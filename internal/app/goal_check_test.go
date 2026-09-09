package app

import "testing"

// LE-7: parseGoalVerdict extracts PASS/FAIL/ambiguous from a goal-check answer
// (case-insensitive, tolerant of leading markdown markers, fail-closed aggregation).
func TestParseGoalVerdict(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"rationale here\nGOAL-CHECK: PASS\n", "PASS"},
		{"GOAL-CHECK: FAIL — criterion 3 unmet", "FAIL"},
		{"- GOAL-CHECK: pass", "PASS"},                          // marker + lowercase
		{"## GOAL-CHECK: Fail", "FAIL"},                         // heading marker
		{"no verdict line at all", ""},                          // ambiguous → ""
		{"GOAL-CHECK: PASS\nGOAL-CHECK: FAIL\n", "FAIL"},        // FAIL is sticky
		{"`GOAL-CHECK: FAIL`", "FAIL"},                          // CF2: backtick-wrapped token
		{`"GOAL-CHECK: PASS"`, "PASS"},                          // CF2: quote-wrapped token
		{"**GOAL-CHECK:** FAIL", "FAIL"},                        // CF2: bolded marker prefix
		{"GOAL-CHECK: FAIL\nGOAL-CHECK: RE-EVALUATING", "FAIL"}, // FAIL is sticky
		{"GOAL-CHECK: FAIL\nRequired format:\nGOAL-CHECK: PASS", "FAIL"},
		{"GOAL-CHECK: RE-EVALUATING\nGOAL-CHECK: PASS", ""},
		{"GOAL-CHECK: PASS — reservations remain\nGOAL-CHECK: PASS", ""},
		{"GOAL-CHECK: PASS\nGOAL-CHECK: PASS", "PASS"},
		{"GOAL-CHECK: PASSING", ""},
		{"GOAL-CHECK: PASS — reservations remain", ""},
		{"GOAL-CHECK: PASS\nGOAL-CHECK: PASS?", ""},
	}
	for _, c := range cases {
		if got := parseGoalVerdict(c.in); got != c.want {
			t.Fatalf("parseGoalVerdict(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
