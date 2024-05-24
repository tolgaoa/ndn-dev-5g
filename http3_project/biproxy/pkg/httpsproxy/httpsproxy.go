package httpsproxy

import (
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "net/http"
)

// StartHTTPSProxy starts the HTTPS proxy
func StartHTTPSProxy() {
    addr := ":11096"
    fmt.Println("Starting HTTPS to HTTP proxy on", addr)
    startHTTPSProxy(addr)
}

func startHTTPSProxy(addr string) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handleHTTPSRequest)

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

func handleHTTPSRequest(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Proxy received HTTPS request: %s %s\n", r.Method, r.URL)
    response, err := forwardRequestToHTTP(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Println("Forwarding HTTPS response headers and body")
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

func forwardRequestToHTTP(r *http.Request) (*http.Response, error) {
    client := &http.Client{}

    // Construct the full URL for forwarding the request
    targetURL := fmt.Sprintf("http://%s%s", r.Host, r.RequestURI)
    log.Printf("Forwarding HTTPS request to URL: %s with method: %s\n", targetURL, r.Method)

    req, err := http.NewRequest(r.Method, targetURL, r.Body)
    if err != nil {
        log.Printf("Error creating HTTP request: %v\n", err)
        return nil, err
    }
    req.Header = r.Header

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending HTTP request: %v\n", err)
        return nil, err
    }
    return resp, nil
}

