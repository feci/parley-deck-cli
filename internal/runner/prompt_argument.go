package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocolpacket"
)

// Keep each prompt argument below Linux's 128 KiB per-string exec limit.
// This bound also applies on other hosts, so the same launch uses the same delivery.
const inlinePromptArgumentLimit = 120 << 10

func expandHeadlessArgs(resolved []string, root, prompt string) []string {
	args := make([]string, 0, len(resolved))
	for _, arg := range resolved {
		switch arg {
		case "{root}":
			args = append(args, root)
			continue
		case "{prompt}":
			args = append(args, prompt)
			continue
		}
		// Replace root first. Text inside the prompt is data, never a placeholder.
		arg = strings.ReplaceAll(arg, "{root}", root)
		args = append(args, strings.ReplaceAll(arg, "{prompt}", prompt))
	}
	return args
}

// headlessPromptArgs changes only oversized argv delivery. The private file keeps
// the complete prepared prompt, including the live protocol and its attestation.
// A file and a successful launch do not prove that a model read those bytes.
func headlessPromptArgs(resolved []string, root, prompt, invocationDir string) ([]string, error) {
	args := expandHeadlessArgs(resolved, root, prompt)
	needsFile := false
	for i, arg := range args {
		if len(arg) < inlinePromptArgumentLimit {
			continue
		}
		if !strings.Contains(resolved[i], "{prompt}") {
			return nil, fmt.Errorf("headless argument %d exceeds the %d-byte bound without a prompt placeholder", i, inlinePromptArgumentLimit)
		}
		needsFile = true
	}
	if !needsFile {
		return args, nil
	}

	// The launch has already created this private invocation directory. Never
	// overwrite an earlier launch's file, and retain successful writes for audit.
	file, err := os.CreateTemp(invocationDir, "headless-prompt-*.md")
	if err != nil {
		return nil, fmt.Errorf("create headless prompt file: %w", err)
	}
	path := file.Name()
	complete := false
	defer func() {
		_ = file.Close()
		if !complete {
			_ = os.Remove(path)
		}
	}()
	if _, err := file.WriteString(prompt); err != nil {
		return nil, fmt.Errorf("write headless prompt file: %w", err)
	}
	if err := fsutil.SyncFile(file); err != nil {
		return nil, fmt.Errorf("sync headless prompt file: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close headless prompt file: %w", err)
	}
	complete = true
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve headless prompt file: %w", err)
	}
	instruction := fmt.Sprintf("The complete Parley launch prompt is in the private UTF-8 file below. It includes the full live protocol, its attestation, and the task. Before any task action, read the entire file and verify its SHA256. Follow that exact prompt. If the file is missing, unreadable, or has a different hash, report a blocker and stop.\nPrompt file (quoted path): %q\nPrompt SHA256: %s\n", abs, protocolpacket.Hash(prompt))
	args = expandHeadlessArgs(resolved, root, instruction)
	for i, arg := range args {
		if len(arg) >= inlinePromptArgumentLimit {
			return nil, fmt.Errorf("headless argument %d still exceeds the %d-byte bound after file delivery", i, inlinePromptArgumentLimit)
		}
	}
	return args, nil
}
