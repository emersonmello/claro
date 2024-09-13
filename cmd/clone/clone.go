// Package clone
package clone

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emersonmello/claro/internal"
	"github.com/emersonmello/claro/internal/tui"
	"github.com/spf13/cobra"
)

// Clone represents the clone command
func Clone() *cobra.Command {
	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone all students assignments from a GitHub Classroom",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !tui.GitHubCliInstalled {
				tui.UserGitHubPAT = internal.GetAndSaveToken()
			}
			if _, err := tea.NewProgram(internal.NewCloneModel()).Run(); err != nil {
				fmt.Println("Error running program:", err)
			}
			return nil
		},
	}
	return cloneCmd
}
