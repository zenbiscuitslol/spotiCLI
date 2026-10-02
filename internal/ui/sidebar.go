package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/zenbiscuitslol/spotiCLI/internal/spotify"
)

// sidebarView renders the scrolling playlist list on the left.
func (m Model) sidebarView(w, h int) string {
	inner, rows := w-2, h-2
	focused := m.focus == paneSidebar

	var lines []string
	title := "Playlists"
	switch {
	case m.plErr != "":
		lines = wrap(m.plErr, styleErr, inner)
	case !m.plLoaded:
		lines = []string{styleMuted.Render(" Loading playlists…")}
	default:
		title = fmt.Sprintf("Playlists · %d", len(m.playlists))
		end := min(m.plOffset+rows, len(m.playlists))
		for i := m.plOffset; i < end; i++ {
			lines = append(lines, playlistRow(m.playlists[i], i == m.plCursor, focused, inner))
		}
	}
	return panel(title, lines, w, h, focused)
}

func playlistRow(p spotify.Playlist, cursor, focused bool, w int) string {
	count := fmt.Sprintf("%d ", p.Tracks)
	if !cursor {
		return spread(styleText.Render("  "+p.Name), styleMuted.Render(count), w)
	}

	// Cursor row: highlight bar spanning the full width, including the gap.
	// Dimmer when the sidebar doesn't have focus.
	fg := colGreen
	if !focused {
		fg = colText
	}
	bg := lipgloss.NewStyle().Background(colSelBg)
	cnt := bg.Foreground(colMuted).Render(count)
	name := ansi.Truncate("▌ "+p.Name, w-ansi.StringWidth(count)-1, "…")
	name = bg.Foreground(fg).Bold(true).Render(name)
	gap := max(w-ansi.StringWidth(name)-ansi.StringWidth(cnt), 0)
	return name + bg.Render(strings.Repeat(" ", gap)) + cnt
}
