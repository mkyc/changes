package cli

import (
	"fmt"

	"github.com/mkyc/changes/internal/app"
)

// Run builds the root command, executes it with args, and returns the
// process exit code. On failure it writes the D011 error contract
// ("error: <message>") to deps.Stderr before returning 1.
func Run(deps *app.Deps, args []string) int {
	root := NewRootCmd(deps)
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintf(deps.Stderr, "error: %s\n", err)
		return 1
	}

	return 0
}
