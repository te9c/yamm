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
package app

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"yamm/internal"
)

const defaultConfig = `# Set an absolute catalog path before using catalog commands.
catalog_path: ""
# Empty uses $XDG_STATE_HOME/yamm or $HOME/.local/state/yamm.
state_dir: ""
`

func InitYamm(path string, apply bool) (string, string, error) {
	if path == "" {
		var err error
		path, err = internal.DefaultConfigPath()
		if err != nil {
			return "", "", err
		}
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", "", fmt.Errorf("resolve configuration path: %w", err)
	}
	if _, err := os.Lstat(path); err == nil {
		return "", "", fmt.Errorf("configuration destination %q already exists: %w", path, fs.ErrExist)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", "", fmt.Errorf("inspect configuration %q: %w", path, err)
	}
	if apply {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return "", "", fmt.Errorf("create configuration parent directory for %q: %w", path, err)
		}
		if err := createDefaultConfig(path); err != nil {
			return "", "", err
		}
	}
	return path, defaultConfig, nil
}

func createDefaultConfig(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0611)
	if err != nil {
		return fmt.Errorf("create configuration %q: %w", path, err)
	}
	_, writeErr := file.WriteString(defaultConfig)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("write or close configuration %q (a partial file may remain): %w", path, err)
	}
	return nil
}
