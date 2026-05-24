package aggregator

import (
	"testing"
	"time"
)

func snapshot(fp, sample string, rate float64) BucketSnapshot {
	return BucketSnapshot{
		Fingerprint: fp,
		Sample:      sample,
		Count:       100,
		Rate:        rate,
		FirstSeen:   time.Now().Add(-time.Minute),
		LastSeen:    time.Now(),
	}
}

func TestEvaluator_NoAlertBelowWarn(t *testing.T) {
	e := NewEvaluator(DefaultAlertRule)
	snaps := []BucketSnapshot{snapshot("fp1", "msg", 5.0)}
	alerts := e.Evaluate(snaps)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestEvaluator_WarnAlert(t *testing.T) {
	e := NewEvaluator(DefaultAlertRule)
	snaps := []BucketSnapshot{snapshot("fp2", "warn msg", 15.0)}
	alerts := e.Evaluate(snaps)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != AlertWarn {
		t.Errorf("expected AlertWarn, got %v", alerts[0].Level)
	}
	if alerts[0].Fingerprint != "fp2" {
		t.Errorf("unexpected fingerprint: %s", alerts[0].Fingerprint)
	}
}

func TestEvaluator_CriticalAlert(t *testing.T) {
	e := NewEvaluator(DefaultAlertRule)
	snaps := []BucketSnapshot{snapshot("fp3", "crit msg", 60.0)}
	alerts := e.Evaluate(snaps)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != AlertCritical {
		t.Errorf("expected AlertCritical, got %v", alerts[0].Level)
	}
	if alerts[0].Threshold != DefaultAlertRule.CriticalRate {
		t.Errorf("unexpected threshold: %f", alerts[0].Threshold)
	}
}

func TestEvaluator_MultipleSnapshots(t *testing.T) {
	e := NewEvaluator(DefaultAlertRule)
	snaps := []BucketSnapshot{
		snapshot("fp4", "ok", 1.0),
		snapshot("fp5", "warn", 20.0),
		snapshot("fp6", "crit", 55.0),
	}
	alerts := e.Evaluate(snaps)
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestEvaluator_CustomRule(t *testing.T) {
	rule := AlertRule{WarnRate: 2.0, CriticalRate: 5.0}
	e := NewEvaluator(rule)
	snaps := []BucketSnapshot{snapshot("fp7", "msg", 3.0)}
	alerts := e.Evaluate(snaps)
	if len(alerts) != 1 || alerts[0].Level != AlertWarn {
		t.Errorf("expected 1 warn alert with custom rule")
	}
}
