package dashboard

import (
	"fmt"
	"strings"

	"github.com/user/logfold/internal/aggregator"
)

const (
	alertWarnPrefix     = "[WARN] "
	alertCriticalPrefix = "[CRIT] "
	alertMaxRows        = 10
)

// AlertsView renders a list of active alerts as a plain string block
// suitable for embedding in the terminal dashboard.
func AlertsView(alerts []aggregator.Alert) string {
	if len(alerts) == 0 {
		return "  No active alerts\n"
	}

	var sb strings.Builder
	limit := len(alerts)
	if limit > alertMaxRows {
		limit = alertMaxRows
	}

	for _, a := range alerts[:limit] {
		prefix := alertWarnPrefix
		if a.Level == aggregator.AlertCritical {
			prefix = alertCriticalPrefix
		}
		line := fmt.Sprintf("  %s%-40s rate=%.1f/s (threshold %.1f/s)\n",
			prefix,
			truncate(a.Message, 40),
			a.Rate,
			a.Threshold,
		)
		sb.WriteString(line)
	}

	if len(alerts) > alertMaxRows {
		sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(alerts)-alertMaxRows))
	}

	return sb.String()
}
