package aggregator

import (
	"sync"
)

// Store manages a collection of Buckets keyed by fingerprint.
type Store struct {
	mu      sync.RWMutex
	buckets map[string]*Bucket
}

// NewStore creates an initialised Store.
func NewStore() *Store {
	return &Store{
		buckets: make(map[string]*Bucket),
	}
}

// Ingest normalises the entry, resolves its bucket, and records it.
func (s *Store) Ingest(entry LogEntry) {
	fp := fingerprint(entry.Level, entry.Message)

	s.mu.Lock()
	b, ok := s.buckets[fp]
	if !ok {
		b = &Bucket{Fingerprint: fp}
		s.buckets[fp] = b
	}
	s.mu.Unlock()

	b.Record(entry)
}

// Snapshots returns a slice of all current bucket snapshots.
func (s *Store) Snapshots() []BucketSnapshot {
	s.mu.RLock()
	keys := make([]string, 0, len(s.buckets))
	for k := range s.buckets {
		keys = append(keys, k)
	}
	s.mu.RUnlock()

	snaps := make([]BucketSnapshot, 0, len(keys))
	for _, k := range keys {
		s.mu.RLock()
		b := s.buckets[k]
		s.mu.RUnlock()
		snaps = append(snaps, b.Snapshot())
	}
	return snaps
}

// Len returns the number of distinct buckets.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.buckets)
}
