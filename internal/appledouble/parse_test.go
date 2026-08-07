package appledouble_test

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/dotclean/internal/appledouble"
)

func TestParseMinimalFinderInfo(t *testing.T) {
	// header + 1 entry (FinderInfo at offset 38, length 32)
	// 26 header + 12 entry = 38
	data := make([]byte, 38+32)
	binary.BigEndian.PutUint32(data[0:4], appledouble.Magic)
	binary.BigEndian.PutUint32(data[4:8], appledouble.Version2)
	binary.BigEndian.PutUint16(data[24:26], 1)
	binary.BigEndian.PutUint32(data[26:30], appledouble.EntryFinderInfo)
	binary.BigEndian.PutUint32(data[30:34], 38)
	binary.BigEndian.PutUint32(data[34:38], 32)
	for i := 0; i < 32; i++ {
		data[38+i] = byte(i)
	}

	f, err := appledouble.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !f.HasFinderInfo || len(f.FinderInfo) != 32 || f.FinderInfo[3] != 3 {
		t.Fatalf("%+v", f)
	}
}

func TestParseRejectsBadMagic(t *testing.T) {
	_, err := appledouble.Parse([]byte("not apple double!!!!!!!!!!"))
	if err != appledouble.ErrNotAppleDouble {
		t.Fatalf("got %v", err)
	}
}

func TestParseRejectsV1(t *testing.T) {
	data := make([]byte, 26)
	binary.BigEndian.PutUint32(data[0:4], appledouble.Magic)
	binary.BigEndian.PutUint32(data[4:8], appledouble.Version1)
	_, err := appledouble.Parse(data)
	if !errors.Is(err, appledouble.ErrAppleDoubleV1) {
		t.Fatalf("got %v", err)
	}
}

func TestParseRejectsUnknownVersion(t *testing.T) {
	data := make([]byte, 26)
	binary.BigEndian.PutUint32(data[0:4], appledouble.Magic)
	binary.BigEndian.PutUint32(data[4:8], 0x00030000)
	_, err := appledouble.Parse(data)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseGoldenFixtures(t *testing.T) {
	dir := "testdata"
	manifestPath := filepath.Join(dir, "manifest.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest []struct {
		File      string   `json:"file"`
		AttrNames []string `json:"attr_names"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	for _, entry := range manifest {
		t.Run(entry.File, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(dir, entry.File))
			if err != nil {
				t.Fatal(err)
			}
			f, err := appledouble.Parse(data)
			if err != nil {
				t.Fatal(err)
			}
			got := make([]string, 0, len(f.Attrs))
			for _, a := range f.Attrs {
				got = append(got, a.Name)
				if len(a.Value) == 0 {
					t.Fatalf("empty value for %q", a.Name)
				}
			}
			if len(got) != len(entry.AttrNames) {
				t.Fatalf("attrs=%v want %v", got, entry.AttrNames)
			}
			for i := range got {
				if got[i] != entry.AttrNames[i] {
					t.Fatalf("attrs[%d]=%q want %q (full got=%v)", i, got[i], entry.AttrNames[i], got)
				}
			}
		})
	}
}
