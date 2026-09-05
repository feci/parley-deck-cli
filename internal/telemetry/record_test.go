package telemetry

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func readRecord(t *testing.T, path string) Record {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record Record
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	return record
}

func TestLifecycleAndUnknowns(t *testing.T) {
	inv, err := Begin(t.TempDir(), Metadata{Agent: "agent-1", Adapter: "test", RequestedModel: String("model-exact"), AttemptOrdinal: 1})
	if err != nil {
		t.Fatal(err)
	}
	requested := readRecord(t, filepath.Join(inv.Dir, "requested.json"))
	if requested.InvocationID == "" || requested.Type != "invocation.requested" || requested.StartedAt != nil || requested.Outcome != nil {
		t.Fatalf("bad request: %+v", requested)
	}
	if err := inv.Started(123); err != nil {
		t.Fatal(err)
	}
	if err := inv.Finish(Outcome{Status: "finished"}); err != nil {
		t.Fatal(err)
	}
	terminal := readRecord(t, filepath.Join(inv.Dir, "terminal.json"))
	if terminal.StartedAt == nil || terminal.CompletedAt == nil || terminal.DurationMS == nil || *terminal.PID != 123 {
		t.Fatalf("missing lifecycle: %+v", terminal)
	}
	usage := terminal.Outcome.Usage
	if usage.CostUSD != nil || usage.InputTokens != nil || usage.ReportedModel != nil || usage.Source != "unavailable" {
		t.Fatalf("invented observation: %+v", usage)
	}
	data, _ := os.ReadFile(filepath.Join(inv.Dir, "terminal.json"))
	if !strings.Contains(string(data), `"cost_usd": null`) || !strings.Contains(string(data), `"reported_model": null`) {
		t.Fatal("unknowns not explicit")
	}
	if terminal.Metadata.RequestedModel == nil || *terminal.Metadata.RequestedModel != "model-exact" {
		t.Fatal("requested identity lost")
	}
	info, _ := os.Stat(filepath.Join(inv.Dir, "terminal.json"))
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("record readable by others: %o", info.Mode().Perm())
	}
}

func TestFailedStartHasNoInventedProcess(t *testing.T) {
	inv, err := Begin(t.TempDir(), Metadata{Agent: "missing-agent"})
	if err != nil {
		t.Fatal(err)
	}
	if err := inv.Finish(Outcome{Status: "failed", FailureClass: String("start-error")}); err != nil {
		t.Fatal(err)
	}
	value := readRecord(t, filepath.Join(inv.Dir, "terminal.json"))
	if value.StartedAt != nil || value.PID != nil || value.Outcome.ExitCode != nil {
		t.Fatal("failed start invented a process")
	}
	if _, err := os.Stat(filepath.Join(inv.Dir, "started.json")); !os.IsNotExist(err) {
		t.Fatal("unexpected start record")
	}
}

