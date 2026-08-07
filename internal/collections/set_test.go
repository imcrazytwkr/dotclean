package collections_test

import (
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/collections"
)

func TestSetAddContains(t *testing.T) {
	var s collections.Set[string]
	if !s.IsEmpty() || s.Size() != 0 {
		t.Fatal("zero set should be empty")
	}
	if s.Contains("a") {
		t.Fatal("missing key")
	}
	s.Add("a")
	s.Add("a")
	if !s.Contains("a") || s.Size() != 1 || s.IsEmpty() {
		t.Fatalf("after Add: contains=%v size=%d empty=%v", s.Contains("a"), s.Size(), s.IsEmpty())
	}
	if s.Contains("b") {
		t.Fatal("unexpected key")
	}
}
