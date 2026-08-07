//go:build darwin

package merge

import (
	"github.com/imcrazytwkr/dotclean/internal/collections"
	"golang.org/x/sys/unix"
)

func listXattr(path string) (collections.Set[string], error) {
	var out collections.Set[string]

	size, err := unix.Llistxattr(path, nil)
	if err != nil || size == 0 {
		return out, err
	}

	buf := make([]byte, size)
	n, err := unix.Llistxattr(path, buf)
	if err != nil {
		return out, err
	}

	start := 0
	for i := range n {
		if buf[i] != 0 {
			continue
		}

		if i > start {
			out.Add(string(buf[start:i]))
		}

		start = i + 1
	}

	return out, nil
}

func setXattr(path, name string, val []byte) error {
	return unix.Lsetxattr(path, name, val, 0)
}
