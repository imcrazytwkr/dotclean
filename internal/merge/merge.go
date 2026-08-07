package merge

import (
	"os"

	"github.com/imcrazytwkr/dotclean/internal/appledouble"
	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/collections"
)

const (
	xattrFinderInfo   = "com.apple.FinderInfo"
	xattrResourceFork = "com.apple.ResourceFork"
)

// Apply merges AppleDouble sidecar into native path according to keep mode.
func Apply(sidecar, native string, keep cli.KeepMode) error {
	if keep == cli.KeepNative {
		return nil
	}

	raw, err := os.ReadFile(sidecar)
	if err != nil {
		return err
	}

	double, err := appledouble.Parse(raw)
	if err != nil {
		return err
	}

	var existing collections.Set[string]
	if keep == cli.KeepMostRecent {
		existing, err = listXattr(native)
		if err != nil {
			return err
		}
	}

	set := func(name string, val []byte) error {
		if keep == cli.KeepMostRecent && existing.Contains(name) {
			return nil // prefer native when present
		}

		return setXattr(native, name, val)
	}

	if double.HasFinderInfo && len(double.FinderInfo) == appledouble.FinderInfoSize {
		if err := set(xattrFinderInfo, double.FinderInfo); err != nil {
			return err
		}
	}

	if len(double.ResourceFork) > 0 {
		if err := set(xattrResourceFork, double.ResourceFork); err != nil {
			return err
		}
	}

	for _, a := range double.Attrs {
		if len(a.Name) == 0 {
			continue
		}
		if err := set(a.Name, a.Value); err != nil {
			return err
		}
	}
	return nil
}
