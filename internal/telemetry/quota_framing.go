package telemetry

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Node's SDK error dump is parsed as a bounded grammar, never searched for
// interesting fields. Unknown prose, errors, truncated inspect values, duplicate
// fields and incomplete containers invalidate the entire stderr capture.
type sdkRecord struct {
	body, header, message string
	end                   int
}
type sdkDump struct {
	raw         string
	pos         int
	records     []sdkRecord
	fingerprint string
}
type sdkUndefined struct{}

var sdkHeader = regexp.MustCompile(`^(APICallError \[AI_APICallError\]|AI_APICallError|RetryError \[AI_RetryError\]|AI_RetryError): (.+)$`)
var sdkStack = regexp.MustCompile(`^at (?:async )?(?:[^\s()]+(?: \[as [^\s\]]+\])? \([^\r\n()]+:[0-9]+:[0-9]+\)|[^\s()]+:[0-9]+:[0-9]+)(?: \{)?$`)
var sdkAttempts = regexp.MustCompile(`^Failed after ([1-9][0-9]*) attempts\. Last error: (.+)$`)

func parseSDKCapture(raw string) ([]sdkRecord, error) {
	p := &sdkDump{raw: strings.TrimRight(raw, "\r\n")}
	for {
		p.space()
		if p.pos >= len(p.raw) {
			return nil, fmt.Errorf("missing terminal turn-failure line")
		}
		if zcodeTerminal.MatchString(p.raw[p.pos:]) {
			if len(p.records) == 0 {
				return nil, fmt.Errorf("missing provider record")
			}
			return p.records, nil
		}
		if _, err := p.failure(0); err != nil {
			return nil, err
		}
	}
}
func (p *sdkDump) space() {
	for p.pos < len(p.raw) && unicode.IsSpace(rune(p.raw[p.pos])) {
		p.pos++
	}
}
func (p *sdkDump) line() string {
	start := p.pos
	for p.pos < len(p.raw) && p.raw[p.pos] != '\n' {
		p.pos++
	}
	s := strings.TrimSpace(p.raw[start:p.pos])
	if p.pos < len(p.raw) {
		p.pos++
	}
	return s
}
func (p *sdkDump) failure(depth int) (string, error) {
	if depth > 8 {
		return "", fmt.Errorf("SDK nesting limit")
	}
	p.space()
	line := p.line()
	m := sdkHeader.FindStringSubmatch(line)
	if m == nil {
		return "", fmt.Errorf("unknown SDK error header")
	}
	retry := strings.Contains(m[1], "RetryError")
	message := m[2]
	opened := strings.HasSuffix(message, " {")
	if opened {
		message = strings.TrimSuffix(message, " {")
	}
	for !opened {
		if p.pos >= len(p.raw) {
			return "", fmt.Errorf("incomplete SDK stack")
		}
		line = p.line()
		if !sdkStack.MatchString(line) {
			return "", fmt.Errorf("unknown SDK stack framing")
		}
		opened = strings.HasSuffix(line, " {")
	}
	fields := map[string]any{}
	bodies := 0
	bodyEnd := 0
	retryMessages := []string{}
	last := ""
	var errorFingerprints []string
	lastFingerprint := ""
	for {
		p.space()
		if p.pos >= len(p.raw) {
			return "", fmt.Errorf("incomplete SDK object")
		}
		if p.raw[p.pos] == '}' {
			p.pos++
			break
		}
		key, err := p.key()
		if err != nil {
			return "", err
		}
		if _, ok := fields[key]; ok {
			return "", fmt.Errorf("duplicate SDK field")
		}
		fields[key] = nil
		p.space()
		switch {
		case retry && key == "errors":
			if !p.take('[') {
				return "", fmt.Errorf("missing retry list")
			}
			for {
				p.space()
				if p.take(']') {
					break
				}
				msg, e := p.failure(depth + 1)
				if e != nil {
					return "", e
				}
				retryMessages = append(retryMessages, msg)
				errorFingerprints = append(errorFingerprints, p.fingerprint)
				p.space()
				if p.take(',') {
					continue
				}
				if !p.take(']') {
					return "", fmt.Errorf("unfinished retry list")
				}
				break
			}
			fields[key] = errorFingerprints
		case retry && key == "lastError":
			recordStart := len(p.records)
			last, err = p.failure(depth + 1)
			if err != nil {
				return "", err
			}
			// lastError is a fingerprint-checked duplicate, not a later observation.
			p.records = p.records[:recordStart]
			lastFingerprint = p.fingerprint
			fields[key] = lastFingerprint
		default:
			allowed := map[string]bool{"cause": true, "Symbol(vercel.ai.error)": true}
			if retry {
				allowed["reason"] = true
				allowed["Symbol(vercel.ai.error.AI_RetryError)"] = true
			} else {
				for _, k := range []string{"url", "requestBodyValues", "statusCode", "responseHeaders", "responseBody", "isRetryable", "data", "Symbol(vercel.ai.error.AI_APICallError)"} {
					allowed[k] = true
				}
			}
			if !allowed[key] {
				return "", fmt.Errorf("unknown SDK field %s", key)
			}
			v, e := p.value(depth + 1)
			if e != nil {
				return "", e
			}
			fields[key] = v
			if key == "responseBody" {
				bodies++
				bodyEnd = p.pos
			}
		}
		p.space()
		if p.take(',') {
			continue
		}
		if !p.take('}') {
			return "", fmt.Errorf("missing SDK delimiter")
		}
		break
	}
	if fields["Symbol(vercel.ai.error)"] != true {
		return "", fmt.Errorf("missing SDK identity")
	}
	// A nested cause cannot be ignored, even if the outer response is quota.
	if cause, exists := fields["cause"]; !exists || cause != (sdkUndefined{}) && cause != nil {
		return "", fmt.Errorf("non-provider cause")
	}
	if retry {
		match := sdkAttempts.FindStringSubmatch(message)
		if match == nil || fields["reason"] != "maxRetriesExceeded" || fields["Symbol(vercel.ai.error.AI_RetryError)"] != true || len(retryMessages) == 0 || last == "" {
			return "", fmt.Errorf("incomplete retry wrapper")
		}
		n, _ := strconv.Atoi(match[1])
		if n != len(retryMessages) || last != match[2] || len(errorFingerprints) == 0 || lastFingerprint != errorFingerprints[len(errorFingerprints)-1] {
			return "", fmt.Errorf("contradictory retry wrapper")
		}
		encoded, _ := json.Marshal(fields)
		p.fingerprint = message + "\n" + string(encoded)
		return last, nil
	}
	if bodies != 1 || fields["statusCode"] != float64(429) || fields["isRetryable"] != true || fields["Symbol(vercel.ai.error.AI_APICallError)"] != true {
		return "", fmt.Errorf("incomplete provider record")
	}
	body, ok := fields["responseBody"].(string)
	if !ok {
		return "", fmt.Errorf("missing provider body")
	}
	var envelope map[string]any
	if !uniqueJSONKeys([]byte(body)) || json.Unmarshal([]byte(body), &envelope) != nil {
		return "", fmt.Errorf("malformed provider body")
	}
	er, ok := envelope["error"].(map[string]any)
	if !ok || er["message"] != message {
		return "", fmt.Errorf("provider header/body mismatch")
	}
	if data, ok := fields["data"]; !ok || data != nil && !reflect.DeepEqual(data, sdkUndefined{}) && !reflect.DeepEqual(data, envelope) {
		return "", fmt.Errorf("mixed SDK data")
	}
	headers, ok := fields["responseHeaders"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("incomplete response headers")
	}
	header := ""
	for k, v := range headers {
		if strings.EqualFold(k, "retry-after") {
			if header != "" {
				return "", fmt.Errorf("duplicate retry header")
			}
			s, ok := v.(string)
			if !ok {
				return "", fmt.Errorf("invalid retry header")
			}
			header = s
		} else {
			if _, ok := v.(string); !ok {
				return "", fmt.Errorf("invalid response header")
			}
		}
	}
	if _, ok := fields["url"].(string); !ok {
		return "", fmt.Errorf("missing SDK URL")
	}
	if _, ok := fields["requestBodyValues"]; !ok {
		return "", fmt.Errorf("missing SDK request framing")
	}
	p.records = append(p.records, sdkRecord{body, header, message, bodyEnd})
	encoded, _ := json.Marshal(fields)
	p.fingerprint = message + "\n" + string(encoded)
	return message, nil
}
func (p *sdkDump) take(c byte) bool {
	p.space()
	if p.pos < len(p.raw) && p.raw[p.pos] == c {
		p.pos++
		return true
	}
	return false
}
func (p *sdkDump) key() (string, error) {
	p.space()
	if p.pos >= len(p.raw) {
		return "", fmt.Errorf("missing key")
	}
	var key string
	if strings.ContainsRune("'\"`", rune(p.raw[p.pos])) {
		v, e := p.quoted()
		if e != nil {
			return "", e
		}
		key = v
	} else {
		start := p.pos
		for p.pos < len(p.raw) && p.raw[p.pos] != ':' && p.raw[p.pos] != '\n' {
			p.pos++
		}
		key = strings.TrimSpace(p.raw[start:p.pos])
		key = strings.TrimPrefix(strings.TrimSuffix(key, "]"), "[")
		if key == "" {
			return "", fmt.Errorf("empty key")
		}
	}
	if !p.take(':') {
		return "", fmt.Errorf("missing key delimiter")
	}
	return key, nil
}
func (p *sdkDump) quoted() (string, error) {
	q := p.raw[p.pos]
	p.pos++
	var b strings.Builder
	for p.pos < len(p.raw) {
		c := p.raw[p.pos]
		p.pos++
		if c == q {
			return b.String(), nil
		}
		if c == '\n' || c == '\r' {
			return "", fmt.Errorf("multiline JS string")
		}
		if c != '\\' {
			b.WriteByte(c)
			continue
		}
		if p.pos >= len(p.raw) {
			break
		}
		e := p.raw[p.pos]
		p.pos++
		switch e {
		case '\\', '\'', '"', '`':
			b.WriteByte(e)
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		default:
			return "", fmt.Errorf("unsupported JS escape")
		}
	}
	return "", fmt.Errorf("unterminated JS string")
}
func (p *sdkDump) value(depth int) (any, error) {
	p.space()
	if depth > 16 || p.pos >= len(p.raw) {
		return nil, fmt.Errorf("incomplete SDK value")
	}
	switch p.raw[p.pos] {
	case '\'', '"', '`':
		return p.quoted()
	case '{':
		p.pos++
		m := map[string]any{}
		for {
			if p.take('}') {
				return m, nil
			}
			k, e := p.key()
			if e != nil {
				return nil, e
			}
			if _, ok := m[k]; ok {
				return nil, fmt.Errorf("duplicate SDK value key")
			}
			v, e := p.value(depth + 1)
			if e != nil {
				return nil, e
			}
			m[k] = v
			if p.take(',') {
				continue
			}
			if !p.take('}') {
				return nil, fmt.Errorf("incomplete SDK map")
			}
			return m, nil
		}
	case '[':
		p.pos++
		a := []any{}
		for {
			if p.take(']') {
				return a, nil
			}
			v, e := p.value(depth + 1)
			if e != nil {
				return nil, e
			}
			a = append(a, v)
			if p.take(',') {
				continue
			}
			if !p.take(']') {
				return nil, fmt.Errorf("incomplete SDK array")
			}
			return a, nil
		}
	}
	start := p.pos
	for p.pos < len(p.raw) && !strings.ContainsRune(",]} \t\r\n", rune(p.raw[p.pos])) {
		p.pos++
	}
	s := p.raw[start:p.pos]
	switch s {
	case "undefined":
		return sdkUndefined{}, nil
	case "null":
		return nil, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	n, e := strconv.ParseFloat(s, 64)
	if e == nil {
		return n, nil
	}
	return nil, fmt.Errorf("unknown SDK value")
}
