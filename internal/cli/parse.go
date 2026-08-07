package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

// ParseOpts parses argv (typically os.Args).
// If the returned code is ExitContinue, opts is ready to run.
// Any other code is a process exit status (opts may be nil).
func ParseOpts(args []string) (*Options, int) {
	opts := &Options{
		Keep:  KeepMostRecent,
		Argv0: "dotclean",
	}

	if len(args) > 0 {
		opts.Argv0 = args[0]
	}

	fs := pflag.NewFlagSet("dotclean", pflag.ContinueOnError)
	fs.SetOutput(io.Discard)

	fs.BoolVarP(&opts.Flat, "flat", "f", false, "flat merge; do not recurse into subdirectories")
	fs.BoolVarP(&opts.AlwaysDelete, "always-delete", "m", false, "always delete AppleDouble (._*) files")
	fs.BoolVarP(&opts.CleanupOrphans, "cleanup", "n", false, "delete AppleDouble file if there is no matching native file")
	fs.BoolVarP(&opts.Preserve, "preserve", "p", false, "preserve AppleDouble file after handling")
	fs.BoolVarP(&opts.FollowSymlinks, "follow-symlinks", "s", false, "follow symbolic links to AppleDouble files")
	fs.BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose output")
	fs.BoolVarP(&opts.DryRun, "dry-run", "N", false, "list deletion targets only; do not merge or delete")

	var help bool
	fs.BoolVarP(&help, "help", "h", false, "print help and exit")

	keep := fs.String("keep", string(KeepMostRecent), "mostrecent|dotbar|native")

	fs.Usage = func() {
		PrintHelp(os.Stderr)
	}

	argv := []string{}
	if len(args) > 1 {
		argv = args[1:]
	}

	err := fs.Parse(argv)
	if err != nil {
		if err == pflag.ErrHelp {
			PrintHelp(os.Stdout)
			return nil, ExitOK
		}

		PrintHelp(os.Stderr)
		fmt.Fprintln(os.Stderr, err)
		return nil, ExitUsage
	}

	if help {
		PrintHelp(os.Stdout)
		return nil, ExitOK
	}

	var ok bool
	opts.Keep, ok = ParseKeepMode(*keep)
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid --keep value %q (want mostrecent|dotbar|native)\n", *keep)
		return nil, ExitUsage
	}

	opts.Dirs = append([]string(nil), fs.Args()...)
	if len(opts.Dirs) == 0 {
		PrintHelp(os.Stderr)
		return nil, ExitFailure
	}

	// -m overrides -p
	if opts.AlwaysDelete {
		opts.Preserve = false
	}

	return opts, ExitContinue
}
