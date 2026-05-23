package dashboard

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logfold/internal/aggregator"
)

func makeEntry(fp, level, msg string) aggregator.Entry {
	return aggregator.Entry{
		Timestamp:   time.Now(),
		Level:       level,
		Message:     msg,
		Fingerprint: fp,
		Raw:         msg,
	}
}

func TestModel_RefreshPopulatesRows(t *testing.T) {
	s := aggregator.NewStore(30 * time.Second)
	s.Ingest(makeEntry("fp1", "INFO", "hello world"))
	s.Ingest(makeEntry("fp1", "INFO", "hello world"))
	s.Ingest(makeEntry("fp2", "ERROR", "something failed"))

	m := NewModel(s)
	m.Refresh()

	if len(m.Rows()) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(m.Rows()))
	}
}

func TestModel_RefreshPausedDoesNotUpdate(t *testing.T) {
	s := aggregator.NewStore(30 * time.Second)
	s.Ingest(makeEntry("fp1", "INFO", "hello"))

	m := NewModel(s)
	m.Refresh() // initial

	m.Paused = true
	s.Ingest(makeEntry("fp2", "WARN", "new entry"))
	m.Refresh() // should be skipped

	if len(m.Rows()) != 1 {
		t.Fatalf("expected 1 row while paused, got %d", len(m.Rows()))
	}
}

func TestModel_DeltaTracking(t *testing.T) {
	s := aggregator.NewStore(30 * time.Second)
	s.Ingest(makeEntry("fp1", "INFO", "hello"))

	m := NewModel(s)
	m.Refresh()

	s.Ingest(makeEntry("fp1", "INFO", "hello"))
	s.Ingest(makeEntry("fp1", "INFO", "hello"))
	m.Refresh()

	rows := m.Rows()
	if len(rows) == 0 {
		t.Fatal("no rows returned")
	}
	if rows[0].Delta != 2 {
		t.Errorf("expected delta 2, got %d", rows[0].Delta)
	}
}

func TestRender_ContainsHeaders(t *testing.T) {
	s := aggregator.NewStore(30 * time.Second)
	s.Ingest(makeEntry("fp1", "ERROR", "disk full"))

	m := NewModel(s)
	m.Refresh()

	var buf bytes.Buffer
	Render(&buf, m)

	out := buf.String()
	for _, want := range []string{"FINGERPRINT", "LEVEL", "SAMPLE", "COUNT", "RATE/s"} {
		if !strings.Contains(out, want) {
			t.Errorf("render output missing %q", want)
		}
	}
}
