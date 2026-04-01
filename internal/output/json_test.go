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

func TestPrintJSONDoesNotHTMLEscape(t *testing.T) {
	out := captureStdout(func() { PrintJSON(map[string]any{"text": "a<b>&c"}) })
	if strings.Contains(out, "\\u003c") || strings.Contains(out, "\\u003e") || strings.Contains(out, "\\u0026") {
		t.Fatalf("expected unescaped html chars, got: %s", out)
	}
	if !strings.Contains(out, "a<b>&c") {
		t.Fatalf("expected raw html chars, got: %s", out)
	}
}
