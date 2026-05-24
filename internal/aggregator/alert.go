package aggregator

import "time"

// AlertLevel classifies the severity of a detected anomaly.
type AlertLevel int

const (
	AlertWarn AlertLevel = iota
	AlertCritical
)

// Alert represents an anomaly detected within a bucket.
type Alert struct {
	Fingerprint string
	Message     string
	Level       AlertLevel
	Rate        float64
	Threshold   float64
	DetectedAt  time.Time
}

// AlertRule defines thresholds for triggering alerts on a bucket.
type AlertRule struct {
	// WarnRate triggers a warning alert when events/sec exceeds this value.
	WarnRate float64
	// CriticalRate triggers a critical alert when events/sec exceeds this value.
	CriticalRate float64
}

// DefaultAlertRule is used when no custom rule is configured.
var DefaultAlertRule = AlertRule{
	WarnRate:     10.0,
	CriticalRate: 50.0,
}

// Evaluator checks store snapshots against alert rules and emits Alerts.
type Evaluator struct {
	rule AlertRule
}

// NewEvaluator creates an Evaluator with the given rule.
func NewEvaluator(rule AlertRule) *Evaluator {
	return &Evaluator{rule: rule}
}

// Evaluate inspects each bucket snapshot and returns any triggered alerts.
func (e *Evaluator) Evaluate(snapshots []BucketSnapshot) []Alert {
	var alerts []Alert
	for _, s := range snapshots {
		var al *Alert
		switch {
		case s.Rate >= e.rule.CriticalRate:
			al = &Alert{
				Fingerprint: s.Fingerprint,
				Message:     s.Sample,
				Level:       AlertCritical,
				Rate:        s.Rate,
				Threshold:   e.rule.CriticalRate,
				DetectedAt:  time.Now(),
			}
		case s.Rate >= e.rule.WarnRate:
			al = &Alert{
				Fingerprint: s.Fingerprint,
				Message:     s.Sample,
				Level:       AlertWarn,
				Rate:        s.Rate,
				Threshold:   e.rule.WarnRate,
				DetectedAt:  time.Now(),
			}
		}
		if al != nil {
			alerts = append(alerts, *al)
		}
	}
	return alerts
}
