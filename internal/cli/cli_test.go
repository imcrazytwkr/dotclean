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

func TestParseOptsSetQuarantine(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-Q", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if !opts.SetQuarantine {
		t.Fatal("expected SetQuarantine from -Q")
	}
	opts, code = cli.ParseOpts([]string{"dotclean", "--set-quarantine", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if !opts.SetQuarantine {
		t.Fatal("expected SetQuarantine from --set-quarantine")
	}
	opts, code = cli.ParseOpts([]string{"dotclean", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if opts.SetQuarantine {
		t.Fatal("SetQuarantine should default false")
	}
}

func TestParseOptsDeepSpotlight(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-D", "-S", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if !opts.Deep || !opts.Spotlight {
		t.Fatalf("deep/spotlight: %+v", opts)
	}
	opts, code = cli.ParseOpts([]string{"dotclean", "--deep", "--spotlight", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if !opts.Deep || !opts.Spotlight {
		t.Fatal("long flags")
	}
	opts, code = cli.ParseOpts([]string{"dotclean", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if opts.Deep || opts.Spotlight {
		t.Fatal("deep/spotlight should default off")
	}
}

func TestParseOptsPreserveDisablesDeep(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-D", "-p", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if opts.Deep {
		t.Fatal("-p should clear Deep")
	}
	if !opts.Preserve {
		t.Fatal("preserve should remain set")
	}
}

func TestParseOptsPreserveDisablesSpotlight(t *testing.T) {
	opts, code := cli.ParseOpts([]string{"dotclean", "-S", "-p", "/tmp"})
	if code != cli.ExitContinue {
		t.Fatalf("want ExitContinue, got %d", code)
	}
	if opts.Spotlight {
		t.Fatal("-p should clear Spotlight")
	}
	if !opts.Preserve {
		t.Fatal("preserve should remain set")
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
