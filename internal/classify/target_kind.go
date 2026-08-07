package classify

// TargetKind identifies junk/sidecar classes.
type TargetKind int

const (
	KindAppleDouble TargetKind = iota
	KindDeepJunk
	KindSpotlightJunk
)
