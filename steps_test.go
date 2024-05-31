package main

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func setup(cmds ...Step) Steps {
	return Steps{
		currentStep:   0,
		steps:         cmds,
		viewportWidth: 0,
	}
}

func TestInitReturnsStart(t *testing.T) {
	t.Parallel()

	s := setup(newStep("first command"))

	cmd := s.Init()
	start := cmd()

	msg, ok := start.(startMsg)
	if !ok {
		t.Fatalf("expect msg to be a startMsg")
	}

	if msg.id != 0 {
		t.Fatalf("expected startMsg{id: 0}")
	}
}

func TestUpdateStartMsg(t *testing.T) {
	var cmd tea.Cmd

	t.Parallel()

	steps := setup(
		newStep("first command"),
		newStep("second command"),
	)
	if steps.currentStep != 0 {
		t.Fatalf("expected currentStep to be 0")
	}

	if steps.steps[steps.currentStep].state != Pending {
		t.Fatalf("expected step to have state pending")
	}

	steps, cmd = steps.Update(startMsg{id: 0})
	if steps.steps[steps.currentStep].state != Started {
		t.Fatalf("expected step to have state started")
	}

	tick := cmd().(tea.BatchMsg)[0] //nolint:forcetypeassert
	_ = tick().(tickMsg)            //nolint:forcetypeassert
	exit := cmd().(tea.BatchMsg)[1] //nolint:forcetypeassert
	_ = exit().(exitMsg)            //nolint:forcetypeassert
}

func TestUpdateTickMsgReturnsTickOnlyIfStarted(t *testing.T) {
	var cmd tea.Cmd

	t.Parallel()

	steps := setup(
		newStep("first command"),
	)
	steps.steps[0].state = Started
	steps, cmd = steps.Update(tickMsg{})

	_, ok := cmd().(tickMsg)
	if !ok {
		t.Fatalf("expected cmd to be a tickMsg")
	}

	steps.steps[0].state = Exited0

	_, cmd = steps.Update(tickMsg{})
	if cmd != nil {
		t.Fatalf("Expected ticking to stop after Exited0")
	}
}

func TestUpdateExitMsg(t *testing.T) {
	var cmd tea.Cmd

	t.Parallel()

	steps := setup(
		newStep("first command"),
		newStep("second command"),
	)
	if steps.currentStep != 0 {
		t.Fatalf("expected currentStep to be 0")
	}

	steps, cmd = steps.Update(exitMsg{
		id:     0,
		output: "",
		err:    nil,
	})
	if steps.currentStep != 1 {
		t.Fatalf("expected currentStep to be 1")
	}

	tick := cmd().(tea.BatchMsg)[0] //nolint:forcetypeassert
	_ = tick().(tickMsg)            //nolint:forcetypeassert
	exit := cmd().(tea.BatchMsg)[1] //nolint:forcetypeassert
	// starts the next one
	_ = exit().(startMsg) //nolint:forcetypeassert
}

func TestUpdateExitMsgErr(t *testing.T) {
	t.Parallel()

	steps := setup(
		newStep("first command"),
		newStep("second command"),
	)
	steps, _ = steps.Update(exitMsg{
		id:     0,
		output: "",
		err:    errors.New("an error occurred"), //nolint:err113
	})

	if steps.steps[0].state != Exited1 {
		t.Fatalf("Expected exitMsg to update state to Exited 1")
	}

	if steps.steps[1].state != Skipped {
		t.Fatalf("Expected remaining steps to be skipped")
	}
}

func TestReset(t *testing.T) {
	t.Parallel()

	steps := setup(
		newStep("first command"),
	)
	steps.steps[0].state = Skipped
	steps.reset()

	if steps.steps[0].state != Pending {
		t.Fatalf("Expected reset to change state to Pending")
	}
}

func TestKeyMsgR(t *testing.T) {
	t.Parallel()

	steps := setup(
		newStep("first command"),
	)
	steps.steps[0].state = Exited1
	keyMsg := tea.KeyMsg{
		Alt:   false,
		Paste: false,
		Type:  tea.KeyRunes,
		Runes: []rune("r"),
	}

	steps, _ = steps.Update(keyMsg)
	if steps.steps[0].state != Pending {
		t.Fatalf("Expected 'r' keypress to reset")
	}
}

func TestView(t *testing.T) {
	t.Parallel()

	steps := setup(
		newStep("first command"),
		newStep("second command"),
	)
	steps.steps[0].state = Exited1
	steps.steps[1].state = Skipped
	steps.View()
	steps.viewportWidth = 100
	steps.View()
}
