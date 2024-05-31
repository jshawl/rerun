package main

import (
	"math/rand"
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

func TestUpdateTickMsgUpdatesCounter(t *testing.T) {
	s := Steps{steps: []Step{
		newStep("first command"),
	}}
	times := rand.Intn(1000000)
	for i := 0; i < times; i++ {
		s, _ = s.Update(tickMsg{})
	}
	if s.steps[0].counter != times {
		t.Fatalf("expected current step counter to increase %d", times)
	}
}
