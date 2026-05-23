// Package tail provides a file-tailing mechanism that emits new lines
// as they are appended to a log file.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// PollInterval is how often the tailer checks for new data.
const PollInterval = 200 * time.Millisecond

// Tailer reads lines from a file, following new content as it is appended.
type Tailer struct {
	path  string
	Lines chan string
}

// New creates a Tailer for the given file path.
func New(path string) *Tailer {
	return &Tailer{
		path:  path,
		Lines: make(chan string, 256),
	}
}

// Run opens the file, seeks to the end, and streams new lines until ctx is
// cancelled. It polls the file at PollInterval when no data is available.
func (t *Tailer) Run(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only tail new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			// Partial read or no new data — wait and retry.
			time.Sleep(PollInterval)
			continue
		}

		if line != "" {
			// Strip trailing newline before sending.
			if len(line) > 0 && line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			select {
			case t.Lines <- line:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
