package cli

import (
	"github.com/netease/zephyr-mail/internal/imap"
	"github.com/netease/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
)

func newListMailboxesCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "list-mailboxes",
		Short:         "List available mailboxes",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runListMailboxes,
	}
}

func runListMailboxes(cmd *cobra.Command, _ []string) error {
	clientCfg, err := loadIMAPConfig()
	if err != nil {
		return err
	}

	client, err := connectIMAPClient(clientCfg)
	if err != nil {
		return err
	}
	defer func() { _ = client.Logout() }()

	result, err := imap.ListMailboxes(client)
	if err != nil {
		return err
	}

	output.PrintJSON(result)
	return nil
}
