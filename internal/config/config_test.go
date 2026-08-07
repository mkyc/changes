package config_test

import (
	"strings"
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

func defaultConventional() map[string][]string {
	return map[string][]string{
		"major": {"feat!", "BREAKING CHANGE"},
		"minor": {"feat"},
		"patch": {"fix", "perf", "docs", "style", "refactor", "test", "build", "chore"},
		"none":  {"ci"},
	}
}

func assertConventional(t *testing.T, got map[string][]string, want map[string][]string) {
	t.Helper()
	for key, vals := range want {
		gotVals, ok := got[key]
		if !ok {
			t.Fatalf("Conventional[%q] missing", key)
		}
		if len(gotVals) != len(vals) {
			t.Fatalf("Conventional[%q] = %v, want %v", key, gotVals, vals)
		}
		for i := range vals {
			if gotVals[i] != vals[i] {
				t.Fatalf("Conventional[%q] = %v, want %v", key, gotVals, vals)
			}
		}
	}
}

func writeConfigFile(t *testing.T, fsys afero.Fs, yaml string) {
	t.Helper()
	if err := fsys.MkdirAll(".changes", 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := afero.WriteFile(fsys, config.ConfigFilePath, []byte(yaml), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
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

	want := defaultConventional()
	assertConventional(t, cfg.Conventional, want)
}

func TestLoad_PartialConfigFallsBackToDefaults(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "changelog: HISTORY.md\n")

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
	assertConventional(t, cfg.Conventional, defaultConventional())
}

func TestLoad_PartialConventionalFallsBackToDefaults(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional:\n  major:\n    - custom\n")

	cfg, err := config.Load(fsys, newFlagSet())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := defaultConventional()
	want["major"] = []string{"custom"}
	assertConventional(t, cfg.Conventional, want)
}

func TestLoad_EmptyConventionalBlockReturnsError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional:\n")

	_, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for empty conventional block")
	}
	if !strings.Contains(err.Error(), "conventional must contain at least one key") {
		t.Fatalf("error = %v, want conventional empty-block message", err)
	}
}

func TestLoad_EmptyMappingConventionalReturnsError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional: {}\n")

	_, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for empty mapping conventional block")
	}
	if !strings.Contains(err.Error(), "conventional must contain at least one key") {
		t.Fatalf("error = %v, want conventional empty-block message", err)
	}
}

func TestLoad_UnknownConventionalKeyReturnsError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional:\n  unknown:\n    - x\n")

	_, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for unknown conventional key")
	}
	if !strings.Contains(err.Error(), `unknown conventional key "unknown"`) {
		t.Fatalf("error = %v, want unknown key message", err)
	}
}

func TestLoad_NestedWrongTypeConventionalReturnsSingleLineError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional:\n  major: feat\n")

	_, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for nested wrong-type conventional value")
	}
	if !strings.Contains(err.Error(), "conventional.major must be a list of strings") {
		t.Fatalf("error = %v, want nested wrong-type message", err)
	}
	if strings.Contains(err.Error(), "\n") {
		t.Fatalf("error = %q, want no embedded newlines", err.Error())
	}
}

func TestLoad_SequenceConventionalReturnsError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional:\n  - a\n  - b\n")

	_, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for sequence conventional value")
	}
	if !strings.Contains(err.Error(), "conventional must be a mapping of keys to lists") {
		t.Fatalf("error = %v, want non-mapping message", err)
	}
}

func TestLoad_ScalarConventionalReturnsError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional: scalar\n")

	_, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for scalar conventional value")
	}
	if !strings.Contains(err.Error(), "conventional must be a mapping of keys to lists") {
		t.Fatalf("error = %v, want non-mapping message", err)
	}
}

func TestLoad_AliasToMappingConventionalSucceeds(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "defs:\n  conv: &conv\n    major:\n      - feat\nconventional: *conv\n")

	cfg, err := config.Load(fsys, newFlagSet())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := defaultConventional()
	want["major"] = []string{"feat"}
	assertConventional(t, cfg.Conventional, want)
}

func TestLoad_AliasToSequenceConventionalReturnsError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "defs:\n  conv: &conv\n    - a\n    - b\nconventional: *conv\n")

	_, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for alias to sequence conventional value")
	}
	if !strings.Contains(err.Error(), "conventional must be a mapping of keys to lists") {
		t.Fatalf("error = %v, want non-mapping message", err)
	}
}

func TestLoad_EmptyConventionalSliceHonored(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "conventional:\n  major: []\n")

	cfg, err := config.Load(fsys, newFlagSet())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := defaultConventional()
	want["major"] = []string{}
	assertConventional(t, cfg.Conventional, want)
}

func TestLoad_CLIFlagWinsOverConfigFile(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "since: develop\n")

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

func TestLoad_EmptyCLIFlagFallsBackToDefault(t *testing.T) {
	fsys := afero.NewMemMapFs()
	flags := newFlagSet()
	if err := flags.Set("since", ""); err != nil {
		t.Fatalf("Set: %v", err)
	}

	cfg, err := config.Load(fsys, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Since != "main" {
		t.Errorf("Since = %q, want default %q", cfg.Since, "main")
	}
}

func TestLoad_EmptyCLIFlagFallsBackToConfigFile(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "since: develop\n")

	flags := newFlagSet()
	if err := flags.Set("since", ""); err != nil {
		t.Fatalf("Set: %v", err)
	}

	cfg, err := config.Load(fsys, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Since != "develop" {
		t.Errorf("Since = %q, want config value %q", cfg.Since, "develop")
	}
}

func TestLoad_EmptyConfigValueFallsBackToDefault(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "since: \"\"\n")

	cfg, err := config.Load(fsys, newFlagSet())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Since != "main" {
		t.Errorf("Since = %q, want default %q", cfg.Since, "main")
	}
}

func TestLoad_InvalidConfigReturnsError(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "changelog: [unterminated\n")

	cfg, err := config.Load(fsys, newFlagSet())
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
	if cfg != nil {
		t.Errorf("expected nil Config on error, got %+v", cfg)
	}
}

func TestLoad_BareNullConfigValueFallsBackToDefault(t *testing.T) {
	fsys := afero.NewMemMapFs()
	writeConfigFile(t, fsys, "since:\n")

	cfg, err := config.Load(fsys, newFlagSet())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Since != "main" {
		t.Errorf("Since = %q, want default %q", cfg.Since, "main")
	}
}
