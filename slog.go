package logradar

import (
	"context"
	"log/slog"
)

// SlogHandler is a slog.Handler that forwards every record to a Sentinel.
// It never writes to the underlying stream — pair it with another handler
// (or a MultiHandler) if you also want normal output.
type SlogHandler struct {
	sentinel *Sentinel
	attrs    []slog.Attr
	group    string
}

// NewSlogHandler wraps a Sentinel in an slog.Handler.
//
//	sentinel := logradar.New(detectors, nil)
//	logger := slog.New(logradar.NewSlogHandler(sentinel))
func NewSlogHandler(s *Sentinel) *SlogHandler {
	return &SlogHandler{sentinel: s}
}

func (h *SlogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Accept every level — the Sentinel's detectors decide what matters.
	return true
}

func (h *SlogHandler) Handle(_ context.Context, r slog.Record) error {
	fields := make(map[string]string, len(h.attrs)+r.NumAttrs())

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.String()
	}
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.String()
		return true
	})

	h.sentinel.Process(LogEntry{
		Level:     r.Level.String(),
		Message:   r.Message,
		Timestamp: r.Time,
		Fields:    fields,
	})
	return nil
}

func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := *h
	clone.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &clone
}

func (h *SlogHandler) WithGroup(name string) slog.Handler {
	clone := *h
	if h.group != "" {
		clone.group = h.group + "." + name
	} else {
		clone.group = name
	}
	return &clone
}