package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Palette, loosely based on Spotify's brand colors.
var (
	colGreen  = lipgloss.Color("#1DB954")
	colText   = lipgloss.Color("#E6E6E6")
	colMuted  = lipgloss.Color("#7A7A7A")
	colBorder = lipgloss.Color("#3A3A3A")
	colSelBg  = lipgloss.Color("#262626")

	// Visualizer gradient: low -> mid -> high amplitude.
	colVizLow  = "#1DB954"
	colVizMid  = "#E8C547"
	colVizHigh = "#E8475F"
)

var (
	styleText   = lipgloss.NewStyle().Foreground(colText)
	styleMuted  = lipgloss.NewStyle().Foreground(colMuted)
	styleAccent = lipgloss.NewStyle().Foreground(colGreen)
	styleBold   = lipgloss.NewStyle().Foreground(colText).Bold(true)
)

// fit truncates s to w cells and pads it with spaces so it is exactly w wide.
func fit(s string, w int) string {
	if w <= 0 {
		return ""
	}
	s = ansi.Truncate(s, w, "…")
	if pad := w - ansi.StringWidth(s); pad > 0 {
		s += strings.Repeat(" ", pad)
	}
	return s
}

// block normalizes lines to exactly w x h (clipping or padding as needed).
func block(lines []string, w, h int) []string {
	out := make([]string, h)
	for i := range out {
		if i < len(lines) {
			out[i] = fit(lines[i], w)
		} else {
			out[i] = strings.Repeat(" ", w)
		}
	}
	return out
}

// spread places left and right on one line of width w, right-aligned.
func spread(left, right string, w int) string {
	gap := w - ansi.StringWidth(left) - ansi.StringWidth(right)
	if gap < 1 {
		return ansi.Truncate(left, w-ansi.StringWidth(right)-1, "…") + " " + right
	}
	return left + strings.Repeat(" ", gap) + right
}

// panel draws a rounded box of exactly w x h with a title in the top border.
func panel(title string, body []string, w, h int) string {
	if w < 4 || h < 2 {
		return ""
	}
	border := lipgloss.NewStyle().Foreground(colBorder)
	titleSt := lipgloss.NewStyle().Foreground(colGreen).Bold(true)

	title = ansi.Truncate(title, w-6, "…")
	fill := w - 2 - 1 - (ansi.StringWidth(title) + 2)
	top := border.Render("╭─") + titleSt.Render(" "+title+" ") +
		border.Render(strings.Repeat("─", max(fill, 0))+"╮")

	var sb strings.Builder
	sb.WriteString(top + "\n")
	side := border.Render("│")
	for _, l := range block(body, w-2, h-2) {
		sb.WriteString(side + l + side + "\n")
	}
	sb.WriteString(border.Render("╰" + strings.Repeat("─", w-2) + "╯"))
	return sb.String()
}

func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }
