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
// deep/spotlight enable extra junk kinds; AppleDouble is always considered.
func ClassifyName(dir, name string, deep, spotlight bool) (Candidate, bool) {
	path := filepath.Join(dir, name)

	if strings.HasPrefix(name, "._") {
		nativeBase := name[2:]
		if len(nativeBase) == 0 {
			return Candidate{}, false
		}
		return Candidate{
			Kind:   KindAppleDouble,
			Path:   path,
			Native: filepath.Join(dir, nativeBase),
		}, true
	}

	if deep {
		switch name {
		case ".DS_Store", ".AppleDouble", ".Trashes", ".TemporaryItems":
			return Candidate{Kind: KindDeepJunk, Path: path}, true
		}
	}

	if spotlight {
		if name == ".fseventsd" || strings.HasPrefix(name, ".Spotlight") {
			return Candidate{Kind: KindSpotlightJunk, Path: path}, true
		}
	}

	return Candidate{}, false
}
