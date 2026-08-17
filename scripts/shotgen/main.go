// Command shotgen prints a view's rendered content, for making screenshots.
//
// It renders through the real models rather than screen-scraping a terminal.
// A pty capture is a stream of repaints and cursor moves, and reconstructing
// the finished screen from it needs a terminal emulator; the models already
// know what the screen says.
//
//	go run ./scripts/shotgen <board-path> <view> <cols> <rows>
//
// Views: home, board, pipeline, signoff, edit.
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/BeMuCa/jaira/internal/tui"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "usage: shotgen <board-path> <view> <cols> <rows>")
		os.Exit(2)
	}
	dir, view := os.Args[1], os.Args[2]
	cols, _ := strconv.Atoi(os.Args[3])
	rows, _ := strconv.Atoi(os.Args[4])
	size := tea.WindowSizeMsg{Width: cols, Height: rows}

	if view == "home" {
		h, err := tui.NewHome(nil)
		if err != nil {
			die(err)
		}
		h.Update(size)
		fmt.Print(h.View().Content)
		return
	}

	s, err := ticket.Discover(dir)
	if err != nil {
		die(err)
	}
	m, err := tui.New(s)
	if err != nil {
		die(err)
	}
	m.Update(size)

	key := func(r rune) { m.Update(tea.KeyPressMsg{Code: r, Text: string(r)}) }
	switch view {
	case "board":
	case "pipeline":
		key('v')
	case "signoff":
		// Walk right until opening the ticket lands on the sign-off screen,
		// rather than counting lanes: a hardcoded count breaks the moment a
		// lane is added to the default board.
		for i := 0; i < 16; i++ {
			m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			if strings.Contains(m.View().Content, "sign-off") {
				break
			}
			m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
			key('l')
		}
	case "edit":
		m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		key('e')
	default:
		die(fmt.Errorf("unknown view %q", view))
	}
	fmt.Print(m.View().Content)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "shotgen:", err)
	os.Exit(1)
}
