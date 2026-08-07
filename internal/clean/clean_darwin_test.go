//go:build darwin

package clean_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/clean"
	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/fstype"
	"github.com/imcrazytwkr/dotclean/internal/testenv"
	"golang.org/x/sys/unix"
)

func TestRunMergesAndDeletesSidecar(t *testing.T) {
	dir := testenv.TempDir(t, "dotclean-clean-*")

	native := filepath.Join(dir, "doc.pdf")
	sidecar := filepath.Join(dir, "._doc.pdf")
	if err := os.WriteFile(native, []byte("payload"), 0o644); err != nil {
		t.Fatal(err)
	}
	golden, err := os.ReadFile(filepath.Join("..", "appledouble", "testdata", "rich.pdf.appledouble"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sidecar, golden, 0o644); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Lstat(native)
	if err != nil {
		t.Fatal(err)
	}
	if !fstype.NewCache().SupportsMergeXattrs(fi, native) {
		t.Skip("TEST_DIR/TMPDIR volume is not APFS/HFS+; set TEST_DIR to an APFS path")
	}

	opts := &cli.Options{
		Dirs: []string{dir},
		Keep: cli.KeepMostRecent,
	}
	if err := clean.Run(opts); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatal("sidecar should be deleted after merge")
	}
	size, err := unix.Lgetxattr(native, "com.apple.quarantine", nil)
	if err != nil || size == 0 {
		t.Fatalf("expected quarantine on native after merge: err=%v size=%d", err, size)
	}
}
