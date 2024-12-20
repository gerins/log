# 📜 Go Server Logging Middleware

[![Go Version](https://img.shields.io/badge/Go-v1.22.0-blue.svg)](https://golang.org/doc/go1.22)[![Status](https://img.shields.io/badge/status-development-green.svg)](https://shields.io/)[![Release](https://img.shields.io/badge/release-v1.1.0-yellow.svg)](https://shields.io/)

Go Server Logging Middleware is a structured logging library for Go servers using frameworks like [Echo](https://echo.labstack.com) or [Fiber](https://gofiber.io). Its key feature is **sub-logging**, enabling you to capture the entire lifecycle of a request, including logs generated across handlers, use cases, and repositories. This makes it easier to trace and debug issues efficiently.



## 🌟 Features

- **Unified Request Logs**: Consolidates logs for each request into a single structured entry.
- **Sub-logging Support**: Tracks logs from handlers, use cases, and repositories for detailed traceability.
- **Customizable Logging**: Configure log levels, formats, and output targets to suit your needs.
- **Framework Integration**: Prebuilt examples for Echo, Fiber, and GORM.
- **JSON Structured Logs**: Easy integration with log aggregation and monitoring tools like Elasticsearch and Kibana.



## 📦 Installation

Add the package to your Go project:

```bash
go get -u github.com/gerins/log
```



## 🔧 Usage

### Example Integrations

Check out the examples directory for ready-to-use implementations:

- **[Echo Example](https://github.com/gerins/log/blob/main/example/echo/main.go)**  
- **[Fiber Example](https://github.com/gerins/log/blob/main/example/fiber/main.go)**  
- **[Gorm Example](https://github.com/gerins/log/blob/main/example/gorm/main.go)** For integrate Gorm log output to sublogging
Sub-log

### Example syntax
```go
import "github.com/gerins/log"

// General logging
log.Debug("Testing Global Log Debug")
log.Info("Testing Global Log Info")
log.Warn("Testing Global Log Warn")
log.Error("Testing Global Log Error")

// This log will append to Sub-logging
log.Context(ctx).Debug("Testing Log Request Debug")
log.Context(ctx).Info("Testing Log Request Info")
log.Context(ctx).Warn("Testing Log Request Warn")
log.Context(ctx).Error("Testing Log Request Error")
```

## 📊 Sample Logs

### General Log Entry
```json
{
    "time": "2024-06-28T20:02:49.496462+07:00",
    "level": "DEBUG",
    "msg": "Testing Global Log Debug"
}
```

### Detailed Request Log with Sub-logging
```json
{
    "time": "2024-06-28T20:00:02.089362+07:00",
    "level": "TRACE",
    "caller": "log/request.go:66",
    "processID": "oWCEjmzbdw7AMuob17wa",
    "ip": "127.0.0.1",
    "method": "GET",
    "url": "localhost:8080/",
    "statusCode": 200,
    "requestDuration": 104,
    "requestHeader": {
        "Accept": ["*/*"],
        "User-Agent": ["curl/8.6.0"]
    },
    "requestBody": {},
    "responseHeader": {
        "Content-Type": ["text/plain; charset=UTF-8"]
    },
    "responseBody": null,
    "extraData": {
        "userData": {
            "Name": "Bob",
            "Age": 29
        }
    },
    "subLog": [
        {
            "level": "[DEBUG] echo/main.go:43",
            "message": "Testing Log Request Debug"
        },
        {
            "level": "[INFO] echo/main.go:44",
            "message": "Testing Log Request Info"
        },
        {
            "level": "[WARN] echo/main.go:45",
            "message": "Testing Log Request Warn"
        },
        {
            "level": "[ERROR] echo/main.go:46",
            "message": "Testing Log Request Error"
        },
        {
            "level": "[DURATION] echo/main.go:54",
            "message": "[104.193ms] handler total process duration"
        },
        {
            "level": "[DATABASE] repository/person.go:45",
            "message": "record not found [82.751ms] [rows:0] SELECT * FROM \"person\" WHERE id = 1 ORDER BY \"person\".\"id\" LIMIT 1"
        },
    ]
}
```


## 💡 Why Use This Library?

1. **Simplified Debugging**: View all relevant logs for a single request in one place.
2. **Traceability**: Easily identify issues at any layer of your application.
3. **Flexibility**: Integrates seamlessly with popular frameworks and tools.
4. **Structured Output**: Well-formatted JSON logs for modern logging systems.



## 📘 Documentation

`WIP`


## 🙌 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you'd like to change.


## 🔒 License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). See the `LICENSE` file for details.


## ✍️ Author

**Garin Prakoso** 
[GitHub](https://github.com/gerins) | [LinkedIn](https://www.linkedin.com/in/garin-prakoso-60244b1a2/)
Feel free to contact me if you need help or have any feedback.

