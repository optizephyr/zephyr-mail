package cli

import (
	"github.com/optizephyr/zephyr-mail/internal/imap"
	"github.com/optizephyr/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type searchFlagOptions struct {
	Mailbox  string
	Unseen   bool
	Seen     bool
	Flagged  bool
	Answered bool
	From     string
	To       string
	Subject  string
	Recent   string
	Since    string
	Before   string
	UID      string
	Limit    int
}

func parseSearchFlags(args []string) searchFlagOptions {
	fs := pflag.NewFlagSet("search", pflag.ContinueOnError)
	mailbox := fs.String("mailbox", "", "")
	unseen := fs.Bool("unseen", false, "")
	seen := fs.Bool("seen", false, "")
	flagged := fs.Bool("flagged", false, "")
	answered := fs.Bool("answered", false, "")
	from := fs.String("from", "", "")
	to := fs.String("to", "", "")
	subject := fs.String("subject", "", "")
	recent := fs.String("recent", "", "")
	since := fs.String("since", "", "")
	before := fs.String("before", "", "")
	uid := fs.String("uid", "", "")
	limit := fs.Int("limit", 100, "")
	_ = fs.Parse(args)

	return searchFlagOptions{
		Mailbox:  *mailbox,
		Unseen:   *unseen,
		Seen:     *seen,
		Flagged:  *flagged,
		Answered: *answered,
		From:     *from,
		To:       *to,
		Subject:  *subject,
		Recent:   *recent,
		Since:    *since,
		Before:   *before,
		UID:      *uid,
		Limit:    *limit,
	}
}

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search",
		Short:         "Search for messages",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runSearch,
	}

	cmd.Flags().String("mailbox", "", "Mailbox to inspect")
	cmd.Flags().Bool("unseen", false, "Only unseen messages")
	cmd.Flags().Bool("seen", false, "Only seen messages")
	cmd.Flags().Bool("flagged", false, "Only flagged messages")
	cmd.Flags().Bool("answered", false, "Only answered messages")
	cmd.Flags().String("from", "", "Search sender")
	cmd.Flags().String("to", "", "Search recipient")
	cmd.Flags().String("subject", "", "Search subject")
	cmd.Flags().String("recent", "", "Relative recent window like 2h or 7d")
	cmd.Flags().String("since", "", "Search messages since date")
	cmd.Flags().String("before", "", "Search messages before date")
	cmd.Flags().String("uid", "", "Search by UID or UID range")
	cmd.Flags().Int("limit", 100, "Maximum messages to return")

	return cmd
}

func runSearch(cmd *cobra.Command, _ []string) error {
	mailbox, err := cmd.Flags().GetString("mailbox")
	if err != nil {
		return err
	}
	unseen, err := cmd.Flags().GetBool("unseen")
	if err != nil {
		return err
	}
	seen, err := cmd.Flags().GetBool("seen")
	if err != nil {
		return err
	}
	flagged, err := cmd.Flags().GetBool("flagged")
	if err != nil {
		return err
	}
	answered, err := cmd.Flags().GetBool("answered")
	if err != nil {
		return err
	}
	from, err := cmd.Flags().GetString("from")
	if err != nil {
		return err
	}
	to, err := cmd.Flags().GetString("to")
	if err != nil {
		return err
	}
	subject, err := cmd.Flags().GetString("subject")
	if err != nil {
		return err
	}
	recent, err := cmd.Flags().GetString("recent")
	if err != nil {
		return err
	}
	since, err := cmd.Flags().GetString("since")
	if err != nil {
		return err
	}
	before, err := cmd.Flags().GetString("before")
	if err != nil {
		return err
	}
	uid, err := cmd.Flags().GetString("uid")
	if err != nil {
		return err
	}
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return err
	}
	if limit <= 0 {
		limit = 100
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

	result, err := imap.Search(client, imap.SearchOptions{
		Mailbox:  mailbox,
		Unseen:   unseen,
		Seen:     seen,
		Flagged:  flagged,
		Answered: answered,
		From:     from,
		To:       to,
		Subject:  subject,
		Recent:   recent,
		Since:    since,
		Before:   before,
		UID:      uid,
		Limit:    limit,
	})
	if err != nil {
		return err
	}

	if len(result.Messages) > limit {
		result.Messages = result.Messages[:limit]
	}

	output.PrintJSON(result)
	return nil
}
