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
package tui

import (
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// interrupted
var ErrCancelled = errors.New("terminal UI cancelled")

func Run() error {
	app := tview.NewApplication()

	page := tview.NewTextView().
		SetText("Hello tview!").
		SetWrap(false).
		SetScrollable(false)
	app.SetRoot(page, true)

	var cancelled bool
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlC:
			cancelled = true
			app.Stop()
			return nil
		case event.Key() == tcell.KeyRune && event.Rune() == 'q':
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.Run(); err != nil {
		return err
	}
	if cancelled {
		return ErrCancelled
	}
	return nil
}
