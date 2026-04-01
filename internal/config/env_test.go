package config

import "testing"

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
