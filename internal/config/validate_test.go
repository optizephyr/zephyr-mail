package config

import (
	"testing"
)

func TestValidateSMTPMissingHost(t *testing.T) {
	cfg := Config{SMTPUser: "a", SMTPPass: "b"}
	err := ValidateSMTP(cfg)
	const want = "Missing SMTP configuration. Please set SMTP_HOST, SMTP_USER, and SMTP_PASS in .env"
	if err == nil || err.Error() != want {
		t.Fatalf("unexpected err: %v", err)
	}
}
