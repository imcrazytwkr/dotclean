package clean_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/clean"
	"github.com/imcrazytwkr/dotclean/internal/cli"
)

func TestDryRunNoMutation(t *testing.T) {
	dir := t.TempDir()
	native := filepath.Join(dir, "a b.jpg")
	sidecar := filepath.Join(dir, "._a b.jpg")
	if err := os.WriteFile(native, []byte("img"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sidecar, []byte("AD"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &cli.Options{DryRun: true, Dirs: []string{dir}, Keep: cli.KeepMostRecent}
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := clean.Run(opts)
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatal(err)
	}
	out, _ := io.ReadAll(r)
	if !bytes.Contains(out, []byte(sidecar)) {
		t.Fatalf("dry-run output missing sidecar: %s", out)
	}
	if _, err := os.Stat(sidecar); err != nil {
		t.Fatal("sidecar should still exist after dry-run")
	}
}

func TestDeletePairedDefault(t *testing.T) {
	dir := t.TempDir()
	native := filepath.Join(dir, "photo.jpg")
	sidecar := filepath.Join(dir, "._photo.jpg")
	os.WriteFile(native, []byte("x"), 0o644)
	os.WriteFile(sidecar, []byte("y"), 0o644)

	opts := &cli.Options{Dirs: []string{dir}, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatal("sidecar should be deleted")
	}
	if _, err := os.Stat(native); err != nil {
		t.Fatal("native should remain")
	}
}

func TestPreserveKeepsSidecar(t *testing.T) {
	dir := t.TempDir()
	native := filepath.Join(dir, "photo.jpg")
	sidecar := filepath.Join(dir, "._photo.jpg")
	os.WriteFile(native, []byte("x"), 0o644)
	os.WriteFile(sidecar, []byte("y"), 0o644)

	opts := &cli.Options{Dirs: []string{dir}, Preserve: true, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecar); err != nil {
		t.Fatal("sidecar should be preserved")
	}
}

func TestOrphanNeedsFlag(t *testing.T) {
	dir := t.TempDir()
	sidecar := filepath.Join(dir, "._lonely.jpg")
	os.WriteFile(sidecar, []byte("y"), 0o644)

	opts := &cli.Options{Dirs: []string{dir}, Keep: cli.KeepMostRecent}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecar); err != nil {
		t.Fatal("orphan should remain without -n/-m")
	}

	opts.CleanupOrphans = true
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatal("orphan should be removed with -n")
	}
}

func TestFlatDoesNotRecurse(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0o755)
	os.WriteFile(filepath.Join(dir, "root.jpg"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(dir, "._root.jpg"), []byte("b"), 0o644)
	os.WriteFile(filepath.Join(sub, "nested.jpg"), []byte("c"), 0o644)
	os.WriteFile(filepath.Join(sub, "._nested.jpg"), []byte("d"), 0o644)

	opts := &cli.Options{Dirs: []string{dir}, Flat: true, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "._root.jpg")); !os.IsNotExist(err) {
		t.Fatal("root sidecar should be gone")
	}
	if _, err := os.Stat(filepath.Join(sub, "._nested.jpg")); err != nil {
		t.Fatal("nested sidecar should remain with -f")
	}
}

