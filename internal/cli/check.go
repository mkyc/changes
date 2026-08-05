package cli

import (
	"github.com/spf13/cobra"

	"github.com/mkyc/changes/internal/app"
	"github.com/mkyc/changes/internal/config"
)

func newCheckCmd(deps *app.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Validate changesets and changelog drift",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load(deps.FS, cmd.Flags())
			return err
		},
	}
}
