package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/user/logfold/internal/tail"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	_, err := f.WriteString(line + "\n")
	if err != nil {
		t.Fatalf("writeLine: %v", err)
	}
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logfold-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tr := tail.New(f.Name())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- tr.Run(ctx)
	}()

	// Give the tailer time to open the file and seek to end.
	time.Sleep(50 * time.Millisecond)

	wantLines := []string{
		"2024-01-01T00:00:00Z INFO hello world",
		"2024-01-01T00:00:01Z ERROR something failed",
	}
	for _, l := range wantLines {
		writeLine(t, f, l)
	}

	received := make([]string, 0, len(wantLines))
	timeout := time.After(2 * time.Second)
	for len(received) < len(wantLines) {
		select {
		case line := <-tr.Lines:
			received = append(received, line)
		case <-timeout:
			t.Fatalf("timed out waiting for lines; got %d/%d", len(received), len(wantLines))
		}
	}

	for i, want := range wantLines {
		if received[i] != want {
			t.Errorf("line %d: got %q, want %q", i, received[i], want)
		}
	}

	cancel()
	if err := <-errCh; err != context.Canceled {
		t.Errorf("Run returned %v, want context.Canceled", err)
	}
}

func TestTailer_MissingFile_ReturnsError(t *testing.T) {
	tr := tail.New("/nonexistent/path/logfold.log")
	ctx := context.Background()
	err := tr.Run(ctx)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
