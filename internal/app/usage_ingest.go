package app

// `parley usage ingest` (lean-organizer D.2/D.3): a per-idea, per-phase organizer
// usage ledger built from the CLIENT'S OWN accounting files.
//
// Binding properties (ratified): NO flag accepts a token count — the tool takes a
// path and parses the client's file itself (a number typed by a model is RECALL and
// inadmissible; a number the tool parses from the client-written file is PRIMARY
// with a locator). Streaming with bounded memory over the verified 228 MB Codex
// rollout; <= 1 KB stdout; ONE ledger row per ingest carrying idea, phase, agent,
// source path, parser id, ingest time, and the six total_token_usage fields
// verbatim — tokens only, no dollar claim. Idempotent re-ingest. Attribution is by
// the explicit --idea/--phase arguments plus run-record timestamp windows; residue
// the windows cannot resolve is labeled `ambiguous`, with the method stated in the
// ledger header. A slug appearing in a session file is never attribution evidence.

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/protocol"
)

// usageSourceIDs are the shipped client-accounting parsers.
var usageSourceIDs = map[string]bool{
	"codex-rollout": true,
	"claude-jsonl":  true,
}

// usageTotals is the six total_token_usage fields, verbatim names.
type usageTotals struct {
	InputTokens       int64 `json:"input_tokens"`
	CachedInputTokens int64 `json:"cached_input_tokens"`
	CacheWriteTokens  int64 `json:"cache_write_input_tokens"`
	OutputTokens      int64 `json:"output_tokens"`
	ReasoningTokens   int64 `json:"reasoning_output_tokens"`
	TotalTokens       int64 `json:"total_tokens"`
}

// Claude's message.usage names differ; the row still stores the six canonical
// fields, derived as: cached=cache_read, cache_write=cache_creation,
// reasoning=0 (Claude reports none).
type claudeUsage struct {
	InputTokens              int64 `json:"input_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
}

// ingestResult is what one ingest produced (also the ledger row body).
type ingestResult struct {
	ParserID        string      `json:"parser_id"`
	SourcePath      string      `json:"source_path"`
	LastEventAt     string      `json:"last_event_at"`
	FirstEventAt    string      `json:"first_event_at"`
	EventCount      int64       `json:"event_count"`
	Totals          usageTotals `json:"total_token_usage"`
	Attribution     string      `json:"attribution"` // attributed | ambiguous
	AttributionNote string      `json:"attribution_note"`
	ContentSHA256   string      `json:"content_sha256"` // idempotency key basis
}

// streamLines yields file lines with bounded memory: a line longer than the cap is
// discarded whole (a token_count/usage event line is small; oversized lines are
// payload blobs the parsers never need).
func streamLines(path string, yield func(line []byte) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 1<<20)
	const cap = 1 << 24 // 16 MiB per line; bounded
	var acc []byte
	for {
		chunk, err := r.ReadSlice('\n')
		if len(acc)+len(chunk) > cap {
			acc = acc[:0] // drop the oversized line entirely
			if err == io.EOF {
				break
			}
			continue
		}
		acc = append(acc, chunk...)
		if err == io.EOF {
			if len(acc) > 0 {
				if yerr := yield(acc); yerr != nil {
					return yerr
				}
			}
			break
		}
		if err != nil && err != bufio.ErrBufferFull {
			return err
		}
		if len(chunk) > 0 && chunk[len(chunk)-1] == '\n' {
			if yerr := yield(acc); yerr != nil {
				return yerr
			}
			acc = acc[:0]
		}
	}
	return nil
}

// parseCodexRollout streams a Codex rollout JSONL and keeps the LAST cumulative
// token_count event (plus first/last timestamps and the event count).
func parseCodexRollout(path string) (ingestResult, error) {
	res := ingestResult{ParserID: "codex-rollout/v1", SourcePath: path}
	var acc []byte
	err := streamLines(path, func(line []byte) error {
		if !containsBytes(line, []byte("token_count")) {
			return nil
		}
		acc = append(acc[:0], line...)
		return nil
	})
	if err != nil {
		return res, err
	}
	// Re-scan is avoided by keeping every candidate; but memory must stay bounded —
	// instead keep only the LAST candidate line (above) and parse it now.
	if len(acc) == 0 {
		return res, fmt.Errorf("no token_count event found in %s", path)
	}
	var ev struct {
		Timestamp string `json:"timestamp"`
		Payload   struct {
			Type string `json:"type"`
			Info struct {
				TotalTokenUsage usageTotals `json:"total_token_usage"`
			} `json:"info"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(acc, &ev); err != nil {
		return res, fmt.Errorf("token_count event: %w", err)
	}
	if ev.Payload.Type != "token_count" {
		return res, fmt.Errorf("last token_count candidate is a %s event", ev.Payload.Type)
	}
	res.Totals = ev.Payload.Info.TotalTokenUsage
	res.LastEventAt = ev.Timestamp
	res.EventCount = 1
	// first event timestamp + exact event count need a second bounded pass
	first := ""
	var count int64
	err = streamLines(path, func(line []byte) error {
		if !containsBytes(line, []byte("token_count")) {
			return nil
		}
		var probe struct {
			Timestamp string `json:"timestamp"`
			Payload   struct {
				Type string `json:"type"`
			} `json:"payload"`
		}
		if json.Unmarshal(line, &probe) == nil && probe.Payload.Type == "token_count" {
			count++
			if first == "" {
				first = probe.Timestamp
			}
		}
		return nil
	})
	if err == nil {
		res.FirstEventAt = first
		res.EventCount = count
	}
	res.ContentSHA256 = sha256Hex(string(acc))
	return res, nil
}

