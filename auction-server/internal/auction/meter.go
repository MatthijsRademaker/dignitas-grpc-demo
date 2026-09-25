package auction

import "time"

// meter counts events over a sliding window. Not safe for concurrent use: the
// House guards it with its mutex.
type meter struct {
	window time.Duration
	times  []time.Time
}

func (m *meter) mark(now time.Time) {
	m.times = append(m.times, now)
}

func (m *meter) perSecond(now time.Time) float64 {
	m.prune(now)
	return float64(len(m.times)) / m.window.Seconds()
}

func (m *meter) count(now time.Time) int {
	m.prune(now)
	return len(m.times)
}

func (m *meter) prune(now time.Time) {
	cutoff := now.Add(-m.window)
	i := 0
	for i < len(m.times) && !m.times[i].After(cutoff) {
		i++
	}
	m.times = m.times[i:]
}
