package main

import (
	"fmt"
	"os"

	"github.com/mkyc/changes/internal/app"
	"github.com/mkyc/changes/internal/cli"
)

func main() {
	deps := app.NewRealDeps()
	root := cli.NewRootCmd(deps)

	if err := root.Execute(); err != nil {
		fmt.Fprintf(deps.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
