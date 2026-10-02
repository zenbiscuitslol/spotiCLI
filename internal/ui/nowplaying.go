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

	// Cover is square in pixels: a cell is ~1:2, so cols = 2*rows. Shrink it
	// if there isn't room left for the info column.
	coverRows := innerH
	coverCols := coverRows * 2
	if maxCols := inner / 2; coverCols > maxCols {
		coverCols = maxCols &^ 1
		coverRows = coverCols / 2
	}
	cover := m.cover.Render(coverCols, coverRows)
	if cover == nil { // no art yet: dim placeholder block
		ph := styleMuted.Render(strings.Repeat("░", coverCols))
		for i := 0; i < coverRows; i++ {
			cover = append(cover, ph)
		}
	}

	const gap = 2
	rightW := inner - coverCols - gap - 1 // 1 col of right padding

	var right []string
	if s := m.state; s != nil {
		status := "❚❚ Paused"
		if s.Playing {
			status = "▶ Playing"
		}
		shuffle, repeat := "off", s.Repeat
		if s.Shuffle {
			shuffle = "on"
		}
		if repeat == "" {
			repeat = "off"
		}
		right = []string{
			styleBold.Render(s.Track.Title),
			styleText.Render(s.Track.Artist),
			styleMuted.Render(s.Track.Album),
			progressBar(m.elapsed, s.Track.Duration, rightW),
			spread(styleAccent.Render(status),
				styleMuted.Render(fmt.Sprintf("⇄ %s  ↻ %s  vol %d%%", shuffle, repeat, s.Volume)), rightW),
			styleMuted.Render(s.DeviceName),
		}
	} else {
		right = []string{
			styleBold.Render("Nothing playing"),
			styleMuted.Render("Open a playlist, move to a track and press enter."),
			styleMuted.Render("Spotify must be open on one of your devices."),
			"", "", "",
		}
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
	return panel("Now Playing", lines, w, h, false)
}

// progressBar renders "1:23 ━━━━●────── 3:45" in exactly w cells.
func progressBar(elapsed, total time.Duration, w int) string {
	elapsed = max(min(elapsed, total), 0)
	l, r := fmtDur(elapsed), fmtDur(total)
	barW := w - len(l) - len(r) - 2
	if barW < 4 || total <= 0 {
		return styleMuted.Render(l + " / " + r)
	}
	pos := clamp(int(float64(barW-1)*float64(elapsed)/float64(total)), 0, barW-1)
	bar := styleAccent.Render(strings.Repeat("━", pos)) +
		styleBold.Render("●") +
		styleMuted.Render(strings.Repeat("─", barW-pos-1))
	return styleMuted.Render(l) + " " + bar + " " + styleMuted.Render(r)
}
