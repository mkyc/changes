package cli_test

import (
	"regexp"
	"strings"
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

	if got := len(root.Commands()); got != 4 {
		names := make([]string, 0, got)
		for _, cmd := range root.Commands() {
			names = append(names, cmd.Name())
		}
		t.Errorf("expected exactly 4 subcommands, got %d: %v", got, names)
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

func TestPropose_SilentSuccessWithSinceFlag(t *testing.T) {
	deps, stdout, stderr := app.NewFakeDeps(time.Now())
	root := cli.NewRootCmd(deps)
	root.SetArgs([]string{"propose", "--since", "HEAD~1"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute(propose --since HEAD~1): unexpected error: %v", err)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
}

func TestSubcommands_RejectExtraPositionalArgs(t *testing.T) {
	for _, name := range []string{"init", "propose", "apply", "check"} {
		t.Run(name, func(t *testing.T) {
			deps, _, _ := app.NewFakeDeps(time.Now())
			root := cli.NewRootCmd(deps)
			root.SetArgs([]string{name, "extra-arg"})

			err := root.Execute()
			if err == nil {
				t.Fatalf("Execute(%s, extra-arg): expected error, got nil", name)
			}
			if err.Error() == "" {
				t.Errorf("Execute(%s, extra-arg): expected non-empty error message", name)
			}
		})
	}
}

func TestRun_UnknownSubcommandFollowsD011Contract(t *testing.T) {
	deps, stdout, stderr := app.NewFakeDeps(time.Now())

	code := cli.Run(deps, []string{"bogus"})

	if code != 1 {
		t.Errorf("Run(bogus): expected exit code 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf("Run(bogus): expected empty stdout, got %q", stdout.String())
	}
	if !regexp.MustCompile(`^error: `).MatchString(stderr.String()) {
		t.Errorf("Run(bogus): expected stderr to match %q, got %q", `^error: `, stderr.String())
	}
}

func TestRun_UnknownSubcommandStderrIsSingleLine(t *testing.T) {
	deps, stdout, stderr := app.NewFakeDeps(time.Now())

	code := cli.Run(deps, []string{"bogus"})

	if code != 1 {
		t.Errorf("Run(bogus): expected exit code 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf("Run(bogus): expected empty stdout, got %q", stdout.String())
	}
	if !regexp.MustCompile(`^error: .+\n$`).MatchString(stderr.String()) {
		t.Errorf("Run(bogus): expected stderr to match %q, got %q", `^error: .+\n$`, stderr.String())
	}
	if got := strings.Count(stderr.String(), "\n"); got != 1 {
		t.Errorf("Run(bogus): expected exactly one newline in stderr, got %d in %q", got, stderr.String())
	}
}
