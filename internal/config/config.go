// Package config resolves the changes CLI configuration by layering
// built-in defaults, an optional .changes/config.yaml, and CLI flags,
// per D010 (CLI flags > config file > built-in defaults).
package config

import (
	"bytes"
	"errors"
	"io/fs"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ConfigFilePath is the location of the optional project config file,
// relative to the working directory.
const ConfigFilePath = ".changes/config.yaml"

// Config is the fully resolved configuration used by every subcommand.
type Config struct {
	Package      string
	Changelog    string
	Since        string
	TagPrefix    string
	Conventional map[string][]string
}

// defaultConventional matches docs/DESIGN.md's documented default
// conventional-commit increment mappings.
func defaultConventional() map[string][]string {
	return map[string][]string{
		"major": {"feat!", "BREAKING CHANGE"},
		"minor": {"feat"},
		"patch": {"fix", "perf", "docs", "style", "refactor", "test", "build", "chore"},
		"none":  {"ci"},
	}
}

// Load resolves a Config by layering built-in defaults, then
// .changes/config.yaml if present, then CLI flag overrides. A missing
// config file is not an error.
func Load(fsys afero.Fs, flags *pflag.FlagSet) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	v.SetDefault("package", ".")
	v.SetDefault("changelog", "CHANGELOG.md")
	v.SetDefault("since", "main")
	v.SetDefault("tag_prefix", "v")
	v.SetDefault("conventional", defaultConventional())

	data, err := afero.ReadFile(fsys, ConfigFilePath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	if err == nil {
		if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
			return nil, err
		}
	}

	if flags != nil {
		if since, err := flags.GetString("since"); err == nil && flags.Changed("since") {
			v.Set("since", since)
		}
	}

	conventional := map[string][]string{}
	for key, val := range v.GetStringMapStringSlice("conventional") {
		conventional[key] = val
	}

	return &Config{
		Package:      v.GetString("package"),
		Changelog:    v.GetString("changelog"),
		Since:        v.GetString("since"),
		TagPrefix:    v.GetString("tag_prefix"),
		Conventional: conventional,
	}, nil
}
