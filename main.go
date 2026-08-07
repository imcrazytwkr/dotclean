package main

import (
	"fmt"
	"os"

	"github.com/imcrazytwkr/dotclean/internal/clean"
	"github.com/imcrazytwkr/dotclean/internal/cli"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	opts, code := cli.ParseOpts(args)
	if code != cli.ExitContinue {
		return code
	}

	err := clean.Run(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return cli.ExitFailure
	}

	return cli.ExitOK
}
