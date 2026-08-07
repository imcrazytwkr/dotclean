package testenv

import (
	"os"
	"testing"
)

// Dir returns the base directory for integration tests:
// TEST_DIR, then TMPDIR, then /tmp.
func Dir() string {
	d := os.Getenv("TEST_DIR")
	if len(d) > 0 {
		return d
	}

	d = os.Getenv("TMPDIR")
	if len(d) > 0 {
		return d
	}

	return "/tmp"
}

func TempDir(t testing.TB, pattern string) string {
	t.Helper()

	dir, err := os.MkdirTemp(Dir(), pattern)
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}

	t.Cleanup(func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("TempDir cleanup: %v", err)
		}
	})

	return dir
}
