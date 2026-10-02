// Package mock holds placeholder data so the UI can be designed before any
// Spotify API calls exist. Replace these with real API results later.
package mock

import (
	"image"
	"image/color"
	"math"
	"time"
)

type Playlist struct {
	Name   string
	Tracks int
}

type Track struct {
	Title    string
	Artist   string
	Album    string
	Duration time.Duration
}

func Playlists() []Playlist {
	return []Playlist{
		{"Liked Songs", 312},
		{"Discover Weekly", 30},
		{"Release Radar", 24},
		{"Daily Mix 1", 50},
		{"Daily Mix 2", 50},
		{"Late Night Drive", 87},
		{"Focus Flow", 120},
		{"Gym Hype", 64},
		{"Chill Lo-fi", 203},
		{"Road Trip", 95},
		{"Throwback 2000s", 140},
		{"Indie Mix", 58},
		{"Jazz Classics", 76},
		{"Study Beats", 160},
		{"Sunday Morning", 41},
		{"Acoustic Evenings", 33},
	}
}

func Tracks() []Track {
	m := func(min, sec int) time.Duration {
		return time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	}
	return []Track{
		{"Neon Horizons", "The Midnight Cartel", "Afterglow", m(3, 42)},
		{"Paper Planes in July", "Wren & the Tides", "Small Hours", m(4, 8)},
		{"Static Bloom", "Halcyon Arcade", "Static Bloom", m(3, 15)},
		{"Run Slow", "Marlowe Park", "Run Slow", m(2, 58)},
		{"Glass Cathedral", "Ember Choir", "Vespers", m(5, 21)},
		{"Honey Overdrive", "Velvet Static", "Honey Overdrive", m(3, 33)},
		{"Low Tide Letters", "Saoirse Vale", "Low Tide", m(4, 47)},
		{"Concrete Garden", "North of Nowhere", "Concrete Garden", m(3, 5)},
		{"Satellites", "Juno Lark", "Orbit", m(3, 51)},
		{"Fever Dream Motel", "The Paper Kites Club", "Vacancy", m(4, 12)},
		{"Slow Burn", "Marlowe Park", "Run Slow", m(3, 27)},
		{"Moonlit Arcade", "Halcyon Arcade", "Static Bloom", m(2, 49)},
		{"Wildflower Radio", "Wren & the Tides", "Small Hours", m(3, 38)},
		{"Telescope", "Juno Lark", "Orbit", m(4, 2)},
		{"Golden Hour Static", "Velvet Static", "Honey Overdrive", m(3, 19)},
		{"Undertow", "Ember Choir", "Vespers", m(5, 44)},
		{"City of Lanterns", "North of Nowhere", "Concrete Garden", m(3, 56)},
		{"Last Train Home", "Saoirse Vale", "Low Tide", m(4, 30)},
		{"Daydream Parade", "The Midnight Cartel", "Afterglow", m(3, 11)},
		{"Soft Machines", "Halcyon Arcade", "Static Bloom", m(4, 6)},
	}
}

// Cover returns a procedurally generated square "album cover" so the cover
// renderer has something to draw without a network request.
func Cover() image.Image {
	const size = 256
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	lerp := func(a, b, t float64) float64 { return a + (b-a)*t }
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			fx, fy := float64(x)/size, float64(y)/size
			t := (fx*0.35 + fy*0.65)
			// Dusk gradient: deep indigo -> magenta -> amber.
			var r, g, b float64
			if t < 0.55 {
				u := t / 0.55
				r, g, b = lerp(28, 190, u), lerp(18, 52, u), lerp(84, 130, u)
			} else {
				u := (t - 0.55) / 0.45
				r, g, b = lerp(190, 255, u), lerp(52, 168, u), lerp(130, 70, u)
			}
			// Sun disc.
			d := math.Hypot(fx-0.5, fy-0.42)
			if d < 0.2 {
				r, g, b = 255, lerp(214, 150, d/0.2), lerp(120, 70, d/0.2)
			}
			// Horizon stripes cutting into the lower sun / ground.
			if fy > 0.55 && fy < 0.75 && d < 0.2 {
				if int((fy-0.55)*60)%2 == 0 {
					r, g, b = 28, 18, 84
				}
			}
			if fy >= 0.75 {
				r, g, b = lerp(20, 40, fy), lerp(12, 20, fy), lerp(60, 90, fy)
			}
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return img
}
