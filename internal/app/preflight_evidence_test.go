package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadinessRejectsUntrustedOrErrorPONG(t *testing.T) {
	for name, output := range map[string]string{
		"user-echo":                  `{"role":"user","content":"PONG"}`,
		"unrecognized-object":        `{"answer":"PONG"}`,
		"error-with-success-subtype": `{"type":"result","is_error":true,"subtype":"success","result":"PONG"}`,
		"wrapped-user-echo":          `{"data":{"role":"user","content":"PONG"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			obs := classifyReadiness(output, "", 0, false, false, "")
			if obs.Ready || obs.Class == ClassReady {
				t.Fatalf("non-assistant/error envelope became ready: %+v", obs)
			}
		})
	}
}

func TestReadinessSchemaAndCaptureCounterexamples(t *testing.T) {
	for name, output := range map[string]string{
		"bare-content":           `{"content":"PONG"}`,
		"bare-result":            `{"result":"PONG"}`,
		"bare-message":           `{"message":"PONG"}`,
		"bare-text":              `{"text":"PONG"}`,
		"unattributed-wrappers":  `{"result":{"data":{"content":"PONG"}}}`,
		"escaped-duplicate":      `{"role":"user","r\u006fle":"assistant","content":"PONG"}`,
		"sole-role-alias":        `{"Role":"assistant","type":"result","result":"PONG"}`,
		"sole-error-alias":       `{"role":"assistant","isError":true,"content":"PONG"}`,
		"duplicate-wrapper":      `{"data":{"role":"user","content":"PONG"},"data":{"role":"assistant","content":"PONG"}}`,
		"null-role":              `{"role":null,"type":"result","result":"PONG"}`,
		"string-error":           `{"role":"assistant","is_error":"true","content":"PONG"}`,
		"nested-string-error":    `{"data":{"role":"assistant","is_error":"true","content":"PONG"}}`,
		"deep-error":             `{"data":{"payload":{"role":"assistant","is_error":true,"content":"PONG"}}}`,
		"unknown-type-with-role": `{"role":"assistant","type":"unknown","content":"PONG"}`,
		"competing-wrappers":     `{"data":{"role":"assistant","content":"PONG"},"response":{"role":"assistant","content":"FAIL"}}`,
		"competing-content":      `{"role":"assistant","content":"PONG","text":"FAIL"}`,
	} {
		t.Run(name, func(t *testing.T) {
			obs := classifyReadiness(output, "", 0, false, false, "")
			if obs.Ready || obs.Class == ClassReady || obs.Class == "" {
				t.Fatalf("invalid readiness observation: %+v", obs)
			}
		})
	}
	for _, output := range []string{`{"role":"assistant","content":"PONG"}`, `{"type":"result","is_error":false,"subtype":"success","result":"PONG"}`, `{"data":{"role":"assistant","content":"PONG"}}`} {
		if obs := classifyReadiness(output, "", 0, false, false, ""); !obs.Ready {
			t.Fatalf("recognized positive refused: %+v", obs)
		}
	}
	for _, tc := range []struct {
		text      string
		truncated bool
	}{{" \n", false}, {"", true}, {" \t", true}} {
		obs := classifyReadiness(tc.text, "", -1, true, tc.truncated, "overflow")
		if obs.Class != ClassDeadlineAfterOutput || obs.Ready {
			t.Fatalf("lost actual output observation: %+v", obs)
		}
	}
	var buf bytes.Buffer
	truncated := false
	w := boundedWriter{buf: &buf, max: 4, truncated: &truncated}
	_, _ = w.Write([]byte("PONG"))
	_, _ = w.Write(nil)
	if truncated {
		t.Fatal("empty write at cap invented truncation")
	}
	_, _ = w.Write([]byte("!"))
	if !truncated {
		t.Fatal("dropped bytes were not observed")
	}
}

func TestReadinessEscapedCredentialDiagnostics(t *testing.T) {
	for _, raw := range []string{
		`{"token":"SYNTHETIC_SECRET_VALUE"}`,
		`{"to\u006ben":"SYNTHETIC_SECRET_VALUE"}`,
		`{"data":{"access_token":["SYNTHETIC_SECRET_VALUE"]}}`,
		`{"token":"prefix\"SYNTHETIC_SECRET_VALUE"}`,
		`{"to\u006ben":"SYNTHETIC_SECRET_VALUE`,
		`token="SYNTHETIC_SECRET_VALUE"`,
	} {
		if tail := sanitizeTail(raw); strings.Contains(tail, "SYNTHETIC_SECRET_VALUE") {
			t.Fatalf("synthetic credential leaked: %s", tail)
		}
	}
}
