//server.go
package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "9001" // fallback
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from backend on port %s!\n", port)
    })

    fmt.Printf("Starting backend on port %s...\n", port)
    http.ListenAndServe(":"+port, nil)
}
