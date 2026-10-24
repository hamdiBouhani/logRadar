package logradar

import (
	"testing"
	"time"
)

func TestKeywordDetector(t *testing.T) {
	d := NewKeywordDetector("error", "failed")
	if ok, _ := d.Detect(LogEntry{Message: "database FAILED"}); !ok {
		t.Fatal("expected match")
	}
	if ok, _ := d.Detect(LogEntry{Message: "all good"}); ok {
		t.Fatal("expected no match")
	}
}

func TestLevelDetector(t *testing.T) {
	d := NewLevelDetector("ERROR", "FATAL")
	if ok, _ := d.Detect(LogEntry{Level: "error"}); !ok {
		t.Fatal("expected match")
	}
}

func TestSpikeDetector(t *testing.T) {
	d := NewSpikeDetector(3, time.Second)
	now := time.Now()
	for i := 0; i < 2; i++ {
		if ok, _ := d.Detect(LogEntry{Level: "ERROR", Timestamp: now}); ok {
			t.Fatal("fired too early")
		}
	}
	if ok, _ := d.Detect(LogEntry{Level: "ERROR", Timestamp: now}); !ok {
		t.Fatal("expected spike")
	}
}

func TestRegexDetector(t *testing.T) {
	d := NewRegexDetector(`(?i)panic|segfault`)
	if ok, _ := d.Detect(LogEntry{Message: "kernel PANIC"}); !ok {
		t.Fatal("expected regex match")
	}
	if ok, _ := d.Detect(LogEntry{Message: "everything fine"}); ok {
		t.Fatal("expected no match")
	}
}
