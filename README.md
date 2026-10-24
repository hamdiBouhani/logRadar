# logRadar

## Built-in Detectors

| Detector | Purpose |
|---|---|
| `KeywordDetector` | Matches literal keywords in the message. |
| `LevelDetector` | Fires on specified log levels (ERROR, FATAL, ...). |
| `SpikeDetector` | Fires when N errors occur within a time window. |
| `RegexDetector` | Fires when a compiled regex matches the message. |

## `log/slog` Integration

```go
sentinel := logradar.New(detectors, nil)
logger := slog.New(logradar.NewSlogHandler(sentinel))
logger.Error("database connection timeout")