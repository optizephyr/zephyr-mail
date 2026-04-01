package imap

import "testing"

func TestNewClientMissingCredentialsReturnsError(t *testing.T) {
	_, err := NewClient(ClientConfig{Username: "", Password: ""})
	if err == nil {
		t.Fatal("expected error for missing credentials")
	}

	if err.Error() != "Missing IMAP_USER or IMAP_PASS environment variables" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewClientAppliesDefaults(t *testing.T) {
	client, err := NewClient(ClientConfig{
		Username: "user@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}

	cfg := client.Config()
	if cfg.Host != "127.0.0.1" {
		t.Fatalf("expected default host, got %q", cfg.Host)
	}
	if cfg.Port != 993 {
		t.Fatalf("expected default port, got %d", cfg.Port)
	}
	if cfg.Mailbox != "INBOX" {
		t.Fatalf("expected default mailbox, got %q", cfg.Mailbox)
	}
}

func TestNewClientPreservesProvidedValues(t *testing.T) {
	client, err := NewClient(ClientConfig{
		Host:     "imap.example.com",
		Port:     1993,
		Mailbox:  "Archive",
		Username: "user@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}

	cfg := client.Config()
	if cfg.Host != "imap.example.com" {
		t.Fatalf("expected provided host, got %q", cfg.Host)
	}
	if cfg.Port != 1993 {
		t.Fatalf("expected provided port, got %d", cfg.Port)
	}
	if cfg.Mailbox != "Archive" {
		t.Fatalf("expected provided mailbox, got %q", cfg.Mailbox)
	}
}
