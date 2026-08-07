//go:build darwin

package merge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/fstype"
	"github.com/imcrazytwkr/dotclean/internal/merge"
	"github.com/imcrazytwkr/dotclean/internal/testenv"
	"golang.org/x/sys/unix"
)

func TestApplyDotbarWritesXattrs(t *testing.T) {
	dir := requireMergeDir(t)
	native, sidecar := stageGolden(t, dir, "doc.pdf", "rich.pdf.appledouble")

	if err := merge.Apply(sidecar, native, &cli.Options{Keep: cli.KeepDotbar}); err != nil {
		t.Fatal(err)
	}

	names := listXattrNames(t, native)
	if names["com.apple.quarantine"] {
		t.Fatal("quarantine should be skipped without SetQuarantine")
	}
	if !names["com.apple.lastuseddate#PS"] {
		t.Fatalf("missing lastuseddate; have %v", names)
	}
}

func TestApplySetQuarantineWritesQuarantine(t *testing.T) {
	dir := requireMergeDir(t)
	native, sidecar := stageGolden(t, dir, "doc.pdf", "rich.pdf.appledouble")

	opts := &cli.Options{Keep: cli.KeepDotbar, SetQuarantine: true}
	if err := merge.Apply(sidecar, native, opts); err != nil {
		t.Fatal(err)
	}

	names := listXattrNames(t, native)
	if !names["com.apple.quarantine"] {
		t.Fatalf("missing quarantine; have %v", names)
	}
	val := getXattr(t, native, "com.apple.quarantine")
	if len(val) == 0 {
		t.Fatal("empty quarantine value")
	}
}

func TestApplyMostRecentKeepsExisting(t *testing.T) {
	dir := requireMergeDir(t)
	native, sidecar := stageGolden(t, dir, "doc.pdf", "rich.pdf.appledouble")

	keep := []byte("keep-me")
	if err := unix.Lsetxattr(native, "com.apple.quarantine", keep, 0); err != nil {
		t.Fatal(err)
	}

	opts := &cli.Options{Keep: cli.KeepMostRecent, SetQuarantine: true}
	if err := merge.Apply(sidecar, native, opts); err != nil {
		t.Fatal(err)
	}

	got := getXattr(t, native, "com.apple.quarantine")
	if string(got) != string(keep) {
		t.Fatalf("quarantine overwritten: %q", got)
	}
	names := listXattrNames(t, native)
	if !names["com.apple.lastuseddate#PS"] {
		t.Fatalf("expected lastuseddate from sidecar; have %v", names)
	}
}

func TestApplySkipsMaclAndProvenance(t *testing.T) {
	dir := requireMergeDir(t)
	native, sidecar := stageGolden(t, dir, "doc.epub", "rich.epub.appledouble")

	if err := merge.Apply(sidecar, native, &cli.Options{Keep: cli.KeepDotbar}); err != nil {
		t.Fatal(err)
	}

	names := listXattrNames(t, native)
	if names["com.apple.macl"] || names["com.apple.provenance"] {
		t.Fatalf("SIP-ish attrs should be skipped; have %v", names)
	}
	if !names["com.apple.lastuseddate#PS"] {
		t.Fatalf("expected lastuseddate from sidecar; have %v", names)
	}
}

func requireMergeDir(t *testing.T) string {
	t.Helper()
	dir := testenv.TempDir(t, "dotclean-merge-*")

	probe := filepath.Join(dir, "probe")
	if err := os.WriteFile(probe, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Lstat(probe)
	if err != nil {
		t.Fatal(err)
	}
	if !fstype.NewCache().SupportsMergeXattrs(fi, probe) {
		t.Skip("TEST_DIR/TMPDIR volume is not APFS/HFS+; set TEST_DIR to an APFS path")
	}
	return dir
}

func stageGolden(t *testing.T, dir, nativeBase, golden string) (native, sidecar string) {
	t.Helper()
	native = filepath.Join(dir, nativeBase)
	sidecar = filepath.Join(dir, "._"+nativeBase)
	if err := os.WriteFile(native, []byte("payload"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join("..", "appledouble", "testdata", golden)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sidecar, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return native, sidecar
}

func listXattrNames(t *testing.T, path string) map[string]bool {
	t.Helper()
	size, err := unix.Llistxattr(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, size)
	n, err := unix.Llistxattr(path, buf)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string]bool)
	start := 0
	for i := 0; i < n; i++ {
		if buf[i] != 0 {
			continue
		}
		if i > start {
			out[string(buf[start:i])] = true
		}
		start = i + 1
	}
	return out
}

func getXattr(t *testing.T, path, name string) []byte {
	t.Helper()
	size, err := unix.Lgetxattr(path, name, nil)
	if err != nil {
		t.Fatalf("getxattr %s: %v", name, err)
	}
	buf := make([]byte, size)
	n, err := unix.Lgetxattr(path, name, buf)
	if err != nil {
		t.Fatal(err)
	}
	return buf[:n]
}
