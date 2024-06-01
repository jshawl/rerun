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

func assertMsgType[T interface{}](t *testing.T, cmd tea.Cmd) T {
	t.Helper()

	msg := cmd()
	typed, ok := msg.(T)

	if !ok {
		t.Fatalf("Expected msg to be of type %T, got %T", *new(T), msg)
	}

	return typed
}

func TestInitReturnsStart(t *testing.T) {
	t.Parallel()

	s := setup(newStep("first command"))

	cmd := s.Init()
	msg := assertMsgType[startMsg](t, cmd)

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

	msgs := assertMsgType[tea.BatchMsg](t, cmd)
	assertMsgType[tickMsg](t, msgs[0])
	assertMsgType[exitMsg](t, msgs[1])
}

func TestUpdateTickMsgReturnsTickOnlyIfStarted(t *testing.T) {
	var cmd tea.Cmd

	t.Parallel()

	steps := setup(
		newStep("first command"),
	)
	steps.steps[0].state = Started
	steps, cmd = steps.Update(tickMsg{})
	assertMsgType[tickMsg](t, cmd)

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
		output: "",
		err:    nil,
	})
	if steps.currentStep != 1 {
		t.Fatalf("expected currentStep to be 1")
	}

	msgs := assertMsgType[tea.BatchMsg](t, cmd)
	assertMsgType[tickMsg](t, msgs[0])
	assertMsgType[startMsg](t, msgs[1])
}

func TestUpdateExitMsgErr(t *testing.T) {
	t.Parallel()

	steps := setup(
		newStep("first command"),
		newStep("second command"),
	)
	steps, _ = steps.Update(exitMsg{
		output: "",
		err:    errors.New("an error occurred"), //nolint:err113
	})

	if steps.steps[0].state != Exited1 {
		t.Fatalf("Expected exitMsg to update state to Exited 1")
	}

	if steps.steps[1].state != Skipped {
		t.Fatalf("Expected remaining steps to be skipped")
	}

	if steps.steps[0].output != "an error occurred" {
		t.Fatalf("Expected error message to be in output")
	}

	steps, _ = steps.Update(exitMsg{
		output: "display the output instead of the error",
		err:    errors.New("an error occurred"), //nolint:err113
	})

	if steps.steps[0].output != "display the output instead of the error" {
		t.Fatalf("Expected error output to be in output")
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
