package imap

import "fmt"

type Client struct {
	config ClientConfig
}

func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("Missing IMAP_USER or IMAP_PASS environment variables")
	}

	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}

	if cfg.Port == 0 {
		cfg.Port = 993
	}

	if cfg.Mailbox == "" {
		cfg.Mailbox = DefaultMailbox
	}

	return &Client{config: cfg}, nil
}

func (c *Client) Config() ClientConfig {
	return c.config
}
