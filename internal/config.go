/*
Copyright © 2026 Roman Dordzheev <roman.te9c@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

type cfg struct {
	configPath  string
	catalogPath string
	stateDir    string
}

type configFile struct {
	CatalogPath string `yaml:"catalog_path"`
	StateDir    string `yaml:"state_dir"`
}

func DefaultConfigPath() (string, error) {
	dir, err := defaultDir("XDG_CONFIG_HOME", ".config")
	if err != nil {
		return "", fmt.Errorf("resolve default configuration path: %w", err)
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func defaultStateDir() (string, error) {
	return defaultDir("XDG_STATE_HOME", filepath.Join(".local", "state"))
}

// fallback is relative to the home directory
func defaultDir(variable, fallback string) (string, error) {
	base := os.Getenv(variable)
	if base != "" {
		if err := validatePath(variable, base); err != nil {
			return "", err
		}
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory for %s fallback: %w", variable, err)
		}
		if err := validatePath("home directory for "+variable+" fallback", home); err != nil {
			return "", err
		}
		base = filepath.Join(home, fallback)
	}
	return filepath.Join(base, "yamm"), nil
}

func validatePath(field, path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("%s must be an absolute path: %q", field, path)
	}
	return nil
}

func LoadConfig(path string) (*cfg, error) {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}
	configPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve configuration path %q: %w", path, err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read configuration %q: %w", configPath, err)
	}
	file, err := parseConfig(data)
	if err != nil {
		return nil, fmt.Errorf("parse configuration %q: %w", configPath, err)
	}
	if file.CatalogPath != "" {
		if err := validatePath("catalog_path", file.CatalogPath); err != nil {
			return nil, fmt.Errorf("validate configuration %q: %w", configPath, err)
		}
	}
	if file.StateDir == "" {
		file.StateDir, err = defaultStateDir()
		if err != nil {
			return nil, fmt.Errorf("resolve state_dir for configuration %q: %w", configPath, err)
		}
	}
	if err := validatePath("state_dir", file.StateDir); err != nil {
		return nil, fmt.Errorf("validate configuration %q: %w", configPath, err)
	}
	return &cfg{configPath: configPath, catalogPath: file.CatalogPath, stateDir: file.StateDir}, nil
}

func (c *cfg) ConfigPath() string { return c.configPath }

func (c *cfg) CatalogPath() string { return c.catalogPath }

func (c *cfg) StateDir() string { return c.stateDir }

func parseConfig(data []byte) (configFile, error) {
	var document yaml.Node
	if err := yaml.Load(data, &document); err != nil {
		return configFile{}, err
	}
	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 {
		return configFile{}, fmt.Errorf("configuration must contain one non-null mapping document")
	}
	root := document.Content[0]
	if root.Kind != yaml.MappingNode || root.Tag != "!!map" || root.Anchor != "" {
		return configFile{}, configNodeError(root, "configuration must be a mapping without anchors or custom tags")
	}

	var file configFile
	seen := make(map[string]bool, 2)
	for i := 0; i < len(root.Content); i += 2 {
		key, value := root.Content[i], root.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" || key.Anchor != "" {
			return configFile{}, configNodeError(key, "configuration keys must be strings without anchors or custom tags")
		}
		if seen[key.Value] {
			return configFile{}, configNodeError(key, "duplicate key %q", key.Value)
		}
		seen[key.Value] = true
		var field *string
		switch key.Value {
		case "catalog_path":
			field = &file.CatalogPath
		case "state_dir":
			field = &file.StateDir
		default:
			return configFile{}, configNodeError(key, "unknown key %q", key.Value)
		}
		if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || value.Anchor != "" {
			return configFile{}, configNodeError(value, "%s must be a string without anchors, aliases, or custom tags", key.Value)
		}
		*field = value.Value
	}
	return file, nil
}

func configNodeError(node *yaml.Node, format string, args ...any) error {
	return fmt.Errorf("line %d, column %d: %s", node.Line, node.Column, fmt.Sprintf(format, args...))
}
