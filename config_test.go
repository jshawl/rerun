package main

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	conf := parseConfig()
	if conf["steps"][0] != "go mod tidy" {
		t.Fatalf("expected first step to be go mod tidy")
	}
}
