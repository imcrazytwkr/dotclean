package classify

import (
	"path/filepath"
	"strings"
)

// Candidate is a classified path.
type Candidate struct {
	Kind   TargetKind
	Path   string // sidecar or junk path
	Native string // for AppleDouble: native path without ._
}

// ClassifyName classifies a single file/dir base name within dir.
// Returns ok=false if the name is not a target under current (normal) mode.
func ClassifyName(dir, name string) (Candidate, bool) {
	if !strings.HasPrefix(name, "._") {
		return Candidate{}, false
	}

	nativeBase := name[2:]
	if len(nativeBase) == 0 {
		return Candidate{}, false
	}

	native := filepath.Join(dir, nativeBase)
	return Candidate{
		Kind:   KindAppleDouble,
		Path:   filepath.Join(dir, name),
		Native: native,
	}, true
}
