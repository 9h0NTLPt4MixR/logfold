package aggregator

import (
	"testing"
	"time"
)

func TestWindow_RecordAndSum(t *testing.T) {
	w := NewWindow(10*time.Second, 10)
	now := time.Now()

	for i := 0; i < 5; i++ {
		w.Record(now)
	}

	if got := w.Sum(now); got != 5 {
		t.Errorf("expected sum 5, got %d", got)
	}
}

func TestWindow_StaleSlotsClearedOnAdvance(t *testing.T) {
	w := NewWindow(2*time.Second, 2)
	past := time.Now().Add(-5 * time.Second)

	for i := 0; i < 10; i++ {
		w.Record(past)
	}

	now := time.Now()
	if got := w.Sum(now); got != 0 {
		t.Errorf("expected stale slots to be cleared, got %d", got)
	}
}

func TestWindow_Rate(t *testing.T) {
	w := NewWindow(10*time.Second, 10)
	now := time.Now()

	for i := 0; i < 20; i++ {
		w.Record(now)
	}

	rate := w.Rate(now)
	// 20 events over 10 seconds = 2.0 eps
	if rate < 1.9 || rate > 2.1 {
		t.Errorf("expected rate ~2.0, got %f", rate)
	}
}

func TestWindow_MultipleSlots(t *testing.T) {
	w := NewWindow(4*time.Second, 4)
	base := time.Now().Truncate(time.Second)

	w.Record(base)
	w.Record(base.Add(1 * time.Second))
	w.Record(base.Add(2 * time.Second))

	now := base.Add(2 * time.Second)
	if got := w.Sum(now); got != 3 {
		t.Errorf("expected sum 3 across slots, got %d", got)
	}
}

func TestWindow_RateZeroOnEmpty(t *testing.T) {
	w := NewWindow(10*time.Second, 10)
	if r := w.Rate(time.Now()); r != 0 {
		t.Errorf("expected 0 rate on empty window, got %f", r)
	}
}
