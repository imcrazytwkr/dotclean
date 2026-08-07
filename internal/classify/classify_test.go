package classify_test

import (
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/classify"
)

func TestClassifyAppleDouble(t *testing.T) {
	c, ok := classify.ClassifyName("/vol", "._photo.jpg")
	if !ok {
		t.Fatal("expected match")
	}
	if c.Path != filepath.Join("/vol", "._photo.jpg") {
		t.Fatal(c.Path)
	}
	if c.Native != filepath.Join("/vol", "photo.jpg") {
		t.Fatal(c.Native)
	}
}

func TestClassifyIgnoresDSStore(t *testing.T) {
	if _, ok := classify.ClassifyName("/vol", ".DS_Store"); ok {
		t.Fatal(".DS_Store must not be a target")
	}
}

func TestClassifyRejectsEmptyNative(t *testing.T) {
	if _, ok := classify.ClassifyName("/vol", "._"); ok {
		t.Fatal(`"._" must not be a target`)
	}
}

func TestClassifySpacedName(t *testing.T) {
	c, ok := classify.ClassifyName("/my vol", "._file name.jpg")
	if !ok {
		t.Fatal("expected match")
	}
	if c.Native != filepath.Join("/my vol", "file name.jpg") {
		t.Fatal(c.Native)
	}
}
