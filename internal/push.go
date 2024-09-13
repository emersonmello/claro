// Package internal
package internal

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emersonmello/claro/internal/tui"
)

type statePush int

const (
	initialPush = iota
	pushDir
)

type pair struct {
	repository    os.DirEntry
	gradeFilename os.DirEntry
}

type repo struct {
	repoMap      map[string]pair
	repositories []os.DirEntry
}

type PushModel struct {
	state                statePush
	submissionsDirectory string
	repos                repo
	totalPushed          int
	index                int
	styles               tui.ClaroStyles
	keyMap               *tui.KeyMap
	help                 help.Model
	progress             progress.Model
	done                 bool
	width                int
	height               int
}

func NewPushModel(directory string) PushModel {
	styles := tui.CreateDefaultStyles()
	keys := tui.ClaroKeyMap()
	h := help.New()
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
		progress.WithoutPercentage(),
	)
	return PushModel{
		state:                initialPush,
		submissionsDirectory: directory,
		totalPushed:          0,
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

func (m PushModel) Init() tea.Cmd {
	return getRepositoriesAndGradeFiles(m.submissionsDirectory)
}

func (m PushModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case initialPush:
		return initialPushUpdate(msg, m)
	case pushDir:
		return pushUpdate(msg, m)
	}
	return m, nil
}

func initialPushUpdate(msg tea.Msg, m PushModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case repo:
		m.repos = msg
		if len(m.repos.repoMap) > 0 {
			m.state = pushDir
			fullpath, _ := filepath.Abs(filepath.Join(m.submissionsDirectory, m.repos.repoMap[m.repos.repositories[m.index].Name()].repository.Name()))
			return m, tea.Sequence(tea.Printf("Grading submissions\n"), gitPullCmd(fullpath, m.repos.repoMap[m.repos.repositories[m.index].Name()].repository.Name()))
		} else {
			return m, tea.Sequence(tea.Printf(tui.ErrorStyle.Render(fmt.Sprintf("No repositories or grade files found in %s\nPlease see the help for more information\n", m.submissionsDirectory))), tea.Quit)
		}
	case tui.AssignmentDirError:
		cmd := tea.Printf("%s", msg)
		return m, tea.Sequence(cmd, tea.Quit)

	}
	return m, nil
}

func pushUpdate(msg tea.Msg, m PushModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle successful push message
	case tui.SuccessfullMsg:
		m.totalPushed++
		cmd := tea.Printf("%s %s ", tui.CheckMark, msg)
		if m.index >= len(m.repos.repositories)-1 {
			// If all repositories have been processed, mark as done and quit
			m.done = true
			return m, tea.Sequence(cmd, tea.Quit)
		}
		// Move to the next repository
		m.index++
		fullpath, _ := filepath.Abs(filepath.Join(m.submissionsDirectory, m.repos.repoMap[m.repos.repositories[m.index].Name()].repository.Name()))
		return m, tea.Sequence(cmd, gitPullCmd(fullpath, m.repos.repoMap[m.repos.repositories[m.index].Name()].repository.Name()))
	// Handle successful pull message
	case tui.SuccessfullPullMsg:
		fullpath, _ := filepath.Abs(filepath.Join(m.submissionsDirectory, m.repos.repoMap[m.repos.repositories[m.index].Name()].repository.Name()))
		return m, gitCommitAndPushCmd(fullpath, m.repos.repositories[m.index].Name(), m.repos.repoMap[m.repos.repositories[m.index].Name()])
	case tui.ErrorMsg:
		reason := lipgloss.NewStyle().Foreground(lipgloss.Color("#783D38")).Italic(true).Render(string(msg))
		cmd := tea.Printf("%s %s %s", tui.ErrorMark, m.repos.repositories[m.index].Name(), reason)
		if m.index >= len(m.repos.repositories)-1 {
			m.done = true
			return m, tea.Sequence(cmd, tea.Quit)
		}
		m.index++
		fullpath, _ := filepath.Abs(filepath.Join(m.submissionsDirectory, m.repos.repoMap[m.repos.repositories[m.index].Name()].repository.Name()))
		return m, tea.Sequence(cmd, gitPullCmd(fullpath, m.repos.repoMap[m.repos.repositories[m.index].Name()].repository.Name()))
	case tui.ErrorPullMsg:
		return m, tea.Sequence(tea.Printf(m.styles.ErrorText.Render(string(msg))), tea.Quit)
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m PushModel) View() string {
	if m.state == pushDir {
		return m.PushView()
	}
	return ""
}

func (m PushModel) PushView() string {
	if m.done {
		return tui.DoneStyle.Render(fmt.Sprintf("Graded %d submissions\n", m.totalPushed))
	}
	n := len(m.repos.repositories)
	w := lipgloss.Width(fmt.Sprintf("%d", n))
	count := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)
	per := float64(m.index) / float64(len(m.repos.repositories)-1)
	prog := m.progress.ViewAs(per)

	repository := tui.CurrentRepositoryStyle.Render(m.repos.repositories[m.index].Name())
	info := lipgloss.NewStyle().Render(fmt.Sprintf("%s %s ", tui.BowtieMark, repository))
	newLine := lipgloss.NewStyle().Render("\n")
	return info + newLine + "  " + prog + count + newLine
}

// getRepositoriesAndGradeFiles returns a tea.Cmd that lists all directories and grade files in the given source directory.
// It matches grade files with the corresponding repository directories based on a specific filename pattern.
//
// Parameters:
// - sourceDirectory: the path to the directory to list.
//
// Returns:
// - tea.Cmd: a command that, when executed, returns a tea.Msg containing a repo struct with matched repositories and grade files.
func getRepositoriesAndGradeFiles(sourceDirectory string) tea.Cmd {
	return func() tea.Msg {
		fileNamePattern := "^(grade-).*(\\.md)$"
		regexpPattern, _ := regexp.Compile(fileNamePattern)

		r := repo{
			repoMap:      make(map[string]pair),
			repositories: []os.DirEntry{},
		}

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

		regularEntries := make([]os.DirEntry, 0, len(entries))
		directories := make([]os.DirEntry, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				if !checkIfDirectoryIsAGitRepo(entry, sourceDirectory) {
					continue
				}
				directories = append(directories, entry)
			} else {
				if info, e := entry.Info(); e == nil {
					if info.Mode().IsRegular() {
						regularEntries = append(regularEntries, entry)
					}
				}
			}
		}

		for _, entry := range directories {
			for _, file := range regularEntries {
				if regexpPattern.MatchString(file.Name()) {
					fN := strings.Split(file.Name(), "grade-")
					if len(fN) == 2 {
						nN := strings.Split(fN[1], ".md")
						if len(nN) == 2 {
							key := nN[0]
							if key == entry.Name() {
								r.repoMap[key] = pair{repository: entry, gradeFilename: file}
								r.repositories = append(r.repositories, entry)
								break
							}
						}
					}
				}
			}
		}
		return r
	}
}
