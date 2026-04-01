package cli

import (
	"fmt"

	"github.com/netease/zephyr-mail/internal/imap"
	"github.com/netease/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
)

func newMarkUnreadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "mark-unread <uid> [uid2...]",
		Short:         "Mark messages as unread",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runMarkUnread,
	}

	cmd.Flags().String("mailbox", "", "Mailbox to inspect")
	return cmd
}

func runMarkUnread(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("UID(s) required: node imap.js mark-unread <uid> [uid2...]")
	}

	mailbox, err := cmd.Flags().GetString("mailbox")
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

	result, err := imap.MarkUnread(client, args, mailbox)
	if err != nil {
		return err
	}

	output.PrintJSON(result)
	return nil
}
