package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func runCLI(args ...string) (int, string, string) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "/tmp/zephyr-mail-test", "./cmd/zephyr-mail")
	buildCmd.Dir = "../.."
	if output, err := buildCmd.CombinedOutput(); err != nil {
		panic(string(output))
	}

	// Run the binary
	cmd := exec.Command("/tmp/zephyr-mail-test", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return exitCode, stdout.String(), stderr.String()
}

func TestRootUnknownCommandExitCode(t *testing.T) {
	code, stdout, stderr := runCLI("unknown-command")
	if code != 1 {
		t.Fatalf("want exit 1, got %d", code)
	}
	if stdout != "" {
		t.Fatalf("want empty stdout, got: %s", stdout)
	}
	if !strings.Contains(stderr, "Unknown command") {
		t.Fatalf("unexpected stderr: %s", stderr)
	}
}