func TestUniqueConcurrentAttemptsAndRetryLink(t *testing.T) {
	dir := t.TempDir()
	var wg sync.WaitGroup
	ids := make(chan string, 40)
	for range 40 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i, err := Begin(dir, Metadata{Agent: "same-agent", AttemptOrdinal: 1, SegmentID: "same-segment"})
			if err != nil {
				t.Error(err)
				return
			}
			ids <- i.ID
			if err := i.Finish(Outcome{Status: "failed"}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	close(ids)
	seen := map[string]bool{}
	previous := ""
	for id := range ids {
		if seen[id] {
			t.Fatal("duplicate invocation ID")
		}
		seen[id] = true
		previous = id
	}
	if len(seen) != 40 {
		t.Fatalf("got %d attempts", len(seen))
	}
	retry, err := Begin(dir, Metadata{Agent: "same-agent", AttemptOrdinal: 2, RetryOf: String(previous)})
	if err != nil {
		t.Fatal(err)
	}
	if retry.ID == previous || *retry.Snapshot().Metadata.RetryOf != previous {
		t.Fatal("retry identity collapsed")
	}
}

func TestUnsafeMetadataNeverPersisted(t *testing.T) {
	secret := "npm_FAKE_CANARY_NOT_A_REAL_TOKEN_1234567890"
	inv, err := Begin(t.TempDir(), Metadata{Agent: secret, RequestedModel: String(secret), RequestedEffort: String("max"), Context: Context{Mode: "full-fallback", FallbackReason: String("https://user:fake-password@example.test/path")}})
	if err != nil {
		t.Fatal(err)
	}
	if err := inv.Started(123); err != nil {
		t.Fatal(err)
	}
	if err := inv.Finish(Outcome{Status: "failed", FailureClass: String(secret), Usage: Usage{ReportedModel: String(secret), ReportedModels: []string{secret, "model-safe"}, Source: secret}}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"requested.json", "started.json", "terminal.json"} {
		data, err := os.ReadFile(filepath.Join(inv.Dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), secret) || strings.Contains(string(data), "fake-password") || strings.Contains(string(data), "example.test") {
			t.Fatalf("secret leaked in %s", name)
		}
	}
	if len(inv.Snapshot().Warnings) < 2 {
		t.Fatal("missing omission diagnostics")
	}
}

func TestCallerCannotMutateValidatedMetadata(t *testing.T) {
	model := "model-safe"
	i, err := Begin(t.TempDir(), Metadata{RequestedModel: &model})
	if err != nil {
		t.Fatal(err)
	}
	model = "npm_FAKE_CANARY"
	snapshot := i.Snapshot()
	*snapshot.Metadata.RequestedModel = "ghp_FAKE_CANARY"
	if err := i.Started(123); err != nil {
		t.Fatal(err)
	}
	actual := readRecord(t, filepath.Join(i.Dir, "started.json"))
	if *actual.Metadata.RequestedModel != "model-safe" {
		t.Fatal("caller changed validated metadata")
	}
}

func TestInvalidUsageAndMonetaryValuesBecomeUnknown(t *testing.T) {
	negative, excessive, cost := int64(-1), int64(2_000_000_000_000), math.NaN()
	u := CleanUsage(Usage{InputTokens: &negative, OutputTokens: &excessive, CostUSD: &cost})
	if u.InputTokens != nil || u.OutputTokens != nil || u.CostUSD != nil {
		t.Fatalf("invalid usage retained: %+v", u)
	}
	if _, err := json.Marshal(u); err != nil {
		t.Fatal(err)
	}
}

func TestLegitimateTaskIdentifiersAreNotCredentialPrefixes(t *testing.T) {
	for _, value := range []string{"task-C01", "meta-protocol-change-phase-packet-and-fixup-budget-long-slug", "kimi-code/k3"} {
		if SafeLabel(value) == nil {
			t.Fatalf("legitimate label rejected: %s", value)
		}
	}
	for _, value := range []string{"sk-FAKE_CANARY", "provider/sk-FAKE_CANARY", strings.Repeat("a", 64)} {
		if SafeLabel(value) != nil {
			t.Fatal("secret-shaped metadata accepted")
		}
	}
}

func TestTerminalIsIdempotentAndNotOverwritten(t *testing.T) {
	i, err := Begin(t.TempDir(), Metadata{Agent: "agent-1"})
	if err != nil {
		t.Fatal(err)
	}
	if err := i.Finish(Outcome{Status: "failed"}); err != nil {
		t.Fatal(err)
	}
	first, _ := os.ReadFile(filepath.Join(i.Dir, "terminal.json"))
	if err := i.Finish(Outcome{Status: "finished"}); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(filepath.Join(i.Dir, "terminal.json"))
	if string(first) != string(second) {
		t.Fatal("terminal result overwritten")
	}
	if err := i.Started(123); err == nil {
		t.Fatal("started after terminal")
	}
}

func TestWriteFailuresAreReturned(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "not-a-directory")
	if err := os.WriteFile(file, []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Begin(file, Metadata{}); err == nil {
		t.Fatal("failed request persistence accepted")
	}
	i, err := Begin(dir, Metadata{})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(i.Dir, "terminal.json"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := i.Finish(Outcome{Status: "finished"}); err == nil {
		t.Fatal("terminal persistence failure hidden")
	}
	if err := i.Finish(Outcome{Status: "finished"}); err == nil {
		t.Fatal("repeated finish hid original failure")
	}
}
