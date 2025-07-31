### Go Project Structure

You MUST strictly follow the directory structure outlined below to ensure a consistent, scalable, and maintainable codebase.

#### 1. Core Principles
Separation of Concerns: Clearly separate the public-facing API from the internal implementation. This is enforced by the internal directory.

Dependency Rule: Dependencies must always flow inward. The internal implementation must not depend on the publicly exposed packages.

Explicit API: Packages located at the project root are the explicit public API of this project.

#### 2. Directory Structure
The following is a standard directory structure for a server application. Use this structure as a baseline and adapt as necessary.

```plaintext
/
├── go.mod
├── go.sum
├── cmd/
│   └── <app-name>/
│       └── main.go         # Application entry point
├── internal/
│   ├── types/              # Shared data structures, enums, AND ALL INTERFACES for contracts
│   ├── server/             # Server logic (HTTP routing, handlers, middleware)
│   ├── client/             # External service API client implementations
│   ├── core/               # Core business logic and domain models, consumes interfaces from /types
│   └── persistence/        # Concrete implementations of persistence interfaces defined in /types
├── <pkg-name>/             # Publicly exposed package (e.g., user, order)
│   └── <pkg-name>.go       # Public interfaces and constructor functions
└── ... (other public packages)
```

#### 3. Package-Specific Roles

### `/cmd/<app-name>`

**Role**: The entry point for an executable application.

**Responsibilities**:

*   Load configuration.
*   Initialize dependencies, such as database connections.
*   Assemble components from the internal packages.
*   Start the server or execute the command.

### `/<pkg-name>` (Public Root Package)

**Role**: The project's public API. It can be imported by other projects.

**Responsibilities**:

*   Define the interfaces and types to be exposed externally.
*   Provide constructor functions that create instances of the internal implementation and return them as the defined interface type.
*   Consumers of this package should not need to know about the existence or structure of the `internal` directory.

### `/internal` (Internal Implementation Package)

**Role**: Contains all internal implementation for the project. The Go compiler prevents it from being imported by external projects.

#### Sub-package Responsibilities:

*   **`/types`**: This package is the central repository for **all shared data structures (structs, enums) AND all interfaces** that define contracts between different `internal` packages. This includes interfaces for persistence, services, external clients, etc. By centralizing interfaces here, it strictly enforces the Dependency Inversion Principle, ensuring `core` and other packages depend on abstractions, and prevents circular imports. Any type or interface that needs to be consumed or implemented by more than one other `internal` package **MUST** be defined here.

*   **`/server`**: Contains everything related to the HTTP server. It is responsible for router setup, middleware, and HTTP request/response handling logic (handlers). Handlers call the business logic in the core layer, consuming interfaces as defined in `/types`.

*   **`/client`**: Contains client code for communicating with external services that this project depends on (e.g., other microservices, third-party APIs). If these clients expose interfaces for contract, those interfaces should also be defined in `/types`.

*   **`/core`**: This package contains the application's most critical business logic and domain models. It **consumes interfaces defined in `/internal/types`** (e.g., `persistence.UserRepository`, `client.NotificationSender`) and is responsible for orchestrating business workflows. It receives concrete implementations of these interfaces via Dependency Injection. This package should be a pure collection of business rules, independent of specific technologies (like HTTP, SQL, or concrete caching implementations).

*   **`/persistence`**: This package is dedicated to providing **concrete implementations of the persistence interfaces defined in `/internal/types`**. It manages specific database connections (e.g., SQL, NoSQL) and implements CRUD operations. It is also responsible for providing mock implementations for testing purposes. It must only import interfaces from the `/internal/types` package and relevant database drivers, ensuring a clean and unidirectional dependency flow.

### 4. Dependency Rules

*   The direction of dependency flows from `cmd` → root packages → `internal`.

*   Within `internal`, dependencies must strictly flow from outer layers to inner layers, adhering to the Dependency Inversion Principle and centralizing contracts in `/types`:
    *   `server` can depend on `core`.
    *   `core` depends **only** on `types` (for interfaces and shared data structures).
    *   `persistence` depends **only** on `types` (for interfaces it implements and shared data structures).
    *   `client` depends **only** on `types` (for interfaces it implements/consumes and shared data structures, if applicable).
    *   `cmd` is responsible for initializing concrete implementations (from `persistence`, `client`) and injecting them into `core` and `server` through the interfaces defined in `types`.
    *   `core` must never depend on `server`, `client`, or concrete implementations of `persistence`. It only interacts with abstractions (interfaces) defined in `types`.

*   **No Circular Dependencies**: This structure prevents circular dependencies by ensuring that all contracts (interfaces) are defined in a single, well-defined, and universally referencable `/types` package. All other `internal` packages (`core`, `persistence`, `server`, `client`) depend exclusively on `/types` for their abstract definitions, and concrete implementations are injected from `cmd`.

