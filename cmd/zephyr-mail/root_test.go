package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args ...string) (int, string, string) {
	// Build the binary to a temp location
	binaryPath := t.TempDir() + "/zephyr-mail-test"
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "github.com/optizephyr/zephyr-mail/cmd/zephyr-mail")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %s", string(output))
	}

	// Run the binary
	cmd := exec.Command(binaryPath, args...)

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
	code, stdout, stderr := runCLI(t, "unknown-command")
	if code != 1 {
		t.Fatalf("want exit 1, got %d", code)
	}
	if stdout != "" {
		t.Fatalf("want empty stdout, got: %s", stdout)
	}
	if !strings.Contains(stderr, "Unknown command") {
		t.Fatalf("unexpected stderr: %s", stderr)
	}
	if strings.Contains(stderr, "Error:") {
		t.Fatalf("unexpected error prefix for unknown command: %s", stderr)
	}
}
