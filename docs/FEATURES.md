# Features

This library focuses on request-centric logging with sub-logs. The main features are:

- Unified request logs stored as a single JSON entry.
- Sub-logging to collect logs across handlers, use cases, and repositories.
- Structured JSON output via `slog` go standart library with custom levels.
- File rotation with `file-rotatelogs`.
- Framework middleware for Echo, Fiber, Gin, and gRPC.
- HTTP trace logging for outbound calls.
- Optional masking for sensitive fields using struct tags.

## Log Levels

Custom levels are mapped to string output in the JSON log:

- `DEBUG`
- `INFO`
- `WARN`
- `ERROR`
- `FATAL`
- `TRACE` (HTTP trace)
- `REQUEST` (request summary with sub-log)

## Sub-Logging

Sub-logs are captured under the `subLog` field in a request log entry. Use `log.Context(ctx)` to append:

```go
log.Context(ctx).Info("repository call")
log.Context(ctx).Error("db error")
```

You can disable sub-log aggregation and force global output via `Config.DisableSubLogs`.

## Sensitive Data Masking

When `Config.HideSensitiveData` is enabled, fields tagged with `log:"hide"` are masked.

```go
type LoginRequest struct {
    Email    string
    Password string `log:"hide"`
}
```

## HTTP Trace

Use `log.NewTrace` to record outbound HTTP requests. You can log to trace output or attach it to `ExtraData`.

```go
trace := log.NewTrace("GET", url, reqHeader, reqBody, false)
// ... make request and read response body
trace.RawRespBody = rawBody
trace.Save(ctx, resp)
```
