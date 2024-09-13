// Package internal
package internal

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emersonmello/claro/internal/tui"
	"github.com/github/gh-classroom/pkg/classroom"
)

type state int

const (
	initial state = iota
	listClassrooms
	fetchAssignmentsList
	listAssignments
	fetchRepositoriesList
	cloningAssignment
)

// CloneModel represents the model for the clone command
type CloneModel struct {
	state           state
	cL              tui.ClassroomList
	aL              tui.AssignmentsList
	repoL           []classroom.AcceptedAssignment
	totalCloned     int
	index           int
	classroomList   list.Model
	assignmentsList list.Model
	spinner         spinner.Model
	progress        progress.Model
	styles          tui.ClaroStyles
	keyMap          *tui.KeyMap
	help            help.Model
	done            bool
	width           int
	height          int
	credentialSet   bool
}

// NewCloneModel creates a new CloneModel
func NewCloneModel() CloneModel {
	styles := tui.CreateDefaultStyles()
	keys := tui.ClaroKeyMap()
	h := help.New()
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
		progress.WithoutPercentage(),
	)
	return CloneModel{
		state:         initial,
		styles:        styles,
		keyMap:        keys,
		help:          h,
		spinner:       sp,
		progress:      p,
		index:         0,
		totalCloned:   0,
		credentialSet: false,
	}
}

func (m CloneModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, restGetClassrooms(0))
}

func (m CloneModel) View() string {
	switch m.state {
	case initial:
		return m.classroomView()
	case listClassrooms:
		return m.classroomList.View()
	case fetchAssignmentsList:
		return m.assignmentsView()
	case listAssignments:
		return m.assignmentsList.View()
	case fetchRepositoriesList:
		return m.acceptedAssignmentsView()
	case cloningAssignment:
		return m.cloneView()
	default:
		return ""
	}
}

func (m CloneModel) classroomView() string {
	str := fmt.Sprintf("%s fetching your classrooms", m.spinner.View())
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.NewStyle().MarginLeft(1).Render(str))
}

func (m CloneModel) assignmentsView() string {
	str := fmt.Sprintf("%s retrieving your assignments", m.spinner.View())
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.NewStyle().MarginLeft(1).Render(str))
}

func (m CloneModel) acceptedAssignmentsView() string {
	str := fmt.Sprintf("%s Retrieving accepted assignments.", m.spinner.View())
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.NewStyle().MarginLeft(1).Render(str))
}

func (m CloneModel) cloneView() string {
	n := len(m.repoL)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return tui.DoneStyle.Render(fmt.Sprintf("Cloned %d repositories\n", m.totalCloned))
	}
	count := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)
	per := float64(m.index) / float64(len(m.repoL)-1)
	prog := m.progress.ViewAs(per)
	cellsAvail := max(0, m.width-lipgloss.Width(prog+count))

	repository := tui.CurrentRepositoryStyle.Render(m.repoL[m.index].Repository.Name)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render(fmt.Sprintf("%s %s ", tui.BowtieMark, repository))
	newLine := lipgloss.NewStyle().Render("\n")
	//cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+count)-2)
	//gap := strings.Repeat(" ", cellsRemaining)
	//return info + gap + prog + count + newLine
	return info + newLine + "  " + prog + count + newLine
}

func (m CloneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case initial:
		return listUpdate(msg, m)
	case listClassrooms:
		return classroomUpdate(msg, m)
	case fetchAssignmentsList:
		return listUpdate(msg, m)
	case listAssignments:
		return assignmentUpdate(msg, m)
	case fetchRepositoriesList:
		return listUpdate(msg, m)
	case cloningAssignment:
		return cloneUpdate(msg, m)
	default:
		return m, nil
	}
}

func classroomUpdate(msg tea.Msg, m CloneModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.classroomList.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "right":
			var i, ok = m.classroomList.SelectedItem().(tui.Item)
			if ok {
				m.state = fetchAssignmentsList
				return m, tea.Batch(m.spinner.Tick, restGetAssignments(i.Id, 0, 0))
			}
		}
	}
	var cmd tea.Cmd
	m.classroomList, cmd = m.classroomList.Update(msg)
	return m, cmd
}

