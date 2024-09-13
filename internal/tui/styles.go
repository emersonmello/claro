// Package tui provides the text user interface for the claro CLI tool.
package tui

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorLight             = "#43BF6D"
	colorDark              = "#73F59F"
	CurrentRepositoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	DoneStyle              = lipgloss.NewStyle().Margin(1, 2)
	TextStyle              = lipgloss.NewStyle().Margin(1, 2)
	ErrorStyle             = lipgloss.NewStyle().Margin(1, 2).Foreground(lipgloss.AdaptiveColor{Light: "#FF2D27", Dark: "#FF644E"})
	CheckMark              = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	BowtieMark             = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8BA00")).SetString("⧖")
	ErrorMark              = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF2D27")).SetString("𐄂")
)

// ClaroStyles represents the styles used in the Claro TUI
type ClaroStyles struct {
	Title        lipgloss.Style
	Item         lipgloss.Style
	SelectedItem lipgloss.Style
	Pagination   lipgloss.Style
	Help         lipgloss.Style
	QuitText     lipgloss.Style
	ErrorText    lipgloss.Style
	NormalText   lipgloss.Style
}

// CreateDefaultStyles creates the default styles for the Claro TUI
func CreateDefaultStyles() (s ClaroStyles) {
	s.Title = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#1B42A3"))
	s.Item = lipgloss.NewStyle().PaddingLeft(4)
	s.SelectedItem = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: colorLight, Dark: colorDark}).PaddingLeft(2)
	s.Pagination = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	s.Help = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.QuitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	s.NormalText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	s.ErrorText = lipgloss.NewStyle().Margin(1, 0, 2, 4).Foreground(lipgloss.AdaptiveColor{Light: "#FF2D27", Dark: "#FF644E"})
	return s
}
