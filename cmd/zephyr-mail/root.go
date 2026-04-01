package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/netease/zephyr-mail/internal/config"
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		// Normalize error output to show "Unknown command" for unknown subcommands
		// Cobra's default behavior varies, so we ensure consistent output
		errMsg := err.Error()
		if strings.Contains(errMsg, "unknown command") || strings.Contains(errMsg, "Unknown command") {
			fmt.Fprintf(os.Stderr, "Unknown command\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
