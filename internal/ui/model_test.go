package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

// Every rendered frame must be exactly w x h, otherwise the terminal scrolls
// or the layout tears.
func TestViewFillsTerminal(t *testing.T) {
	sizes := [][2]int{{60, 19}, {80, 24}, {120, 40}, {200, 50}, {100, 30}}
	for _, sz := range sizes {
		m := New()
		var tm tea.Model = m
		tm, _ = tm.Update(tea.WindowSizeMsg{Width: sz[0], Height: sz[1]})
		for i := 0; i < 30; i++ { // let the visualizer run a bit
			tm, _ = tm.Update(tickMsg{})
		}
		lines := strings.Split(tm.View(), "\n")
		if len(lines) != sz[1] {
			t.Errorf("%dx%d: got %d lines", sz[0], sz[1], len(lines))
		}
		for i, l := range lines {
			if w := ansi.StringWidth(l); w > sz[0] {
				t.Errorf("%dx%d: line %d is %d wide", sz[0], sz[1], i, w)
			}
		}
	}
}

func TestTooSmall(t *testing.T) {
	var tm tea.Model = New()
	tm, _ = tm.Update(tea.WindowSizeMsg{Width: 30, Height: 10})
	if !strings.Contains(tm.View(), "too small") {
		t.Error("expected too-small message")
	}
}
