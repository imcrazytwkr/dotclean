package cli

// Options holds parsed CLI flags.
type Options struct {
	Flat           bool
	AlwaysDelete   bool
	CleanupOrphans bool
	Preserve       bool
	FollowSymlinks bool
	Verbose        bool
	DryRun         bool
	SetQuarantine  bool // apply com.apple.quarantine from AppleDouble (off by default)
	Deep           bool // remove deep junk (.DS_Store, .Trashes, …)
	Spotlight      bool // remove Spotlight/fsevents junk
	Keep           KeepMode
	Dirs           []string
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

// DeepEnabled reports whether deep junk removal is active (-D and not -p).
func (o *Options) DeepEnabled() bool {
	return o.Deep && !o.Preserve
}

// SpotlightEnabled reports whether spotlight junk removal is active (-S and not -p).
func (o *Options) SpotlightEnabled() bool {
	return o.Spotlight && !o.Preserve
}
