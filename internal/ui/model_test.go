package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/zenbiscuitslol/spotiCLI/internal/mock"
)

// settle runs a command and feeds its resulting messages back into the model
// until nothing is left. Animation ticks are dropped so it terminates.
func settle(tm tea.Model, cmd tea.Cmd) tea.Model {
	if cmd == nil {
		return tm
	}
	switch msg := cmd().(type) {
	case tickMsg:
		return tm
	case tea.BatchMsg:
		for _, c := range msg {
			tm = settle(tm, c)
		}
		return tm
	default:
		tm, next := tm.Update(msg)
		return settle(tm, next)
	}
}

func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case " ":
		return tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func press(tm tea.Model, keys ...string) tea.Model {
	for _, k := range keys {
		var cmd tea.Cmd
		tm, cmd = tm.Update(key(k))
		tm = settle(tm, cmd)
	}
	return tm
}

func init() { refreshDelay = 0 }

// boot returns a model of the given size that has finished its initial load.
func boot(w, h int) tea.Model {
	m := New(mock.New())
	var tm tea.Model = m
	tm, _ = tm.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return settle(tm, m.Init())
}

func TestInitialLoadShowsPlaylistsAndTracks(t *testing.T) {
	v := boot(110, 32).View()
	for _, want := range []string{"Playlists · 16", "Liked Songs", "Discover Weekly", "Neon Horizons", "Nothing playing"} {
		if !strings.Contains(v, want) {
			t.Errorf("view is missing %q", want)
		}
	}
}

// Every rendered frame must fit the terminal, otherwise it scrolls or tears.
func TestViewFitsTerminal(t *testing.T) {
	for _, sz := range [][2]int{{60, 19}, {80, 24}, {120, 40}, {200, 50}} {
		tm := boot(sz[0], sz[1])
		tm = press(tm, "tab", "enter", "down", "down") // start playing
		for i := 0; i < 30; i++ {
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
	var tm tea.Model = New(mock.New())
	tm, _ = tm.Update(tea.WindowSizeMsg{Width: 30, Height: 10})
	if !strings.Contains(tm.View(), "too small") {
		t.Error("expected too-small message")
	}
}

func TestPlayPauseNextPrevious(t *testing.T) {
	tm := boot(110, 32)
	tm = press(tm, "tab", "down", "down", "enter") // third track
	v := tm.View()
	if !strings.Contains(v, "▶ Playing") || !strings.Contains(v, "Static Bloom") {
		t.Fatalf("expected Static Bloom playing:\n%s", v)
	}

	tm = press(tm, " ")
	if !strings.Contains(tm.View(), "Paused") {
		t.Error("space should pause")
	}
	tm = press(tm, " ")
	if !strings.Contains(tm.View(), "▶ Playing") {
		t.Error("space should resume")
	}

	tm = press(tm, "n")
	if !strings.Contains(tm.View(), "Run Slow") {
		t.Error("n should skip to the next track")
	}
	tm = press(tm, "p", "p")
	if !strings.Contains(tm.View(), "Paper Planes in July") {
		t.Error("p should go back")
	}
}

func TestOpenOtherPlaylistReloadsTracks(t *testing.T) {
	tm := boot(110, 32)
	tm = press(tm, "j", "enter")
	v := tm.View()
	if !strings.Contains(v, "Discover Weekly · 30 tracks") {
		t.Errorf("opened playlist title missing:\n%s", v)
	}
}

func TestSpaceWithNothingPlayingExplains(t *testing.T) {
	tm := press(boot(110, 32), " ")
	if !strings.Contains(tm.View(), "Nothing playing — pick a track") {
		t.Error("expected hint in footer")
	}
}

func TestSidebarScrollsToKeepCursorVisible(t *testing.T) {
	tm := boot(80, 20) // too short to show all 16 playlists
	for i := 0; i < 15; i++ {
		tm = press(tm, "j")
	}
	if v := tm.View(); !strings.Contains(v, "Acoustic Even") {
		t.Errorf("last playlist should have scrolled into view:\n%s", v)
	}
}
