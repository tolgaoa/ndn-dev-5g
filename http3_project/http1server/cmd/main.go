package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from HTTP/1 server at path: %s", r.URL.Path)
        log.Printf("Handled request on path: %s", r.URL.Path)
    })

    log.Println("Starting HTTP/1 server on :8082...")
    log.Fatal(http.ListenAndServe(":8082", nil))
}

