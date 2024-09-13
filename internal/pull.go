// Package internal
package internal

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emersonmello/claro/internal/tui"
)

type statePull int

const (
	initialPull statePull = iota
	pullDir
)

type PullModel struct {
	state                statePull
	submissionsDirectory string
	repositories         []os.DirEntry
	totalPulled          int
	index                int
	styles               tui.ClaroStyles
	keyMap               *tui.KeyMap
	help                 help.Model
	progress             progress.Model
	done                 bool
	width                int
	height               int
}

func NewPullModel(directory string) PullModel {
	styles := tui.CreateDefaultStyles()
	keys := tui.ClaroKeyMap()
	h := help.New()
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
		progress.WithoutPercentage(),
	)
	return PullModel{
		state:                initialPull,
		submissionsDirectory: directory,
		totalPulled:          0,
		index:                0,
		styles:               styles,
		keyMap:               keys,
		help:                 h,
		progress:             p,
		done:                 false,
		width:                0,
		height:               0,
	}
}

func (m PullModel) Init() tea.Cmd {
	cmd := getReposDirectoryList(m.submissionsDirectory)
	return cmd
}

func (m PullModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	}
	switch m.state {
	case initialPull:
		return initialPullUpdate(msg, m)
	case pullDir:
		return pullUpdate(msg, m)
	}
	return m, nil
}

func initialPullUpdate(msg tea.Msg, m PullModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.AssignmentDirError:
		cmd := tea.Printf("%s", msg)
		return m, tea.Sequence(cmd, tea.Quit)
	case []os.DirEntry:
		m.repositories = msg
		if len(m.repositories) > 0 {
			m.state = pullDir
			m.index = 0
			fullpath, _ := filepath.Abs(filepath.Join(m.submissionsDirectory, m.repositories[m.index].Name()))
			return m, tea.Sequence(tea.Printf("Pulling %d repositories\n", len(m.repositories)), gitPullCmd(fullpath, m.repositories[m.index].Name()))
		} else {
			return m, tea.Sequence(tea.Printf(tui.ErrorStyle.Render(fmt.Sprintf("No repositories found in %s\n", m.submissionsDirectory))), tea.Quit)
		}

	}
	return m, nil
}

func pullUpdate(msg tea.Msg, m PullModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.SuccessfullMsg, tui.SuccessfullPullMsg:
		m.totalPulled++
		cmd := tea.Printf("%s %s", tui.CheckMark, msg)
		if m.index >= len(m.repositories)-1 {
			m.done = true
			return m, tea.Sequence(cmd, tea.Quit)
		}
		m.index++
		fullpath, _ := filepath.Abs(filepath.Join(m.submissionsDirectory, m.repositories[m.index].Name()))
		return m, tea.Sequence(cmd, gitPullCmd(fullpath, m.repositories[m.index].Name()))
	case tui.ErrorMsg:
		reason := lipgloss.NewStyle().Foreground(lipgloss.Color("#783D38")).Italic(true).SetString(string(msg)).Render()
		cmd := tea.Printf("%s %s %s", tui.ErrorMark, m.repositories[m.index].Name(), reason)
		if m.index >= len(m.repositories)-1 {
			m.done = true
			return m, tea.Sequence(cmd, tea.Quit)
		}
		m.index++
		fullpath, _ := filepath.Abs(filepath.Join(m.submissionsDirectory, m.repositories[m.index].Name()))
		return m, tea.Sequence(cmd, gitPullCmd(fullpath, m.repositories[m.index].Name()))
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

// View renders the current view of the PullModel based on its state.
//
// If the state is `pullDir` and repositories are available, it will display
// the current repository being processed. If all repositories have been pulled,
// it will display a message indicating the total number of repositories pulled.
//
// Returns a string representing the current view of the PullModel.
func (m PullModel) View() string {
	if m.state == pullDir {
		return m.PullView()
	}
	return ""
}

func (m PullModel) PullView() string {
	if m.done {
		return tui.DoneStyle.Render(fmt.Sprintf("Pulled %d repositories\n", m.totalPulled))
	}
	n := len(m.repositories)
	w := lipgloss.Width(fmt.Sprintf("%d", n))
	count := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)
	per := float64(m.index) / float64(len(m.repositories)-1)
	prog := m.progress.ViewAs(per)

	repository := tui.CurrentRepositoryStyle.Render(m.repositories[m.index].Name())
	info := lipgloss.NewStyle().Render(fmt.Sprintf("%s %s ", tui.BowtieMark, repository))
	newLine := lipgloss.NewStyle().Render("\n")
	return info + newLine + "  " + prog + count + newLine
}

// getReposDirectoryList returns a tea.Cmd that lists all directories in the given source directory.
// If the source directory starts with "~", it will be expanded to the user's home directory.
// If the source directory does not exist, it returns a message indicating the error.
//
// Parameters:
// - sourceDirectory: the path to the directory to list.
//
// Returns:
// - tea.Cmd: a command that, when executed, returns a tea.Msg containing a list of directories.
func getReposDirectoryList(sourceDirectory string) tea.Cmd {
	return func() tea.Msg {
		if strings.HasPrefix(sourceDirectory, "~") {
			dirname, _ := os.UserHomeDir()
			sourceDirectory = filepath.Join(dirname, sourceDirectory[1:])
		}
		if _, err := os.Stat(sourceDirectory); os.IsNotExist(err) {
			return tui.AssignmentDirError(tui.ErrorStyle.Render("The assignment directory does not exist. Please run the clone command first.\n"))
		}
		entries, err := os.ReadDir(sourceDirectory)
		if err != nil {
			return tui.AssignmentDirError(tui.ErrorStyle.Render(fmt.Sprintf("Failed to read the directory: %v\n", err)))
		}
		onlyDirs := entries[:0]
		for _, entry := range entries {
			if entry.IsDir() {
				onlyDirs = append(onlyDirs, entry)
			}
		}
		return onlyDirs
	}
}
