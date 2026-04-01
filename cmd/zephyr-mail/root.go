package main

import (
	"fmt"
	"os"

	"github.com/netease/zephyr-mail/internal/cli"
	"github.com/netease/zephyr-mail/internal/common"
	"github.com/netease/zephyr-mail/internal/config"
	"github.com/netease/zephyr-mail/internal/output"
	"github.com/spf13/cobra"
)

var appConfig config.Config

var rootCmd = &cobra.Command{
	Use:   "zephyr-mail",
	Short: "Zephyr Mail - IMAP/SMTP email CLI tool",
	Long:  `A CLI tool for sending and receiving email via IMAP and SMTP protocols.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If we reach here, no valid subcommand was provided
		if len(args) == 0 {
			return fmt.Errorf("no command specified")
		}
		return fmt.Errorf("Unknown command")
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	var err error
	appConfig, err = config.LoadFromEnv()
	if err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
	cli.Register(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		normalized := common.NormalizeCLIError(err)
		if common.IsUnknownCommandError(normalized) {
			fmt.Fprintln(os.Stderr, "Unknown command")
		} else {
			output.PrintError(normalized)
		}
		os.Exit(common.ExitCode(normalized))
	}
}
