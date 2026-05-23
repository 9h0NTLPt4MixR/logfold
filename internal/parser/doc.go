// Package parser provides log line parsing for logfold.
//
// It supports multiple common log formats via a priority-ordered list of
// [LineFormat] matchers. When no format matches, a best-effort [Entry] is
// returned with [LevelUnknown] so that no log line is silently dropped.
//
// Supported formats (in matching order):
//
//   - ISO 8601 timestamp + level + message
//     e.g. "2024-01-15T12:34:56Z INFO request completed"
//
//   - Bracketed timestamp and level
//     e.g. "[2024-01-15 12:34:56] [ERROR] something went wrong"
//
//   - Level-colon prefix (no timestamp)
//     e.g. "ERROR: connection refused"
//
// Usage:
//
//	p := parser.New()
//	entry, err := p.Parse(rawLine)
//	if err != nil {
//		// line was empty
//	}
//	fmt.Println(entry.Level, entry.Message)
package parser
