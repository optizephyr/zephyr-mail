package output

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = original

	data, _ := io.ReadAll(r)
	_ = r.Close()

	return string(data)
}

func captureStderr(fn func()) string {
	original := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	_ = w.Close()
	os.Stderr = original

	data, _ := io.ReadAll(r)
	_ = r.Close()

	return string(data)
}

func TestPrintErrorUsesCompatibilityPrefix(t *testing.T) {
	stderr := captureStderr(func() { PrintError(errors.New("boom")) })
	if !strings.Contains(stderr, "Error: boom") {
		t.Fatal(stderr)
	}
}

func TestPrintJSONPretty(t *testing.T) {
	out := captureStdout(func() { PrintJSON(map[string]any{"ok": true}) })
	if !strings.Contains(out, "\n") {
		t.Fatal("expected pretty json")
	}
}
