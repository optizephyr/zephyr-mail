package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	IMAPHost               string
	IMAPPort               int
	IMAPUser               string
	IMAPPass               string
	IMAPTLS                bool
	IMAPRejectUnauthorized bool
	IMAPMailbox            string
	SMTPHost               string
	SMTPPort               int
	SMTPUser               string
	SMTPPass               string
	SMTPSecure             bool
	SMTPRejectUnauthorized bool
	SMTPFrom               string
}

func LoadFromEnv() (Config, error) {
	_ = godotenv.Load(".env")

	cfg := Config{
		IMAPHost:               getEnvDefault("IMAP_HOST", "127.0.0.1"),
		IMAPPort:               parseIntDefault("IMAP_PORT", 993),
		IMAPUser:               os.Getenv("IMAP_USER"),
		IMAPPass:               os.Getenv("IMAP_PASS"),
		IMAPTLS:                os.Getenv("IMAP_TLS") == "true",
		IMAPRejectUnauthorized: os.Getenv("IMAP_REJECT_UNAUTHORIZED") != "false",
		IMAPMailbox:            getEnvDefault("IMAP_MAILBOX", "INBOX"),
		SMTPHost:               os.Getenv("SMTP_HOST"),
		SMTPPort:               parseIntDefault("SMTP_PORT", 587),
		SMTPUser:               os.Getenv("SMTP_USER"),
		SMTPPass:               os.Getenv("SMTP_PASS"),
		SMTPSecure:             os.Getenv("SMTP_SECURE") == "true",
		SMTPRejectUnauthorized: os.Getenv("SMTP_REJECT_UNAUTHORIZED") != "false",
		SMTPFrom:               os.Getenv("SMTP_FROM"),
	}

	return cfg, nil
}

func getEnvDefault(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func parseIntDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(v)
	if err != nil || parsed == 0 {
		return fallback
	}

	return parsed
}
