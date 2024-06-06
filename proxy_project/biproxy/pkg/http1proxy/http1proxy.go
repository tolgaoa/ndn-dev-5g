package http1proxy

import (
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

// StartHTTP1Proxy starts the HTTP/1.1 proxy, with optional TLS support
func StartHTTP1Proxy() {
    useTLS := os.Getenv("USE_TLS") == "true"
    addr := ":11095"

    if useTLS {
        fmt.Println("Starting HTTP/1 to HTTP/1 proxy on", addr, "with TLS")
        startHTTP1ProxyWithTLS(addr)
    } else {
        fmt.Println("Starting HTTP/1 to HTTP/1 proxy on", addr)
        startHTTP1Proxy(addr)
    }
}

func startHTTP1Proxy(addr string) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handleHTTP1Request)
    server := &http.Server{
        Addr:    addr,
        Handler: mux,
    }
    log.Fatal(server.ListenAndServe())
}

func startHTTP1ProxyWithTLS(addr string) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handleHTTP1Request)

    // Load TLS certificates
    certFile := "certs/server.crt"
    keyFile := "certs/server.key"
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to load TLS certificates: %v", err)
    }

    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
    }

    server := &http.Server{
        Addr:      addr,
        Handler:   mux,
        TLSConfig: tlsConfig,
    }

    log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}

func handleHTTP1Request(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Proxy received HTTP/1 request: %s %s\n", r.Method, r.URL)
    response, err := forwardRequestToHTTP1(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Println("Forwarding HTTP/1 response headers and body")
    for key, values := range response.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(response.StatusCode)
    _, err = io.Copy(w, response.Body)
    if err != nil {
        log.Printf("Error copying response body: %v\n", err)
    }
}

func forwardRequestToHTTP1(r *http.Request) (*http.Response, error) {
    client := &http.Client{}

    // Construct the full URL for forwarding the request
    targetURL := fmt.Sprintf("http://%s%s", r.Host, r.RequestURI)
    log.Printf("Forwarding HTTP/1 request to URL: %s with method: %s\n", targetURL, r.Method)

    req, err := http.NewRequest(r.Method, targetURL, r.Body)
    if err != nil {
        log.Printf("Error creating HTTP/1 request: %v\n", err)
        return nil, err
    }
    req.Header = r.Header

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending HTTP/1 request: %v\n", err)
        return nil, err
    }
    return resp, nil
}