func assignmentUpdate(msg tea.Msg, m CloneModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.assignmentsList.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "right":
			var i, ok = m.assignmentsList.SelectedItem().(tui.Item)
			if ok {
				m.state = fetchRepositoriesList
				return m, tea.Batch(m.spinner.Tick, restGetAcceptedAssignmentsList(i.Id, 0, 0))
			}
		case "left":
			m.state = listClassrooms
			var cmd tea.Cmd
			m.classroomList, cmd = m.classroomList.Update(msg)
			return m, cmd
		}
	}
	var cmd tea.Cmd
	m.assignmentsList, cmd = m.assignmentsList.Update(msg)
	return m, cmd
}

func listUpdate(msg tea.Msg, m CloneModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ClassroomList:
		m.cL = msg
		if m.cL != nil {
			if len(m.cL) > 0 {
				m.state = listClassrooms
				styles := tui.CreateDefaultStyles()
				keys := tui.ClaroKeyMap()
				height := min(len(msg)+8, m.height) - 2
				l := list.New(tui.MakeClassroomList(m.cL), tui.NewItemDelegate(&styles, keys), tui.DefaultWidth, height)
				l = tui.FormatList(l, "Select a classroom")
				l.AdditionalShortHelpKeys = m.keyMap.ShortHelp
				m.classroomList = l
				return m, nil
			}
		}
		return m, tea.Sequence(tea.Printf(m.styles.QuitText.Render("You don't have GitHub Classrooms, or you do not have permission to access them.")), tea.Quit)
	case tui.AssignmentsList:
		m.aL = msg
		if m.aL != nil {
			if len(m.aL) > 0 {
				m.state = listAssignments
				styles := tui.CreateDefaultStyles()
				keys := tui.ClaroKeyMap()
				height := min(len(msg)+8, m.height) - 2
				l := list.New(tui.MakeAssignmentsList(m.aL), tui.NewItemDelegate(&styles, keys), tui.DefaultWidth, height)
				l = tui.FormatList(l, "Select an assignment")
				l.AdditionalShortHelpKeys = m.keyMap.ShortHelp
				m.assignmentsList = l
				return m, nil
			}
		}
		return m, tea.Sequence(tea.Printf(m.styles.QuitText.Render("No assignments were found for this classroom, or you do not have permission to access them.")), tea.Quit)
	case []classroom.AcceptedAssignment:
		m.repoL = msg
		if len(m.repoL) > 0 {
			m.state = cloningAssignment
			m.index = 0
			return m, tea.Sequence(tea.Printf("Found %d repositories. Cloning...\n", len(m.repoL)), gitCloneAssignment(m.repoL[m.index]), m.spinner.Tick)
		}
		return m, tea.Sequence(tea.Printf(m.styles.QuitText.Render("No student submissions were found for this assignment, or you do not have permission to access them.")), tea.Quit)
	case tui.ErrorMsg:
		return m, tea.Sequence(tea.Printf(m.styles.ErrorText.Render(string(msg))), tea.Quit)
	}
	return m, nil
}

func cloneUpdate(msg tea.Msg, m CloneModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tui.SuccessfullMsg:
		m.totalCloned++
		cmd := tea.Printf("%s %s", tui.CheckMark, m.repoL[m.index].Repository.Name)
		if m.index >= len(m.repoL)-1 {
			m.done = true
			return m, tea.Sequence(cmd, tea.Quit)
		}
		m.index++
		return m, tea.Sequence(cmd, gitCloneAssignment(m.repoL[m.index]))
	case tui.ErrorMsg:
		reason := lipgloss.NewStyle().Foreground(lipgloss.Color("#783D38")).Italic(true).SetString(string(msg))
		cmd := tea.Printf("%s %s %s", tui.ErrorMark, m.repoL[m.index].Repository.Name, reason)
		if m.index >= len(m.repoL)-1 {
			m.done = true
			return m, tea.Sequence(cmd, tea.Quit)
		}
		m.index++
		return m, tea.Sequence(cmd, gitCloneAssignment(m.repoL[m.index]))
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}
