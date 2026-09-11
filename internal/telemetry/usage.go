package telemetry

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

const parserLimit = 256 * 1024

// Collector interprets only explicitly structured adapter output. Text that a
// model prints in a plain-text transport is not provider usage evidence.
type Collector struct {
	stdoutBytes, stderrBytes int64
	mu                       sync.Mutex
	adapter                  string
	structured               bool
	began                    time.Time
	observation              Observation
	line                     []byte
	tail                     []byte
	dropping                 bool
	usage                    Usage
	failure                  string
	steps                    map[string]Usage
	ambiguousSteps           bool
}

type streamWriter struct {
	collector *Collector
	stream    string
}

func NewCollector(adapter string, structured bool) *Collector {
	return &Collector{adapter: adapter, structured: structured, began: time.Now(), steps: map[string]Usage{}}
}

func (c *Collector) Writer(stream string) io.Writer {
	return streamWriter{collector: c, stream: stream}
}

func (w streamWriter) Write(data []byte) (int, error) {
	c := w.collector
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(data) > 0 && c.observation.FirstActivityMS == nil {
		elapsed := time.Since(c.began).Milliseconds()
		c.observation.FirstActivityMS = &elapsed
	}
	if w.stream == "stderr" {
		c.stderrBytes += int64(len(data))
		return len(data), nil
	}
	c.stdoutBytes += int64(len(data))
	if !c.structured {
		return len(data), nil
	}
	if len(data) >= parserLimit {
		c.tail = append(c.tail[:0], data[len(data)-parserLimit:]...)
	} else {
		if len(c.tail)+len(data) > parserLimit {
			c.tail = c.tail[len(c.tail)+len(data)-parserLimit:]
		}
		c.tail = append(c.tail, data...)
	}
	for _, b := range data {
		if b == '\n' {
			if !c.dropping {
				c.consume(c.line)
			}
			c.line = c.line[:0]
			c.dropping = false
			continue
		}
		if c.dropping {
			continue
		}
		if len(c.line) >= parserLimit {
			c.dropping = true
			c.observation.TruncatedInput = true
			c.line = c.line[:0]
			continue
		}
		c.line = append(c.line, b)
	}
	return len(data), nil
}

func object(value any) map[string]any { result, _ := value.(map[string]any); return result }
func word(value any) string           { result, _ := value.(string); return result }
func count(value any) *int64 {
	number, ok := value.(json.Number)
	if !ok {
		return nil
	}
	valueInt, err := number.Int64()
	if err != nil || valueInt < 0 || valueInt > 1_000_000_000_000 {
		return nil
	}
	return &valueInt
}
func money(value any) *float64 {
	number, ok := value.(json.Number)
	if !ok {
		return nil
	}
	valueFloat, err := number.Float64()
	if err != nil || math.IsNaN(valueFloat) || math.IsInf(valueFloat, 0) || valueFloat < 0 || valueFloat > 1_000_000 {
		return nil
	}
	return &valueFloat
}

func reportedUsage(raw map[string]any, source string) Usage {
	u := Usage{Source: source, CostBasis: "unavailable", Coverage: "reported",
		InputTokens: count(raw["input_tokens"]), OutputTokens: count(raw["output_tokens"]),
		CacheReadTokens: count(raw["cache_read_input_tokens"]), CacheWriteTokens: count(raw["cache_creation_input_tokens"]),
		TotalTokens: count(raw["total_tokens"])}
	if u.InputTokens == nil {
		u.InputTokens = count(raw["prompt_tokens"])
	}
	if u.OutputTokens == nil {
		u.OutputTokens = count(raw["completion_tokens"])
	}
	if u.CacheReadTokens == nil {
		u.CacheReadTokens = count(raw["cached_input_tokens"])
	}
	return u
}

