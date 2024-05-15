package main

import (
    "crypto/tls"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/quic-go/quic-go/http3"
)

func main() {
    // Print the current working directory
    wd, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error getting current directory: %v", err)
    }
    fmt.Println("Current working directory:", wd)

    server := &http3.Server{
        Addr:      ":8443",
        Handler:   http.HandlerFunc(handleRequest),
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{mustLoadCertificate()},
            NextProtos:   []string{"h3"},
        },
    }

    // Start the HTTP/3 server
    fmt.Println("Starting HTTP3 Server")
    log.Fatal(server.ListenAndServeTLS("certs/server.crt", "certs/server.key"))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Received over HTTP/3: " + r.URL.Path))
}

func mustLoadCertificate() tls.Certificate {
    certFile := "certs/server.crt"
    keyFile := "certs/server.key"
    fmt.Println("Loading certificate from:", certFile)
    fmt.Println("Loading key from:", keyFile)

    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to load certificate and key from %s and %s: %v", certFile, keyFile, err)
    }
    return cert
}

