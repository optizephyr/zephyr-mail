package config

import "fmt"

func ValidateIMAP(cfg Config) error {
	if cfg.IMAPUser == "" || cfg.IMAPPass == "" {
		return fmt.Errorf("Missing IMAP_USER or IMAP_PASS environment variables")
	}
	return nil
}

func ValidateSMTP(cfg Config) error {
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" || cfg.SMTPPass == "" {
		return fmt.Errorf("Missing SMTP configuration. Please set SMTP_HOST, SMTP_USER, and SMTP_PASS in .env")
	}
	return nil
}
