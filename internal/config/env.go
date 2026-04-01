package config

import (
	"os"
	"path/filepath"
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
	if err := loadDotEnvWithFallback(); err != nil {
		return Config{}, err
	}

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

func loadDotEnvWithFallback() error {
	paths := []string{".env"}

	if executablePath, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(executablePath), ".env"))
	}

	for _, envPath := range paths {
		if _, err := os.Stat(envPath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		if err := godotenv.Load(envPath); err != nil {
			return err
		}
		return nil
	}

	return nil
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
