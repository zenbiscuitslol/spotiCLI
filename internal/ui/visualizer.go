package ui

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// numBands matches spotify-player: 64 log-scale frequency bands, bass on the
// left and treble on the right.
const numBands = 64

var eighths = []rune(" ▁▂▃▄▅▆▇█")

// visualizer is a bar-chart spectrum display. For now the band levels are
// simulated in Step; later they should be fed by a real FFT of the audio.
type visualizer struct {
	bands [numBands]float64 // 0..1
	t     float64
	rng   *rand.Rand
}

func newVisualizer() *visualizer {
	return &visualizer{rng: rand.New(rand.NewSource(7))}
}

// Step advances the simulation by dt seconds. While inactive (paused or
// idle) the bars fall back to zero.
func (v *visualizer) Step(dt float64, active bool) {
	v.t += dt
	// ~120 BPM kick drum envelope.
	beat := math.Pow(math.Max(0, math.Sin(v.t*2*math.Pi*2)), 6)
	for i := range v.bands {
		p := float64(i) / (numBands - 1)
		base := 0.12 + 0.55*math.Pow(1-p, 1.6)
		wob := 0.5 + 0.5*math.Sin(v.t*(1.5+p*6)+float64(i)*0.9)
		kick := beat * 0.5 * math.Pow(1-p, 2.5)
		jitter := v.rng.Float64() * 0.1 * (1 - p*0.5)
		target := math.Min(base*wob+kick+jitter, 1)
		if !active {
			target = 0
		}
		if target > v.bands[i] {
			v.bands[i] += (target - v.bands[i]) * 0.6 // fast attack
		} else {
			v.bands[i] = v.bands[i]*0.82 + target*0.18 // slow decay
		}
	}
}

// View renders the spectrum into exactly w x h cells.
func (v *visualizer) View(w, h int) []string {
	if w < 2 || h < 1 {
		return nil
	}
	n := min(numBands, w/2)
	slot := w / n
	barW := max(slot-1, 1)
	lead := (w - n*slot) / 2

	// Resample the 64 bands down to n bars.
	bars := make([]float64, n)
	for j := range bars {
		lo, hi := j*numBands/n, (j+1)*numBands/n
		hi = max(hi, lo+1)
		var sum float64
		for _, b := range v.bands[lo:hi] {
			sum += b
		}
		bars[j] = sum / float64(hi-lo)
	}

	// One color per row, from the bottom (low) to the top (high).
	rowSt := make([]lipgloss.Style, h)
	for r := 0; r < h; r++ {
		frac := (float64(h-1-r) + 0.5) / float64(h)
		rowSt[r] = lipgloss.NewStyle().Foreground(lipgloss.Color(gradient(frac)))
	}

	out := make([]string, h)
	for r := 0; r < h; r++ {
		var sb strings.Builder
		sb.WriteString(strings.Repeat(" ", lead))
		for _, lvl := range bars {
			fill := clamp(int(lvl*float64(h*8))-(h-1-r)*8, 0, 8)
			cell := strings.Repeat(string(eighths[fill]), barW)
			if fill > 0 {
				cell = rowSt[r].Render(cell)
			}
			sb.WriteString(cell + strings.Repeat(" ", slot-barW))
		}
		out[r] = sb.String()
	}
	return out
}

// gradient maps 0..1 onto low -> mid -> high.
func gradient(t float64) string {
	if t < 0.5 {
		return mix(colVizLow, colVizMid, t/0.5)
	}
	return mix(colVizMid, colVizHigh, (t-0.5)/0.5)
}

func mix(a, b string, t float64) string {
	var ar, ag, ab, br, bg, bb int
	fmt.Sscanf(a, "#%02x%02x%02x", &ar, &ag, &ab)
	fmt.Sscanf(b, "#%02x%02x%02x", &br, &bg, &bb)
	l := func(x, y int) int { return int(float64(x) + float64(y-x)*t) }
	return fmt.Sprintf("#%02x%02x%02x", l(ar, br), l(ag, bg), l(ab, bb))
}
