// Package tui provides the text user interface for the claro CLI tool.
package tui

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left  key.Binding
	Right key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right}
}

// ClaroKeyMap returns the keybindings for the Claro TUI
func ClaroKeyMap() *KeyMap {
	return &KeyMap{
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "to go back"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "enter"),
			key.WithHelp("→/enter", "confirm"),
		),
	}
}
