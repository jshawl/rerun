package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialModel(t *testing.T) {
	t.Parallel()

	m := initialModel()

	if m.content != "" {
		t.Fatalf("Expected initial model content to be empty")
	}
}

func TestInit(t *testing.T) {
	t.Parallel()

	m := initialModel()
	cmd := m.Init()
	msg, ok := cmd().(startMsg)

	if !ok {
		t.Fatalf("expected msg to be a startMsg")
	}

	if msg.id != 0 {
		t.Fatalf("Expected init to send startMsg with id 0")
	}
}

func TestUpdateKeyMsgQuit(t *testing.T) {
	var cmd tea.Cmd

	t.Parallel()

	model := initialModel()
	keyMsg := tea.KeyMsg{
		Alt:   false,
		Paste: false,
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	}
	_, cmd = model.Update(keyMsg)
	_, ok := cmd().(tea.QuitMsg)

	if !ok {
		t.Fatalf("expected msg to be a tea.QuitMsg")
	}
}

func TestUpdateWindowSizeMsg(t *testing.T) {
	t.Parallel()

	m := initialModel()
	nextModel, _ := m.Update(tea.WindowSizeMsg{Height: 100, Width: 100})
	model2, correctType := nextModel.(model)

	if !correctType {
		t.Fatalf("Expected model type assertion")
	}

	if model2.steps.viewportWidth != 100 {
		t.Fatalf("Expected viewport width to be stored on steps model")
	}

	m3, _ := model2.Update(tea.WindowSizeMsg{Height: 100, Width: 99})
	model3, ok := m3.(model)

	if !ok {
		t.Fatalf("Expected model type assertion")
	}

	if model3.steps.viewportWidth != 99 {
		t.Fatalf("Expected resize to update steps model")
	}
}

func TestMainView(t *testing.T) {
	t.Parallel()

	m := initialModel()

	m.View()
}
