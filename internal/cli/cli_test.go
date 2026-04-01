package cli

import "testing"

func TestCheckBareUnseenDoesNotEnableUnseen(t *testing.T) {
	opts := parseCheckFlags([]string{"--unseen"})
	if opts.UnseenRaw == "true" {
		t.Fatal("bare --unseen must not be treated as true")
	}
}

func TestSearchPresenceBooleans(t *testing.T) {
	opts := parseSearchFlags([]string{"--unseen", "--flagged"})
	if !opts.Unseen || !opts.Flagged {
		t.Fatal("presence booleans expected")
	}
}
