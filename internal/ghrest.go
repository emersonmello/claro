// Package internal
package internal

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"errors"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/emersonmello/claro/internal/tui"
	"github.com/github/gh-classroom/pkg/classroom"
)

func restGetClassrooms(page int) tea.Cmd {
	return func() tea.Msg {
		client, er := getAPIRESTClient()
		if client == nil {
			return er
		}
		var classroomList tui.ClassroomList
		var path = "classrooms"
		if page != 0 {
			path = fmt.Sprintf("%s?page=%d", path, page)
		}
		if e := client.Get(path, &classroomList); e != nil {
			return checkIfBadCredentialError(e, "Failed to retrieve the classrooms list")
		}

		return classroomList
		//return generateRandomClassrooms(30)
	}
}
func restGetAssignments(classroomId string, page int, perPage int) tea.Cmd {
	return func() tea.Msg {
		client, er := getAPIRESTClient()
		if client == nil {
			return er
		}
		var assignments tui.AssignmentsList
		var path = fmt.Sprintf("classrooms/%s/assignments", classroomId)
		if page != 0 {
			path += fmt.Sprintf("?page=%v", page)
		}
		if perPage != 0 {
			path += fmt.Sprintf("&per_page=%v", perPage)
		}
		if e := client.Get(path, &assignments); e != nil {
			return checkIfBadCredentialError(e, "Failed to retrieve the assignments list")
		}
		return assignments
		//return generateRandomAssignments(30)
	}
}

func restGetAcceptedAssignmentsList(assignmentId string, page int, perPage int) tea.Cmd {
	return func() tea.Msg {
		client, er := getAPIRESTClient()
		if client == nil {
			return er
		}
		var repos = make([]classroom.AcceptedAssignment, 0)
		var path = fmt.Sprintf("assignments/%v/accepted_assignments", assignmentId)
		if page != 0 {
			path += fmt.Sprintf("?page=%v", page)
		}
		if perPage != 0 {
			path += fmt.Sprintf("&per_page=%v", perPage)
		}
		if e := client.Get(path, &repos); e != nil {
			return checkIfBadCredentialError(e, "Failed to retrieve the accepted assignments list")
		}
		return repos
	}
}

func getAPIRESTClient() (*api.RESTClient, tui.ErrorMsg) {
	var client *api.RESTClient
	var err error
	// If you have GitHub CLI installed, you can use it to connect to the GitHub REST API.
	//if tui.GitHubCliInstalled {
	//	client, err = api.DefaultRESTClient()
	//} else {
	// Ok, no problem. Since I'm not using GitHub CLI, I need to have access to a Personal Access Token
	opts := api.ClientOptions{AuthToken: tui.UserGitHubPAT}
	client, err = api.NewRESTClient(opts)
	//}
	var errorMsg tui.ErrorMsg
	if client == nil {
		errorMsg = tui.ErrorMsg(fmt.Sprintf("An error occurred while retrieving the GitHub REST API client. %s", err))
	}
	return client, errorMsg
}

func checkIfBadCredentialError(e error, msg string) tui.ErrorMsg {
	var hE *api.HTTPError
	s := "claro config"
	if errors.As(e, &hE) {
		if hE.StatusCode == http.StatusUnauthorized {
			return tui.ErrorMsg(fmt.Sprintf("%s => HTTP %d: %s."+
				"\nVisit https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens to generate a new PAT"+
				"\nThen, execute '%s' to update your GitHub Personal Access Token.", msg, hE.StatusCode, hE.Message, s))
		}
	}
	return tui.ErrorMsg(fmt.Sprintf("%s. %s", msg, e))
}
