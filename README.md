# spotiCLI

A Spotify TUI written in Go (Bubble Tea). Playlists on the left, tracks in the
middle, now-playing with album art and an audio visualizer at the bottom.

## Setup

1. Create an app at <https://developer.spotify.com/dashboard>.
2. Add `http://127.0.0.1:8888/callback` to its **Redirect URIs**.
3. `export SPOTIFY_CLIENT_ID=<your client id>` (no client secret needed — PKCE).
4. `go run .` — a browser window opens for the one-time login. The token is
   cached in your user config dir (`~/.config/spotiCLI/token.json`).

Playback needs Spotify Premium and Spotify open on one of your devices.

`go run . --demo` runs the UI on sample data; `go run . --logout` forgets the login.

## Keys

| key | action |
| --- | --- |
| `↑`/`↓`, `j`/`k` | move |
| `tab`, `h`/`l` | switch pane |
| `enter` | open playlist / play track |
| `space` | play / pause |
| `n` / `p` | next / previous |
| `q` | quit |
