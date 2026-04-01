package cli

import "github.com/spf13/cobra"

func Register(root *cobra.Command) {
	root.AddCommand(
		newCheckCmd(),
		newFetchCmd(),
		newDownloadCmd(),
		newSearchCmd(),
		newMarkReadCmd(),
		newMarkUnreadCmd(),
		newListMailboxesCmd(),
		newSendCmd(),
		newTestCmd(),
	)
}
