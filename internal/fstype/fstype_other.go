//go:build !darwin

package fstype

const mergeXattrsPossible = false

func supportsMergeXattrs(path string) bool {
	// Signature parity with darwin; unused on !darwin.
	_ = path
	return false
}
