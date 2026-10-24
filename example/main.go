package main

import (
	"fmt"
	"time"

	logsentinel "github.com/hamdiBouhani/logRadar"
)

func main() {
	sentinel := logsentinel.New(
		[]logsentinel.Detector{
			logsentinel.NewLevelDetector("ERROR", "FATAL"),
			logsentinel.NewKeywordDetector("panic", "timeout"),
			logsentinel.NewSpikeDetector(3, 5*time.Second),
		},
		func(a logsentinel.Alert) {
			fmt.Printf("🚨 [%s] %s | %s\n",
				a.Entry.Level, a.Reason, a.Entry.Message)
		},
	)

	sentinel.Info("server started")
	sentinel.Error("database connection timeout") // keyword + level
	sentinel.Error("payment failed")
	sentinel.Error("another error") // 3rd error in window → spike
	sentinel.Info("recovered")
}