*   **No Circular Dependencies**: A structure where an `internal` package imports a public root package is not allowed.

### 5. Additional Guidelines

*   **No `pkg` Directory**: Do not create a `pkg` directory at the project root. It adds unnecessary nesting, and the distinction between public and private packages is already clear based on their location in the root or `internal` directory.
*   **Interface-Driven Design**: The public API must be provided through interfaces. This is key to facilitating Dependency Injection (DI) and writing testable code.
*   **Clear Package Names**: Package names should clearly describe their roles. Avoid vague names like `util` or `common`.
### **Defining and Using Sentinel Errors**

In Go, errors are values. For common, predictable error conditions, you should define package-level "sentinel errors." These are pre-declared error values that functions can return to signal a specific, well-known state.

This pattern allows calling code to programmatically check for and handle these specific errors, making your application more robust and reliable.

### Guiding Principles

1.  **Define Sentinel Errors as Package-Level Variables:** Use the `errors.New` function to create package-level variables for each distinct error condition. This ensures there is a single, consistent value for each error that callers can reference.

2.  **Follow Naming and Formatting Conventions:**
    *   **Naming:** Start the variable name with `Err` (e.g., `ErrNotFound`). This is a strong, widely-followed convention in the Go community.
    *   **Message:** Prefix the error message with the package name (e.g., `"database: item not found"`). This provides immediate context in logs and error chains, clarifying where the error originated.
    *   **Documentation:** Add a doc comment to each error variable explaining the condition under which it is returned.

3.  **Use `errors.Is` for Checking:** Callers should use the `errors.Is` function to check if a returned error matches a specific sentinel error. This is the idiomatic and most robust method, as it correctly handles wrapped errors.

### Example Implementation

Here is how you would define and use sentinel errors in a `database` package.

**1. Define the Errors in `database/database.go`**

```go
// Package database provides functions for interacting with the data store.
package database

import "errors"

var (
	// ErrNotFound is returned when a requested record is not found.
	ErrNotFound = errors.New("database: record not found")

	// ErrDuplicateKey is returned when attempting to insert a record
	// with a primary key that already exists.
	ErrDuplicateKey = errors.New("database: duplicate key")

    // ErrInvalidID is returned when a provided ID is malformed or empty.
    ErrInvalidID = errors.New("database: invalid ID")
)

// GetUser retrieves a user by their ID.
func GetUser(id string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	// ... logic to query the database for the user ...

	if userWasNotFoundInDB {
		// Return the predefined sentinel error.
		return nil, ErrNotFound
	}
    
	return &user, nil
}
```

**2. Handle the Errors in `main.go`**

The true power of sentinel errors is realized in the calling code, which can now react intelligently to different failure modes.

```go
package main

import (
	"database"
	"errors"
	"fmt"
	"log"
)

func processUser(userID string) {
	user, err := database.GetUser(userID)
	if err != nil {
		// Use errors.Is to check for a specific, recoverable error.
		if errors.Is(err, database.ErrNotFound) {
			fmt.Printf("User with ID '%s' not found. Let's create a new profile for them.\n", userID)
			// ... logic to create a new user ...
			return
		}

		// Handle other known errors if necessary.
		if errors.Is(err, database.ErrInvalidID) {
			log.Printf("Error: The provided user ID '%s' is invalid.\n", userID)
			return
		}

		// For all other, unexpected errors, log fatal.
		log.Fatalf("An unexpected error occurred while processing user '%s': %v", userID, err)
	}

	fmt.Printf("Successfully processed user: %s\n", user.Name)
}

func main() {
    // This call will be handled gracefully.
	processUser("user-123") // Assuming this user doesn't exist.

    // This call will also be handled.
    processUser("")
}
```
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

### httprouter

When you start a new HTTP API project, it's strongly recommended to use `httprouter` (github.com/julienschmidt/httprouter) as your router.

This router is fast, lightweight, and provides a clean API for defining routes.

Here's an example of how to use `httprouter`:

```go
package main

import (
    "fmt"
    "net"
    "net/http"
    "github.com/rs/zerolog/log"

    "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)

    server := &http.Server{
        Handler: router,
    }

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to listen")
    }
    defer ln.Close()

    go func() {
        log.Info().Msg("Server started")
        err := http.Serve(ln, server)
        if err != nil && err != http.ErrServerClosed {
            log.Fatal().Err(err).Msg("Failed to serve HTTP")
        }
        log.Info().Msg("Server stopped")
    }()

    stopHTTPCh := make(chan struct{})
    // some logic to handle graceful shutdown, e.g., signal handling
    <-stopHTTPCh // Block forever or implement graceful shutdown
    log.Info().Msg("Shutting down server")
    err := server.Close()
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to close server")
    }
    log.Info().Msg("Server gracefully stopped")
}
```

