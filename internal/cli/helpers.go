package cli

import (
	"crypto/tls"
	"net"
	"strconv"

	imapclient "github.com/emersion/go-imap/client"
	"github.com/optizephyr/zephyr-mail/internal/config"
	imapsvc "github.com/optizephyr/zephyr-mail/internal/imap"
	smtpsvc "github.com/optizephyr/zephyr-mail/internal/smtp"
)

func loadIMAPConfig() (imapsvc.ClientConfig, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return imapsvc.ClientConfig{}, err
	}

	client, err := imapsvc.NewClient(imapsvc.ClientConfig{
		Host:               cfg.IMAPHost,
		Port:               cfg.IMAPPort,
		Username:           cfg.IMAPUser,
		Password:           cfg.IMAPPass,
		TLS:                cfg.IMAPTLS,
		RejectUnauthorized: cfg.IMAPRejectUnauthorized,
		Mailbox:            cfg.IMAPMailbox,
	})
	if err != nil {
		return imapsvc.ClientConfig{}, err
	}

	return client.Config(), nil
}

func connectIMAPClient(cfg imapsvc.ClientConfig) (*imapclient.Client, error) {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	var client *imapclient.Client
	var err error
	if cfg.TLS {
		client, err = imapclient.DialTLS(addr, &tls.Config{
			ServerName:         cfg.Host,
			InsecureSkipVerify: !cfg.RejectUnauthorized,
		})
	} else {
		client, err = imapclient.Dial(addr)
	}
	if err != nil {
		return nil, err
	}

	if err := client.Login(cfg.Username, cfg.Password); err != nil {
		_ = client.Logout()
		return nil, err
	}

	return client, nil
}

func loadSMTPClient() (*smtpsvc.Client, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, err
	}

	return smtpsvc.NewClient(smtpsvc.Config{
		Host:               cfg.SMTPHost,
		Port:               cfg.SMTPPort,
		Secure:             cfg.SMTPSecure,
		RejectUnauthorized: cfg.SMTPRejectUnauthorized,
		Username:           cfg.SMTPUser,
		Password:           cfg.SMTPPass,
		From:               cfg.SMTPFrom,
	})
}

func resolveMailbox(flagMailbox, defaultMailbox string) string {
	if flagMailbox != "" {
		return flagMailbox
	}
	return defaultMailbox
}
