package cli

type KeepMode string

const (
	KeepMostRecent KeepMode = "mostrecent"
	KeepDotbar     KeepMode = "dotbar"
	KeepNative     KeepMode = "native"
)

func ParseKeepMode(value string) (KeepMode, bool) {
	v := KeepMode(value)
	switch v {
	case KeepMostRecent, KeepDotbar, KeepNative:
		return v, true
	default:
		return KeepMostRecent, false
	}
}
