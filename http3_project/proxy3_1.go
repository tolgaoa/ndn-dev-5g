package main

import (
    "bytes"
    "crypto/tls"
    "io"
    "io/ioutil"
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
    log.Println("Current working directory:", wd)

    server := &http3.Server{
        Addr:      ":8443",
        Handler:   http.HandlerFunc(forwardToHTTP1),
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{mustLoadCertificate()},
            NextProtos:   []string{"h3"},
        },
    }

    // Start the HTTP/3 server
    log.Println("Starting HTTP3 Server as Proxy to HTTP/1")
    log.Fatal(server.ListenAndServeTLS("certs/server.crt", "certs/server.key"))
}

func forwardToHTTP1(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received HTTP/3 request for: %s %s", r.Method, r.URL.Path)

    response, err := forwardRequest(r)
    if err != nil {
        http.Error(w, "Failed to forward request: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Copy the response headers
    for key, values := range response.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(response.StatusCode)
    
    // Correct handling of the Body for writing to the response
    if response.Body != nil {
        _, err = io.Copy(w, response.Body)
        if err != nil {
            log.Printf("Failed to write response body: %v", err)
        }
    }
    log.Printf("Forwarded to HTTP/1 and responded with: %d", response.StatusCode)
}

func forwardRequest(r *http.Request) (*http.Response, error) {
    client := &http.Client{}
    reqBody, _ := ioutil.ReadAll(r.Body)
    req, err := http.NewRequest(r.Method, "http://localhost:8082"+r.URL.String(), bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }
    req.Header = r.Header

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

func mustLoadCertificate() tls.Certificate {
    certFile := "certs/server.crt"
    keyFile := "certs/server.key"
    log.Println("Loading certificate from:", certFile)
    log.Println("Loading key from:", keyFile)

    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to load certificate and key from %s and %s: %v", certFile, keyFile, err)
    }
    return cert
}

