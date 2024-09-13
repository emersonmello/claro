// Package tui provides the text user interface for the claro CLI tool.
package tui

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"
	legal "io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Item represents an item in the list
type Item struct {
	Id   string
	Name string
	Url  string
}

// ItemDelegate represents the delegate for the list
type ItemDelegate struct {
	styles *ClaroStyles
	keys   *KeyMap
}

// NewItemDelegate creates a new ItemDelegate
func NewItemDelegate(styles *ClaroStyles, keys *KeyMap) *ItemDelegate {
	return &ItemDelegate{
		styles: styles,
		keys:   keys,
	}
}

func (i Item) FilterValue() string                             { return i.Name }
func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w legal.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%s", i.Name)

	fn := d.styles.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.SelectedItem.Render("> " + strings.Join(s, " "))
		}
	}
	_, err := fmt.Fprint(w, fn(str))
	if err != nil {
		return
	}
}
