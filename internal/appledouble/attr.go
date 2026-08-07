package appledouble

import (
	"encoding/binary"
	"fmt"
)

// Attr is a named extended attribute from the ATTR blob after FinderInfo.
type Attr struct {
	Name  string
	Value []byte
}

// parseATTR parses the Mac OS X ATTR blob that follows FinderInfo inside entry 9.
// finderInfoOff is the file offset of the FinderInfo entry data.
// ATTR header/entry offsets (data_start, adx_offset) are file-absolute (XNU/Samba).
func parseATTR(file []byte, finderInfoOff int) ([]Attr, error) {
	start := finderInfoOff + FinderInfoSize
	if start > len(file) {
		return nil, nil
	}

	// Skip up to 2 pad bytes before ATTR magic.
	i := start
	end := min(start+2, len(file))
	for i < end && !(len(file)-i >= 4 && binary.BigEndian.Uint32(file[i:i+4]) == AttrMagic) {
		i++
	}

	if len(file)-i < 36 {
		return nil, nil
	}

	if binary.BigEndian.Uint32(file[i:i+4]) != AttrMagic {
		return nil, nil
	}

	base := i
	hdr := file[base:]
	// Header (36 bytes):
	// 0 magic, 4 debug, 8 total_size, 12 data_start, 16 data_length,
	// 20 reserved[3], 32 flags, 34 num_attrs
	numAttrs := int(binary.BigEndian.Uint16(hdr[34:36]))

	attrs := make([]Attr, 0, numAttrs)
	off := 36 // relative to ATTR header
	for n := 0; n < numAttrs && base+off+11 <= len(file); n++ {
		entryOff := int(binary.BigEndian.Uint32(file[base+off : base+off+4]))
		entryLen := int(binary.BigEndian.Uint32(file[base+off+4 : base+off+8]))
		nameLen := int(file[base+off+10])
		off += 11

		if nameLen < 1 || base+off+nameLen > len(file) {
			return nil, fmt.Errorf("%w: bad ATTR name", ErrCorrupt)
		}

		nameBytes := file[base+off : base+off+nameLen]
		off += nameLen
		for off%4 != 0 {
			off++
		}

		name := string(nameBytes)
		if nameLen > 0 && nameBytes[nameLen-1] == 0 {
			name = string(nameBytes[:nameLen-1])
		}

		if entryOff < 0 || entryOff+entryLen > len(file) {
			return nil, fmt.Errorf("%w: ATTR data OOB for %q", ErrCorrupt, name)
		}

		val := append([]byte(nil), file[entryOff:entryOff+entryLen]...)
		if len(name) != 0 {
			attrs = append(attrs, Attr{Name: name, Value: val})
		}
	}

	return attrs, nil
}
