package cli

import (
	"fmt"

	"github.com/optizephyr/zephyr-mail/internal/output"
	smtpsvc "github.com/optizephyr/zephyr-mail/internal/smtp"
	"github.com/spf13/cobra"
)

func newSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "send",
		Short:         "Send an email via SMTP",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runSend,
	}

	cmd.Flags().String("from", "", "Override sender address")
	cmd.Flags().String("to", "", "Recipient address")
	cmd.Flags().String("cc", "", "CC recipients")
	cmd.Flags().String("bcc", "", "BCC recipients")
	cmd.Flags().String("subject", "", "Email subject")
	cmd.Flags().String("subject-file", "", "Read subject from file")
	cmd.Flags().String("body", "", "Plain text body")
	cmd.Flags().String("body-file", "", "Read body from file")
	cmd.Flags().String("html-file", "", "Read HTML body from file")
	cmd.Flags().Bool("html", false, "Treat body as HTML")
	cmd.Flags().String("attach", "", "Comma-separated attachment paths")

	return cmd
}

func runSend(cmd *cobra.Command, _ []string) error {
	from, err := cmd.Flags().GetString("from")
	if err != nil {
		return err
	}
	to, err := cmd.Flags().GetString("to")
	if err != nil {
		return err
	}
	if to == "" {
		return fmt.Errorf("Missing required option: --to <email>")
	}
	cc, err := cmd.Flags().GetString("cc")
	if err != nil {
		return err
	}
	bcc, err := cmd.Flags().GetString("bcc")
	if err != nil {
		return err
	}
	subject, err := cmd.Flags().GetString("subject")
	if err != nil {
		return err
	}
	subjectFile, err := cmd.Flags().GetString("subject-file")
	if err != nil {
		return err
	}
	if subject == "" && subjectFile == "" {
		return fmt.Errorf("Missing required option: --subject <text> or --subject-file <file>")
	}
	body, err := cmd.Flags().GetString("body")
	if err != nil {
		return err
	}
	bodyFile, err := cmd.Flags().GetString("body-file")
	if err != nil {
		return err
	}
	htmlFile, err := cmd.Flags().GetString("html-file")
	if err != nil {
		return err
	}
	html, err := cmd.Flags().GetBool("html")
	if err != nil {
		return err
	}
	attach, err := cmd.Flags().GetString("attach")
	if err != nil {
		return err
	}

	client, err := loadSMTPClient()
	if err != nil {
		return err
	}

	result, err := client.Send(smtpsvc.SendRequest{
		From:        from,
		To:          to,
		Cc:          cc,
		Bcc:         bcc,
		Subject:     subject,
		SubjectFile: subjectFile,
		Body:        body,
		BodyFile:    bodyFile,
		HTMLFile:    htmlFile,
		HTML:        html,
		Attach:      attach,
	})
	if err != nil {
		return err
	}

	output.PrintJSON(result)
	return nil
}
