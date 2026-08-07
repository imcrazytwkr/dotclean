package fstype

import "os"

// Cache avoids repeated Statfs per mount (keyed by st_dev).
type Cache struct {
	m map[uint64]bool
}

func NewCache() *Cache {
	return &Cache{m: make(map[uint64]bool)}
}

// SupportsMergeXattrs reports whether path's volume can store native Mac xattrs
// (HFS+/APFS only). fi must be from Lstat(path). Nested mounts are distinct
// cache entries via st_dev. Non-Darwin builds always return false.
func (c *Cache) SupportsMergeXattrs(fi os.FileInfo, path string) bool {
	if !mergeXattrsPossible {
		return false
	}

	dev, ok := fileDev(fi)
	if !ok {
		return supportsMergeXattrs(path)
	}

	v, ok := c.m[dev]
	if ok {
		return v
	}

	v = supportsMergeXattrs(path)
	c.m[dev] = v
	return v
}
