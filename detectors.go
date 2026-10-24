package logradar

import (
	"strings"
	"sync"
	"time"
)

// ─────────────────────────────────────────────────────────────
// 1. KeywordDetector — matches keywords in the message.
// ─────────────────────────────────────────────────────────────
type KeywordDetector struct {
	keywords []string // stored uppercased
}

func NewKeywordDetector(keywords ...string) *KeywordDetector {
	upper := make([]string, len(keywords))
	for i, k := range keywords {
		upper[i] = strings.ToUpper(k)
	}
	return &KeywordDetector{keywords: upper}
}

func (d *KeywordDetector) Detect(entry LogEntry) (bool, string) {
	msg := strings.ToUpper(entry.Message)
	for _, kw := range d.keywords {
		if strings.Contains(msg, kw) {
			return true, "keyword match: " + kw
		}
	}
	return false, ""
}

// ─────────────────────────────────────────────────────────────
// 2. LevelDetector — fires on specific log levels.
// ─────────────────────────────────────────────────────────────
type LevelDetector struct {
	levels map[string]bool
}

func NewLevelDetector(levels ...string) *LevelDetector {
	m := make(map[string]bool)
	for _, l := range levels {
		m[strings.ToUpper(l)] = true
	}
	return &LevelDetector{levels: m}
}

func (d *LevelDetector) Detect(entry LogEntry) (bool, string) {
	if d.levels[strings.ToUpper(entry.Level)] {
		return true, "level: " + entry.Level
	}
	return false, ""
}

// ─────────────────────────────────────────────────────────────
// 3. SpikeDetector — fires when N errors happen within a window.
// ─────────────────────────────────────────────────────────────
type SpikeDetector struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	hits      []time.Time // timestamps of recent error-level entries
}

func NewSpikeDetector(threshold int, window time.Duration) *SpikeDetector {
	return &SpikeDetector{threshold: threshold, window: window}
}

func (d *SpikeDetector) Detect(entry LogEntry) (bool, string) {
	lvl := strings.ToUpper(entry.Level)
	if lvl != "ERROR" && lvl != "FATAL" && lvl != "CRITICAL" {
		return false, ""
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := entry.Timestamp
	if now.IsZero() {
		now = time.Now()
	}

	// Drop timestamps older than the window.
	cutoff := now.Add(-d.window)
	kept := d.hits[:0]
	for _, t := range d.hits {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	d.hits = append(kept, now)

	if len(d.hits) >= d.threshold {
		d.hits = d.hits[:0] // reset after firing
		return true, "spike: errors per window exceeded"
	}
	return false, ""
}
