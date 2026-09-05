package app

import "testing"

func TestReadinessRejectsUntrustedOrErrorPONG(t *testing.T) {
	for name, output := range map[string]string{
		"user-echo":                  `{"role":"user","content":"PONG"}`,
		"unrecognized-object":        `{"answer":"PONG"}`,
		"error-with-success-subtype": `{"type":"result","is_error":true,"subtype":"success","result":"PONG"}`,
		"wrapped-user-echo":          `{"data":{"role":"user","content":"PONG"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			obs := classifyReadiness(output, "", 0, false)
			if obs.Ready || obs.Class == ClassReady {
				t.Fatalf("non-assistant/error envelope became ready: %+v", obs)
			}
		})
	}
}
