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
	"go.yaml.in/yaml/v3"
)

const (
	// ConfigFilePath is the location of the optional project config file,
	// relative to the working directory.
	ConfigFilePath = ".changes/config.yaml"
	defaultSince   = "main"
)

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

// fileConventionalFromYAML parses the conventional block directly instead of
// reading it from v, because Viper's typed getters can't distinguish an
// absent key from one that's present but empty — exactly the distinction
// the empty-block validation below needs.
func fileConventionalFromYAML(data []byte) (map[string][]string, bool, error) {
	var doc struct {
		Conventional yaml.Node `yaml:"conventional"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, false, err
	}
	if doc.Conventional.Kind == 0 {
		return nil, false, nil
	}
	resolved := &doc.Conventional
	for resolved.Kind == yaml.AliasNode && resolved.Alias != nil {
		resolved = resolved.Alias
	}
	if resolved.Kind != yaml.MappingNode && resolved.Tag != "!!null" {
		return nil, true, errors.New("conventional must be a mapping of keys to lists")
	}

	fileConv := make(map[string][]string, len(resolved.Content)/2)
	for i := 0; i < len(resolved.Content); i += 2 {
		key := resolved.Content[i].Value
		var val []string
		if err := resolved.Content[i+1].Decode(&val); err != nil {
			return nil, true, fmt.Errorf("conventional.%s must be a list of strings", key)
		}
		fileConv[key] = val
	}
	if len(fileConv) == 0 {
		return nil, true, errors.New("conventional must contain at least one key")
	}
	return fileConv, true, nil
}

func mergeConventional(data []byte, configFileRead bool) (map[string][]string, error) {
	merged := defaultConventional()
	if !configFileRead {
		return merged, nil
	}

	fileConv, present, err := fileConventionalFromYAML(data)
	if err != nil {
		return nil, err
	}
	if !present {
		return merged, nil
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
	v.SetDefault("since", defaultSince)
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
		if since, err := flags.GetString("since"); err == nil && flags.Changed("since") && since != "" {
			v.Set("since", since)
		}
	}
	since := v.GetString("since")
	if since == "" {
		since = defaultSince
	}

	conventional, err := mergeConventional(data, configFileRead)
	if err != nil {
		return nil, err
	}

	return &Config{
		Package:      v.GetString("package"),
		Changelog:    v.GetString("changelog"),
		Since:        since,
		TagPrefix:    v.GetString("tag_prefix"),
		Conventional: conventional,
	}, nil
}
