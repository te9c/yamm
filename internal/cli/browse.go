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
package cli

import (
	"errors"
	"fmt"
	"os"

	"yamm/internal/tui"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newBrowseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "browse",
		Short: "Open TUI for browsing your library",
		Long:  "Open the TUI browser. Press q to quit or Ctrl-C to cancel.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.NoArgs(cmd, args); err != nil {
				return &exitError{code: 2, err: err}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
				return &exitError{code: 2, err: errors.New("browse requires terminal input and output")}
			}
			if err := tui.Run(); err != nil {
				if errors.Is(err, tui.ErrCancelled) {
					cmd.SilenceErrors = true
					return &exitError{code: 130, err: err}
				}
				return &exitError{code: 6, err: fmt.Errorf("start terminal UI: %w", err)}
			}
			return nil
		},
	}
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return &exitError{code: 2, err: err}
	})
	return cmd
}

func init() {
	rootCmd.AddCommand(newBrowseCmd())
}
