package dashboard

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	colFP      = 14
	colLevel   = 7
	colSample  = 38
	colCount   = 8
	colRate    = 9
	colLastSeen = 10
)

var levelColors = map[string]string{
	"ERROR": "\033[31m",
	"WARN":  "\033[33m",
	"INFO":  "\033[32m",
	"DEBUG": "\033[36m",
}

const resetColor = "\033[0m"

// Render writes a full dashboard frame to w.
func Render(w io.Writer, m *Model) {
	fmt.Fprintf(w, "\033[H\033[2J") // clear screen, move cursor home

	fmt.Fprintf(w, "  logfold  —  %s    (q quit  p pause)\n\n",
		m.LastUpdate.Format(time.RFC3339))

	header := fmt.Sprintf("  %-*s  %-*s  %-*s  %*s  %*s  %*s\n",
		colFP, "FINGERPRINT",
		colLevel, "LEVEL",
		colSample, "SAMPLE",
		colCount, "COUNT",
		colRate, "RATE/s",
		colLastSeen, "LAST SEEN",
	)
	fmt.Fprint(w, header)
	fmt.Fprintln(w, "  "+strings.Repeat("─", len(header)-3))

	for _, row := range m.Rows() {
		color := levelColors[strings.ToUpper(row.Level)]
		if color == "" {
			color = ""
		}

		sample := row.Sample
		if len(sample) > colSample {
			sample = sample[:colSample-1] + "…"
		}

		deltaStr := ""
		if row.Delta > 0 {
			deltaStr = fmt.Sprintf("+%d", row.Delta)
		}

		fmt.Fprintf(w, "  %-*s  %s%-*s%s  %-*s  %*s%s  %*.2f  %*s\n",
			colFP, truncate(row.Fingerprint, colFP),
			color, colLevel, truncate(row.Level, colLevel), resetColor,
			colSample, sample,
			colCount, fmt.Sprintf("%d", row.Count), deltaStr,
			colRate, row.Rate,
			colLastSeen, row.LastSeen.Format("15:04:05"),
		)
	}

	if m.Paused {
		fmt.Fprintln(w, "\n  [PAUSED]")
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