func TestAlwaysDeleteOrphan(t *testing.T) {
	dir := t.TempDir()
	sidecar := filepath.Join(dir, "._lonely.jpg")
	if err := os.WriteFile(sidecar, []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &cli.Options{Dirs: []string{dir}, AlwaysDelete: true, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatal("orphan should be removed with -m")
	}
}

func TestAlwaysDeleteOverridesPreserve(t *testing.T) {
	dir := t.TempDir()
	native := filepath.Join(dir, "photo.jpg")
	sidecar := filepath.Join(dir, "._photo.jpg")
	if err := os.WriteFile(native, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sidecar, []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &cli.Options{
		Dirs:         []string{dir},
		AlwaysDelete: true,
		Preserve:     true,
		Keep:         cli.KeepNative,
	}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatal("sidecar should be deleted when -m overrides -p")
	}
	if _, err := os.Stat(native); err != nil {
		t.Fatal("native should remain")
	}
}

func TestDeepRemovesJunkLeavesAppleDoubleRules(t *testing.T) {
	dir := t.TempDir()
	ds := filepath.Join(dir, ".DS_Store")
	trashes := filepath.Join(dir, ".Trashes")
	native := filepath.Join(dir, "photo.jpg")
	sidecar := filepath.Join(dir, "._photo.jpg")
	if err := os.WriteFile(ds, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(trashes, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(trashes, "inner"), []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(native, []byte("n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sidecar, []byte("s"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &cli.Options{Dirs: []string{dir}, Deep: true, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ds); !os.IsNotExist(err) {
		t.Fatal(".DS_Store should be removed with -D")
	}
	if _, err := os.Stat(trashes); !os.IsNotExist(err) {
		t.Fatal(".Trashes should be removed with -D")
	}
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatal("paired sidecar should still be deleted by default")
	}
	if _, err := os.Stat(native); err != nil {
		t.Fatal("native should remain")
	}
}

func TestWithoutDeepJunkRemains(t *testing.T) {
	dir := t.TempDir()
	ds := filepath.Join(dir, ".DS_Store")
	if err := os.WriteFile(ds, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	opts := &cli.Options{Dirs: []string{dir}, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ds); err != nil {
		t.Fatal(".DS_Store should remain without -D")
	}
}

func TestPreserveDisablesDeep(t *testing.T) {
	dir := t.TempDir()
	ds := filepath.Join(dir, ".DS_Store")
	native := filepath.Join(dir, "photo.jpg")
	sidecar := filepath.Join(dir, "._photo.jpg")
	if err := os.WriteFile(ds, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(native, []byte("n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sidecar, []byte("s"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &cli.Options{Dirs: []string{dir}, Deep: true, Preserve: true, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ds); err != nil {
		t.Fatal(".DS_Store should remain when -p disables deep")
	}
	if _, err := os.Stat(sidecar); err != nil {
		t.Fatal("sidecar should be preserved with -p")
	}
}

func TestPreserveDisablesSpotlight(t *testing.T) {
	dir := t.TempDir()
	spot := filepath.Join(dir, ".Spotlight-V100")
	if err := os.Mkdir(spot, 0o755); err != nil {
		t.Fatal(err)
	}
	opts := &cli.Options{Dirs: []string{dir}, Spotlight: true, Preserve: true, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(spot); err != nil {
		t.Fatal(".Spotlight-V100 should remain when -p disables spotlight")
	}
}

func TestSpotlightRemovesSpotlightDir(t *testing.T) {
	dir := t.TempDir()
	spot := filepath.Join(dir, ".Spotlight-V100")
	if err := os.Mkdir(spot, 0o755); err != nil {
		t.Fatal(err)
	}
	opts := &cli.Options{Dirs: []string{dir}, Spotlight: true, Keep: cli.KeepNative}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(spot); !os.IsNotExist(err) {
		t.Fatal(".Spotlight-V100 should be removed with -S")
	}
}

func TestDeepDryRunListsJunk(t *testing.T) {
	dir := t.TempDir()
	ds := filepath.Join(dir, ".DS_Store")
	if err := os.WriteFile(ds, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	opts := &cli.Options{Dirs: []string{dir}, Deep: true, DryRun: true, Keep: cli.KeepNative}
	err := clean.Run(opts)
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatal(err)
	}
	out, _ := io.ReadAll(r)
	if !bytes.Contains(out, []byte(ds)) {
		t.Fatalf("dry-run missing .DS_Store: %s", out)
	}
	if _, err := os.Stat(ds); err != nil {
		t.Fatal(".DS_Store should remain after dry-run")
	}
}
