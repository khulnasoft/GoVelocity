---
id: logger
---

# Logger

Logger middleware for [Velocity](https://github.com/khulnasoft/velocity) that logs HTTP request/response details.

## Signatures

```go
func New(config ...Config) velocity.Handler
```

## Examples

Import the middleware package that is part of the Velocity web framework

```go
import (
    "github.com/khulnasoft/velocity"
    "github.com/khulnasoft/velocity/middleware/logger"
)
```

:::tip
The order of registration plays a role. Only all routes that are registered after this one will be logged.
The middleware should therefore be one of the first to be registered.
:::

After you initiate your Velocity app, you can use the following possibilities:

```go
// Initialize default config
app.Use(logger.New())

// Or extend your config for customization
// Logging remote IP and Port
app.Use(logger.New(logger.Config{
    Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
}))

// Logging Request ID
app.Use(requestid.New())
app.Use(logger.New(logger.Config{
    // For more options, see the Config section
    Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
}))

// Changing TimeZone & TimeFormat
app.Use(logger.New(logger.Config{
    Format:     "${pid} ${status} - ${method} ${path}\n",
    TimeFormat: "02-Jan-2006",
    TimeZone:   "America/New_York",
}))

// Custom File Writer
file, err := os.OpenFile("./123.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
if err != nil {
    log.Fatalf("error opening file: %v", err)
}
defer file.Close()
app.Use(logger.New(logger.Config{
    Output: file,
}))

// Add Custom Tags
app.Use(logger.New(logger.Config{
    CustomTags: map[string]logger.LogFunc{
        "custom_tag": func(output logger.Buffer, c velocity.Ctx, data *logger.Data, extraParam string) (int, error) {
            return output.WriteString("it is a custom tag")
        },
    },
}))

// Callback after log is written
app.Use(logger.New(logger.Config{
    TimeFormat: time.RFC3339Nano,
    TimeZone:   "Asia/Shanghai",
    Done: func(c velocity.Ctx, logString []byte) {
        if c.Response().StatusCode() != velocity.StatusOK {
            reporter.SendToSlack(logString) 
        }
    },
}))

// Disable colors when outputting to default format
app.Use(logger.New(logger.Config{
    DisableColors: true,
}))
```

### Use Logger Middleware with Other Loggers

In order to use Velocity logger middleware with other loggers such as zerolog, zap, logrus; you can use `LoggerToWriter` helper which converts Velocity logger to a writer, which is compatible with the middleware.

```go
package main

import (
    "github.com/khulnasoft/contrib/velocityzap/v2"
    "github.com/khulnasoft/velocity"
    "github.com/khulnasoft/velocity/log"
    "github.com/khulnasoft/velocity/middleware/logger"
)

func main() {
    // Create a new Velocity instance
    app := velocity.New()

    // Create a new zap logger which is compatible with Velocity AllLogger interface
    zap := velocityzap.NewLogger(velocityzap.LoggerConfig{
        ExtraKeys: []string{"request_id"},
    })

    // Use the logger middleware with zerolog logger
    app.Use(logger.New(logger.Config{
        Output: logger.LoggerToWriter(zap, log.LevelDebug),
    }))

    // Define a route
    app.Get("/", func(c velocity.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Start server on http://localhost:3000
    app.Listen(":3000")
}
```

:::tip
Writing to os.File is goroutine-safe, but if you are using a custom Output that is not goroutine-safe, make sure to implement locking to properly serialize writes.
:::

## Config

### Config

| Property         | Type                       | Description                                                                                                                      | Default                                                               |
|:-----------------|:---------------------------|:---------------------------------------------------------------------------------------------------------------------------------|:----------------------------------------------------------------------|
| Next             | `func(velocity.Ctx) bool`    | Next defines a function to skip this middleware when returned true.                                                              | `nil`                                                                 |
| Done             | `func(velocity.Ctx, []byte)` | Done is a function that is called after the log string for a request is written to Output, and pass the log string as parameter. | `nil`                                                                 |
| CustomTags       | `map[string]LogFunc`       | tagFunctions defines the custom tag action.                                                                                      | `map[string]LogFunc`                                                  |
| Format           | `string`                   | Format defines the logging tags.                                                                                                 | `[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n` |
| TimeFormat       | `string`                   | TimeFormat defines the time format for log timestamps.                                                                           | `15:04:05`                                                            |
| TimeZone         | `string`                   | TimeZone can be specified, such as "UTC" and "America/New_York" and "Asia/Chongqing", etc                                        | `"Local"`                                                             |
| TimeInterval     | `time.Duration`            | TimeInterval is the delay before the timestamp is updated.                                                                       | `500 * time.Millisecond`                                              |
| Output           | `io.Writer`                | Output is a writer where logs are written.                                                                                       | `os.Stdout`                                                           |
| LoggerFunc | `func(c velocity.Ctx, data *Data, cfg Config) error` | Custom logger function for integration with logging libraries (Zerolog, Zap, Logrus, etc). Defaults to Velocity's default logger if not defined. | `see default_logger.go defaultLoggerInstance` |
| DisableColors    | `bool`                     | DisableColors defines if the logs output should be colorized.                                                                    | `false`                                                               |
| enableColors     | `bool`                     | Internal field for enabling colors in the log output. (This is not a user-configurable field)                                    | -                                                                     |
| enableLatency    | `bool`                     | Internal field for enabling latency measurement in logs. (This is not a user-configurable field)                                 | -                                                                     |
| timeZoneLocation | `*time.Location`           | Internal field for the time zone location. (This is not a user-configurable field)                                               | -                                                                     |

## Default Config

```go
var ConfigDefault = Config{
    Next:          nil,
    Done:          nil,
    Format:        "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
    TimeFormat:    "15:04:05",
    TimeZone:      "Local",
    TimeInterval:  500 * time.Millisecond,
    Output:        os.Stdout,
    DisableColors: false,
    LoggerFunc:    defaultLoggerInstance,
}
```

## Constants

```go
// Logger variables
const (
    TagPid               = "pid"
    TagTime              = "time"
    TagReferer           = "referer"
    TagProtocol          = "protocol"
    TagPort              = "port"
    TagIP                = "ip"
    TagIPs               = "ips"
    TagHost              = "host"
    TagMethod            = "method"
    TagPath              = "path"
    TagURL               = "url"
    TagUA                = "ua"
    TagLatency           = "latency"
    TagStatus            = "status"         // response status
    TagResBody           = "resBody"        // response body
    TagReqHeaders        = "reqHeaders"
    TagQueryStringParams = "queryParams"    // request query parameters
    TagBody              = "body"           // request body
    TagBytesSent         = "bytesSent"
    TagBytesReceived     = "bytesReceived"
    TagRoute             = "route"
    TagError             = "error"
    // DEPRECATED: Use TagReqHeader instead
    TagHeader            = "header:"        // request header
    TagReqHeader         = "reqHeader:"     // request header
    TagRespHeader        = "respHeader:"    // response header
    TagQuery             = "query:"         // request query
    TagForm              = "form:"          // request form
    TagCookie            = "cookie:"        // request cookie
    TagLocals            = "locals:"
    // colors
    TagBlack             = "black"
    TagRed               = "red"
    TagGreen             = "green"
    TagYellow            = "yellow"
    TagBlue              = "blue"
    TagMagenta           = "magenta"
    TagCyan              = "cyan"
    TagWhite             = "white"
    TagReset             = "reset"
)
```
