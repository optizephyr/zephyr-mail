package cli

import (
	"fmt"

	"github.com/optizephyr/zephyr-mail/internal/imap"
	"github.com/optizephyr/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
)

func newFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "fetch <uid>",
		Short:         "Fetch a message by UID",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runFetch,
	}

	cmd.Flags().String("mailbox", "", "Mailbox to inspect")
	return cmd
}

func runFetch(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("UID required: node imap.js fetch <uid>")
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

	result, err := imap.Fetch(client, args[0], mailbox)
	if err != nil {
		return err
	}

	output.PrintJSON(result)
	return nil
}
