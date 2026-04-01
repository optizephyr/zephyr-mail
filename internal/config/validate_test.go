package config

import (
	"strings"
	"testing"
)

func TestValidateSMTPMissingHost(t *testing.T) {
	cfg := Config{SMTPUser: "a", SMTPPass: "b"}
	err := ValidateSMTP(cfg)
	if err == nil || !strings.Contains(err.Error(), "Missing SMTP configuration") {
		t.Fatalf("unexpected err: %v", err)
	}
}
