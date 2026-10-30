package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCaptureOnlyMachineConfiguration(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	other := t.TempDir()
	t.Setenv("PARLEY_HOME", other)
	t.Setenv("PARLEY_CAPTURE_ZCODE_FAILURES", "true")
	os.WriteFile(filepath.Join(other, "agents.toml"), []byte("[defaults]\ncapture_zcode_failures = true\n"), 0600)
	enabled, dir, err := MachineZcodeCapture()
	if err != nil || enabled || dir != filepath.Join(home, ".parley", "captures", "zcode") {
		t.Fatalf("%v %s %v", enabled, dir, err)
	}
	os.Mkdir(filepath.Join(home, ".parley"), 0700)
	path := filepath.Join(home, ".parley", "agents.toml")
	for _, value := range []string{"true", "false", "'true'"} {
		os.WriteFile(path, []byte("[defaults]\ncapture_zcode_failures = "+value+"\n"), 0600)
		enabled, _, err = MachineZcodeCapture()
		if value == "'true'" {
			if err == nil {
				t.Fatal("string enabled capture")
			}
			continue
		}
		if err != nil || enabled != (value == "true") {
			t.Fatalf("%s %v %v", value, enabled, err)
		}
	}
}
