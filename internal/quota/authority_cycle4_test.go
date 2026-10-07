package quota_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
)

func TestQuotaCycle4AuthorityWorkingCopiesOnlyGateBinding(t *testing.T) {
	for _, archived := range []bool{false, true} {
		t.Run(map[bool]string{false: "live", true: "archived"}[archived], func(t *testing.T) {
			root := t.TempDir()
			a := quotatest.Authority(t, root, "idea", "answer", "Owner decision.")
			path := filepath.Join(root, a.Path)
			if archived {
				dst := filepath.Join(filepath.Dir(path), "archived", filepath.Base(path))
				if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(path, dst); err != nil {
					t.Fatal(err)
				}
				path = dst
			}
			for _, edit := range []string{"\n## Owner update\nA later answer.\n", "Arbitrary replacement owned by the recipient.\n"} {
				if err := os.WriteFile(path, []byte(edit), 0600); err != nil {
					t.Fatal(err)
				}
				if err := quota.ValidateAuthority(root, "idea", a); err != nil {
					t.Fatal("bound object invalidated", err)
				}
				if _, err := quota.BindAuthority(root, "idea", a.Path, a.Commit, a.Quote); err == nil || !strings.Contains(err.Error(), "working copy changed") {
					t.Fatal("superseded answer newly bound", err)
				}
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := quota.ValidateAuthority(root, "idea", a); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestQuotaCycle4AuthorityImmutableEvidenceStillGates(t *testing.T) {
	for _, field := range []string{"missing-commit", "blob", "digest", "quote", "missing-blob"} {
		t.Run(field, func(t *testing.T) {
			root := t.TempDir()
			a := quotatest.Authority(t, root, "idea", "answer", "Owner decision.")
			switch field {
			case "missing-commit":
				a.Commit = strings.Repeat("1", 40)
			case "blob":
				a.Blob = strings.Repeat("1", 40)
			case "digest":
				a.SHA256 = strings.Repeat("1", 64)
			case "quote":
				a.Quote = "Never authorized"
			case "missing-blob":
				if err := os.Remove(filepath.Join(root, ".git/objects", a.Blob[:2], a.Blob[2:])); err != nil {
					t.Fatal(err)
				}
			}
			if err := quota.ValidateAuthority(root, "idea", a); err == nil {
				t.Fatal("invalid immutable evidence accepted", field)
			}
		})
	}
}
