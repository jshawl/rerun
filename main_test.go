package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialModel(t *testing.T) {
	m := initialModel()
	if m.content != "" {
		t.Fatalf("Expected initial model content to be empty")
	}
}

func TestInit(t *testing.T) {
	m := initialModel()
	cmd := m.Init()
	msg := cmd().(startMsg)
	if msg.id != 0 {
		t.Fatalf("Expected init to send startMsg with id 0")
	}
}

func TestUpdateKeyMsgQuit(t *testing.T) {
	var cmd tea.Cmd
	m := initialModel()
	k := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd = m.Update(k)
	_ = cmd().(tea.QuitMsg)
}

func TestUpdateWindowSizeMsg(t *testing.T) {
	m := initialModel()
	m2, _ := m.Update(tea.WindowSizeMsg{Height: 100, Width: 100})

	if m2.(model).steps.viewportWidth != 100 {
		t.Fatalf("Expected viewport width to be stored on steps model")
	}
	m3, _ := m2.Update(tea.WindowSizeMsg{Height: 100, Width: 99})
	if m3.(model).steps.viewportWidth != 99 {
		t.Fatalf("Expected resize to update steps model")
	}
}

func TestMainView(t *testing.T) {
	m := initialModel()
	m.View()
}
