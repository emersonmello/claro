// Package pull
package pull

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emersonmello/claro/internal"
	"github.com/emersonmello/claro/internal/tui"
	"github.com/spf13/cobra"
)

// Pull represents the pull command
func Pull() *cobra.Command {
	pullCmd := &cobra.Command{
		Use:   "pull <directory-with-student-submissions>",
		Short: "Incorporate changes from students' remote repositories into local copy",
		Long:  tui.LongHelpMsg("Incorporate changes from students' remote repositories into local copy"),
		//Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New(tui.UseErrorMsg("pull"))
			}
			if _, err := tea.NewProgram(internal.NewPullModel(args[0])).Run(); err != nil {
				fmt.Println("Error running program:", err)
			}
			return nil
		},
	}
	return pullCmd
}
