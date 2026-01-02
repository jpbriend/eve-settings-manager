package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "esm",
	Short: "Eve Settings Manager - Manage Eve Online character settings",
	Long: `Eve Settings Manager (esm) is a CLI tool to manage Eve Online character settings.

It supports listing, copying, backing up, and restoring character-specific settings
(core_char_*.dat files) across different accounts and installations.

Works with both Steam and non-Steam versions on Windows and Linux.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(restoreCmd)
}
