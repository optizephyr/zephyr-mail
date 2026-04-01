package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
