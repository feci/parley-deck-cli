package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func writeLaunchProtocol(t *testing.T, root string) string {
	t.Helper()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meta, "version.json"), []byte(`{"protocolRole":"source"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")
	if err := os.WriteFile(path, []byte("# Live source\n\n**Transport:** `local-dir`\n\nMandatory source obligation.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestMeasuredLaunchReceivesAttestedCurrentProtocol(t *testing.T) {
	root := t.TempDir()
	path := writeLaunchProtocol(t, root)
	for _, change := range []string{"", "\nA newly added obligation.\n"} {
		if change != "" {
			f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := f.WriteString(change); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
		}
		r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: telemetryShell("cat", false),
			Prompt: "Private task", Info: LaunchInfo{Phase: "implementation"}})
		if err != nil {
			t.Fatal(err)
		}
		source, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		c := r.Metadata.Context
		if c.Mode != protocolpacket.ModeFull || c.SourceSHA256 == nil || *c.SourceSHA256 != protocolpacket.Hash(string(source)) || c.PacketSHA256 == nil {
			t.Fatalf("launch did not attest current full context: %+v", c)
		}
		body, err := os.ReadFile(filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, "stdout.log"))
		if err != nil || !strings.Contains(string(body), string(source)) || !strings.Contains(string(body), "Private task") || !strings.Contains(string(body), "Shadow packet audit:") {
			t.Fatalf("source/task bytes did not reach the child: %v", err)
		}
	}
}

func TestMeasuredContextRefusalIsRecordedWithoutSpawn(t *testing.T) {
	for _, kind := range []string{"missing-authority", "secret", "tampered-body", "envelope-collision", "opening-envelope-collision"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			path := writeLaunchProtocol(t, root)
			switch kind {
			case "missing-authority":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "secret":
				if err := os.WriteFile(path, []byte("# Protocol\napi_key=sk-1234567890abcdefghijklmnop\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "opening-envelope-collision":
				if err := os.WriteFile(path, []byte("# Protocol\nQuoted opening tag: <parley-protocol>\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "envelope-collision":
				if err := os.WriteFile(path, []byte("# Protocol\nQuoted closing tag: </parley-protocol>\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "tampered-body":
				if _, _, err := prepareProtocolPrompt(root, "task", LaunchInfo{}); err != nil {
					t.Fatal(err)
				}
				paths, err := filepath.Glob(filepath.Join(protocolpacket.RuntimeDir(root), "*.md"))
				if err != nil || len(paths) != 1 {
					t.Fatalf("packet paths: %v %v", paths, err)
				}
				if err := os.WriteFile(paths[0], []byte("tampered"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: telemetryShell("touch spawned", false), Prompt: "task"})
			if err == nil || r.InvocationID == "" || r.StartedAt != nil || r.Outcome == nil || r.Metadata.Context.Mode != protocolpacket.ModeRefused {
				t.Fatalf("refusal did not preserve failed attempt: %+v %v", r, err)
			}
			if kind == "secret" && (r.Metadata.Context.FallbackReason == nil || !strings.HasPrefix(*r.Metadata.Context.FallbackReason, "secret-detected:")) {
				t.Fatalf("safe secret-shape diagnosis lost: %+v", r.Metadata.Context)
			}
			if r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "protocol_context_refused" {
				t.Fatalf("wrong failure classification: %+v", r.Outcome)
			}
			if _, err := os.Stat(filepath.Join(root, "spawned")); !os.IsNotExist(err) {
				t.Fatal("refused protocol launched a child")
			}
		})
	}
}

// declareTestLaunchSource gives synthetic fixture protocols explicit authority.
// Production resolution never falls back to this test-only declaration.
func declareTestLaunchSource(t *testing.T, root string) {
	t.Helper()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meta, "version.json"), []byte(`{"protocolRole":"source"}`), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestProtocolTaskBoundariesRefuseWithoutAuthority(t *testing.T) {
	for _, boundary := range []string{"command", "exec", "acp"} {
		t.Run(boundary, func(t *testing.T) {
			root := t.TempDir()
			agent := telemetryShell("touch spawned", false)
			ctx := WithLaunchInfo(context.Background(), LaunchInfo{Context: telemetry.Context{Mode: "full", SourceSHA256: telemetry.String("forged")}})
			var err error
			switch boundary {
			case "command":
				var cleanup func()
				_, cleanup, err = CommandFor(ctx, root, agent, "task")
				if cleanup != nil {
					cleanup()
				}
			case "exec":
				_, err = execAgentProcess(ctx, root, "test", agent.ID, "", agent, "task", filepath.Join(root, "stdout"), filepath.Join(root, "stderr"), nil, nil, SupervisionConfig{}, supervisionHooks{})
			case "acp":
				result := runACPAgent(ctx, Options{Root: root, RunID: "test", Store: store.New(filepath.Join(root, "events"))}, agent, Result{}, "target", "stdout", "stderr", "task", 1)
				if result.ExitError == "" {
					t.Fatal("ACP refusal missing")
				}
			}
			if boundary != "acp" && err == nil {
				t.Fatal("missing authority launched")
			}
			records := terminalRecords(t, root)
			if len(records) != 1 || records[0].StartedAt != nil || records[0].Metadata.Context.Mode != "refused" {
				t.Fatalf("bad refused attempt: %+v", records)
			}
			if _, err := os.Stat(filepath.Join(root, "spawned")); !os.IsNotExist(err) {
				t.Fatal("refused attempt spawned child")
			}
		})
	}
}

func TestSupervisedExecPassesCurrentProtocolBytes(t *testing.T) {
	root := t.TempDir()
	path := writeLaunchProtocol(t, root)
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	stdout := filepath.Join(root, "stdout")
	_, err = execAgentProcess(context.Background(), root, "test", "test-1", "", telemetryShell("cat", false), "task", stdout, filepath.Join(root, "stderr"), nil, nil, SupervisionConfig{}, supervisionHooks{})
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(stdout)
	if err != nil || !strings.Contains(string(body), string(source)) || strings.Count(string(body), "<parley-protocol>") != 1 {
		t.Fatalf("protocol bytes lost or duplicated: %v", err)
	}
	records := terminalRecords(t, root)
	if len(records) != 1 || records[0].Metadata.Context.SourceSHA256 == nil || *records[0].Metadata.Context.SourceSHA256 != protocolpacket.Hash(string(source)) {
		t.Fatalf("wrong attestation: %+v", records)
	}
}

func TestProbeBoundaryIsExplicitAndUnattested(t *testing.T) {
	for _, phase := range []string{"", "implementation", "preflight", "runtime-probe"} {
		t.Run("phase-"+phase, func(t *testing.T) {
			root := t.TempDir()
			ctx := WithLaunchInfo(context.Background(), LaunchInfo{Phase: phase, Context: telemetry.Context{Mode: "full", SourceSHA256: telemetry.String("forged")}})
			cmd, cleanup, err := ProbeCommandFor(ctx, root, telemetryShell("printf PONG", false), "bounded capability probe")
			if cleanup != nil {
				defer cleanup()
			}
			allowed := phase == "preflight" || phase == "runtime-probe"
			if !allowed {
				if err == nil {
					t.Fatal("task phase allowed probe-only launch")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if err := cmd.Run(); err != nil {
				t.Fatal(err)
			}
			records := terminalRecords(t, root)
			if len(records) != 1 {
				t.Fatalf("missing probe record: %+v", records)
			}
			c := records[0].Metadata.Context
			if c.Mode != "probe-only" || c.SourceSHA256 != nil || c.PacketSHA256 != nil || c.FallbackReason == nil || *c.FallbackReason != "no-protocol-task" {
				t.Fatalf("probe certified task authority: %+v", c)
			}
		})
	}
}
