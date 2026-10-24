package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	logradar "github.com/hamdiBouhani/logRadar"
)

func main() {
	// 1. Build the Sentinel.
	sentinel := logradar.New(
		[]logradar.Detector{
			logradar.NewLevelDetector("ERROR", "FATAL"),
			logradar.NewRegexDetector(`(?i)panic|segfault|oom`),
			logradar.NewSpikeDetector(3, 5*time.Second),
		},
		func(a logradar.Alert) {
			fmt.Printf("🚨 [%s] %s | %s\n",
				a.Entry.Level, a.Reason, a.Entry.Message)
		},
	)

	// 2. Wire it into log/slog.
	logger := slog.New(logradar.NewSlogHandler(sentinel))

	logger.Info("server started")
	logger.Error("database connection timeout")
	logger.Error("payment failed")
	logger.Warn("memory pressure")
	logger.Error("another error") // 3rd error in window → spike
	logger.Info("recovered")
	_ = context.Background()
}
