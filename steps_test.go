package main

import (
	"math/rand"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitReturnsStart(t *testing.T) {
	s := Steps{steps: []Step{
		{command: "first command"},
	}}
	cmd := s.Init()
	start := cmd()
	if start.(startMsg).id != 0 {
		t.Fatalf("expected startMsg{id: 0}")
	}
}

func TestUpdateStartMsg(t *testing.T) {
	var cmd tea.Cmd
	s := Steps{steps: []Step{
		{command: "first command"},
		{command: "second command"},
	}}
	if s.currentStep != 0 {
		t.Fatalf("expected currentStep to be 0")
	}
	s, cmd = s.Update(startMsg{id: 1})
	_, _ = cmd().(tickMsg)
	if s.currentStep != 1 {
		t.Fatalf("expected currentStep to be 1")
	}
}

func TestUpdateTickMsgUpdatesCounter(t *testing.T) {
	s := Steps{steps: []Step{
		{command: "first command"},
	}}
	times := rand.Intn(1000000)
	for i := 0; i < times; i++ {
		s, _ = s.Update(tickMsg{})
	}
	if s.steps[0].counter != times {
		t.Fatalf("expected current step counter to increase %d", times)
	}
}
