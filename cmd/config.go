package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure claro's properties (github token, commit message, etc.)",
}

func init() {
	rootCmd.AddCommand(configCmd)
}
