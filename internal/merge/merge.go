package merge

import (
	"fmt"
	"os"

	"github.com/imcrazytwkr/dotclean/internal/appledouble"
	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/collections"
)

const (
	xattrFinderInfo   = "com.apple.FinderInfo"
	xattrResourceFork = "com.apple.ResourceFork"
	xattrQuarantine   = "com.apple.quarantine"
	xattrMacl         = "com.apple.macl"
	xattrProvenance   = "com.apple.provenance"
)

// Apply merges AppleDouble sidecar into native path according to opts.
// A nil opts is treated as KeepMostRecent with SetQuarantine false.
func Apply(sidecar, native string, opts *cli.Options) error {
	keep := cli.KeepMostRecent
	setQuarantine := false
	verbose := false
	if opts != nil {
		keep = opts.Keep
		setQuarantine = opts.SetQuarantine
		verbose = opts.Verbose
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
		if len(a.Name) == 0 || skipAttr(a.Name, setQuarantine) {
			continue
		}
		if err := set(a.Name, a.Value); err != nil {
			if isIgnoredXattrError(err) {
				if verbose {
					fmt.Fprintf(os.Stderr, "skip xattr %s on %s: %v\n", a.Name, native, err)
				}
				continue
			}
			return err
		}
	}
	return nil
}

// skipAttr reports ATTR names that must not be written from AppleDouble.
func skipAttr(name string, setQuarantine bool) bool {
	switch name {
	case xattrMacl, xattrProvenance:
		return true
	case xattrQuarantine:
		return !setQuarantine
	default:
		return false
	}
}
