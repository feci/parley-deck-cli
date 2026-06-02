package pipeline

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReviewAgreedFixes reads the machine-readable Phase-7 contract from a
// review/consensus.md frontmatter: `outstanding_agreed_fixes: <int>` (the
// authority for the fix-up loop) and optional `blocked: true`. `found` is false
// when the field is absent, so the auto driver can FAIL CLOSED rather than
// guess "zero fixes" from prose (§12 / decided 2026-06-02).
func ReviewAgreedFixes(path string) (count int, blocked bool, found bool, err error) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return 0, false, false, fmt.Errorf("read review consensus: %w", readErr)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return 0, false, false, nil
	}
	for _, line := range lines[1:] {
		t := strings.TrimSpace(line)
		if t == "---" {
			break
		}
		key, val, ok := strings.Cut(t, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		switch key {
		case "outstanding_agreed_fixes":
			n, convErr := strconv.Atoi(val)
			if convErr != nil {
				return 0, false, false, fmt.Errorf("outstanding_agreed_fixes=%q is not an int: %w", val, convErr)
			}
			count, found = n, true
		case "blocked":
			blocked = val == "true" || val == "yes"
		}
	}
	return count, blocked, found, nil
}

// Phase8Action is the next step the unattended fix-up loop should take.
type Phase8Action string

const (
	Phase8Complete Phase8Action = "complete"      // zero outstanding fixes -> block done
	Phase8Blocked  Phase8Action = "blocked"       // a reviewer BLOCKed -> stop for human
	Phase8Fixup    Phase8Action = "fixup"          // outstanding>0 and under the cap -> apply fixes, re-review
	Phase8Maxed    Phase8Action = "max-cycles"     // hit the safety cap -> stop
	Phase8NoData   Phase8Action = "missing-contract" // review consensus lacked the machine field -> fail closed
)

// Phase8Decision is the pure control logic for the unattended fix-up loop. It is
// kept separate from agent I/O so it is fully unit-testable.
func Phase8Decision(outstanding int, blocked, foundContract bool, cycle, maxCycles int) Phase8Action {
	if !foundContract {
		return Phase8NoData
	}
	if blocked {
		return Phase8Blocked
	}
	if outstanding <= 0 {
		return Phase8Complete
	}
	if cycle >= maxCycles {
		return Phase8Maxed
	}
	return Phase8Fixup
}