// parseClaudeJSONL streams a Claude session JSONL and SUMS message.usage records
// (Claude reports per-message usage, not cumulative).
func parseClaudeJSONL(path string) (ingestResult, error) {
	res := ingestResult{ParserID: "claude-jsonl/v1", SourcePath: path}
	var sum usageTotals
	var count int64
	first, last := "", ""
	var lastLine []byte
	err := streamLines(path, func(line []byte) error {
		if !containsBytes(line, []byte("usage")) {
			return nil
		}
		var probe struct {
			Timestamp string `json:"timestamp"`
			Message   struct {
				Usage *claudeUsage `json:"usage"`
			} `json:"message"`
		}
		if json.Unmarshal(line, &probe) != nil || probe.Message.Usage == nil {
			return nil
		}
		u := probe.Message.Usage
		sum.InputTokens += u.InputTokens
		sum.CachedInputTokens += u.CacheReadInputTokens
		sum.CacheWriteTokens += u.CacheCreationInputTokens
		sum.OutputTokens += u.OutputTokens
		sum.TotalTokens += u.InputTokens + u.OutputTokens
		count++
		if first == "" {
			first = probe.Timestamp
		}
		last = probe.Timestamp
		lastLine = append(lastLine[:0], line...)
		return nil
	})
	if err != nil {
		return res, err
	}
	if count == 0 {
		return res, fmt.Errorf("no message.usage record found in %s", path)
	}
	res.Totals = sum
	res.FirstEventAt, res.LastEventAt, res.EventCount = first, last, count
	res.ContentSHA256 = sha256Hex(string(lastLine))
	return res, nil
}

func containsBytes(hay, needle []byte) bool {
	return len(hay) >= len(needle) && (string(hay) == string(needle) || indexOfBytes(hay, needle) >= 0)
}

