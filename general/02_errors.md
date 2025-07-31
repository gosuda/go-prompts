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