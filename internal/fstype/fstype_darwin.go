//go:build darwin

package fstype

import (
	"strings"

	"golang.org/x/sys/unix"
)

const mergeXattrsPossible = true

func supportsMergeXattrs(path string) bool {
	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		return false
	}

	switch fstypeName(st.Fstypename) {
	case "apfs", "hfs":
		return true
	default:
		return false
	}
}

// @TODO: consider re-implementing if gopls starts supporting inline-disable
// @SEE: https://github.com/golang/go/issues/50764
func fstypeName(b [16]byte) string {
	i := 0

	// @NOTE: Using classic `for` is necessary bere because while I would love
	// to use `bytes.IndexByte` here, gopls insists on me using wasteful and
	// largely redundant in this case `bytes.Cut`
	for b[i] != 0 && i < len(b) {
		i++
	}

	// Cheaper in allocations if the underlying bytes are expected to be
	// lowercase already
	return strings.ToLower(string(b[:i]))
}
