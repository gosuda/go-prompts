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
