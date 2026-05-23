package parser

import (
	"regexp"
	"time"
)

// LineFormat is an interface for a specific log format matcher.
type LineFormat interface {
	TryParse(line string) (*Entry, bool)
}

// defaultFormats returns the ordered list of built-in formats.
// More specific formats should come first.
func defaultFormats() []LineFormat {
	return []LineFormat{
		&regexFormat{
			// e.g. 2024-01-15T12:34:56.789Z INFO some message here
			name:    "iso8601-level-msg",
			re:      regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?)\s+(\w+)\s+(.+)$`),
			tsFmt:   time.RFC3339,
			tsIdx:   1,
			lvlIdx:  2,
			msgIdx:  3,
		},
		&regexFormat{
			// e.g. [2024-01-15 12:34:56] [ERROR] some message
			name:    "bracketed-ts-level",
			re:      regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]\s*\[(\w+)\]\s*(.+)$`),
			tsFmt:   "2006-01-02 15:04:05",
			tsIdx:   1,
			lvlIdx:  2,
			msgIdx:  3,
		},
		&regexFormat{
			// e.g. ERROR: some message (no timestamp)
			name:   "level-colon-msg",
			re:     regexp.MustCompile(`^(DEBUG|INFO|WARN|WARNING|ERROR|FATAL|CRITICAL|TRACE):\s+(.+)$`),
			tsFmt:  "",
			tsIdx:  0,
			lvlIdx: 1,
			msgIdx: 2,
		},
	}
}

type regexFormat struct {
	name   string
	re     *regexp.Regexp
	tsFmt  string
	tsIdx  int
	lvlIdx int
	msgIdx int
}

func (f *regexFormat) TryParse(line string) (*Entry, bool) {
	m := f.re.FindStringSubmatch(line)
	if m == nil {
		return nil, false
	}
	var ts time.Time
	if f.tsIdx > 0 && f.tsFmt != "" {
		parsed, err := time.Parse(f.tsFmt, m[f.tsIdx])
		if err != nil {
			ts = time.Now()
		} else {
			ts = parsed
		}
	} else {
		ts = time.Now()
	}
	return &Entry{
		Timestamp: ts,
		Level:     normalizeLevel(m[f.lvlIdx]),
		Message:   m[f.msgIdx],
	}, true
}
