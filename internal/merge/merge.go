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
	xattrQuarantine   = "com.apple.quarantine"
)

// Apply merges AppleDouble sidecar into native path according to opts.
// A nil opts is treated as KeepMostRecent with SetQuarantine false.
func Apply(sidecar, native string, opts *cli.Options) error {
	keep := cli.KeepMostRecent
	setQuarantine := false
	if opts != nil {
		keep = opts.Keep
		setQuarantine = opts.SetQuarantine
	}

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
		if a.Name == xattrQuarantine && !setQuarantine {
			continue
		}
		if err := set(a.Name, a.Value); err != nil {
			return err
		}
	}
	return nil
}
