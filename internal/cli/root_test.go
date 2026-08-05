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

func TestSubcommands_SilentSuccessWithNoConfigFile(t *testing.T) {
	for _, name := range []string{"init", "propose", "apply", "check"} {
		t.Run(name, func(t *testing.T) {
			deps, stdout, stderr := app.NewFakeDeps(time.Now())
			root := cli.NewRootCmd(deps)
			root.SetArgs([]string{name})

			if err := root.Execute(); err != nil {
				t.Fatalf("Execute(%s): unexpected error: %v", name, err)
			}
			if stdout.Len() != 0 {
				t.Errorf("Execute(%s): expected empty stdout, got %q", name, stdout.String())
			}
			if stderr.Len() != 0 {
				t.Errorf("Execute(%s): expected empty stderr, got %q", name, stderr.String())
			}
		})
	}
}
