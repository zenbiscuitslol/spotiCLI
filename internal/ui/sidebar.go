package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/zenbiscuitslol/spotiCLI/internal/mock"
)

// sidebarView renders the playlist list that sits on the left of the screen.
// Items that don't fit are clipped (scrolling comes with navigation later).
func (m Model) sidebarView(w, h int) string {
	inner := w - 2
	var lines []string
	for i, p := range m.playlists {
		lines = append(lines, playlistRow(p, i == m.selected, inner))
	}
	title := fmt.Sprintf("Playlists · %d", len(m.playlists))
	return panel(title, lines, w, h)
}

func playlistRow(p mock.Playlist, selected bool, w int) string {
	count := fmt.Sprintf("%d ", p.Tracks)
	if !selected {
		return spread(styleText.Render("  "+p.Name), styleMuted.Render(count), w)
	}

	// Selected row: highlight bar spanning the full width, including the gap.
	bg := lipgloss.NewStyle().Background(colSelBg)
	cnt := bg.Foreground(colMuted).Render(count)
	name := ansi.Truncate("▌ "+p.Name, w-ansi.StringWidth(count)-1, "…")
	name = bg.Foreground(colGreen).Bold(true).Render(name)
	gap := max(w-ansi.StringWidth(name)-ansi.StringWidth(cnt), 0)
	return name + bg.Render(strings.Repeat(" ", gap)) + cnt
}
