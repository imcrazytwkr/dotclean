package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/cli"
)

func TestAlwaysDeleteOverridesPreserve(t *testing.T) {
	opts := &cli.Options{
		AlwaysDelete: true,
		Preserve:     false,
	}
	if !opts.ShouldDeletePaired() {
		t.Fatal("should delete paired with AlwaysDelete")
	}
}

func TestPreserveKeepsPaired(t *testing.T) {
	opts := &cli.Options{Preserve: true}
	if opts.ShouldDeletePaired() {
		t.Fatal("preserve should keep paired")
	}
}

func TestOrphanNeedsFlag(t *testing.T) {
	opts := &cli.Options{}
	if opts.ShouldDeleteOrphan() {
		t.Fatal("orphan should not delete by default")
	}
	opts.CleanupOrphans = true
	if !opts.ShouldDeleteOrphan() {
		t.Fatal("CleanupOrphans should delete orphans")
	}
}

func TestParseOptsNoArgs(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean"})
	if code != cli.ExitFailure {
		t.Fatalf("want ExitFailure, got %d", code)
	}
	if opts != nil {
		t.Fatal("opts should be nil")
	}
}

func TestParseOptsHelp(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-h"})
	if code != cli.ExitOK {
		t.Fatalf("want ExitOK, got %d", code)
	}
	if opts != nil {
		t.Fatal("opts should be nil")
	}
}

func TestParseOptsRejectsSingleDashLong(t *testing.T) {
	_, code := cli.ParseOpts([]string{"dotclean", "-flat", "/tmp"})
	if code != cli.ExitUsage {
		t.Fatalf("want ExitUsage, got %d", code)
	}
}

func TestParseOptsShortAndLong(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-f", "--dry-run", "dir a"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if !opts.Flat || !opts.DryRun {
		t.Fatalf("flags: %+v", opts)
	}
	if len(opts.Dirs) != 1 || opts.Dirs[0] != "dir a" {
		t.Fatalf("dirs: %#v", opts.Dirs)
	}
}

func TestParseOptsClusteredShorts(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-mvN", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if !opts.AlwaysDelete || !opts.Verbose || !opts.DryRun {
		t.Fatalf("clustered: %+v", opts)
	}
	if opts.Preserve {
		t.Fatal("-m should clear preserve")
	}
}

func TestParseOptsAlwaysDeleteOverridesPreserve(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-m", "-p", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if opts.Preserve {
		t.Fatal("preserve should be cleared when -m set")
	}
	if !opts.ShouldDeletePaired() {
		t.Fatal("should delete paired with -m")
	}
}

func TestParseOptsPreserve(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-p", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if opts.ShouldDeletePaired() {
		t.Fatal("preserve should keep paired")
	}
}

func TestParseOptsInvalidKeep(t *testing.T) {
	_, code := cli.ParseOpts([]string{"dotclean", "--keep=nope", "/tmp"})
	if code != cli.ExitUsage {
		t.Fatalf("want ExitUsage, got %d", code)
	}
}

func TestParseOptsSpacedDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "dir a")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	opts, code := cli.ParseOpts([]string{"dotclean", "-N", dir})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if len(opts.Dirs) != 1 || opts.Dirs[0] != dir {
		t.Fatalf("dirs: %#v", opts.Dirs)
	}
}
