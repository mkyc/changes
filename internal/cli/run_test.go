package cli_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/spf13/afero"

	"github.com/mkyc/changes/internal/app"
	"github.com/mkyc/changes/internal/cli"
	"github.com/mkyc/changes/internal/config"
)

func TestRun_InvalidNestedConventionalFollowsD011Contract(t *testing.T) {
	deps, stdout, stderr := app.NewFakeDeps(time.Now())
	if err := deps.FS.MkdirAll(".changes", 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := afero.WriteFile(deps.FS, config.ConfigFilePath, []byte("conventional:\n  major: feat\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	code := cli.Run(deps, []string{"check"})

	if code != 1 {
		t.Errorf("Run(check): expected exit code 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf("Run(check): expected empty stdout, got %q", stdout.String())
	}
	if !regexp.MustCompile(`^error: conventional\.major must be a list of strings\n$`).MatchString(stderr.String()) {
		t.Errorf("Run(check): expected single-line D011 stderr, got %q", stderr.String())
	}
}
