package cli

import (
	"fmt"
	"io"
	"strings"
)

func PrintHelp(w io.Writer) {
	fmt.Fprintln(w, strings.TrimSpace(`
usage: dotclean [options] directory ...

Merge AppleDouble (._*) sidecars with native files on HFS+/APFS, or discard
and remove them on other filesystems (e.g. FAT32/exFAT).

Options:
  -f, --flat              Do not recurse into subdirectories
  -h, --help              Print this help and exit
  -m, --always-delete     Always delete AppleDouble files
  -n, --cleanup           Delete AppleDouble if there is no matching native file
  -p, --preserve          Preserve AppleDouble file after handling
  -s, --follow-symlinks   Follow symbolic links to AppleDouble files
  -v, --verbose           Verbose output
  -N, --dry-run           List deletion targets only; do not merge or delete
      --keep=MODE         mostrecent (default), dotbar, or native

See docs/dotclean.1.md or 'man dotclean' for the full manual.
`))
}