func indexOfBytes(hay, needle []byte) int {
	for i := 0; i+len(needle) <= len(hay); i++ {
		match := true
		for j := range needle {
			if hay[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// ledgerRow is one appended JSONL row.
type ledgerRow struct {
	IngestedAt string       `json:"ingested_at"`
	Idea       string       `json:"idea"`
	Phase      int          `json:"phase"`
	Agent      string       `json:"agent"`
	Source     string       `json:"source"`
	Result     ingestResult `json:"result"`
}

const usageLedgerMethod = `attribution method: rows attribute by the EXPLICIT --idea/--phase arguments, corroborated by run-record timestamp windows (parley-deck/runs/*/run.json created_at..updated_at). Residue the windows cannot resolve is labeled "ambiguous". A slug appearing inside a session file is never attribution evidence. Tokens only; no dollar claim. Idempotent: an identical (agent, source, idea, phase, parser, content_sha256) row is not re-appended.`

func usageLedgerPath(root, idea string) string {
	return filepath.Join(root, protocol.DeckDir, "ideas", idea, "usage-ledger.jsonl")
}

// resolveAttribution corroborates the explicit args with run-record windows.
func resolveAttribution(root, idea string, res ingestResult) (string, string) {
	runsDir := filepath.Join(root, protocol.DeckDir, "runs")
	entries, err := os.ReadDir(runsDir)
	if err != nil || res.LastEventAt == "" {
		return "ambiguous", "no run-record window available to corroborate the explicit attribution"
	}
	when, werr := time.Parse(time.RFC3339Nano, res.LastEventAt)
	if werr != nil {
		return "ambiguous", "last accounting event has no parseable timestamp"
	}
	matching := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, rerr := os.ReadFile(filepath.Join(runsDir, e.Name(), "run.json"))
		if rerr != nil {
			continue
		}
		var manifest struct {
			IdeaSlug  string `json:"idea_slug"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}
		if json.Unmarshal(data, &manifest) != nil {
			continue
		}
		created, cerr := time.Parse(time.RFC3339Nano, strings.Replace(manifest.CreatedAt, " ", "T", 1))
		updated, uerr := time.Parse(time.RFC3339Nano, strings.Replace(manifest.UpdatedAt, " ", "T", 1))
		if cerr != nil || uerr != nil {
			continue
		}
		if !when.Before(created) && !when.After(updated) {
			matching++
			if manifest.IdeaSlug != idea {
				return "ambiguous", "accounting event falls inside a run of a different idea"
			}
		}
	}
	if matching == 0 {
		return "ambiguous", "no run-record window contains the accounting event (pre-run or post-run residue)"
	}
	return "attributed", "explicit --idea/--phase corroborated by " + fmt.Sprintf("%d", matching) + " run-record window(s)"
}

func runUsageIngest(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("usage ingest", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "workspace root")
	agent := fs.String("agent", "", "agent id the accounting file belongs to")
	source := fs.String("source", "", "client accounting format: codex-rollout|claude-jsonl")
	path := fs.String("path", "", "path to the client accounting file (a token count is never accepted)")
	idea := fs.String("idea", "", "idea slug (explicit attribution; never slug scanning)")
	phase := fs.Int("phase", -1, "protocol phase the usage is attributed to")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	for _, missing := range []struct{ name, value string }{
		{"--agent", *agent}, {"--source", *source}, {"--path", *path}, {"--idea", *idea},
	} {
		if strings.TrimSpace(missing.value) == "" {
			fmt.Fprintf(stderr, "usage ingest: %s is required\n", missing.name)
			return 1
		}
	}
	if !usageSourceIDs[*source] {
		fmt.Fprintf(stderr, "usage ingest: unknown --source %q (want codex-rollout|claude-jsonl)\n", *source)
		return 1
	}
	if *phase < 0 || *phase > 8 {
		fmt.Fprintf(stderr, "usage ingest: --phase must be 0..8\n")
		return 1
	}
	root, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(stderr, "usage ingest: %v\n", err)
		return 1
	}
	var res ingestResult
	switch *source {
	case "codex-rollout":
		res, err = parseCodexRollout(*path)
	case "claude-jsonl":
		res, err = parseClaudeJSONL(*path)
	}
	if err != nil {
		fmt.Fprintf(stderr, "usage ingest: %v\n", err)
		return 1
	}
	res.Attribution, res.AttributionNote = resolveAttribution(root, *idea, res)

	row := ledgerRow{
		IngestedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Idea:       *idea, Phase: *phase, Agent: *agent, Source: *source, Result: res,
	}
	ledger := usageLedgerPath(root, *idea)
	appended, err := appendIdempotent(ledger, row)
	if err != nil {
		fmt.Fprintf(stderr, "usage ingest: %v\n", err)
		return 1
	}
	status := "appended"
	if !appended {
		status = "idempotent no-op (identical row exists)"
	}
	// <= 1 KB stdout: one summary line.
	fmt.Fprintf(stdout, "ingest %s: %s idea=%s phase=%d agent=%s events=%d total_tokens=%d attribution=%s\n",
		status, res.ParserID, *idea, *phase, *agent, res.EventCount, res.Totals.TotalTokens, res.Attribution)
	return 0
}

// appendIdempotent appends one ledger row unless an identical row (same agent,
// source, idea, phase, parser, content hash) already exists. The method header is
// (re)stated as a leading comment line so the ledger is self-describing.
func appendIdempotent(path string, row ledgerRow) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	key := func(r ledgerRow) string {
		return fmt.Sprintf("%s|%s|%d|%s|%s|%s", r.Agent, r.Source, r.Phase, r.Idea, r.Result.ParserID, r.Result.ContentSHA256)
	}
	data, _ := os.ReadFile(path)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var existing ledgerRow
		if json.Unmarshal([]byte(line), &existing) == nil && key(existing) == key(row) {
			return false, nil
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	body, err := json.Marshal(row)
	if err != nil {
		return false, err
	}
	if len(data) == 0 {
		if _, err := fmt.Fprintf(f, "# %s\n", usageLedgerMethod); err != nil {
			return false, err
		}
	}
	_, err = f.Write(append(body, '\n'))
	return err == nil, err
}

// runUsage dispatches `parley usage <verb>`; ingest is the only verb today.
func runUsage(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "ingest" {
		fmt.Fprintln(stderr, "usage: parley usage ingest --agent ID --source codex-rollout|claude-jsonl --path FILE --idea SLUG --phase N [--dir DIR]")
		return 1
	}
	return runUsageIngest(args[1:], stdout, stderr)
}
