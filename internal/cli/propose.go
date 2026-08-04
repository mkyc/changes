package cli

import (
	"github.com/spf13/cobra"

	"github.com/mkyc/changes/internal/app"
	"github.com/mkyc/changes/internal/config"
)

func newProposeCmd(deps *app.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose",
		Short: "Propose a new changeset from recent commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load(deps.FS, cmd.Flags())
			return err
		},
	}
	cmd.Flags().String("since", "", "base ref to diff against when detecting changes")
	return cmd
}
