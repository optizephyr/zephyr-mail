package cli

import (
	"github.com/optizephyr/zephyr-mail/internal/imap"
	"github.com/optizephyr/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type checkFlagOptions struct {
	Mailbox   string
	Limit     int
	Recent    string
	UnseenRaw string
}

func parseCheckFlags(args []string) checkFlagOptions {
	fs := pflag.NewFlagSet("check", pflag.ContinueOnError)
	mailbox := fs.String("mailbox", "", "")
	limit := fs.Int("limit", 10, "")
	recent := fs.String("recent", "", "")
	unseen := fs.String("unseen", "", "")
	fs.Lookup("unseen").NoOptDefVal = "false"
	_ = fs.Parse(args)

	return checkFlagOptions{
		Mailbox:   *mailbox,
		Limit:     *limit,
		Recent:    *recent,
		UnseenRaw: *unseen,
	}
}

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "check",
		Short:         "Check for new or unread mail",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runCheck,
	}

	cmd.Flags().String("mailbox", "", "Mailbox to inspect")
	cmd.Flags().Int("limit", 10, "Maximum messages to return")
	cmd.Flags().String("recent", "", "Relative recent window like 2h or 7d")
	cmd.Flags().String("unseen", "", "Only include unseen messages when set to true")
	cmd.Flags().Lookup("unseen").NoOptDefVal = "false"

	return cmd
}

func runCheck(cmd *cobra.Command, _ []string) error {
	mailbox, err := cmd.Flags().GetString("mailbox")
	if err != nil {
		return err
	}
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return err
	}
	recent, err := cmd.Flags().GetString("recent")
	if err != nil {
		return err
	}
	unseenRaw, err := cmd.Flags().GetString("unseen")
	if err != nil {
		return err
	}

	clientCfg, err := loadIMAPConfig()
	if err != nil {
		return err
	}
	if mailbox == "" {
		mailbox = clientCfg.Mailbox
	}

	client, err := connectIMAPClient(clientCfg)
	if err != nil {
		return err
	}
	defer func() { _ = client.Logout() }()

	result, err := imap.Check(client, imap.CheckOptions{
		Mailbox:   mailbox,
		Limit:     limit,
		Recent:    recent,
		UnseenRaw: unseenRaw,
	})
	if err != nil {
		return err
	}

	output.PrintJSON(result)
	return nil
}
