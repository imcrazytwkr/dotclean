package classify

// TargetKind identifies junk/sidecar classes. Extra kinds reserved for future deep/spotlight.
type TargetKind int

const (
	KindAppleDouble   TargetKind = iota
	KindDeepJunk                 // future
	KindSpotlightJunk            // future
)
