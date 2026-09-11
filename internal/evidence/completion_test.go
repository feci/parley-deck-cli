package evidence

import (
	"encoding/json"
	"strings"
	"testing"
)

const transitionDoc = "---\nidea: idea-x\nstatus: implemented\n---\n\n## Summary of work\n\ndone\n"

func TestTransitionFrontmatterStatus(t *testing.T) {
	out, from, err := TransitionFrontmatterStatus([]byte(transitionDoc), "complete")
	if err != nil || from != "implemented" {
		t.Fatalf("from=%q err=%v", from, err)
	}
	want := strings.Replace(transitionDoc, "status: implemented", "status: complete", 1)
	if string(out) != want {
		t.Fatalf("transformation must be byte-exact outside the status line:\n%q\nwant:\n%q", out, want)
	}
	// Exactly invertible: transforming back restores the original bytes.
	back, _, err := TransitionFrontmatterStatus(out, from)
	if err != nil || string(back) != transitionDoc {
		t.Fatalf("inverse must restore the original bytes: err=%v", err)
	}

	cases := []struct {
		name string
		doc  string
		to   string
	}{
		{"no frontmatter", "body only\n", "complete"},
		{"unterminated", "---\nstatus: implemented\n", "complete"},
		{"missing status", "---\nidea: x\n---\n", "complete"},
		{"duplicate status", "---\nstatus: implemented\nstatus: implemented\n---\n", "complete"},
		{"empty value", "---\nstatus:\n---\n", "complete"},
		{"unnormalized double space", "---\nstatus:  implemented\n---\n", "complete"},
		{"unnormalized no space", "---\nstatus:implemented\n---\n", "complete"},
		{"unnormalized leading space", "---\n status: implemented\n---\n", "complete"},
		{"crlf not normalized", "---\r\nstatus: implemented\r\n---\r\n", "complete"},
		{"multi-token value", "---\nstatus: not run\n---\n", "complete"},
		{"comment value", "---\nstatus: implemented # done\n---\n", "complete"},
		{"already complete", "---\nstatus: complete\n---\n", "complete"},
		{"empty target", transitionDoc, ""},
		{"spaced target", transitionDoc, "not run"},
		{"colon target", transitionDoc, "a:b"},
		{"hash target", transitionDoc, "a#b"},
	}
	for _, tc := range cases {
		if _, _, err := TransitionFrontmatterStatus([]byte(tc.doc), tc.to); err == nil {
			t.Fatalf("%s: must fail closed", tc.name)
		}
	}
}

// attestedTransitionReport returns a report with two independently attested
// records plus an ExtraDigests binding for doc, ready for authorization.
func attestedTransitionReport(doc, path string) *Report {
	r := positiveReport()
	r.ExtraDigests = map[string]string{path: sha256Hex([]byte(doc))}
	return r
}

func TestAuthorizeCompletionTransition(t *testing.T) {
	const path = "parley-deck/ideas/idea-x/IMPLEMENTATION.md"
	r := attestedTransitionReport(transitionDoc, path)
	if err := AuthorizeCompletionTransition(r, path, []byte(transitionDoc), "codex-1"); err != nil {
		t.Fatalf("the attesting verifier must be able to authorize: %v", err)
	}
	tr := r.CompletionTransition
	if tr == nil || tr.Path != path || tr.FromStatus != "implemented" || tr.ToStatus != "complete" || tr.AuthorizedBy != "codex-1" {
		t.Fatalf("transition record wrong: %+v", tr)
	}
	if tr.BeforeSHA256 != sha256Hex([]byte(transitionDoc)) {
		t.Fatal("before digest must be the actual current content's digest")
	}
	transformed, _, _ := TransitionStatusToComplete([]byte(transitionDoc))
	if tr.AfterSHA256 != sha256Hex(transformed) {
		t.Fatal("after digest must be recomputed by applying the transformation, never caller-supplied")
	}
	// A second authorization on the same recorded report is refused.
	if err := AuthorizeCompletionTransition(r, path, []byte(transitionDoc), "codex-1"); err == nil {
		t.Fatal("a second completion transition must be refused")
	}
}

func TestAuthorizeCompletionTransitionRefusals(t *testing.T) {
	const path = "parley-deck/ideas/idea-x/IMPLEMENTATION.md"
	cases := []struct {
		name string
		mut  func(r *Report) (p string, doc string, verifier string)
	}{
		{"missing verifier", func(r *Report) (string, string, string) { return path, transitionDoc, "" }},
		{"unbound path", func(r *Report) (string, string, string) { return "other/path.md", transitionDoc, "codex-1" }},
		{"drifted content", func(r *Report) (string, string, string) {
			return path, strings.Replace(transitionDoc, "done", "edited", 1), "codex-1"
		}},
		{"zero records", func(r *Report) (string, string, string) { r.Records = nil; return path, transitionDoc, "codex-1" }},
		{"unattested record", func(r *Report) (string, string, string) {
			r.Records[0].Provenance.Verifier = ""
			r.Records[0].Provenance.VerifierRerun = nil
			return path, transitionDoc, "codex-1"
		}},
		{"retained rerun missing", func(r *Report) (string, string, string) {
			r.Records[1].Provenance.VerifierRerun = nil
			return path, transitionDoc, "codex-1"
		}},
		{"not the attesting verifier", func(r *Report) (string, string, string) { return path, transitionDoc, "opencode-1" }},
		{"self authorization", func(r *Report) (string, string, string) { return path, transitionDoc, "kimi-1" }},
		{"malformed document", func(r *Report) (string, string, string) {
			doc := strings.Replace(transitionDoc, "status: implemented", "status:  implemented", 1)
			r.ExtraDigests[path] = sha256Hex([]byte(doc))
			return path, doc, "codex-1"
		}},
	}
	for _, tc := range cases {
		r := attestedTransitionReport(transitionDoc, path)
		p, doc, v := tc.mut(r)
		if err := AuthorizeCompletionTransition(r, p, []byte(doc), v); err == nil {
			t.Fatalf("%s: must be refused", tc.name)
		}
		if r.CompletionTransition != nil {
			t.Fatalf("%s: a refused authorization must not record a transition", tc.name)
		}
	}
	if err := AuthorizeCompletionTransition(nil, path, []byte(transitionDoc), "codex-1"); err == nil {
		t.Fatal("nil report must be refused")
	}
}

