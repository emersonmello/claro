// Package internal
package internal

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emersonmello/claro/internal/tui"
	"github.com/github/gh-classroom/pkg/classroom"
	"github.com/spf13/viper"
)

func gitCloneAssignment(assignment classroom.AcceptedAssignment) tea.Cmd {
	var directory string

	if strings.HasPrefix(directory, "~") {
		dirname, _ := os.UserHomeDir()
		directory = filepath.Join(dirname, directory[1:])
	}

	fullPath, _ := filepath.Abs(filepath.Join(directory, assignment.Assignment.Slug+"-submissions"))

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return func() tea.Msg {
				return tui.ErrorMsg(fmt.Sprintf("Error creating directory: %s", fullPath))
			}
		}
	}
	clonePath := filepath.Join(fullPath, assignment.Repository.Name)
	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "-q", assignment.Repository.HtmlUrl, clonePath)
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return tui.ErrorMsg(fmt.Sprintf("Error '%s' encountered while cloning: %s", err, assignment.Repository.FullName))
			}
			// Getting the commit hash to be used in the grade file
			cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
			commit, _ := executeCommand(cmd, clonePath)
			// Getting commit date
			cmd = exec.Command("git", "show", "-s", "--format=%ci")
			commitDate, _ := executeCommand(cmd, clonePath)
			commitStr := fmt.Sprintf("> Commit: %s | %s", strings.ReplaceAll(string(commit), "\n", ""), strings.ReplaceAll(string(commitDate), "\n", ""))
			// Creating grade file .md
			gradeFileName := filepath.Join(fullPath, "grade-"+assignment.Repository.Name+".md")
			if _, err = os.Stat(gradeFileName); os.IsNotExist(err) {
				if f, e := os.Create(gradeFileName); e != nil {
					return tui.ErrorMsg(fmt.Sprintf("Unable to create grade file: %s", e))
				} else {
					mdText := fmt.Sprintf("# %s\n%s\n\n- ...\n- **%s** \n\n", viper.GetString("title"), commitStr, viper.GetString("grade"))
					if _, e = f.WriteString(mdText); e != nil {
						return tui.ErrorMsg(fmt.Sprintf("Unable to write to markdown file: %s", e))
					}
					defer func(f *os.File) {
						_ = f.Close()
					}(f)
				}
			}
			return tui.SuccessfullMsg(assignment.Repository.Name)
		})
	}
	return func() tea.Msg {
		return tui.ErrorMsg("Repository already exists, skipping clone")
	}
}

func gitPullCmd(directory string, repositoryName string) tea.Cmd {
	//pause := time.Duration(rand.Int63n(1000)+3000) * time.Millisecond
	//time.Sleep(pause)

	cmd := exec.Command("git", "reset", "-q")
	if _, e := executeCommand(cmd, directory); e != nil {
		return func() tea.Msg { return tui.ErrorMsg("This is not a git repository") }
	}
	cmd = exec.Command("git", "stash", "-a", "-q")
	_, _ = executeCommand(cmd, directory)
	cmd = exec.Command("git", "rev-list", "--all", "--count")
	totalCommitsBeforePull, _ := executeCommand(cmd, directory)
	cmd = exec.Command("git", "pull", "-q", "--rebase")
	cmd.Dir = directory
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		cmd = exec.Command("git", "rev-list", "--all", "--count")
		totalCommitsAfterPull, _ := executeCommand(cmd, directory)
		cmd = exec.Command("git", "stash", "pop")
		_, _ = executeCommand(cmd, directory)
		if err != nil {
			return tui.ErrorMsg("Failed to execute 'git pull'")
		}
		if !bytes.Equal(totalCommitsAfterPull, totalCommitsBeforePull) {
			str := lipgloss.NewStyle().Foreground(lipgloss.Color("#E9E64D")).Italic(true).SetString("new commits").Render()
			return tui.SuccessfullPullMsg(fmt.Sprintf("%s %s", repositoryName, str))
		}
		return tui.SuccessfullPullMsg(repositoryName)
	})
}

func gitCommitAndPushCmd(directory string, repositoryName string, submission pair) tea.Cmd {
	gradeFileName := viper.GetString("filename")
	parentDir := filepath.Dir(directory)
	srcName, _ := filepath.Abs(filepath.Join(parentDir, submission.gradeFilename.Name()))
	dstName, _ := filepath.Abs(filepath.Join(directory, gradeFileName))
	cmd := exec.Command("git", "reset", "-q")
	_, _ = executeCommand(cmd, directory)
	cmd = exec.Command("git", "stash", "-a", "-q")
	_, _ = executeCommand(cmd, directory)
	if errCopy := copyFile(srcName, dstName); errCopy == nil {
		cmd := exec.Command("git", "add", gradeFileName)
		_, _ = executeCommand(cmd, directory)
		cmd = exec.Command("git", "status", "--porcelain")
		o, _ := executeCommand(cmd, directory)
		cmd = exec.Command("git", "stash", "pop")
		_, _ = executeCommand(cmd, directory)
		var str string
		if string(o) == "" {
			str = lipgloss.NewStyle().Foreground(lipgloss.Color("#E9E64D")).Italic(true).SetString("nothing to commit, working tree clean").Render()
		} else {
			cmd = exec.Command("git", "commit", "-q", "-m", viper.GetString("message"))
			_, _ = executeCommand(cmd, directory)
		}
		cmd = exec.Command("git", "push", "-q")
		cmd.Dir = directory
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return tui.ErrorMsg(fmt.Sprintf("Failed to execute 'git push' for %s", repositoryName))
			}
			return tui.SuccessfullMsg(fmt.Sprintf("%s %s", repositoryName, str))
		})
	}
	cmd = exec.Command("git", "stash", "pop")
	_, _ = executeCommand(cmd, directory)
	return func() tea.Msg {
		return tui.ErrorMsg(fmt.Sprintf("Error copying grade file: %s", gradeFileName))
	}
}

func checkIfDirectoryIsAGitRepo(directory os.DirEntry, sourceDirectory string) bool {
	baseDir, _ := os.Getwd()
	fullpath, _ := filepath.Abs(filepath.Join(sourceDirectory, directory.Name()))
	if err := os.Chdir(fullpath); err != nil {
		return false
	}
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	if out, err := cmd.Output(); err == nil {
		_ = os.Chdir(baseDir)
		if strings.Compare(string(out), fullpath+"\n") == 0 {
			return true
		}
	}
	_ = os.Chdir(baseDir)
	return false
}
