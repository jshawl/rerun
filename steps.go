package main

import tea "github.com/charmbracelet/bubbletea"

type Steps struct {
}

func (m Steps) Init() tea.Cmd {
	return nil
}

func (m Steps) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Steps) View() string {
	return "steps view"
}
