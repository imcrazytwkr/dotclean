package appledouble

const (
	Magic    = 0x00051607
	Version1 = 0x00010000
	Version2 = 0x00020000

	EntryDataFork   = 1
	EntryResource   = 2
	EntryRealName   = 3
	EntryComment    = 4
	EntryFileDates  = 8
	EntryFinderInfo = 9
	EntryMacInfo    = 10

	FinderInfoSize = 32
	AttrMagic      = 0x41545452 // 'ATTR'
)
