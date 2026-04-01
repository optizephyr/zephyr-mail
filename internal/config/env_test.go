package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvWithSetEnvIsDeterministic(t *testing.T) {
	t.Setenv("IMAP_HOST", "imap.example.com")
	t.Setenv("IMAP_PORT", "993")
	t.Setenv("IMAP_USER", "user@example.com")
	t.Setenv("IMAP_PASS", "secret")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_USER", "user@example.com")
	t.Setenv("SMTP_PASS", "secret")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.IMAPHost != "imap.example.com" || cfg.SMTPHost != "smtp.example.com" {
		t.Fatalf("unexpected parsed config: %+v", cfg)
	}
}

func TestLoadEnvDefaultsAndParsingEdges(t *testing.T) {
	t.Setenv("IMAP_HOST", "")
	t.Setenv("IMAP_PORT", "not-a-number")
	t.Setenv("IMAP_TLS", "true")
	t.Setenv("IMAP_REJECT_UNAUTHORIZED", "false")
	t.Setenv("IMAP_MAILBOX", "")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_PORT", "0")
	t.Setenv("SMTP_SECURE", "true")
	t.Setenv("SMTP_REJECT_UNAUTHORIZED", "false")
	t.Setenv("SMTP_FROM", "")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.IMAPHost != "127.0.0.1" {
		t.Fatalf("unexpected IMAP host default: %q", cfg.IMAPHost)
	}
	if cfg.IMAPPort != 993 {
		t.Fatalf("unexpected IMAP port fallback: %d", cfg.IMAPPort)
	}
	if !cfg.IMAPTLS {
		t.Fatalf("unexpected IMAP TLS parse: %v", cfg.IMAPTLS)
	}
	if cfg.IMAPRejectUnauthorized {
		t.Fatalf("unexpected IMAP rejectUnauthorized parse: %v", cfg.IMAPRejectUnauthorized)
	}
	if cfg.IMAPMailbox != "INBOX" {
		t.Fatalf("unexpected IMAP mailbox default: %q", cfg.IMAPMailbox)
	}
	if cfg.SMTPPort != 587 {
		t.Fatalf("unexpected SMTP port fallback: %d", cfg.SMTPPort)
	}
	if !cfg.SMTPSecure {
		t.Fatalf("unexpected SMTP secure parse: %v", cfg.SMTPSecure)
	}
	if cfg.SMTPRejectUnauthorized {
		t.Fatalf("unexpected SMTP rejectUnauthorized parse: %v", cfg.SMTPRejectUnauthorized)
	}
}

func TestLoadEnvParseErrorFromCWDEnvIsReturned(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tempDir := t.TempDir()
	badEnvPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(badEnvPath, []byte("BROKEN='unterminated\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	_, err = LoadFromEnv()
	if err == nil {
		t.Fatal("expected parse error from malformed .env")
	}
}
