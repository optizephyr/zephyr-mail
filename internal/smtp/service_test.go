package smtp

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSubjectFileOverridesSubject(t *testing.T) {
	dir := t.TempDir()
	subjectFile := writeTestFile(t, dir, "subject.txt", "from-file")

	req := SendRequest{Subject: "inline", SubjectFile: subjectFile}
	resolved, err := ResolveSendRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	if resolved.Subject != "from-file" {
		t.Fatalf("got %s", resolved.Subject)
	}
}

func TestBodyFileHtmlDetection(t *testing.T) {
	dir := t.TempDir()
	bodyFile := writeTestFile(t, dir, "body.html", "<p>hello</p>")

	req := SendRequest{BodyFile: bodyFile, HTML: true}
	resolved, err := ResolveSendRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	if resolved.HTMLBody == "" {
		t.Fatal("expected html body")
	}
}

func TestHTMLFileFallback(t *testing.T) {
	dir := t.TempDir()
	htmlFile := writeTestFile(t, dir, "fallback.html", "<p>fallback</p>")

	req := SendRequest{HTMLFile: htmlFile}
	resolved, err := ResolveSendRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	if resolved.HTMLBody != "<p>fallback</p>" {
		t.Fatalf("unexpected html body: %q", resolved.HTMLBody)
	}
}

func TestBodyFallback(t *testing.T) {
	req := SendRequest{Body: "plain text"}
	resolved, err := ResolveSendRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	if resolved.TextBody != "plain text" {
		t.Fatalf("unexpected text body: %q", resolved.TextBody)
	}
}

func TestEmptyBodyFallback(t *testing.T) {
	resolved, err := ResolveSendRequest(SendRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if resolved.TextBody != "" {
		t.Fatalf("expected empty text body, got %q", resolved.TextBody)
	}
}
