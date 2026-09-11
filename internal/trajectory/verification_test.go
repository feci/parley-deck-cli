package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type verificationProcessInput struct {
	Ticket   VerificationTicket
	Criteria []Criterion
	Parent   string
}

// The synthetic helper is a separate executable process, never a model call.
func TestVerificationJournalProcessHelper(t *testing.T) {
	path := os.Getenv("PARLEY_JOURNAL_PROCESS_INPUT")
	if path == "" {
		return
	}
	var input verificationProcessInput
	data, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(data, &input)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(19)
	}
	if err = os.Setenv("PARLEY_VERIFICATION_FIXTURE_PID", strconv.Itoa(os.Getpid())); err != nil {
		os.Exit(19)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err = ExecuteCapturedVerification(ctx, input.Ticket, "verifier-invocation", input.Criteria, input.Parent)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(17)
	}
	os.Exit(0)
}

func journalFixture(t *testing.T, criteria []Criterion, mutate func(string)) (VerificationTicket, string) {
	t.Helper()
	root, b, _ := capturedFixture(t, criteria, mutate)
	ticket, err := PrepareCapturedVerification(context.Background(), root, "fixture", "reviewer", "verification-run")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-verifications", ticket.Request.Charge.EntryKey)
	return ticket, path
}

func reserveJournal(t *testing.T, ticket VerificationTicket) {
	t.Helper()
	if err := ReserveCapturedVerificationLaunch(context.Background(), ticket, "verifier-invocation"); err != nil {
		t.Fatal(err)
	}
}

func journalCommand(t *testing.T, ticket VerificationTicket, criteria []Criterion) (*exec.Cmd, *bytes.Buffer) {
	t.Helper()
	input := verificationProcessInput{ticket, criteria, t.TempDir()}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "helper-input.json")
	if err = os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	t.Cleanup(cancel)
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestVerificationJournalProcessHelper$")
	cmd.Env = append(os.Environ(), "PARLEY_JOURNAL_PROCESS_INPUT="+path)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	return cmd, &out
}

func TestDurableCapturedVerificationProcessRaceAndReplay(t *testing.T) {
	trace := filepath.Join(t.TempDir(), "executions")
	criteria := capturedCriteria()
	criteria[0].Command = "cat source >> '" + strings.ReplaceAll(trace, "'", "'\\''") + "'; " + criteria[0].Command
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	original := snapshotRead(t, filepath.Join(path, "request.json"))
	for _, run := range []string{ticket.RunID, "different-run"} {
		if _, err := PrepareCapturedVerification(context.Background(), ticket.Root, "fixture", "reviewer", run); err == nil {
			t.Fatal("prepared verification was silently replaced or retried")
		}
	}
	reserveJournal(t, ticket)
	if err := ReserveCapturedVerificationLaunch(context.Background(), ticket, "another-invocation"); err == nil {
		t.Fatal("reserved verifier invocation was overwritten")
	}
	// The comparison must still execute the retained original dirty patch.
	snapshotWrite(t, ticket.Root, "source", []byte("later live edit\n"), 0600)
	a, outputA := journalCommand(t, ticket, criteria)
	b, outputB := journalCommand(t, ticket, criteria)
	if err := a.Start(); err != nil {
		t.Fatal(err)
	}
	if err := b.Start(); err != nil {
		_ = a.Process.Kill()
		_ = a.Wait()
		t.Fatal(err)
	}
	errorA, errorB := a.Wait(), b.Wait()
	if (errorA == nil) == (errorB == nil) {
		t.Fatalf("exactly one helper must win: %v %v\n%s\n%s", errorA, errorB, outputA.String(), outputB.String())
	}
	receipt, observation, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation")
	if err != nil || receipt.Steps != 4 || receipt.HelperPID == os.Getpid() || receipt.PreparedSHA256 == "" {
		t.Fatalf("actual independent process receipt unavailable: %+v %v", receipt, err)
	}
	assessment, err := AssessCaptured(ticket.Request, observation)
	if err != nil || assessment.Outcome != Regression {
		t.Fatalf("actual paired regression not recovered: %+v %v", assessment, err)
	}
	for _, root := range []string{receipt.BeforeRoot, receipt.AfterRoot} {
		if _, err := os.Stat(root); !os.IsNotExist(err) {
			t.Fatal("completed helper leaked its private source workspace")
		}
	}
	beforeReceipt := snapshotRead(t, filepath.Join(path, "receipt.json"))
	replay, replayOutput := journalCommand(t, ticket, criteria)
	if err = replay.Run(); err == nil {
		t.Fatalf("new process replayed a completed verification: %s", replayOutput)
	}
	if string(snapshotRead(t, trace)) != "original\nbroken\nbroken\noriginal\n" ||
		!bytes.Equal(original, snapshotRead(t, filepath.Join(path, "request.json"))) ||
		!bytes.Equal(beforeReceipt, snapshotRead(t, filepath.Join(path, "receipt.json"))) {
		t.Fatal("competing or replayed helper changed retained execution history")
	}
	if err = RequireResolved(context.Background(), ticket.Root, "fixture"); err == nil {
		t.Fatal("journal self-issued trajectory resolution")
	}
}

