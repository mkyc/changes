// Package config resolves the changes CLI configuration by layering
// built-in defaults, an optional .changes/config.yaml, and CLI flags,
// per D010 (CLI flags > config file > built-in defaults).
package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ConfigFilePath is the location of the optional project config file,
// relative to the working directory.
const ConfigFilePath = ".changes/config.yaml"

var allowedConventionalKeys = map[string]struct{}{
	"major": {},
	"minor": {},
	"patch": {},
	"none":  {},
}

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

func mergeConventional(v *viper.Viper, configFileRead bool) (map[string][]string, error) {
	merged := defaultConventional()
	if !configFileRead {
		return merged, nil
	}

	fileConv := v.GetStringMapStringSlice("conventional")
	if len(fileConv) == 0 {
		return nil, errors.New("conventional must contain at least one key")
	}

	for key, val := range fileConv {
		if _, ok := allowedConventionalKeys[key]; !ok {
			return nil, fmt.Errorf("unknown conventional key %q", key)
		}
		merged[key] = val
	}
	return merged, nil
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
	configFileRead := err == nil
	if configFileRead {
		if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
			return nil, err
		}
	}

	if flags != nil {
		if since, err := flags.GetString("since"); err == nil && flags.Changed("since") {
			v.Set("since", since)
		}
	}

	conventional, err := mergeConventional(v, configFileRead)
	if err != nil {
		return nil, err
	}

	return &Config{
		Package:      v.GetString("package"),
		Changelog:    v.GetString("changelog"),
		Since:        v.GetString("since"),
		TagPrefix:    v.GetString("tag_prefix"),
		Conventional: conventional,
	}, nil
}
