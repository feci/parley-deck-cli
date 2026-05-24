package acp

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// inheritedEnvVars is the allow-list of variables we copy from the login
// shell into the child. Mirrors AionUi's SHELL_INHERITED_ENV_VARS — only
// keys that affect tool discovery (PATH) and TLS/CA trust are propagated.
var inheritedEnvVars = []string{
	"PATH",
	"NODE_EXTRA_CA_CERTS",
	"SSL_CERT_FILE",
	"SSL_CERT_DIR",
	"REQUESTS_CA_BUNDLE",
	"CURL_CA_BUNDLE",
}

var (
	shellEnvOnce sync.Once
	shellEnv     map[string]string
)

// LoadShellEnv runs "$SHELL -l -c env" once per process and caches the
// allow-listed variables. On Windows it returns nil because PATH inheritance
// already works through normal environment propagation.
//
// We use -l (login) and explicitly NOT -i (interactive); AionUi notes
// interactive shells call tcsetpgrp() and can break Ctrl+C delivery.
func LoadShellEnv(ctx context.Context) map[string]string {
	shellEnvOnce.Do(func() {
		shellEnv = loadShellEnvImpl(ctx)
	})
	return shellEnv
}

func loadShellEnvImpl(ctx context.Context) map[string]string {
	if runtime.GOOS == "windows" {
		return nil
	}
	shell := strings.TrimSpace(os.Getenv("SHELL"))
	if shell == "" || !strings.HasPrefix(shell, "/") {
		return nil
	}
	probeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(probeCtx, shell, "-l", "-c", "env")
	cmd.Env = append(os.Environ(), "HOME="+os.Getenv("HOME"))
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	result := make(map[string]string, len(inheritedEnvVars))
	allowed := make(map[string]struct{}, len(inheritedEnvVars))
	for _, key := range inheritedEnvVars {
		allowed[key] = struct{}{}
	}
	for _, line := range strings.Split(string(out), "\n") {
		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			continue
		}
		key := line[:eq]
		if _, ok := allowed[key]; !ok {
			continue
		}
		result[key] = line[eq+1:]
	}
	return result
}

// MergedEnv returns os.Environ() augmented with shell-loaded keys where they
// are missing (or merging PATH so user-installed tools become visible to
// processes launched outside a terminal).
func MergedEnv(ctx context.Context, extra map[string]string) []string {
	shell := LoadShellEnv(ctx)
	base := os.Environ()
	merged := make([]string, 0, len(base)+len(shell)+len(extra))
	have := make(map[string]int, len(base))
	for i, kv := range base {
		eq := strings.IndexByte(kv, '=')
		if eq <= 0 {
			continue
		}
		have[kv[:eq]] = i
		merged = append(merged, kv)
	}
	for key, value := range shell {
		if key == "PATH" {
			if idx, ok := have["PATH"]; ok {
				merged[idx] = "PATH=" + mergePaths(merged[idx][len("PATH="):], value)
				continue
			}
		}
		if _, ok := have[key]; ok {
			continue
		}
		merged = append(merged, key+"="+value)
		have[key] = len(merged) - 1
	}
	for key, value := range extra {
		if idx, ok := have[key]; ok {
			merged[idx] = key + "=" + value
			continue
		}
		merged = append(merged, key+"="+value)
		have[key] = len(merged) - 1
	}
	return merged
}

func mergePaths(a, b string) string {
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}
	seen := map[string]struct{}{}
	var out []string
	for _, source := range []string{a, b} {
		for _, p := range strings.Split(source, sep) {
			if p == "" {
				continue
			}
			if _, dup := seen[p]; dup {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	return strings.Join(out, sep)
}
