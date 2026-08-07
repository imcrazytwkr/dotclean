package walkdir

import (
	"os"
	"path/filepath"

	"github.com/imcrazytwkr/dotclean/internal/classify"
	"github.com/imcrazytwkr/dotclean/internal/cli"
)

func resolveSymlink(c classify.Candidate, opts *cli.Options) (classify.Candidate, bool) {
	fi, err := os.Lstat(c.Path)
	if err != nil {
		return c, false
	}

	if !isSymlink(fi) {
		return c, true
	}

	if !opts.FollowSymlinks {
		return c, true // operate on the symlink path itself (Apple: follow only with -s)
	}

	target, err := filepath.EvalSymlinks(c.Path)
	if err != nil {
		return c, false
	}
	c.Path = target
	dir, name := filepath.Split(target)
	if n, ok := classify.ClassifyName(filepath.Clean(dir), name, opts.DeepEnabled(), opts.SpotlightEnabled()); ok {
		c.Native = n.Native
		c.Kind = n.Kind
	}
	return c, true
}

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}
