package classify_test

import (
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/classify"
)

func TestClassifyAppleDouble(t *testing.T) {
	c, ok := classify.ClassifyName("/vol", "._photo.jpg", false, false)
	if !ok {
		t.Fatal("expected match")
	}
	if c.Path != filepath.Join("/vol", "._photo.jpg") {
		t.Fatal(c.Path)
	}
	if c.Native != filepath.Join("/vol", "photo.jpg") {
		t.Fatal(c.Native)
	}
	if c.Kind != classify.KindAppleDouble {
		t.Fatal(c.Kind)
	}
}

func TestClassifyIgnoresDSStoreWithoutDeep(t *testing.T) {
	if _, ok := classify.ClassifyName("/vol", ".DS_Store", false, false); ok {
		t.Fatal(".DS_Store must not be a target without deep")
	}
}

func TestClassifyDeepJunk(t *testing.T) {
	for _, name := range []string{".DS_Store", ".AppleDouble", ".Trashes", ".TemporaryItems"} {
		c, ok := classify.ClassifyName("/vol", name, true, false)
		if !ok || c.Kind != classify.KindDeepJunk {
			t.Fatalf("%s: ok=%v kind=%v", name, ok, c.Kind)
		}
	}
}

func TestClassifySpotlightJunk(t *testing.T) {
	c, ok := classify.ClassifyName("/vol", ".fseventsd", false, true)
	if !ok || c.Kind != classify.KindSpotlightJunk {
		t.Fatalf("fseventsd: ok=%v kind=%v", ok, c.Kind)
	}
	c, ok = classify.ClassifyName("/vol", ".Spotlight-V100", false, true)
	if !ok || c.Kind != classify.KindSpotlightJunk {
		t.Fatalf("Spotlight: ok=%v kind=%v", ok, c.Kind)
	}
	if _, ok := classify.ClassifyName("/vol", ".Spotlight-V100", false, false); ok {
		t.Fatal("Spotlight must not match without -S")
	}
}

func TestClassifyRejectsEmptyNative(t *testing.T) {
	if _, ok := classify.ClassifyName("/vol", "._", false, false); ok {
		t.Fatal(`"._" must not be a target`)
	}
}

func TestClassifySpacedName(t *testing.T) {
	c, ok := classify.ClassifyName("/my vol", "._file name.jpg", false, false)
	if !ok {
		t.Fatal("expected match")
	}
	if c.Native != filepath.Join("/my vol", "file name.jpg") {
		t.Fatal(c.Native)
	}
}
