package cli_test

import (
	"testing"
	"time"

	"github.com/mkyc/changes/internal/app"
	"github.com/mkyc/changes/internal/cli"
)

func TestNewRootCmd_CommandDiscovery(t *testing.T) {
	deps, _, _ := app.NewFakeDeps(time.Now())
	root := cli.NewRootCmd(deps)

	want := map[string]bool{"init": false, "propose": false, "apply": false, "check": false}
	for _, cmd := range root.Commands() {
		if _, ok := want[cmd.Name()]; ok {
			want[cmd.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestNewRootCmd_ProposeHasSinceFlag(t *testing.T) {
	deps, _, _ := app.NewFakeDeps(time.Now())
	root := cli.NewRootCmd(deps)

	proposeCmd, _, err := root.Find([]string{"propose"})
	if err != nil {
		t.Fatalf("Find(propose): %v", err)
	}
	if proposeCmd.Flags().Lookup("since") == nil {
		t.Error("expected propose command to declare a --since flag")
	}
}
