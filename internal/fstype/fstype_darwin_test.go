//go:build darwin

package fstype_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/fstype"
	"github.com/imcrazytwkr/dotclean/internal/testenv"
	"golang.org/x/sys/unix"
)

func TestSupportsMergeXattrsMatchesStatfs(t *testing.T) {
	dir := testenv.TempDir(t, "dotclean-fstype-*")

	path := filepath.Join(dir, "probe")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}

	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		t.Fatal(err)
	}
	name := fstypeName(st.Fstypename)
	want := name == "apfs" || name == "hfs"

	cache := fstype.NewCache()
	got := cache.SupportsMergeXattrs(fi, path)
	if got != want {
		t.Fatalf("SupportsMergeXattrs=%v want %v (fstype %q)", got, want, name)
	}
	if cache.SupportsMergeXattrs(fi, path) != got {
		t.Fatal("cached result mismatch")
	}
}

func fstypeName(b [16]byte) string {
	i := 0
	for i < len(b) && b[i] != 0 {
		i++
	}
	return strings.ToLower(string(b[:i]))
}