func TestVerifyCompletionTransition(t *testing.T) {
	const path = "parley-deck/ideas/idea-x/IMPLEMENTATION.md"
	bound := sha256Hex([]byte(transitionDoc))
	current, _, _ := TransitionStatusToComplete([]byte(transitionDoc))

	authorized := func() *Report {
		r := attestedTransitionReport(transitionDoc, path)
		if err := AuthorizeCompletionTransition(r, path, []byte(transitionDoc), "codex-1"); err != nil {
			t.Fatal(err)
		}
		return r
	}

	if reasons := VerifyCompletionTransition(authorized(), path, bound, current, "codex-1"); len(reasons) != 0 {
		t.Fatalf("the authorized transition must reconcile: %v", reasons)
	}

	cases := []struct {
		name    string
		r       *Report
		path    string
		bound   string
		current []byte
		verif   string
	}{
		{"nil report", nil, path, bound, current, "codex-1"},
		{"no transition recorded", attestedTransitionReport(transitionDoc, path), path, bound, current, "codex-1"},
		{"wrong path", authorized(), "other.md", bound, current, "codex-1"},
		{"not the selected verifier", authorized(), path, bound, current, "opencode-1"},
		{"empty selected verifier", authorized(), path, bound, current, ""},
		{"not the recorded state", authorized(), path, sha256Hex([]byte("other")), current, "codex-1"},
		{"current not the authorized state", authorized(), path, bound, []byte(transitionDoc), "codex-1"},
		{"extra edit beyond the flip", authorized(), path, bound, []byte(strings.Replace(string(current), "done", "edited", 1)), "codex-1"},
		{"malformed current doc", authorized(), path, bound, []byte("no frontmatter\n"), "codex-1"},
	}
	for _, tc := range cases {
		if reasons := VerifyCompletionTransition(tc.r, tc.path, tc.bound, tc.current, tc.verif); len(reasons) == 0 {
			t.Fatalf("%s: must not reconcile", tc.name)
		}
	}

	// Tampered records each fail a recomputation.
	tamper := func(mut func(tr *CompletionTransition)) *Report {
		r := authorized()
		mut(r.CompletionTransition)
		return r
	}
	for name, r := range map[string]*Report{
		"tampered target":       tamper(func(tr *CompletionTransition) { tr.ToStatus = "final" }),
		"tampered source":       tamper(func(tr *CompletionTransition) { tr.FromStatus = "final" }),
		"empty source":          tamper(func(tr *CompletionTransition) { tr.FromStatus = "" }),
		"tampered authorizer":   tamper(func(tr *CompletionTransition) { tr.AuthorizedBy = "kimi-1" }),
		"tampered before digest": tamper(func(tr *CompletionTransition) { tr.BeforeSHA256 = sha256Hex([]byte("x")) }),
		"tampered after digest":  tamper(func(tr *CompletionTransition) { tr.AfterSHA256 = sha256Hex([]byte("x")) }),
	} {
		if reasons := VerifyCompletionTransition(r, path, bound, current, "codex-1"); len(reasons) == 0 {
			t.Fatalf("%s: must not reconcile", name)
		}
	}
}

func TestCompletionTransitionJSONRoundTrip(t *testing.T) {
	// Old reports without a transition serialize without the field and load
	// with a nil transition — never silently upgraded.
	r := positiveReport()
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "completion_transition") {
		t.Fatal("a report without a transition must not serialize the field")
	}
	var old Report
	if err := json.Unmarshal(data, &old); err != nil || old.CompletionTransition != nil {
		t.Fatalf("old report must load with a nil transition: %v", err)
	}

	const path = "p/IMPLEMENTATION.md"
	ar := attestedTransitionReport(transitionDoc, path)
	if err := AuthorizeCompletionTransition(ar, path, []byte(transitionDoc), "codex-1"); err != nil {
		t.Fatal(err)
	}
	data, err = json.Marshal(ar)
	if err != nil {
		t.Fatal(err)
	}
	var loaded Report
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.CompletionTransition == nil || *loaded.CompletionTransition != *ar.CompletionTransition {
		t.Fatalf("transition must round-trip: %+v", loaded.CompletionTransition)
	}
}
