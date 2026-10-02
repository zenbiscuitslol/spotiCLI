package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type recorded struct {
	Method, Path, Query, Body string
}

// fake starts a server answering by "METHOD /path"; it records every request.
func fake(t *testing.T, routes map[string]func(w http.ResponseWriter, r *http.Request)) (*Client, *[]recorded) {
	t.Helper()
	var mu sync.Mutex
	var log []recorded
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		log = append(log, recorded{r.Method, r.URL.Path, r.URL.RawQuery, string(b)})
		mu.Unlock()
		h, ok := routes[r.Method+" "+r.URL.Path]
		if !ok {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	return &Client{http: srv.Client(), base: srv.URL}, &log
}

func js(body string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}
}

func apiErr(status int, reason string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintf(w, `{"error":{"status":%d,"message":"nope","reason":%q}}`, status, reason)
	}
}

func TestPlaylistsAcceptsRenamedItemsField(t *testing.T) {
	c, _ := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"GET /me/playlists": js(`{"items":[
			{"id":"a","name":"New","uri":"spotify:playlist:a","items":{"total":7}},
			{"id":"b","name":"Old","uri":"spotify:playlist:b","tracks":{"total":3}}],"next":null}`),
		"GET /me/tracks": js(`{"total":42,"items":[]}`),
	})
	got, err := c.Playlists(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || !got[0].Liked || got[0].Tracks != 42 {
		t.Fatalf("liked songs should come first: %+v", got)
	}
	if got[1].Tracks != 7 || got[2].Tracks != 3 {
		t.Errorf("track counts wrong: %+v", got)
	}
}

func TestPlaylistsFollowsPagination(t *testing.T) {
	var base string
	c, _ := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"GET /me/playlists": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("offset") == "50" {
				js(`{"items":[{"id":"2","name":"Two"}],"next":null}`)(w, r)
				return
			}
			js(fmt.Sprintf(`{"items":[{"id":"1","name":"One"}],"next":"%s/me/playlists?offset=50"}`, base))(w, r)
		},
		"GET /me/tracks": js(`{"total":0}`),
	})
	base = c.base
	got, err := c.Playlists(context.Background())
	if err != nil || len(got) != 3 {
		t.Fatalf("want liked + 2 playlists, got %+v, %v", got, err)
	}
}

func TestPlaylistTracksUsesItemsEndpoint(t *testing.T) {
	c, log := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"GET /playlists/abc/items": js(`{"items":[
			{"item":{"uri":"spotify:track:1","name":"Song","type":"track","duration_ms":61000,
			  "artists":[{"name":"A"},{"name":"B"}],
			  "album":{"name":"Alb","images":[{"url":"big","width":640},{"url":"mid","width":300},{"url":"small","width":64}]}}},
			{"track":{"uri":"spotify:track:2","name":"OldShape","type":"track","artists":[],"album":{}}},
			{"item":null},
			{"item":{"uri":"spotify:local:x","name":"Local","is_local":true}},
			{"item":{"uri":"spotify:episode:9","name":"Pod","type":"episode"}}
		],"next":null}`),
	})
	got, err := c.Tracks(context.Background(), Playlist{ID: "abc"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 playable tracks, got %+v", got)
	}
	if got[0].Artist != "A, B" || got[0].Duration.Seconds() != 61 || got[0].ImageURL != "mid" {
		t.Errorf("track decoded wrong: %+v", got[0])
	}
	if (*log)[0].Path != "/playlists/abc/items" {
		t.Errorf("hit removed endpoint: %s", (*log)[0].Path)
	}
}

func TestPlaylistTracksForbiddenIsExplained(t *testing.T) {
	c, _ := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"GET /playlists/x/items": apiErr(403, ""),
	})
	_, err := c.Tracks(context.Background(), Playlist{ID: "x"})
	if err == nil || !strings.Contains(err.Error(), "own or collaborate") {
		t.Fatalf("want friendly 403 message, got %v", err)
	}
}

func TestPlaybackIdleIs204(t *testing.T) {
	c, _ := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"GET /me/player": func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(204) },
	})
	st, err := c.Playback(context.Background())
	if st != nil || err != nil {
		t.Fatalf("want nil,nil got %v,%v", st, err)
	}
}

func TestPlaybackDecodes(t *testing.T) {
	c, _ := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"GET /me/player": js(`{"is_playing":true,"progress_ms":5000,"shuffle_state":true,"repeat_state":"context",
			"device":{"name":"Laptop","volume_percent":55},
			"item":{"uri":"spotify:track:1","name":"S","type":"track","duration_ms":9000,"artists":[{"name":"A"}],"album":{"name":"Al"}}}`),
	})
	st, err := c.Playback(context.Background())
	if err != nil || st == nil {
		t.Fatal(st, err)
	}
	if !st.Playing || st.Progress.Seconds() != 5 || st.Volume != 55 || !st.Shuffle || st.Repeat != "context" || st.DeviceName != "Laptop" {
		t.Errorf("%+v", st)
	}
}

func TestPlayBodies(t *testing.T) {
	tracks := []Track{{URI: "u0"}, {URI: "u1"}, {URI: "u2"}}

	c, log := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"PUT /me/player/play": func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(204) },
	})
	if err := c.Play(context.Background(), Playlist{URI: "spotify:playlist:p"}, tracks, 1); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	json.Unmarshal([]byte((*log)[0].Body), &body)
	if body["context_uri"] != "spotify:playlist:p" || body["offset"].(map[string]any)["uri"] != "u1" {
		t.Errorf("playlist body: %s", (*log)[0].Body)
	}

	if err := c.Play(context.Background(), Playlist{Liked: true}, tracks, 1); err != nil {
		t.Fatal(err)
	}
	json.Unmarshal([]byte((*log)[1].Body), &body)
	uris := body["uris"].([]any)
	if len(uris) != 2 || uris[0] != "u1" {
		t.Errorf("liked body should queue from the selected track: %s", (*log)[1].Body)
	}
}

func TestCommandRetriesOnFirstDeviceWhenNoneActive(t *testing.T) {
	calls := 0
	c, log := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"POST /me/player/next": func(w http.ResponseWriter, r *http.Request) {
			calls++
			if r.URL.Query().Get("device_id") == "" {
				apiErr(404, "NO_ACTIVE_DEVICE")(w, r)
				return
			}
			w.WriteHeader(204)
		},
		"GET /me/player/devices": js(`{"devices":[{"id":"locked","is_restricted":true},{"id":"dev1"},{"id":"dev2"}]}`),
	})
	if err := c.Next(context.Background()); err != nil {
		t.Fatal(err)
	}
	if calls != 2 || !strings.Contains((*log)[len(*log)-1].Query, "device_id=dev1") {
		t.Errorf("want retry on dev1, log=%+v", *log)
	}
}

func TestNoDevicesGivesHelpfulError(t *testing.T) {
	c, _ := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"PUT /me/player/pause":   apiErr(404, "NO_ACTIVE_DEVICE"),
		"GET /me/player/devices": js(`{"devices":[]}`),
	})
	err := c.Pause(context.Background())
	if err == nil || !strings.Contains(err.Error(), "open Spotify") {
		t.Fatalf("got %v", err)
	}
}

func TestRateLimitMessage(t *testing.T) {
	c, _ := fake(t, map[string]func(http.ResponseWriter, *http.Request){
		"GET /me/player": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Retry-After", "7")
			w.WriteHeader(429)
		},
	})
	_, err := c.Playback(context.Background())
	if err == nil || !strings.Contains(err.Error(), "retry in 7s") {
		t.Fatalf("got %v", err)
	}
}
