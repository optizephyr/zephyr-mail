package cli

import (
	"github.com/optizephyr/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "test",
		Short:         "Test SMTP connectivity",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runTest,
	}
}

func runTest(cmd *cobra.Command, _ []string) error {
	client, err := loadSMTPClient()
	if err != nil {
		return err
	}

	result, err := client.TestConnection()
	if err != nil {
		return err
	}

	output.PrintJSON(result)
	return nil
}
