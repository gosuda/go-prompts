### Go Project Structure Guidelines

This document defines the guidelines for structuring a Go project. Adherence to these rules is required to build a consistent, scalable, and maintainable codebase.

#### 1. Core Principles
Separation of Concerns: Clearly separate the public-facing API from the internal implementation. This is enforced by the internal directory.

Dependency Rule: Dependencies must always flow inward. The internal implementation must not depend on the publicly exposed packages.

Explicit API: Packages located at the project root are the explicit public API of this project.

#### 2. Directory Structure
The following is a standard directory structure for a server application. Use this structure as a baseline and adapt as necessary.

/
├── go.mod
├── go.sum
│
├── cmd/
│   └── <app-name>/
│       └── main.go         // Application entry point
│
├── internal/
│   ├── types/              // Shared data structures and types
│   ├── server/             // Server logic (HTTP routing, handlers, middleware)
│   ├── client/             // External service API client implementations
│   ├── core/               // Core business logic and domain models
│   │   ├── data/           // Core data structures and models
│   │   ├── rtmp/           // RTMP-related business logic
│   │   ├── webhooks/       // Webhook-related business logic
│   │   └── cache/          // Caching-related interfaces and logic
│   └── persistence/        // Data persistence layer (DB interaction)
│
├── <pkg-name>/             // Publicly exposed package (e.g., user, order)
│   └── <pkg-name>.go       // Public interfaces and constructor functions
│
└── ... (other public packages)

#### 3. Package-Specific Roles

/cmd/<app-name>

Role: The entry point for an executable application.

Responsibilities:

Load configuration.

Initialize dependencies, such as database connections.

Assemble components from the internal packages.

Start the server or execute the command.

/<pkg-name> (Public Root Package)

Role: The project's public API. It can be imported by other projects.

Responsibilities:

Define the interfaces and types to be exposed externally.

Provide constructor functions that create instances of the internal implementation and return them as the defined interface type.

Consumers of this package should not need to know about the existence or structure of the internal directory.

/internal (Internal Implementation Package)

Role: Contains all internal implementation for the project. The Go compiler prevents it from being imported by external projects.

#### Sub-package Responsibilities:

/types: Contains major struct definitions and types shared across different internal packages (e.g., server, core, persistence). This helps avoid circular dependencies by providing a common, neutral location for data structures.

/server: Contains everything related to the HTTP server. It is responsible for router setup, middleware, and HTTP request/response handling logic (handlers). Handlers call the business logic in the core layer.

/client: Contains client code for communicating with external services that this project depends on (e.g., other microservices, third-party APIs).

/core: Contains the application's most critical business logic and domain models. This package should be a pure collection of business rules, independent of specific technologies (like HTTP or SQL).

/persistence: Manages database connections and provides data persistence logic. It implements operations (e.g., CRUD) using the data structures defined in /internal/types. It is also responsible for providing mock implementations for testing purposes. This package does not depend on /core.

#### 4. Dependency Rules

The direction of dependency flows from cmd → root packages → internal.

Within internal, dependencies should point from outer layers to inner layers:

server can depend on core.

core can depend on types.

persistence can depend on types.

core must never depend on server or persistence.

No Circular Dependencies: A structure where an internal package imports a public root package is not allowed.

#### 5. Additional Guidelines

No pkg Directory: Do not create a pkg directory at the project root. It adds unnecessary nesting, and the distinction between public and private packages is already clear based on their location in the root or internal directory.

Interface-Driven Design: The public API must be provided through interfaces. This is key to facilitating Dependency Injection (DI) and writing testable code.

Clear Package Names: Package names should clearly describe their roles. Avoid vague names like util or common.