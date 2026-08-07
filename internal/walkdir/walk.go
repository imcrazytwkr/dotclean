package walkdir

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/imcrazytwkr/dotclean/internal/classify"
	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/collections"
)

// Collect walks roots and returns classified candidates.
func Collect(opts *cli.Options) ([]classify.Candidate, error) {
	var out []classify.Candidate
	var seen collections.Set[string]

	for _, root := range opts.Dirs {
		err := walkOne(root, opts, func(c classify.Candidate) {
			if seen.Contains(c.Path) {
				return
			}

			seen.Add(c.Path)
			out = append(out, c)
		})

		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func walkOne(root string, opts *cli.Options, emit func(classify.Candidate)) error {
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s: not a directory", root)
	}

	if opts.Flat {
		return walkFlat(root, opts, emit)
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.Verbose {
				fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			}

			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		// Never treat the walk root as a junk target.
		if path == root {
			return nil
		}

		name := d.Name()
		c, ok := classify.ClassifyName(filepath.Dir(path), name, opts.DeepEnabled(), opts.SpotlightEnabled())
		if !ok {
			return nil
		}

		if d.IsDir() {
			emit(c)
			return filepath.SkipDir
		}

		c, ok = resolveSymlink(c, opts)
		if !ok {
			return nil
		}

		emit(c)
		return nil
	})
}

func walkFlat(root string, opts *cli.Options, emit func(classify.Candidate)) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, e := range entries {
		c, ok := classify.ClassifyName(root, e.Name(), opts.DeepEnabled(), opts.SpotlightEnabled())
		if !ok {
			continue
		}

		if e.IsDir() {
			emit(c)
			continue
		}

		c, ok = resolveSymlink(c, opts)
		if !ok {
			continue
		}

		emit(c)
	}
	return nil
}
