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