func (c *Collector) consume(data []byte) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || data[0] != '{' || len(data) > parserLimit {
		return
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var event map[string]any
	if decoder.Decode(&event) != nil {
		return
	}
	var trailing any
	if decoder.Decode(&trailing) != io.EOF {
		return
	}
	kind := word(event["type"])
	if event["is_error"] == true || kind == "turn.failed" || kind == "error" {
		c.failure = "provider-error"
		if code := count(event["api_error_status"]); code != nil {
			switch *code {
			case 401, 403:
				c.failure = "auth-error"
			case 429:
				c.failure = "rate-limit"
			}
		}
	}
	switch c.adapter {
	case "claude":
		if kind != "result" {
			return
		}
		if event["is_error"] != true {
			c.failure = ""
		}
		u := reportedUsage(object(event["usage"]), "claude.result")
		u.CostUSD = money(event["total_cost_usd"])
		if u.CostUSD != nil {
			u.CostBasis = "cli-estimate"
		}
		for model := range object(event["modelUsage"]) {
			if safe := safeModel(model); safe != nil {
				u.ReportedModels = append(u.ReportedModels, *safe)
			}
		}
		sort.Strings(u.ReportedModels)
		if len(u.ReportedModels) == 1 {
			u.ReportedModel = String(u.ReportedModels[0])
		}
		c.usage = u
	case "codex":
		if kind != "turn.completed" {
			return
		}
		if event["is_error"] != true {
			c.failure = ""
		}
		c.usage = reportedUsage(object(event["usage"]), "codex.turn-completed")
		c.usage.ReportedModel = safeModel(word(event["model"]))
	case "opencode":
		if kind != "step_finish" {
			return
		}
		part := object(event["part"])
		identity := word(part["id"])
		if identity == "" || len(identity) > 256 || (len(c.steps) >= 4096 && c.steps[identity].Source == "") {
			c.ambiguousSteps = true
			return
		}
		tokens := object(part["tokens"])
		cache := object(tokens["cache"])
		u := Usage{InputTokens: count(tokens["input"]), OutputTokens: count(tokens["output"]),
			CacheReadTokens: count(cache["read"]), CacheWriteTokens: count(cache["write"]), TotalTokens: count(tokens["total"]),
			CostUSD: money(part["cost"]), Source: "opencode.step-finish", CostBasis: "cli-estimate", Coverage: "reported"}
		if prior, exists := c.steps[identity]; exists && !reflect.DeepEqual(prior, u) {
			c.ambiguousSteps = true
		}
		c.steps[identity] = u
	default:
		// A recognized terminal usage envelope is admissible. ACP used/size
		// snapshots deliberately do not match this shape or become billed usage.
		if kind != "result" && kind != "usage" {
			return
		}
		raw := object(event["usage"])
		if raw == nil {
			return
		}
		u := reportedUsage(raw, c.adapter+".reported-usage")
		u.ReportedModel = safeModel(word(event["model"]))
		u.CostUSD = money(event["cost_usd"])
		if u.CostUSD == nil {
			u.CostUSD = money(event["total_cost_usd"])
		}
		if u.CostUSD != nil {
			u.CostBasis = "cli-estimate"
		}
		c.usage = u
	}
}

func (c *Collector) Result() (Usage, Observation, string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.structured {
		if !c.dropping {
			c.consume(c.line)
		}
		c.consume(c.tail)
	}
	u := c.usage
	if c.adapter == "opencode" && c.ambiguousSteps {
		u = Usage{Source: "opencode.step-finish", Coverage: "ambiguous-step-identity"}
	} else if c.adapter == "opencode" && len(c.steps) > 0 {
		u = Usage{Source: "opencode.step-finish-sum", CostBasis: "cli-estimate", Coverage: "reported"}
		keys := make([]string, 0, len(c.steps))
		for key := range c.steps {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		steps := make([]Usage, 0, len(keys))
		for _, key := range keys {
			steps = append(steps, c.steps[key])
		}
		u.InputTokens = sumCounts(steps, func(v Usage) *int64 { return v.InputTokens })
		u.OutputTokens = sumCounts(steps, func(v Usage) *int64 { return v.OutputTokens })
		u.CacheReadTokens = sumCounts(steps, func(v Usage) *int64 { return v.CacheReadTokens })
		u.CacheWriteTokens = sumCounts(steps, func(v Usage) *int64 { return v.CacheWriteTokens })
		u.TotalTokens = sumCounts(steps, func(v Usage) *int64 { return v.TotalTokens })
		cost, known := 0.0, true
		for _, step := range steps {
			if step.CostUSD == nil {
				known = false
			} else {
				cost += *step.CostUSD
			}
		}
		if known {
			u.CostUSD = &cost
		}
	}
	if u.InputTokens == nil && u.OutputTokens == nil && u.CostUSD == nil && u.Coverage == "reported" {
		u.Coverage = "unavailable"
	}
	if u.CostUSD == nil {
		u.CostBasis = "unavailable"
	}
	observation := c.observation
	observation.StdoutBytes = clone(&c.stdoutBytes)
	observation.StderrBytes = clone(&c.stderrBytes)
	observation.FirstActivityMS = clone(observation.FirstActivityMS)
	return CleanUsage(u), observation, c.failure
}

func sumCounts(values []Usage, field func(Usage) *int64) *int64 {
	total := int64(0)
	for _, value := range values {
		amount := field(value)
		if amount == nil || total > 1_000_000_000_000-*amount {
			return nil
		}
		total += *amount
	}
	return &total
}

// StructuredArgs inspects, but never changes, the effective argv template.
func StructuredArgs(args []string) bool {
	for index, arg := range args {
		if arg == "--json" {
			return true
		}
		for _, flag := range []string{"--output-format", "--format"} {
			if arg == flag && index+1 < len(args) && (args[index+1] == "json" || args[index+1] == "stream-json") {
				return true
			}
			if strings.HasPrefix(arg, flag+"=") {
				value := strings.TrimPrefix(arg, flag+"=")
				if value == "json" || value == "stream-json" {
					return true
				}
			}
		}
	}
	return false
}
