// Package spotify is a small client for the parts of the Spotify Web API that
// spotiCLI needs: authentication, playlists, tracks and playback control.
//
// It talks to the API directly instead of using a wrapper library because
// Spotify's February 2026 changes removed endpoints such as
// GET /playlists/{id}/tracks (now /items) that wrappers still call.
package spotify

import "time"

type Playlist struct {
	ID     string
	Name   string
	URI    string
	Tracks int
	Liked  bool // the synthetic "Liked Songs" entry (/me/tracks)
}

type Track struct {
	URI      string
	Title    string
	Artist   string
	Album    string
	Duration time.Duration
	ImageURL string
}

// State is a snapshot of the user's current playback.
type State struct {
	Playing    bool
	Track      Track
	Progress   time.Duration
	Volume     int
	Shuffle    bool
	Repeat     string // "off", "track" or "context"
	DeviceName string
}
