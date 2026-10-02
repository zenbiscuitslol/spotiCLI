package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func fmtDur(d time.Duration) string {
	s := int(d.Seconds())
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// tracksView renders the track table for the selected playlist.
func (m Model) tracksView(w, h int) string {
	inner := w - 2

	// Column widths: fixed #, time; the rest split between title/artist/album.
	const numW, timeW = 4, 6
	rest := max(inner-numW-timeW-4, 6) // 4 single-space gutters
	titleW := rest * 40 / 100
	artistW := rest * 30 / 100
	albumW := rest - titleW - artistW

	row := func(num, title, artist, album, dur string, st lipgloss.Style) string {
		cells := []string{
			fit(num, numW), fit(title, titleW), fit(artist, artistW), fit(album, albumW),
		}
		return st.Render(strings.Join(cells, " ")) + " " + st.Render(fmt.Sprintf("%*s", timeW, dur))
	}

	head := styleMuted.Underline(true)
	lines := []string{row("#", "Title", "Artist", "Album", "Time", head)}
	for i, t := range m.tracks {
		st, num := styleText, fmt.Sprintf("%d", i+1)
		if i == m.nowIdx {
			st, num = styleAccent.Bold(true), "▶"
		}
		lines = append(lines, row(num, t.Title, t.Artist, t.Album, fmtDur(t.Duration), st))
	}

	pl := m.playlists[m.selected]
	title := ansi.Truncate(pl.Name, max(inner-20, 8), "…")
	return panel(fmt.Sprintf("%s · %d tracks", title, pl.Tracks), lines, w, h)
}
