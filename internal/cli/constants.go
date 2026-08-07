package cli

const (
	// ExitOK is returned after successful help (-h / --help).
	ExitOK = 0
	// ExitFailure is returned when no directories are given or clean.Run fails.
	ExitFailure = 1
	// ExitUsage is returned for invalid flags or --keep values.
	ExitUsage = 2
	// ExitContinue means ParseOpts succeeded and the caller should run the tool.
	ExitContinue = -1
)
