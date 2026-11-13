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
	"fmt"

	"yamm/internal/app"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var apply, dryRun bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create the default configuration",
		Long:  "Preview the default configuration and its path. Use --apply to create it without overwriting an existing destination.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			path, contents, err := app.InitYamm(path, apply)
			if err != nil {
				return err
			}
			if apply {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Created configuration: %s\n", path)
			} else {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Would create configuration: %s\n\n%s", path, contents)
			}
			return err
		},
	}
	cmd.Flags().BoolVar(&apply, "apply", false, "Create the default configuration")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without creating files or directories")
	cmd.MarkFlagsMutuallyExclusive("dry-run", "apply")
	return cmd
}

func init() {
	rootCmd.AddCommand(newInitCmd())
}
