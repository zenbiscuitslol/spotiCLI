// Package ui contains the Bubble Tea model and all rendering for spotiCLI.
package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zenbiscuitslol/spotiCLI/internal/mock"
)

const (
	frameRate = 20
	tickEvery = time.Second / frameRate

	minW, minH = 60, 19
)

type tickMsg time.Time

type Model struct {
	w, h int

	// Placeholder state; the API layer will populate these later.
	playlists []mock.Playlist
	tracks    []mock.Track
	selected  int // highlighted playlist
	nowIdx    int // playing track
	elapsed   time.Duration
	playing   bool
	volume    int

	viz   *visualizer
	cover *coverRenderer
}

func New() Model {
	return Model{
		playlists: mock.Playlists(),
		tracks:    mock.Tracks(),
		nowIdx:    2,
		elapsed:   71 * time.Second,
		playing:   true,
		volume:    70,
		viz:       newVisualizer(),
		cover:     newCoverRenderer(mock.Cover()),
	}
}

func tick() tea.Cmd {
	return tea.Tick(tickEvery, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) Init() tea.Cmd { return tick() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyMsg:
		// No navigation keybinds yet; only quit.
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		m.viz.Step(tickEvery.Seconds())
		if m.playing {
			m.elapsed += tickEvery
			if m.elapsed >= m.tracks[m.nowIdx].Duration {
				m.elapsed = 0
			}
		}
		return m, tick()
	}
	return m, nil
}

func (m Model) View() string {
	if m.w == 0 {
		return ""
	}
	if m.w < minW || m.h < minH {
		return lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center,
			styleMuted.Render("Terminal too small — need at least 60x19"))
	}

	bodyH := m.h - 1 // 1 line for the footer
	sideW := clamp(m.w/4, 24, 34)
	mainW := m.w - sideW

	npInner := clamp(bodyH-10, 8, 12)
	npH := npInner + 2
	tracksH := bodyH - npH

	main := lipgloss.JoinVertical(lipgloss.Left,
		m.tracksView(mainW, tracksH),
		m.nowPlayingView(mainW, npH),
	)
	body := lipgloss.JoinHorizontal(lipgloss.Top, m.sidebarView(sideW, bodyH), main)

	footer := styleMuted.Render(" spotiCLI  ·  q quit")
	return body + "\n" + footer
}
