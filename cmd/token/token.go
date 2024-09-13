// Package token
package token

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"github.com/emersonmello/claro/internal"
	"github.com/spf13/cobra"
)

// Token represents the token command
func Token() *cobra.Command {
	tokenCmd := &cobra.Command{
		Use:       "token <add|del>",
		Short:     "add or remove a claro's GitHub Personal Access Token in the OS Keychain",
		ValidArgs: []string{"add", "del"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:      configureToken,
	}
	return tokenCmd
}

func configureToken(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "add":
		internal.AddTokenToKeyring()
	case "del":
		internal.DeleteTokenFromKeyring()
	}
	return nil
}
