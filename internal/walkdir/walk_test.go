package walkdir_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/classify"
	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/walkdir"
)

func TestCollectRecurses(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dir, "._root.jpg"), "a")
	write(t, filepath.Join(sub, "._nested.jpg"), "b")

	cands, err := walkdir.Collect(&cli.Options{Dirs: []string{dir}})
	if err != nil {
		t.Fatal(err)
	}
	paths := candPaths(cands)
	if !paths[filepath.Join(dir, "._root.jpg")] || !paths[filepath.Join(sub, "._nested.jpg")] {
		t.Fatalf("expected root and nested sidecars; got %v", paths)
	}
}

func TestCollectFlatSkipsNested(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dir, "._root.jpg"), "a")
	write(t, filepath.Join(sub, "._nested.jpg"), "b")

	cands, err := walkdir.Collect(&cli.Options{Dirs: []string{dir}, Flat: true})
	if err != nil {
		t.Fatal(err)
	}
	paths := candPaths(cands)
	if !paths[filepath.Join(dir, "._root.jpg")] {
		t.Fatal("missing root sidecar")
	}
	if paths[filepath.Join(sub, "._nested.jpg")] {
		t.Fatal("nested sidecar should be skipped with Flat")
	}
}

func TestCollectSpacedName(t *testing.T) {
	dir := t.TempDir()
	sidecar := filepath.Join(dir, "._a b.jpg")
	write(t, sidecar, "ad")

	cands, err := walkdir.Collect(&cli.Options{Dirs: []string{dir}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 {
		t.Fatalf("got %d candidates", len(cands))
	}
	if cands[0].Path != sidecar {
		t.Fatalf("path=%q", cands[0].Path)
	}
	if cands[0].Native != filepath.Join(dir, "a b.jpg") {
		t.Fatalf("native=%q", cands[0].Native)
	}
}

func TestCollectFollowSymlink(t *testing.T) {
	dir := t.TempDir()
	walkRoot := filepath.Join(dir, "walk")
	realDir := filepath.Join(dir, "real")
	if err := os.Mkdir(walkRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(realDir, 0o755); err != nil {
		t.Fatal(err)
	}
	realSidecar := filepath.Join(realDir, "._photo.jpg")
	write(t, realSidecar, "ad")

	link := filepath.Join(walkRoot, "._photo.jpg")
	if err := os.Symlink(realSidecar, link); err != nil {
		t.Fatal(err)
	}

	cands, err := walkdir.Collect(&cli.Options{Dirs: []string{walkRoot}, FollowSymlinks: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 {
		t.Fatalf("got %d candidates", len(cands))
	}
	resolved, err := filepath.EvalSymlinks(link)
	if err != nil {
		t.Fatal(err)
	}
	if cands[0].Path != resolved {
		t.Fatalf("FollowSymlinks path=%q want %q", cands[0].Path, resolved)
	}

	// realDir may sit under a symlinked prefix (e.g. /tmp -> /private/tmp on macOS).
	realResolved, err := filepath.EvalSymlinks(realDir)
	if err != nil {
		t.Fatal(err)
	}

	wantNative := filepath.Join(realResolved, "photo.jpg")
	if cands[0].Native != wantNative {
		t.Fatalf("native=%q want %q", cands[0].Native, wantNative)
	}
}

func TestCollectNoFollowKeepsSymlinkPath(t *testing.T) {
	dir := t.TempDir()
	walkRoot := filepath.Join(dir, "walk")
	realDir := filepath.Join(dir, "real")
	if err := os.Mkdir(walkRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(realDir, 0o755); err != nil {
		t.Fatal(err)
	}
	realSidecar := filepath.Join(realDir, "._photo.jpg")
	write(t, realSidecar, "ad")

	link := filepath.Join(walkRoot, "._photo.jpg")
	if err := os.Symlink(realSidecar, link); err != nil {
		t.Fatal(err)
	}

	cands, err := walkdir.Collect(&cli.Options{Dirs: []string{walkRoot}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 {
		t.Fatalf("got %d candidates", len(cands))
	}
	if cands[0].Path != link {
		t.Fatalf("no-follow path=%q want %q", cands[0].Path, link)
	}
	if cands[0].Native != filepath.Join(walkRoot, "photo.jpg") {
		t.Fatalf("native=%q", cands[0].Native)
	}
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func candPaths(cands []classify.Candidate) map[string]bool {
	out := make(map[string]bool, len(cands))
	for _, c := range cands {
		out[c.Path] = true
	}
	return out
}
