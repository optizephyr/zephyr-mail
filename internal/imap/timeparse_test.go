package imap

import (
	"regexp"
	"testing"
)

func TestParseRelativeTime(t *testing.T) {
	got, err := ParseRelativeTime("2h")
	if err != nil {
		t.Fatal(err)
	}

	if !regexp.MustCompile(`^\d{2}-[A-Z][a-z]{2}-\d{4}$`).MatchString(got) {
		t.Fatalf("unexpected format: %s", got)
	}
}

func TestParseRelativeTimeInvalidFormat(t *testing.T) {
	_, err := ParseRelativeTime("abc")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}

	if err.Error() != "Invalid time format. Use: 30m, 2h, 7d" {
		t.Fatalf("unexpected error: %v", err)
	}
}
