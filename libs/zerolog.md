### **Go Logging Standard: zerolog**

As an expert Go developer, you will adhere to the following standards for logging in all Go code you generate.

#### **Core Mandates**

1.  **Mandatory Library:** `zerolog` (`github.com/rs/zerolog`) is the **exclusive** logging library. Do not use any other logging package, including the standard library's `log` or `fmt`, unless explicitly instructed.
2.  **Structured First:** All logs must be structured with key-value pairs. Use methods like `.Str()`, `.Int()`, and `.Err()` to add context. A final `.Msg()` call provides the human-readable message.
3.  **Correct Error Logging:** When logging an `error` type, **always** use the `.Err(err)` method. This attaches the error to a dedicated `error` field, which is standard practice for structured logging.
4.  **Contextual Loggers:** Prefer creating contextual logger instances over using the global logger, especially within specific components or request handlers.

---

### **1. Global Configuration**

Global configuration should be set once at the start of the `main` function. During development, use `zerolog.ConsoleWriter` for human-readable, color-coded output.

```go
// main.go
package main

import (
    "os"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // For development, use a pretty, color-coded console output.
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    // Set the global log level. Logs with a level of Debug or higher will be written.
    zerolog.SetGlobalLevel(zerolog.DebugLevel)
    
    // --- Application logic starts here ---
    log.Info().Msg("Application starting up.")
    doSomething("my-request-id")
}
```

---

### **2. Logging Levels**

Use the appropriate level for each message. The available levels, from highest to lowest priority, are:

| Level   | Constant              | Value | Use Case                                           |
| :------ | :-------------------- | :---: | :------------------------------------------------- |
| `panic` | `zerolog.PanicLevel`  |   5   | Logs the message and then calls `panic()`.           |
| `fatal` | `zerolog.FatalLevel`  |   4   | Logs the message and then calls `os.Exit(1)`.        |
| `error` | `zerolog.ErrorLevel`  |   3   | For significant errors that require attention.     |
| `warn`  | `zerolog.WarnLevel`   |   2   | For potential issues that don't break functionality. |
| `info`  | `zerolog.InfoLevel`   |   1   | For informational messages about application state.  |
| `debug` | `zerolog.DebugLevel`  |   0   | For detailed debugging information.                |
| `trace` | `zerolog.TraceLevel`  |  -1   | For extremely granular, verbose tracing.           |

To disable logging, use `zerolog.Disabled`.

---

### **3. Usage and Examples**

#### **Basic Structured Logging**

Add contextual fields and finish the chain with a message.

```go
import "github.com/rs/zerolog/log"

log.Debug().
    Str("Scale", "833 cents").
    Float64("Interval", 833.09).
    Msg("Fibonacci is everywhere")
```
**JSON Output:**
`{"level":"debug","Scale":"833 cents","Interval":833.09,"time":"2023-10-27T10:30:00Z","message":"Fibonacci is everywhere"}`

#### **Logging Errors (The Right Way)**

Always use the `.Err()` method to log `error` variables. This ensures they are serialized correctly into a dedicated field.

```go
import "errors"

// ...

err := errors.New("failed to connect to database")
if err != nil {
    log.Error().
        Err(err). // <-- CORRECT: Use the .Err() method
        Str("component", "database").
        Msg("A critical error occurred")
}
```
**JSON Output:**
`{"level":"error","error":"failed to connect to database","component":"database","time":"2023-10-27T10:30:00Z","message":"A critical error occurred"}`

**INCORRECT USAGE (DO NOT DO THIS):**
`log.Error().Msgf("Error connecting to database: %v", err)`
*This loses the structured `error` field, making logs harder to parse and query.*

#### **Creating Contextual Sub-loggers**

Create a sub-logger with pre-populated fields to maintain context across multiple log entries, such as within a request's lifecycle. You can pass this logger via `context.Context` to downstream functions.

```go
// Add a logger with a "component" field to the context.
ctx := log.With().Str("component", "module").Logger().WithContext(ctx)

// Retrieve the logger from the context and use it.
log.Ctx(ctx).Info().Msg("hello world")

// Output: {"level":"info","component":"module","message":"hello world"}
```

#### **Redirecting the Standard Logger**

To redirect output from Go's standard `log` package to `zerolog`, use `stdlog.SetOutput()`. This ensures that logs from third-party libraries using the standard logger are also captured and structured consistently.

```go
import (
    stdlog "log"
    "github.com/rs/zerolog/log"
)

// Redirect standard library log messages to zerolog.
stdlog.SetFlags(0)
stdlog.SetOutput(log.Logger)

// Now, calls to the standard logger will be formatted by zerolog.
stdlog.Print("hello world")

// Output: {"level":"info","message":"hello world"}
```
