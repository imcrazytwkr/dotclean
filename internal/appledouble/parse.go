package appledouble

import (
	"encoding/binary"
	"fmt"
)

// File is a parsed AppleDouble header file.
type File struct {
	FinderInfo    []byte // 32 bytes if present
	ResourceFork  []byte
	Attrs         []Attr
	HasFinderInfo bool
}

type entryDesc struct {
	id, offset, length uint32
}

// Parse parses an AppleDouble v2 sidecar.
func Parse(data []byte) (*File, error) {
	if len(data) < 26 {
		return nil, ErrNotAppleDouble
	}

	magic := binary.BigEndian.Uint32(data[0:4])
	if magic != Magic {
		return nil, ErrNotAppleDouble
	}

	version := binary.BigEndian.Uint32(data[4:8])
	switch version {
	case Version1:
		return nil, ErrAppleDoubleV1
	case Version2:
		// ATTR layout is OS X v2 style.
		break
	default:
		return nil, fmt.Errorf("Unsupported AppleDouble version: %d", version)
	}

	nEntries := binary.BigEndian.Uint16(data[24:26])
	headerEnd := 26 + int(nEntries)*12
	if len(data) < headerEnd {
		return nil, fmt.Errorf("%w: truncated entry table", ErrCorrupt)
	}

	entries := make([]entryDesc, nEntries)
	for i := range nEntries {
		off := 26 + i*12
		e := entryDesc{
			id:     binary.BigEndian.Uint32(data[off : off+4]),
			offset: binary.BigEndian.Uint32(data[off+4 : off+8]),
			length: binary.BigEndian.Uint32(data[off+8 : off+12]),
		}

		if int(e.offset)+int(e.length) > len(data) {
			return nil, fmt.Errorf("%w: entry %d out of bounds", ErrCorrupt, e.id)
		}

		entries[i] = e
	}

	out := &File{}
	for _, e := range entries {
		chunk := data[e.offset : e.offset+e.length]
		switch e.id {
		case EntryFinderInfo:
			if len(chunk) < FinderInfoSize {
				return nil, fmt.Errorf("%w: short FinderInfo", ErrCorrupt)
			}

			out.HasFinderInfo = true
			out.FinderInfo = append([]byte(nil), chunk[:FinderInfoSize]...)

			if len(chunk) > FinderInfoSize {
				attrs, err := parseATTR(data, int(e.offset))
				if err != nil {
					return nil, err
				}

				out.Attrs = attrs
			}
		case EntryResource:
			out.ResourceFork = append([]byte{}, chunk...)
		}
	}
	return out, nil
}
