package config_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"

	"github.com/mkyc/changes/internal/config"
)

func newFlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("since", "", "")
	return fs
}

func TestLoad_Defaults(t *testing.T) {
	fsys := afero.NewMemMapFs()

	cfg, err := config.Load(fsys, newFlagSet())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Package != "." {
		t.Errorf("Package = %q, want %q", cfg.Package, ".")
	}
	if cfg.Changelog != "CHANGELOG.md" {
		t.Errorf("Changelog = %q, want %q", cfg.Changelog, "CHANGELOG.md")
	}
	if cfg.Since != "main" {
		t.Errorf("Since = %q, want %q", cfg.Since, "main")
	}
	if cfg.TagPrefix != "v" {
		t.Errorf("TagPrefix = %q, want %q", cfg.TagPrefix, "v")
	}

	want := map[string][]string{
		"major": {"feat!", "BREAKING CHANGE"},
		"minor": {"feat"},
		"patch": {"fix", "perf", "docs", "style", "refactor", "test", "build", "chore"},
		"none":  {"ci"},
	}
	for key, vals := range want {
		got, ok := cfg.Conventional[key]
		if !ok {
			t.Fatalf("Conventional[%q] missing", key)
		}
		if len(got) != len(vals) {
			t.Fatalf("Conventional[%q] = %v, want %v", key, got, vals)
		}
		for i := range vals {
			if got[i] != vals[i] {
				t.Fatalf("Conventional[%q] = %v, want %v", key, got, vals)
			}
		}
	}
}

func TestLoad_PartialConfigFallsBackToDefaults(t *testing.T) {
	fsys := afero.NewMemMapFs()
	if err := fsys.MkdirAll(".changes", 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	yaml := "changelog: HISTORY.md\n"
	if err := afero.WriteFile(fsys, config.ConfigFilePath, []byte(yaml), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := config.Load(fsys, newFlagSet())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Changelog != "HISTORY.md" {
		t.Errorf("Changelog = %q, want %q", cfg.Changelog, "HISTORY.md")
	}
	if cfg.Package != "." {
		t.Errorf("Package = %q, want default %q", cfg.Package, ".")
	}
	if cfg.Since != "main" {
		t.Errorf("Since = %q, want default %q", cfg.Since, "main")
	}
	if cfg.TagPrefix != "v" {
		t.Errorf("TagPrefix = %q, want default %q", cfg.TagPrefix, "v")
	}
}

func TestLoad_CLIFlagWinsOverConfigFile(t *testing.T) {
	fsys := afero.NewMemMapFs()
	if err := fsys.MkdirAll(".changes", 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	yaml := "since: develop\n"
	if err := afero.WriteFile(fsys, config.ConfigFilePath, []byte(yaml), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	flags := newFlagSet()
	if err := flags.Set("since", "release-branch"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	cfg, err := config.Load(fsys, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Since != "release-branch" {
		t.Errorf("Since = %q, want CLI value %q", cfg.Since, "release-branch")
	}
}
