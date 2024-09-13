// Package push
package push

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

// Push represents the push command
func Push() *cobra.Command {
	pushCmd := &cobra.Command{
		Use:   "push <directory-with-student-submissions>",
		Short: "Add, commit, and push the grading file to each student's remote repository",
		Long:  tui.LongHelpMsg("Use this command to add, commit, and push the grading file for each student's repository to the remote repository"),
		//Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New(tui.UseErrorMsg("push"))
			}
			if _, err := tea.NewProgram(internal.NewPushModel(args[0])).Run(); err != nil {
				fmt.Println("Error running program:", err)
			}
			return nil
		},
	}
	return pushCmd
}