func TestDurableCapturedVerificationCrashPreservesSteps(t *testing.T) {
	criteria := capturedCriteria()
	// Kill only the synthetic helper's explicitly exported PID. The criterion
	// shell exits immediately; no sleeping descendant or unrelated process is used.
	criteria[0].Command = `if [ "$(cat source)" = broken ]; then kill -KILL "$PARLEY_VERIFICATION_FIXTURE_PID"; exit 9; fi; ` + criteria[0].Command
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	cmd, output := journalCommand(t, ticket, criteria)
	if err := cmd.Run(); err == nil {
		t.Fatalf("fixture helper was not interrupted: %s", output)
	}
	_, observation, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation")
	if err == nil || !strings.Contains(err.Error(), "terminal receipt unavailable") || len(observation.Pairs) != 1 || !observation.Pairs[0].Before[0].Complete || observation.Pairs[0].After[0].Complete {
		t.Fatalf("crash lost completed execution or fabricated a terminal: %+v %v", observation, err)
	}
	step := snapshotRead(t, filepath.Join(path, "step-001.json"))
	if _, err = os.Stat(filepath.Join(path, "receipt.json")); !os.IsNotExist(err) {
		t.Fatal("abruptly killed helper claimed a terminal receipt")
	}
	replay, _ := journalCommand(t, ticket, criteria)
	if err = replay.Run(); err == nil || !bytes.Equal(step, snapshotRead(t, filepath.Join(path, "step-001.json"))) {
		t.Fatal("fresh process repeated interrupted work or replaced its evidence")
	}
	// Test-owned parent directories are removed by testing after the helper is
	// terminal. Production recovery must inspect prepared.json before cleanup.
}

func TestDurableCapturedVerificationFailuresRetained(t *testing.T) {
	for _, kind := range []string{"staged-source", "source-drift", "changed-command", "opaque-output"} {
		t.Run(kind, func(t *testing.T) {
			criteria := capturedCriteria()
			mutate := dirtyCapturedChild(t)
			if kind == "staged-source" {
				mutate = func(root string) {
					snapshotWrite(t, root, "source", []byte("broken\n"), 0600)
					gitFixture(t, root, "add", "source")
				}
			} else if kind == "source-drift" {
				criteria[0].Command = "printf changed > source; " + criteria[0].Command
			} else if kind == "opaque-output" {
				criteria[0].Command = "true"
			}
			ticket, path := journalFixture(t, criteria, mutate)
			reserveJournal(t, ticket)
			if kind == "changed-command" {
				criteria[0].Command = "true"
			}
			receipt, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir())
			if err == nil || receipt.FailureStage == "" {
				t.Fatal("failed helper was certified or not retained")
			}
			wantSteps := map[string]int{"staged-source": 0, "changed-command": 0, "source-drift": 1, "opaque-output": 4}[kind]
			read, observation, readErr := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation")
			if readErr == nil || read.Steps != wantSteps || read.FailureStage != receipt.FailureStage {
				t.Fatalf("failure/partial journal not recovered: %+v %v", read, readErr)
			}
			if wantSteps > 0 && len(observation.Pairs) != 1 {
				t.Fatal("partial execution observation lost")
			}
			original := snapshotRead(t, filepath.Join(path, "receipt.json"))
			if _, err = ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir()); err == nil || !bytes.Equal(original, snapshotRead(t, filepath.Join(path, "receipt.json"))) {
				t.Fatal("failed attempt was overwritten by a replay")
			}
		})
	}
}

func TestDurableCapturedVerificationRejectsChangedHistory(t *testing.T) {
	criteria := capturedCriteria()
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	if _, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir()); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"request.json", "launch.json", "claim.json", "prepared.json", "step-001.json", "step-003.json", "receipt.json"} {
		t.Run(name, func(t *testing.T) {
			file := filepath.Join(path, name)
			original := snapshotRead(t, file)
			defer func() {
				if err := os.WriteFile(file, original, 0600); err != nil {
					t.Fatal(err)
				}
			}()
			for _, tamper := range []func() error{
				func() error { return os.WriteFile(file, append(append([]byte{}, original...), '\n'), 0600) },
				func() error { return os.Remove(file) },
			} {
				if err := tamper(); err != nil {
					t.Fatal(err)
				}
				if _, _, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil {
					t.Fatal("changed or missing journal artifact accepted")
				}
			}
		})
	}
	for _, name := range []string{"step-005.json", "unexpected.json"} {
		file := filepath.Join(path, name)
		if err := os.WriteFile(file, []byte("{}\n"), 0600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil {
			t.Fatal("unexpected execution history was hidden")
		}
		if err := os.Remove(file); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err != nil {
		t.Fatal("unchanged history no longer readable", err)
	}
}

func TestDurableCapturedVerificationRetentionFailureStops(t *testing.T) {
	criteria := capturedCriteria()
	criteria[0].Command = `mkdir "$PARLEY_JOURNAL_FIXTURE_DIR/step-001.json"; ` + criteria[0].Command
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	t.Setenv("PARLEY_JOURNAL_FIXTURE_DIR", path)
	reserveJournal(t, ticket)
	receipt, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir())
	if err == nil || receipt.Steps != 0 || receipt.FailureStage != "execution" {
		t.Fatalf("failed execution write allowed progression: %+v %v", receipt, err)
	}
	if _, err = os.Stat(filepath.Join(path, "step-002.json")); !os.IsNotExist(err) {
		t.Fatal("helper kept executing after losing required evidence")
	}
	if _, _, err = ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil {
		t.Fatal("unretained execution was accepted")
	}
}

