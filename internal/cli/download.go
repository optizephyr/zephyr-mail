package cli

import (
	"fmt"

	"github.com/optizephyr/zephyr-mail/internal/imap"
	"github.com/optizephyr/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
)

func newDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "download <uid>",
		Short:         "Download attachment metadata for a message",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runDownload,
	}

	cmd.Flags().String("mailbox", "", "Mailbox to inspect")
	cmd.Flags().String("dir", ".", "Output directory")
	cmd.Flags().String("file", "", "Specific filename to match")
	return cmd
}

func runDownload(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("UID required: node imap.js download <uid>")
	}

	mailbox, err := cmd.Flags().GetString("mailbox")
	if err != nil {
		return err
	}
	outDir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}
	filename, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	clientCfg, err := loadIMAPConfig()
	if err != nil {
		return err
	}
	mailbox = resolveMailbox(mailbox, clientCfg.Mailbox)

	client, err := connectIMAPClient(clientCfg)
	if err != nil {
		return err
	}
	defer func() { _ = client.Logout() }()

	result, err := imap.Download(client, args[0], mailbox, outDir, filename)
	if err != nil {
		return err
	}

	output.PrintJSON(result)
	return nil
}
