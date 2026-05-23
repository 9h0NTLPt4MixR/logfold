package parser

import (
	"fmt"
	"strings"
	"time"
)

// Level represents a log severity level.
type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
	LevelUnknown Level = "UNKNOWN"
)

// Entry is a parsed log line.
type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Raw       string
}

// Parser attempts to parse raw log lines into structured entries.
type Parser struct {
	formats []LineFormat
}

// New returns a Parser pre-loaded with built-in formats.
func New() *Parser {
	return &Parser{
		formats: defaultFormats(),
	}
}

// Parse attempts to parse a raw log line. If no format matches, a best-effort
// entry is returned with LevelUnknown and the raw line as the message.
func (p *Parser) Parse(raw string) (*Entry, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty line")
	}
	for _, f := range p.formats {
		if e, ok := f.TryParse(raw); ok {
			e.Raw = raw
			return e, nil
		}
	}
	return &Entry{
		Timestamp: time.Now(),
		Level:     LevelUnknown,
		Message:   raw,
		Raw:       raw,
	}, nil
}

// normalizeLevel maps common level strings to a canonical Level.
func normalizeLevel(s string) Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG", "DBG", "TRACE":
		return LevelDebug
	case "INFO", "INFORMATION":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR", "ERR":
		return LevelError
	case "FATAL", "CRITICAL", "CRIT":
		return LevelFatal
	default:
		return LevelUnknown
	}
}
