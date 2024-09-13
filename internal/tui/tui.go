// Package tui provides the text user interface for the claro CLI tool.
package tui

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/github/gh-classroom/pkg/classroom"
)

type ClassroomList []classroom.Classroom
type AssignmentsList []classroom.Assignment

type SuccessfullMsg string
type NewCommits string
type ErrorMsg string
type SuccessfullPullMsg string
type ErrorPullMsg string
type AssignmentDirError string

var GitHubCliInstalled bool
var UserGitHubPAT string

const (
	DefaultWidth = 80
	useMsg       = "\nThe directory should have been created using the 'clone' command and should include:" +
		"\n- Subdirectories, each named after a student's repository (e.g., assignment-01-JohnDoeStudent)." +
		"\n- Markdown files, each named with the pattern 'grade-<repository-name>.md' (e.g., grade-assignment-01-JohnDoeStudent.md)."
)

func UseErrorMsg(cmdName string) string {
	return ErrorStyle.Render(fmt.Sprintf("The '%s' command requires a directory containing student repositories and their corresponding grade files.", cmdName)) +
		useMsg +
		fmt.Sprintf("\n\nEnsure that this structure is followed for the '%s' command to work correctly.\n", cmdName)
}

func LongHelpMsg(shortHelpMsg string) string {
	return fmt.Sprintf("%s\n%s", shortHelpMsg, useMsg)
}

// MakeClassroomList creates a list of classrooms
func MakeClassroomList(listC ClassroomList) []list.Item {
	var items []list.Item
	for _, c := range listC {
		if c.Archived == false {
			items = append(items, Item{Id: fmt.Sprintf("%d", c.Id), Name: c.Name, Url: c.Url})
		}
	}
	return items
}

// MakeAssignmentsList creates a list of assignments
func MakeAssignmentsList(listA AssignmentsList) []list.Item {
	var items []list.Item
	for _, c := range listA {
		items = append(items, Item{Id: fmt.Sprintf("%d", c.Id), Name: c.Title})
	}
	return items
}

// FormatList formats a list using the claro default styles
func FormatList(l list.Model, title string) list.Model {
	styles := CreateDefaultStyles()
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.Title
	l.Styles.PaginationStyle = styles.Pagination
	l.Styles.HelpStyle = styles.Help
	return l
}
