package main

import (
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/cli"
)

func TestRunHelpExitOK(t *testing.T) {
	if code := run([]string{"dotclean", "-h"}); code != cli.ExitOK {
		t.Fatalf("want ExitOK, got %d", code)
	}
}

func TestRunNoArgsExitFailure(t *testing.T) {
	if code := run([]string{"dotclean"}); code != cli.ExitFailure {
		t.Fatalf("want ExitFailure, got %d", code)
	}
}
