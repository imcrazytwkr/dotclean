package action_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/action"
	"github.com/imcrazytwkr/dotclean/internal/cli"
)

func TestExecuteKeepsSidecarAfterMergeFailure(t *testing.T) {
	dir := t.TempDir()
	sidecar := filepath.Join(dir, "._photo.jpg")
	native := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(sidecar, []byte("ad"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(native, []byte("img"), 0o644); err != nil {
		t.Fatal(err)
	}

	actions := []action.Action{
		{Kind: action.Merge, Path: sidecar, Native: native},
		{Kind: action.Delete, Path: sidecar},
	}
	err := action.Execute(actions, &cli.Options{}, func(string, string, *cli.Options) error {
		return errors.New("merge boom")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if _, err := os.Stat(sidecar); err != nil {
		t.Fatal("sidecar should remain after merge failure")
	}
}

func TestExecuteDeletesAfterSuccessfulMerge(t *testing.T) {
	dir := t.TempDir()
	sidecar := filepath.Join(dir, "._photo.jpg")
	native := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(sidecar, []byte("ad"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(native, []byte("img"), 0o644); err != nil {
		t.Fatal(err)
	}

	actions := []action.Action{
		{Kind: action.Merge, Path: sidecar, Native: native},
		{Kind: action.Delete, Path: sidecar},
	}
	if err := action.Execute(actions, &cli.Options{}, func(string, string, *cli.Options) error {
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatal("sidecar should be deleted after successful merge")
	}
}
