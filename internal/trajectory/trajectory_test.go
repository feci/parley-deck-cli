package trajectory

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"parley-deck-cli/internal/evidence"
)

func copyTrail(t *testing.T, r []Request, o []Observation) ([]Request, []Observation) {
	t.Helper()
	var out struct {
		R []Request
		O []Observation
	}
	data, err := json.Marshal(struct {
		R []Request
		O []Observation
	}{r, o})
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out.R, out.O
}

func TestTrajectoryRejectsCoverageGapsReplayAndForgedBindings(t *testing.T) {
	requests, observations := fixtureTrail(t, [][2]string{{"good", "good"}, {"bad", "good"}, {"good", "bad"}})
	for _, tc := range []struct {
		name   string
		change func([]Request, []Observation) ([]Request, []Observation)
	}{
		{"missing-observation", func(r []Request, o []Observation) ([]Request, []Observation) { return r, o[:1] }},
		{"omitted-first-patch", func(r []Request, o []Observation) ([]Request, []Observation) { return r[1:], o[1:] }},
		{"sequence-gap", func(r []Request, o []Observation) ([]Request, []Observation) { r[1].Sequence = 3; return r, o }},
		{"ancestry-change", func(r []Request, o []Observation) ([]Request, []Observation) {
			r[1].PreviousSHA256 = strings.Repeat("0", 64)
			return r, o
		}},
		{"replayed-id", func(r []Request, o []Observation) ([]Request, []Observation) { r[1].ID = r[0].ID; return r, o }},
		{"replayed-patch", func(r []Request, o []Observation) ([]Request, []Observation) {
			r[1].Before, r[1].After = r[0].Before, r[0].After
			return r, o
		}},
		{"scope-change", func(r []Request, o []Observation) ([]Request, []Observation) {
			r[1].Criteria[0].Name = "new materiality"
			return r, o
		}},
		{"idea-change", func(r []Request, o []Observation) ([]Request, []Observation) { r[1].Idea = "unrelated"; return r, o }},
		{"self-request", func(r []Request, o []Observation) ([]Request, []Observation) {
			r[0].Verifier = r[0].Implementer
			return r, o
		}},
		{"wrong-observer", func(r []Request, o []Observation) ([]Request, []Observation) { o[0].Verifier = "another"; return r, o }},
		{"self-execution", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].Record.Provenance.Executor = r[0].Implementer
			return r, o
		}},
		{"completion-attestation-reuse", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].Before[0].Record.Provenance.Verifier = r[0].Verifier
			return r, o
		}},
		{"missing-pair", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs = o[0].Pairs[:1]
			return r, o
		}},
		{"incomplete-capture", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].Complete = false
			return r, o
		}},
		{"changed-command", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].Record.Command.CommandSHA256 = strings.Repeat("0", 64)
			return r, o
		}},
		{"changed-tree", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].TreeAfterSHA256 = strings.Repeat("0", 64)
			return r, o
		}},
		{"invalid-count", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].Record.Command.FailedCases = -2
			return r, o
		}},
		{"missing-output", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].Record.Command.OutputSHA256 = ""
			return r, o
		}},
		{"interrupted-process", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].Record.Command.ExitCode = -1
			return r, o
		}},
		{"contradictory-pass", func(r []Request, o []Observation) ([]Request, []Observation) {
			o[0].Pairs[0].After[0].Record.Status = evidence.StatusPass
			return r, o
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, o := copyTrail(t, requests, observations)
			r, o = tc.change(r, o)
			if d, err := Evaluate(r, o); err == nil {
				t.Fatalf("invalid %s admitted: %+v", tc.name, d)
			}
		})
	}
}

func TestReducedPassingCoverageCannotClearARegressionSequence(t *testing.T) {
	r, o := fixtureTrail(t, [][2]string{{"good", "good"}, {"good", "good"}})
	for i := 0; i < 2; i++ {
		o[0].Pairs[0].Before[i].Record.Command.ExecutedCases = 2
	}
	a, err := Assess(r[0], o[0])
	if err != nil {
		t.Fatal(err)
	}
	if a.Outcome != Inconclusive {
		t.Fatalf("reduced coverage was certified clean: %+v", a)
	}
}

func TestFreezeAndVerifyRequireExactIndependentInputs(t *testing.T) {
	roots := fixtureSnapshots(t, [][2]string{{"good", "good"}, {"bad", "good"}})
	base := Request{Idea: "fixture", ID: "patch-1", Sequence: 1, Implementer: "author", Verifier: "verifier"}
	criteria := fixtureCriteria()
	if _, err := Freeze(context.Background(), base, roots[0], roots[0], criteria); err == nil {
		t.Fatal("same snapshot admitted")
	}
	bad := base
	bad.Verifier = base.Implementer
	if _, err := Freeze(context.Background(), bad, roots[0], roots[1], criteria); err == nil {
		t.Fatal("self verification admitted")
	}
	if _, err := Freeze(context.Background(), base, roots[0], roots[1], []Criterion{criteria[0], criteria[0]}); err == nil {
		t.Fatal("duplicate acceptance criterion admitted")
	}
	r, err := Freeze(context.Background(), base, roots[0], roots[1], criteria)
	if err != nil {
		t.Fatal(err)
	}
	changed := append([]Criterion(nil), criteria...)
	changed[0].Command = "true"
	if _, err := Verify(context.Background(), r, roots[0], roots[1], changed); err == nil {
		t.Fatal("changed command admitted")
	}
	if _, err := Verify(context.Background(), r, roots[0], roots[1], criteria[:1]); err == nil {
		t.Fatal("partial scope admitted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Verify(ctx, r, roots[0], roots[1], criteria); err == nil {
		t.Fatal("cancelled verification admitted")
	}
}
