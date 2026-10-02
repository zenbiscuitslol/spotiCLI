package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
)

// nowPlayingView renders the bottom panel: album cover on the left, track
// info + progress + controls + audio visualizer on the right.
func (m Model) nowPlayingView(w, h int) string {
	inner, innerH := w-2, h-2
	t := m.tracks[m.nowIdx]

	// Cover is square in pixels: a cell is ~1:2, so cols = 2*rows. Shrink it
	// if there isn't room left for the info column.
	coverRows := innerH
	coverCols := coverRows * 2
	if maxCols := inner / 2; coverCols > maxCols {
		coverCols = maxCols &^ 1
		coverRows = coverCols / 2
	}
	cover := m.cover.Render(coverCols, coverRows)

	const gap = 2
	rightW := inner - coverCols - gap - 1 // 1 col of right padding
	status := "⏸ Paused"
	if m.playing {
		status = "▶ Playing"
	}
	right := []string{
		styleBold.Render(t.Title),
		styleText.Render(t.Artist),
		styleMuted.Render(t.Album),
		progressBar(m.elapsed, t.Duration, rightW),
		spread(styleAccent.Render(status),
			styleMuted.Render(fmt.Sprintf("⇄ off   ↻ off   vol %d%%", m.volume)), rightW),
		"",
	}
	right = append(right, m.viz.View(rightW, innerH-len(right))...)

	lines := make([]string, innerH)
	for i := range lines {
		left := strings.Repeat(" ", coverCols)
		if i < len(cover) {
			left = cover[i]
		}
		r := ""
		if i < len(right) {
			r = right[i]
		}
		lines[i] = " " + left + strings.Repeat(" ", gap) + ansi.Truncate(r, rightW, "")
	}
	return panel("Now Playing", lines, w, h)
}

// progressBar renders "1:23 ━━━━●────── 3:45" in exactly w cells.
func progressBar(elapsed, total time.Duration, w int) string {
	l, r := fmtDur(elapsed), fmtDur(total)
	barW := w - len(l) - len(r) - 2
	if barW < 4 {
		return styleMuted.Render(l + " / " + r)
	}
	pos := int(float64(barW-1) * float64(elapsed) / float64(total))
	pos = clamp(pos, 0, barW-1)
	bar := styleAccent.Render(strings.Repeat("━", pos)) +
		styleBold.Render("●") +
		styleMuted.Render(strings.Repeat("─", barW-pos-1))
	return styleMuted.Render(l) + " " + bar + " " + styleMuted.Render(r)
}
