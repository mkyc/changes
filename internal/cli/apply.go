package cli

import (
	"github.com/spf13/cobra"

	"github.com/mkyc/changes/internal/app"
	"github.com/mkyc/changes/internal/config"
)

func newApplyCmd(deps *app.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply pending changesets to the changelog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load(deps.FS, cmd.Flags())
			return err
		},
	}
}
