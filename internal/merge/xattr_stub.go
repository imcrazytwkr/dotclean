//go:build !darwin

package merge

import "github.com/imcrazytwkr/dotclean/internal/collections"

func listXattr(path string) (collections.Set[string], error) {
	return collections.Set[string]{}, nil
}

func setXattr(path, name string, val []byte) error {
	return nil
}

func isIgnoredXattrError(err error) bool {
	return false
}
