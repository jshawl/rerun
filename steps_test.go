package main

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitReturnsStart(t *testing.T) {
	s := Steps{steps: []Step{
		newStep("first command"),
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
		newStep("first command"),
		newStep("second command"),
	}}
	if s.currentStep != 0 {
		t.Fatalf("expected currentStep to be 0")
	}
	if s.steps[s.currentStep].state != Pending {
		t.Fatalf("expected step to have state pending")
	}

	s, cmd = s.Update(startMsg{id: 0})
	if s.steps[s.currentStep].state != Started {
		t.Fatalf("expected step to have state started")
	}

	tick := cmd().(tea.BatchMsg)[0]
	_ = tick().(tickMsg)
	exit := cmd().(tea.BatchMsg)[1]
	_ = exit().(exitMsg)
}

func TestUpdateTickMsgReturnsTickOnlyIfStarted(t *testing.T) {
	var cmd tea.Cmd
	s := Steps{steps: []Step{
		newStep("first command"),
	}}
	s.steps[0].state = Started
	s, cmd = s.Update(tickMsg{})
	_ = cmd().(tickMsg)
	s.steps[0].state = Exited0
	s, cmd = s.Update(tickMsg{})
	if cmd != nil {
		t.Fatalf("Expected ticking to stop after Exited0")
	}
}

func TestUpdateExitMsg(t *testing.T) {
	var cmd tea.Cmd
	s := Steps{steps: []Step{
		newStep("first command"),
		newStep("second command"),
	}}
	if s.currentStep != 0 {
		t.Fatalf("expected currentStep to be 0")
	}
	s, cmd = s.Update(exitMsg{id: 0})
	if s.currentStep != 1 {
		t.Fatalf("expected currentStep to be 1")
	}
	tick := cmd().(tea.BatchMsg)[0]
	_ = tick().(tickMsg)
	exit := cmd().(tea.BatchMsg)[1]
	// starts the next one
	_ = exit().(startMsg)
}

func TestUpdateExitMsgErr(t *testing.T) {

	s := Steps{steps: []Step{
		newStep("first command"),
		newStep("second command"),
	}}
	s, _ = s.Update(exitMsg{id: 0, err: errors.New("yikes!")})

	if s.steps[0].state != Exited1 {
		t.Fatalf("Expected exitMsg to update state to Exited 1")
	}

	if s.steps[1].state != Skipped {
		t.Fatalf("Expected remaining steps to be skipped")
	}
}

func TestReset(t *testing.T) {
	s := Steps{steps: []Step{
		newStep("first command"),
	}}
	s.steps[0].state = Skipped
	s.reset()
	if s.steps[0].state != Pending {
		t.Fatalf("Expected reset to change state to Pending")
	}
}

func TestKeyMsgR(t *testing.T) {
	s := Steps{steps: []Step{
		newStep("first command"),
	}}
	s.steps[0].state = Exited1
	k := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
	s, _ = s.Update(k)
	if s.steps[0].state != Pending {
		t.Fatalf("Expected 'r' keypress to reset")
	}
}

func TestView(t *testing.T) {
	s := Steps{steps: []Step{
		newStep("first command"),
		newStep("second command"),
	}}
	s.steps[0].state = Exited1
	s.steps[1].state = Skipped
	s.View()
	s.viewportWidth = 100
	s.View()
}