func TestDurableCapturedVerificationDetectsMidExecutionAuthorityChange(t *testing.T) {
	criteria := capturedCriteria()
	criteria[0].Command = `printf '\n' >> "$PARLEY_JOURNAL_FIXTURE_DIR/launch.json"; ` + criteria[0].Command
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	t.Setenv("PARLEY_JOURNAL_FIXTURE_DIR", path)
	reserveJournal(t, ticket)
	receipt, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir())
	if err == nil || receipt.Steps != 1 || receipt.FailureStage != "execution" {
		t.Fatalf("mid-execution launch change was not retained and refused: %+v %v", receipt, err)
	}
	if _, err = os.Stat(filepath.Join(path, "step-002.json")); !os.IsNotExist(err) {
		t.Fatal("helper continued after its launch authority changed")
	}
	if _, _, err = ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil {
		t.Fatal("changed authority accepted a receipt")
	}
}

func TestDurableCapturedVerificationReceiptWriteFailureNeverAccepts(t *testing.T) {
	criteria := capturedCriteria()
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	if err := os.Mkdir(filepath.Join(path, "receipt.json"), 0700); err != nil {
		t.Fatal(err)
	}
	receipt, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir())
	if err == nil || receipt.Steps != 4 {
		t.Fatalf("terminal write failure hidden or executions lost: %+v %v", receipt, err)
	}
	_, observation, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation")
	if err == nil || len(observation.Pairs) != 1 || !observation.Pairs[0].Before[1].Complete {
		t.Fatal("missing terminal was accepted or actual completed observations disappeared", err)
	}
	if _, err = ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir()); err == nil {
		t.Fatal("receipt failure caused implicit execution retry")
	}
}

func TestDurableCapturedVerificationRequiresReservedAuthority(t *testing.T) {
	criteria := capturedCriteria()
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	if _, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir()); err == nil {
		t.Fatal("helper executed before launch reservation")
	}
	if err := ReserveCapturedVerificationLaunch(context.Background(), ticket, ticket.Request.InvocationID); err == nil {
		t.Fatal("patch invocation was accepted as independent verification")
	}
	reserveJournal(t, ticket)
	for _, changed := range []string{"invocation", "run", "verifier", "scope", "root"} {
		t.Run(changed, func(t *testing.T) {
			data, _ := json.Marshal(ticket)
			var other VerificationTicket
			if err := json.Unmarshal(data, &other); err != nil {
				t.Fatal(err)
			}
			invocation := "verifier-invocation"
			switch changed {
			case "invocation":
				invocation = "different-invocation"
			case "run":
				other.RunID = "different-run"
			case "verifier":
				other.Request.Verifier = "different-verifier"
			case "scope":
				other.Request.Criteria[0].CommandSHA256 = digest([]byte("true"))
			case "root":
				other.Root = t.TempDir()
			}
			if _, err := ExecuteCapturedVerification(context.Background(), other, invocation, criteria, t.TempDir()); err == nil {
				t.Fatal("substituted authority claimed a helper")
			}
			if _, err := os.Stat(filepath.Join(path, "claim.json")); !os.IsNotExist(err) {
				t.Fatal("refused authority consumed or wrote another invocation's helper claim")
			}
		})
	}
}

func TestDurableCapturedVerificationArtifactRefusals(t *testing.T) {
	path := t.TempDir()
	dir, err := os.OpenRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()
	value := struct {
		Version int `json:"version"`
	}{1}
	if _, err = writeVerificationArtifact(dir, "record.json", value); err != nil {
		t.Fatal(err)
	}
	original := snapshotRead(t, filepath.Join(path, "record.json"))
	if _, err = writeVerificationArtifact(dir, "record.json", value); err == nil || !bytes.Equal(original, snapshotRead(t, filepath.Join(path, "record.json"))) {
		t.Fatal("exclusive artifact publication overwrote an existing record")
	}
	for _, malformed := range []string{"{", "{\"version\":1,\"version\":1}\n", "{\"version\":1,\"unknown\":1}\n", strings.Repeat("x", (1<<20)+1)} {
		if err = os.WriteFile(filepath.Join(path, "bad.json"), []byte(malformed), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err = readVerificationArtifact(dir, "bad.json", &value); err == nil {
			t.Fatal("malformed or oversized artifact accepted")
		}
	}
	if err = os.Symlink("record.json", filepath.Join(path, "link.json")); err != nil {
		t.Skip("symlinks unavailable")
	}
	if _, err = readVerificationArtifact(dir, "link.json", &value); err == nil {
		t.Fatal("symlink used as a verification artifact")
	}
}
