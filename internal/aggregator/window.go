package aggregator

import (
	"sync"
	"time"
)

// Window tracks a rolling time window of event counts for rate calculation.
type Window struct {
	mu       sync.Mutex
	slots    []int64
	slotSize time.Duration
	numSlots int
	lastSlot int64
}

// NewWindow creates a Window that covers totalDuration split into numSlots buckets.
func NewWindow(totalDuration time.Duration, numSlots int) *Window {
	return &Window{
		slots:    make([]int64, numSlots),
		slotSize: totalDuration / time.Duration(numSlots),
		numSlots: numSlots,
		lastSlot: time.Now().UnixNano() / int64(totalDuration/time.Duration(numSlots)),
	}
}

// currentSlot returns the slot index for the given time.
func (w *Window) currentSlotIndex(now time.Time) int64 {
	return now.UnixNano() / int64(w.slotSize)
}

// advance clears stale slots up to now.
func (w *Window) advance(now time.Time) {
	current := w.currentSlotIndex(now)
	delta := current - w.lastSlot
	if delta <= 0 {
		return
	}
	clear := delta
	if clear > int64(w.numSlots) {
		clear = int64(w.numSlots)
	}
	for i := int64(1); i <= clear; i++ {
		idx := (w.lastSlot + i) % int64(w.numSlots)
		w.slots[idx] = 0
	}
	w.lastSlot = current
}

// Record increments the count for the current time slot.
func (w *Window) Record(now time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(now)
	idx := w.currentSlotIndex(now) % int64(w.numSlots)
	w.slots[idx]++
}

// Sum returns the total count across all slots in the window.
func (w *Window) Sum(now time.Time) int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(now)
	var total int64
	for _, v := range w.slots {
		total += v
	}
	return total
}

// Rate returns events per second over the window duration.
func (w *Window) Rate(now time.Time) float64 {
	total := w.Sum(now)
	windowSec := float64(w.slotSize) * float64(w.numSlots) / float64(time.Second)
	if windowSec == 0 {
		return 0
	}
	return float64(total) / windowSec
}
