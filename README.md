# logRadar 📡

**A small, readable error-detection library for Go log streams.**

`logRadar` watches your logs in real time and fires an alert the moment something goes wrong — a keyword like `panic`, an `ERROR` level, a burst of failures, a regex hit, or a single request that keeps erroring. You plug it into `log/slog`, pick your detectors, and route alerts to stdout, email, syslog, or your own sink.

### Why logRadar?

- **Composable detectors** — keyword, level, regex, spike, and trace-failure detectors, each ~20 lines and easy to extend.
- **Native `log/slog` integration** — one wrapper, zero changes to your existing logging code.
- **Pluggable sinks** — send alerts to email, syslog, or any custom destination via a tiny `Sink` interface.
- **Zero heavy dependencies** — built on the standard library; no ML, no magic.

### 30-second example

```go
sentinel := logradar.New(
    []logradar.Detector{
        logradar.NewLevelDetector("ERROR", "FATAL"),
        logradar.NewSpikeDetector(5, 10*time.Second),
    },
    logradar.AsAlertHandler(logradar.NewEmailSink(emailOpts)),
)

logger := slog.New(logradar.NewSlogHandler(sentinel))
logger.Error("database connection timeout") // → alert fired
```

**Tagline:** *Point it at your logs. Get alerted when they break.*

---

You can trim this to just the first paragraph + tagline if you want a super-short version for the GitHub "About" field.

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

## Sinks

`logRadar` decouples alert delivery from detection via the `Sink` interface.

```go
sink := logradar.NewMultiSink(
    logradar.NewEmailSink(logradar.EmailSinkOptions{ /* ... */ }),
    // syslogSink, slackSink, ...
)
sentinel := logradar.New(detectors, logradar.AsAlertHandler(sink))