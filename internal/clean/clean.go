package clean

import (
	"fmt"
	"os"

	"github.com/imcrazytwkr/dotclean/internal/action"
	"github.com/imcrazytwkr/dotclean/internal/classify"
	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/fstype"
	"github.com/imcrazytwkr/dotclean/internal/merge"
	"github.com/imcrazytwkr/dotclean/internal/walkdir"
)

// Run executes the clean pipeline for opts.
func Run(options *cli.Options) error {
	candidates, err := walkdir.Collect(options)
	if err != nil {
		return err
	}

	fsTypeCache := fstype.NewCache()
	actions := plan(candidates, options, fsTypeCache)

	if options.DryRun {
		action.DryRunPrint(actions)
		return nil
	}

	return action.Execute(actions, options, merge.Apply)
}

func plan(cands []classify.Candidate, opts *cli.Options, fsCache *fstype.Cache) []action.Action {
	var out []action.Action
	for _, c := range cands {
		if c.Kind != classify.KindAppleDouble {
			continue
		}

		fi, err := os.Lstat(c.Native)
		if err != nil {
			if opts.ShouldDeleteOrphan() {
				out = append(out, action.Action{Kind: action.Delete, Path: c.Path})
				continue
			}

			if opts.Verbose {
				fmt.Fprintf(os.Stderr, "skip orphan %s\n", c.Path)
			}

			continue
		}

		if fsCache.SupportsMergeXattrs(fi, c.Native) && opts.Keep != cli.KeepNative {
			out = append(out, action.Action{Kind: action.Merge, Path: c.Path, Native: c.Native})
		} else if opts.Verbose {
			fmt.Fprintf(os.Stderr, "discard %s (no native xattr store or --keep=native)\n", c.Path)
		}

		if opts.ShouldDeletePaired() {
			out = append(out, action.Action{Kind: action.Delete, Path: c.Path})
		}
	}
	return out
}
