package cli

// Options holds parsed CLI flags. Deep/Spotlight reserved for future modes.
type Options struct {
	Flat           bool
	AlwaysDelete   bool
	CleanupOrphans bool
	Preserve       bool
	FollowSymlinks bool
	Verbose        bool
	DryRun         bool
	SetQuarantine  bool // apply com.apple.quarantine from AppleDouble (off by default)
	Keep           KeepMode
	Dirs           []string
	Deep           bool // future
	Spotlight      bool // future
	Argv0          string
}

func (o *Options) ShouldDeletePaired() bool {
	if o.AlwaysDelete {
		return true
	}

	if o.Preserve {
		return false
	}

	// default: delete after handle
	return true
}

func (o *Options) ShouldDeleteOrphan() bool {
	return o.AlwaysDelete || o.CleanupOrphans
}
