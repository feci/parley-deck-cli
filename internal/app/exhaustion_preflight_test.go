package app

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

func TestExhaustionReadinessUnpinnedDefault(t *testing.T) {
	for _, source := range []string{"legacy", "global-default"} {
		for _, native := range []bool{false, true} {
			name := source + "/transient"
			if native {
				name = source + "/exhaustion"
			}
			t.Run(name, func(t *testing.T) {
				p := quota.NewExhaustionPolicy(nil, nil)
				root, _, _ := fixupAppFixture(t, p)
				ids := []string{"alpha", "beta", "gamma"}
				snapshot := &quota.RoleSnapshot{Version: 1, Order: ids, Organizer: "beta"}
				roles := quota.Roles{Snapshot: snapshot, Assignment: quota.RoleAssignment{
					Organizer: "beta", Implementer: "alpha", ImplementerSource: source,
				}}
				failure := "transient"
				var failurePath string
				if native {
					raw, err := os.ReadFile("../telemetry/testdata/quota/zcode-sdk-source-derived.stderr")
					if err != nil {
						t.Fatal(err)
					}
					path := filepath.Join(root, "failure.stderr")
					if err := os.WriteFile(path, raw, 0600); err != nil {
						t.Fatal(err)
					}
					failure, failurePath = "native", path
				}
				var discovered []agents.Discovery
				for _, id := range ids {
					mode := "success"
					if id == "alpha" {
						mode = failure
					}
					a := agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-test.run=^TestExhaustionReadinessHelper$", "--", "--exhaustion-readiness-helper", mode, failurePath}, PromptMode: agents.PromptStdin}, Found: true, Path: os.Args[0]}
					if native && id == "alpha" {
						a.AdapterID = "zcode"
					}
					discovered = append(discovered, a)
				}
				entries := checkRoster(context.Background(), preflightOptions{Root: root, QuotaPolicy: &p, QuotaRoles: roles, PingTimeout: time.Second}, discovered)
				var members []quota.Member
				for _, e := range entries {
					members = append(members, quota.Member{ID: e.RosterID, Usable: e.Available, Evidence: e.QuotaEvidence})
				}
				d := quota.Evaluate(p, ids, members, roles)
				if !d.Applied || !reflect.DeepEqual(d.After, []string{"beta", "gamma"}) {
					t.Fatalf("decision=%+v entries=%+v", d, entries)
				}
				attempts := 0
				for _, record := range dropoutAppRecords(t, root) {
					if record.Metadata.Agent == "alpha" {
						attempts++
					}
				}
				if native {
					if attempts != 1 || d.RoleTransition == nil || d.RoleTransition.After.Implementer != "beta" {
						t.Fatalf("native attempts=%d decision=%+v", attempts, d)
					}
				} else {
					if attempts != 2 || d.RoleTransition != nil || d.Candidates[0].Evidence.Failure == nil {
						t.Fatalf("transient attempts=%d decision=%+v", attempts, d)
					}
					proof := d.Candidates[0].Evidence.Failure.Attempts
					if proof[1].ObservedAt.Sub(proof[0].ObservedAt) < 5*time.Second {
						t.Fatal("transient retry lost its five-second delay")
					}
				}
				k := quota.NewKickoff("idea", "run", p, d.After, &d, time.Now(), snapshot)
				if err := k.Validate(); err != nil {
					t.Fatal(err)
				}
				if source == "global-default" && !native {
					setCentralDefaults(t, "default_implementer = \"alpha\"\n")
					idea, _, err := protocol.CreateIdeaWithQuota(root, "after default dropout", d.After, nil, "deliberation", "", "next-run", &p, &d, snapshot)
					if err != nil {
						t.Fatal(err)
					}
					ops := designationOps(t, root, idea.Slug, idea.Path, d.After, d.After, io.Discard)
					if got := ops.implementer; got != "beta" || ops.roleErr != "" || ops.dispatchErr != "" {
						t.Fatalf("default did not fall through after dropout: %q", got)
					}
				}
			})
		}
	}
}

func TestExhaustionReadinessHelper(t *testing.T) {
	i := slices.Index(os.Args, "--exhaustion-readiness-helper")
	if i < 0 {
		return
	}
	io.Copy(io.Discard, os.Stdin)
	switch os.Args[i+1] {
	case "success":
		os.Stdout.WriteString("PONG")
		os.Exit(0)
	case "transient":
		os.Stderr.WriteString("HTTP 503")
		os.Exit(9)
	case "native":
		raw, err := os.ReadFile(os.Args[i+2])
		if err != nil {
			t.Fatal(err)
		}
		os.Stderr.Write(raw)
		os.Exit(7)
	default:
		t.Fatal("unknown readiness fixture mode")
	}
}
