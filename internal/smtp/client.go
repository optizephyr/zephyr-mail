package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
	stdsmtp "net/smtp"
	"strconv"
)

type Client struct {
	cfg Config
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("missing SMTP_HOST")
	}
	if cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("missing SMTP_USER or SMTP_PASS")
	}
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	if cfg.From == "" {
		cfg.From = cfg.Username
	}
	return &Client{cfg: cfg}, nil
}

func (c *Client) dial() (*stdsmtp.Client, error) {
	addr := net.JoinHostPort(c.cfg.Host, strconv.Itoa(c.cfg.Port))

	if c.cfg.Secure {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: c.cfg.Host, InsecureSkipVerify: !c.cfg.RejectUnauthorized})
		if err != nil {
			return nil, err
		}
		return stdsmtp.NewClient(conn, c.cfg.Host)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	client, err := stdsmtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: c.cfg.Host, InsecureSkipVerify: !c.cfg.RejectUnauthorized}); err != nil {
			_ = client.Close()
			return nil, err
		}
	}
	return client, nil
}
