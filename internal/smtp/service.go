package smtp

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	stdsmtp "net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

func ResolveSendRequest(req SendRequest) (ResolvedSendRequest, error) {
	resolved := ResolvedSendRequest{From: defaultFrom(req.From)}

	if req.SubjectFile != "" {
		b, err := os.ReadFile(req.SubjectFile)
		if err != nil {
			return resolved, err
		}
		resolved.Subject = strings.TrimSpace(string(b))
	} else if req.Subject != "" {
		resolved.Subject = req.Subject
	} else {
		resolved.Subject = "(no subject)"
	}

	resolved.To = splitCSV(req.To)
	resolved.Cc = splitCSV(req.Cc)
	resolved.Bcc = splitCSV(req.Bcc)

	if req.BodyFile != "" {
		b, err := os.ReadFile(req.BodyFile)
		if err != nil {
			return resolved, err
		}
		content := string(b)
		if strings.HasSuffix(strings.ToLower(req.BodyFile), ".html") || req.HTML {
			resolved.HTMLBody = content
		} else {
			resolved.TextBody = content
		}
	} else if req.HTMLFile != "" {
		b, err := os.ReadFile(req.HTMLFile)
		if err != nil {
			return resolved, err
		}
		resolved.HTMLBody = string(b)
	} else if req.Body != "" {
		resolved.TextBody = req.Body
	}

	for _, path := range splitCSV(req.Attach) {
		resolved.Attachments = append(resolved.Attachments, Attachment{Filename: filepath.Base(path), Path: path})
	}

	return resolved, nil
}

func (c *Client) Send(req SendRequest) (SendResult, error) {
	resolved, err := ResolveSendRequest(req)
	if err != nil {
		return SendResult{}, err
	}
	if len(resolved.To)+len(resolved.Cc)+len(resolved.Bcc) == 0 {
		return SendResult{}, fmt.Errorf("missing recipient")
	}

	client, err := c.dial()
	if err != nil {
		return SendResult{}, err
	}
	defer client.Close()

	if err := client.Auth(stdsmtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)); err != nil {
		return SendResult{}, err
	}

	msg, err := buildMessage(resolved)
	if err != nil {
		return SendResult{}, err
	}

	if err := client.Mail(resolved.From); err != nil {
		return SendResult{}, err
	}
	for _, addr := range append(append(append([]string{}, resolved.To...), resolved.Cc...), resolved.Bcc...) {
		if err := client.Rcpt(addr); err != nil {
			return SendResult{}, err
		}
	}

	w, err := client.Data()
	if err != nil {
		return SendResult{}, err
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return SendResult{}, err
	}
	if err := w.Close(); err != nil {
		return SendResult{}, err
	}

	return SendResult{Success: true, To: resolved.To}, nil
}

func (c *Client) TestConnection() (SendResult, error) {
	from := c.cfg.From
	if from == "" {
		from = c.cfg.Username
	}
	return c.Send(SendRequest{
		From:    from,
		To:      c.cfg.Username,
		Subject: "SMTP Connection Test",
		Body:    "This is a test email from the IMAP/SMTP email skill.",
		HTML:    true,
	})
}

func buildMessage(resolved ResolvedSendRequest) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", resolved.From))
	if len(resolved.To) > 0 {
		buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(resolved.To, ", ")))
	}
	if len(resolved.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(resolved.Cc, ", ")))
	}
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", resolved.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")

	if len(resolved.Attachments) > 0 {
		mw := multipart.NewWriter(&buf)
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%q\r\n\r\n", mw.Boundary()))
		if err := writeBodyPart(mw, resolved); err != nil {
			return nil, err
		}
		for _, attachment := range resolved.Attachments {
			data, err := os.ReadFile(attachment.Path)
			if err != nil {
				return nil, err
			}
			hdr := textproto.MIMEHeader{}
			hdr.Set("Content-Type", contentTypeForFile(attachment.Filename))
			hdr.Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, attachment.Filename))
			part, err := mw.CreatePart(hdr)
			if err != nil {
				return nil, err
			}
			if _, err := part.Write(data); err != nil {
				return nil, err
			}
		}
		if err := mw.Close(); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	if resolved.HTMLBody != "" {
		buf.WriteString("Content-Type: text/html; charset=utf-8\r\n\r\n")
		buf.WriteString(resolved.HTMLBody)
		return buf.Bytes(), nil
	}

	buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	buf.WriteString(resolved.TextBody)
	return buf.Bytes(), nil
}

func writeBodyPart(mw *multipart.Writer, resolved ResolvedSendRequest) error {
	hdr := textproto.MIMEHeader{}
	if resolved.HTMLBody != "" {
		hdr.Set("Content-Type", "text/html; charset=utf-8")
	} else {
		hdr.Set("Content-Type", "text/plain; charset=utf-8")
	}
	part, err := mw.CreatePart(hdr)
	if err != nil {
		return err
	}
	if resolved.HTMLBody != "" {
		_, err = part.Write([]byte(resolved.HTMLBody))
	} else {
		_, err = part.Write([]byte(resolved.TextBody))
	}
	return err
}

func contentTypeForFile(filename string) string {
	if ext := strings.ToLower(filepath.Ext(filename)); ext != "" {
		if ctype := mime.TypeByExtension(ext); ctype != "" {
			return ctype
		}
	}
	return "application/octet-stream"
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func defaultFrom(from string) string {
	if from != "" {
		return from
	}
	if envFrom := os.Getenv("SMTP_FROM"); envFrom != "" {
		return envFrom
	}
	return os.Getenv("SMTP_USER")
}
