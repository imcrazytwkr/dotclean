package fstype

import (
	"os"

	"golang.org/x/sys/unix"
)

func fileDev(fi os.FileInfo) (uint64, bool) {
	st, ok := fi.Sys().(*unix.Stat_t)
	if !ok {
		return 0, false
	}

	// Darwin uses weird int32 encoding for device type
	// but since it fits in 64-bit width, we can just
	// dumb-cast it to use as a cache key
	return uint64(st.Dev), true
}
