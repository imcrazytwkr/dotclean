package action

import (
	"errors"
	"fmt"
	"os"

	"github.com/imcrazytwkr/dotclean/internal/cli"
	"github.com/imcrazytwkr/dotclean/internal/collections"
)

type Kind int

const (
	Merge Kind = iota
	Delete
)

type Action struct {
	Kind   Kind
	Path   string // sidecar path for Delete/Merge
	Native string // native path for Merge
}

func DryRunPrint(actions []Action) {
	for _, a := range actions {
		if a.Kind == Delete {
			fmt.Println(a.Path)
		}
	}
}

func Execute(actions []Action, opts *cli.Options, mergeFn func(sidecar, native string, opts *cli.Options) error) error {
	var failed uint64
	var mergeFailed collections.Set[string]

	for _, a := range actions {
		switch a.Kind {
		case Merge:
			if opts.Verbose {
				fmt.Fprintf(os.Stderr, "merge %s -> %s\n", a.Path, a.Native)
			}

			if mergeFn == nil {
				continue
			}

			err := mergeFn(a.Path, a.Native, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "merge %s: %v\n", a.Path, err)
				mergeFailed.Add(a.Path)
				failed++
			}
		case Delete:
			if mergeFailed.Contains(a.Path) {
				if opts.Verbose {
					fmt.Fprintf(os.Stderr, "keep %s after merge failure\n", a.Path)
				}
				continue
			}

			if opts.Verbose {
				fmt.Fprintf(os.Stderr, "delete %s\n", a.Path)
			}

			err := os.Remove(a.Path)
			if err != nil && !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "delete %s: %v\n", a.Path, err)
				failed++
			}
		}
	}

	switch failed {
	case 0:
		return nil
	case 1:
		return errors.New("one operation failed")
	default:
		return fmt.Errorf("%d operations failed", failed)
	}
}
