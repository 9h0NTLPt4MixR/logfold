package aggregator

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

// Level represents the severity of a log entry.
type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)

// Entry represents a single parsed log line.
type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Fields    map[string]string
	Raw       string
}

// Fingerprint returns a short hash that identifies structurally
// similar log messages, ignoring dynamic tokens like numbers and UUIDs.
func (e *Entry) Fingerprint() string {
	normalized := normalizeMessage(e.Message)
	sum := md5.Sum([]byte(string(e.Level) + ":" + normalized))
	return fmt.Sprintf("%x", sum[:6])
}

// normalizeMessage strips dynamic parts from a message so that
// repeated variants collapse to the same fingerprint.
func normalizeMessage(msg string) string {
	words := strings.Fields(msg)
	for i, w := range words {
		if looksNumeric(w) || looksUUID(w) || looksPath(w) {
			words[i] = "<X>"
		}
	}
	return strings.Join(words, " ")
}

func looksNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && c != '.' && c != '-' {
			return false
		}
	}
	return true
}

func looksUUID(s string) bool {
	// rough heuristic: 8-4-4-4-12 hex groups
	return len(s) == 36 && s[8] == '-' && s[13] == '-'
}

func looksPath(s string) bool {
	return strings.HasPrefix(s, "/") || strings.HasPrefix(s, "./")
}
