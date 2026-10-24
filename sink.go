package logradar

import (
	"context"
	"errors"
)

// Sink receives alerts. Implementations should be safe for concurrent use.
type Sink interface {
	Send(ctx context.Context, alert Alert) error
	Close() error
}

// MultiSink fans an alert out to several sinks.
// Errors from individual sinks are joined; sending continues to the rest.
type MultiSink struct {
	sinks []Sink
}

func NewMultiSink(sinks ...Sink) *MultiSink {
	return &MultiSink{sinks: sinks}
}

func (m *MultiSink) Send(ctx context.Context, alert Alert) error {
	var errs []error
	for _, s := range m.sinks {
		if err := s.Send(ctx, alert); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (m *MultiSink) Close() error {
	var errs []error
	for _, s := range m.sinks {
		if err := s.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// AsAlertHandler adapts a Sink into the func(Alert) expected by New().
// Failures are swallowed (printed) so they never break the log pipeline.
func AsAlertHandler(sink Sink) func(Alert) {
	return func(a Alert) {
		if err := sink.Send(context.Background(), a); err != nil {
			println("[sink error]", err.Error())
		}
	}
}
