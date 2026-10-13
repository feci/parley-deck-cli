package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocolpacket"
)

// This real child consumes the file instruction. It checks the digest before it
// prints bytes. Its behavior is a transport test, not proof of model compliance.
func TestSTEArgumentHelper(t *testing.T) {
	marker := -1
	for i, arg := range os.Args {
		if arg == "--ste-argument-helper" {
			marker = i
			break
		}
	}
	if marker < 0 {
		return
	}
	if os.Args[marker+1] == "stdin" {
		body, err := io.ReadAll(os.Stdin)
		if err != nil {
			os.Exit(2)
		}
		fmt.Print(string(body))
		os.Exit(0)
	}
	arg := strings.TrimPrefix(os.Args[marker+2], "--prompt=")
	var path, hash string
	for _, line := range strings.Split(arg, "\n") {
		if value, ok := strings.CutPrefix(line, "Prompt file (quoted path): "); ok {
			path, _ = strconv.Unquote(value)
		}
		if value, ok := strings.CutPrefix(line, "Prompt SHA256: "); ok {
			hash = value
		}
	}
	if path == "" || !filepath.IsAbs(path) {
		os.Exit(3)
	}
	body, err := os.ReadFile(path)
	if err != nil || protocolpacket.Hash(string(body)) != hash {
		os.Exit(4)
	}
	fmt.Print(string(body))
	os.Exit(0)
}

func TestLargeHeadlessPromptReachesRealChildWithExactAttestedBytes(t *testing.T) {
	for _, boundary := range []string{"measured", "supervised"} {
		for _, template := range []string{"{prompt}", "--prompt={prompt}"} {
			t.Run(boundary+template, func(t *testing.T) {
				root := t.TempDir()
				writeLaunchProtocol(t, root)
				prompt := "-leading option\n> pridaj všetko\r\n---\nste_output: false\n---\n```sh\nprintf '%s' '{root}'\n```\n" + strings.Repeat("Exact evidence.\n", 10000)
				prepared, _, err := prepareProtocolPrompt(root, prompt, LaunchInfo{})
				if err != nil {
					t.Fatal(err)
				}
				if len(prepared) <= 128<<10 {
					t.Fatal("fixture must exceed Linux's single-argument limit")
				}
				agent := agents.Discovery{Spec: agents.Spec{ID: "fixture", PromptMode: agents.PromptArg,
					HeadlessArgs: []string{"-test.run=^TestSTEArgumentHelper$", "--", "--ste-argument-helper", "file", template}}, Path: os.Args[0], Found: true}
				var outputPath string
				if boundary == "measured" {
					r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: agent, Prompt: prompt})
					if err != nil {
						t.Fatal(err)
					}
					if r.Metadata.Context.Mode != "full" || r.StartedAt == nil || r.Outcome.ExitCode == nil || *r.Outcome.ExitCode != 0 {
						t.Fatalf("missing real full-context launch: %+v", r)
					}
					outputPath = filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, "stdout.log")
				} else {
					outputPath = filepath.Join(root, "stdout")
					_, err := execAgentProcess(context.Background(), root, "fixture", agent.ID, "", agent, prompt, outputPath, filepath.Join(root, "stderr"), nil, nil, SupervisionConfig{}, supervisionHooks{})
					if err != nil {
						t.Fatal(err)
					}
				}
				body, err := os.ReadFile(outputPath)
				if err != nil || string(body) != prepared {
					t.Fatalf("child did not read the exact prepared protocol and task: %v", err)
				}
				files, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "headless-prompt-*.md"))
				if err != nil || len(files) != 1 {
					t.Fatalf("expected one retained private prompt: %v %v", files, err)
				}
				info, err := os.Stat(files[0])
				if err != nil || (runtime.GOOS != "windows" && info.Mode().Perm() != 0600) {
					t.Fatalf("prompt file must be private: %v %v", info, err)
				}
			})
		}
	}
}

func TestHeadlessArgumentBoundsPreserveShortInputsAndRefuseUnsafeDelivery(t *testing.T) {
	root := t.TempDir()
	prompt := "-first\n{root}\n{prompt}\n100% 'quoted'"
	args, err := headlessPromptArgs([]string{"--cwd", "{root}", "--prompt={prompt}", "{prompt}"}, "/dir/{prompt}", prompt, root)
	if err != nil || len(args) != 4 || args[1] != "/dir/{prompt}" || args[2] != "--prompt="+prompt || args[3] != prompt {
		t.Fatalf("short inputs or embedded placeholders changed: %q %v", args, err)
	}
	files, _ := os.ReadDir(root)
	if len(files) != 0 {
		t.Fatal("short input created a file")
	}
	large := strings.Repeat("x", 128<<10)
	for _, tc := range []struct {
		name, dir string
		args      []string
	}{
		{"missing-directory", filepath.Join(root, "absent"), []string{"{prompt}"}},
		{"large-literal", root, []string{large}},
		{"large-prefix", root, []string{large + "{prompt}"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args, err := headlessPromptArgs(tc.args, root, large, tc.dir)
			if err == nil || args != nil {
				t.Fatalf("unsafe delivery did not refuse: %v", err)
			}
		})
	}
}

func TestLargeStdinRemainsExactWithInheritedPromptArgument(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	prompt := strings.Repeat("stdin remains exact\n", 8000)
	agent := agents.Discovery{Spec: agents.Spec{ID: "fixture", PromptMode: agents.PromptStdin,
		HeadlessArgs: []string{"-test.run=^TestSTEArgumentHelper$", "--", "--ste-argument-helper", "stdin", "--prompt={prompt}"}}, Path: os.Args[0], Found: true}
	prepared, _, err := prepareProtocolPrompt(root, prompt, LaunchInfo{})
	if err != nil {
		t.Fatal(err)
	}
	cmd, cleanup, err := CommandFor(context.Background(), root, agent, prompt)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	for _, arg := range cmd.Args {
		if len(arg) >= inlinePromptArgumentLimit {
			t.Fatal("oversized inherited argv was not bounded")
		}
	}
	got, err := cmd.Output()
	if err != nil || string(got) != prepared {
		t.Fatalf("stdin bytes changed: %v", err)
	}
}

func TestLargePromptDeliveryCoversEveryBuiltinArgumentTemplate(t *testing.T) {
	for _, spec := range agents.DefaultSpecs() {
		t.Run(spec.ID, func(t *testing.T) {
			dir := t.TempDir()
			resolved, _ := spec.ResolveLaunchArgs()
			usesPrompt := false
			for _, arg := range resolved {
				usesPrompt = usesPrompt || strings.Contains(arg, "{prompt}")
			}
			prompt := strings.Repeat("payload\n", 20000)
			args, err := headlessPromptArgs(resolved, dir, prompt, dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, arg := range args {
				if len(arg) >= inlinePromptArgumentLimit {
					t.Fatal("an adapter retained an oversized argument")
				}
			}
			files, _ := filepath.Glob(filepath.Join(dir, "headless-prompt-*.md"))
			if !usesPrompt {
				if len(files) != 0 || strings.Join(args, "\x00") != strings.Join(expandHeadlessArgs(resolved, dir, prompt), "\x00") {
					t.Fatal("stdin-only adapter changed")
				}
				return
			}
			if len(files) != 1 {
				t.Fatalf("want one complete prompt file, got %v", files)
			}
			body, err := os.ReadFile(files[0])
			if err != nil || string(body) != prompt || !strings.Contains(strings.Join(args, "\n"), protocolpacket.Hash(prompt)) {
				t.Fatalf("adapter lost prompt bytes or digest: %v", err)
			}
		})
	}
}
