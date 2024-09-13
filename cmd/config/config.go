// Package config
package config

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"github.com/emersonmello/claro/internal"
	"github.com/spf13/cobra"
)

// Config represents the clone command
func Config() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure claro's properties (commit message, filename, etc)",
		RunE:  internal.ConfigCmd,
	}
	return configCmd
}
