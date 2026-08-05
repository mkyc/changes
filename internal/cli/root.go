// Package cli builds the changes command tree.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/mkyc/changes/internal/app"
)

// NewRootCmd constructs the root "changes" command and attaches the
// init, propose, apply, and check subcommands.
func NewRootCmd(deps *app.Deps) *cobra.Command {
	root := &cobra.Command{
		Use:           "changes",
		Short:         "changes manages changesets and generates CHANGELOG.md",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(newInitCmd(deps))
	root.AddCommand(newProposeCmd(deps))
	root.AddCommand(newApplyCmd(deps))
	root.AddCommand(newCheckCmd(deps))

	return root
}
