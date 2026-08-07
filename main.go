package main

import (
	"os"

	"github.com/mkyc/changes/internal/app"
	"github.com/mkyc/changes/internal/cli"
)

func main() {
	deps := app.NewRealDeps()
	os.Exit(cli.Run(deps, os.Args[1:]))
}
