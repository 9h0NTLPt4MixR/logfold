package dashboard

import (
	"time"

	"github.com/user/logfold/internal/aggregator"
)

// BucketRow represents a single row in the dashboard table.
type BucketRow struct {
	Fingerprint string
	Level       string
	Sample      string
	Count       int64
	Rate        float64
	LastSeen    time.Time
	Delta       int64 // count change since last refresh
}

// Model holds the state for the terminal dashboard.
type Model struct {
	store      *aggregator.Store
	rows       []BucketRow
	prevCounts map[string]int64
	LastUpdate time.Time
	Paused     bool
}

// NewModel creates a new dashboard Model backed by the given Store.
func NewModel(s *aggregator.Store) *Model {
	return &Model{
		store:      s,
		prevCounts: make(map[string]int64),
	}
}

// Refresh pulls the latest snapshot from the store and updates rows.
func (m *Model) Refresh() {
	if m.Paused {
		return
	}

	snapshots := m.store.Snapshot()
	m.rows = make([]BucketRow, 0, len(snapshots))

	for _, s := range snapshots {
		delta := s.Count - m.prevCounts[s.Fingerprint]
		m.prevCounts[s.Fingerprint] = s.Count

		m.rows = append(m.rows, BucketRow{
			Fingerprint: s.Fingerprint,
			Level:       s.Level,
			Sample:      s.Sample,
			Count:       s.Count,
			Rate:        s.Rate,
			LastSeen:    s.LastSeen,
			Delta:       delta,
		})
	}

	m.LastUpdate = time.Now()
}

// Rows returns the current snapshot rows.
func (m *Model) Rows() []BucketRow {
	return m.rows
}
