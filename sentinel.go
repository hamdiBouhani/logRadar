package logradar

import "time"

// LogEntry is a normalized log record passed to detectors.
type LogEntry struct {
	Level     string            // e.g. "INFO", "ERROR"
	Message   string            // the raw log message
	Timestamp time.Time         // when it happened
	Fields    map[string]string // optional structured fields
}

// Alert is emitted when a detector fires.
type Alert struct {
	Entry  LogEntry
	Reason string // human-readable explanation
}

// Detector is the interface every detection rule implements.
// Return (true, reason) if the entry is an error.
type Detector interface {
	Detect(entry LogEntry) (bool, string)
}

// Sentinel runs a list of detectors on every log entry.
type Sentinel struct {
	detectors []Detector
	onAlert   func(Alert)
}

// New creates a Sentinel with the given detectors.
// onAlert is called for every alert (use nil for a default stdout printer).
func New(detectors []Detector, onAlert func(Alert)) *Sentinel {
	if onAlert == nil {
		onAlert = func(a Alert) {
			println("[ALERT]", a.Reason, "->", a.Entry.Message)
		}
	}
	return &Sentinel{detectors: detectors, onAlert: onAlert}
}

// Process runs all detectors against a single log entry.
// Stops at the first detector that fires (change to "continue" for multiple alerts).
func (s *Sentinel) Process(entry LogEntry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	for _, d := range s.detectors {
		if ok, reason := d.Detect(entry); ok {
			s.onAlert(Alert{Entry: entry, Reason: reason})
			return
		}
	}
}

// Info is a convenience method.
func (s *Sentinel) Info(msg string) {
	s.Process(LogEntry{Level: "INFO", Message: msg})
}

// Error is a convenience method.
func (s *Sentinel) Error(msg string) {
	s.Process(LogEntry{Level: "ERROR", Message: msg})
}
