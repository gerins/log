# 📜 Log Package 
[![Generic badge](https://img.shields.io/badge/Go-v1.17.0-blue.svg)](https://golang.org/doc/go1.17)
[![Generic badge](https://img.shields.io/badge/status-development-green.svg)](https://shields.io/)
[![Generic badge](https://img.shields.io/badge/release-v1.1.3-yellow.svg)](https://shields.io/)


## 📌 Getting Started
```shell
go get -u github.com/gerins/log
```

### Echo 
```go
func main() {
	// Using default configuration, please use InitWithConfig() for production environment.
	log.Init() 
	e := echo.New()

	// Init logging middleware
	e.Use(middlewareLog.SetLogRequest())                       // Mandatory
	e.Use(middleware.BodyDump(middlewareLog.SaveLogRequest())) // Mandatory

	// Init handler
	e.GET("/", func(c echo.Context) error {
		// Get context from echo locals.
		ctx := c.Get("ctx").(context.Context)

		// Assign user id to Log Request model
		// So wen can know who make the request to the server
		log.Context(ctx).UserID = 2020

		// Log Request
		log.Context(ctx).Debug("Testing Log Request Debug")
		log.Context(ctx).Info("Testing Log Request Info")
		log.Context(ctx).Warn("Testing Log Request Warn")
		log.Context(ctx).Error("Testing Log Request Error")

		// Global log
		log.Debug("Testing Global Log Debug")
		log.Info("Testing Global Log Info")
		log.Warn("Testing Global Log Warn")
		log.Error("Testing Global Log Error")

		return c.String(http.StatusOK, "Hello, Log!")
	})

	e.Start("localhost:8080")
}
```

## 🍀 Sample Log Request
```json
{
   "level": "info",
   "time": "2021-10-23T15:02:53.189+0700",
   "caller": "echo/log_request.go:59",
   "msg": "REQUEST_LOG",
   "ProcessID": "JLFyq8YcfN5KOy9knY42",
   "UserID": 2020,
   "IP": "127.0.0.1",
   "Method": "GET",
   "URL": "localhost:8080/",
   "RequestHeader": {
      "Accept": ["*/*"],
      "User-Agent": ["curl/7.68.0"]
   },
   "RequestBody": {},
   "ResponseHeader": {
      "Content-Type": ["text/plain; charset=UTF-8"]
   },
   "ResponseBody": null,
   "StatusCode": 200,
   "RequestDuration": 1,
   "ExtraData": null,
   "SubLog": [
      {
         "level": "DEBUG",
         "message": "Testing Log Request Debug"
      },
      {
         "level": "INFO",
         "message": "Testing Log Request Info"
      },
      {
         "level": "WARN",
         "message": "Testing Log Request Warn"
      },
      {
         "level": "ERROR",
         "message": "Testing Log Request Error"
      }
   ]
}
```
