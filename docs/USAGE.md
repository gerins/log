# Usage

## Installation

```bash
go get -u github.com/gerins/log
```

## Initialize Logger

```go
import "github.com/gerins/log"

func main() {
    log.Init() // With default config
    // or
    log.InitWithConfig(log.Config{
		LogToTerminal:     true,
		LogToFile:         true,
		Location:          "/log/",
		FileLogName:       "server_log",
		FileFormat:        ".%Y-%b-%d-%H-%M.log",
		MaxAge:            30,
		RotationFile:      24,
		Level:             LevelDebug,
		HideSensitiveData: false,
        DisableSubLogs:    false,
    })
}
```

## Global Logging

```go
log.Debug("message")
log.Infof("hello %s", "world")
log.Error("something failed")
```

## Request Logging With Sub-Logs

Attach a request logger to your context and use `log.Context(ctx)` to append sub-logs.

```go
ctx := log.NewRequest().SaveToContext(context.Background())
log.Context(ctx).Info("handler start")
log.Context(ctx).Warn("slow query")
log.Context(ctx).ExtraData["userData"] = user
log.Context(ctx).Save()
```

### Waiting For Goroutines Before Save

`request.Save()` waits on `WaitGroup` before printing. If you log inside goroutines, add them to the request `WaitGroup` so the sub-logs are complete.

```go
req := log.NewRequest()
ctx := req.SaveToContext(context.Background())

req.WaitGroup.Add(1)
go func() {
    defer req.WaitGroup.Done()
    log.Context(ctx).Info("async job finished")
}()

req.Save()
```

## Record Duration

```go
timer := log.Context(ctx).RecordDuration("handler total process duration")
// ... do work
timer.Stop()
```

## HTTP Trace (Outbound)

```go
trace := log.NewTrace("GET", url, reqHeader, reqBody, false)
// ... make request
trace.RawRespBody = rawBody
trace.Save(ctx, resp)
```

## Framework Middleware

### Echo

```go
import (
    logMiddleware "github.com/gerins/log/middleware/echo"
    "github.com/labstack/echo/v4/middleware"
)

e.Use(logMiddleware.SetLogRequest())
e.Use(middleware.BodyDump(logMiddleware.SaveLogRequest()))
```

### Fiber

```go
import logMiddleware "github.com/gerins/log/middleware/fiber"

a.Use(logMiddleware.SaveLogRequest())
```

### Gin

```go
import logMiddleware "github.com/gerins/log/middleware/gin"

r.Use(logMiddleware.SetLogRequest())
r.Use(logMiddleware.SaveLogRequest())
```

### gRPC (Unary Interceptor)

```go
import logMiddleware "github.com/gerins/log/middleware/grpc"

server := grpc.NewServer(
    grpc.UnaryInterceptor(logMiddleware.SaveLogRequest()),
)
```

## GORM Extension

```go
import (
    logGorm "github.com/gerins/log/extension/gorm"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logGorm.Default,
})
if err != nil {
    panic(err)
}
```

## Sensitive Data Masking

```go
type LoginRequest struct {
    Email    string
    Password string `log:"hide"`
}

log.InitWithConfig(log.Config{HideSensitiveData: true})
```
