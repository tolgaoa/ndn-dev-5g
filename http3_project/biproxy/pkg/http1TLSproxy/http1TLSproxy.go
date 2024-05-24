package http1TLSproxy

import (
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "net/http"
)

// StartHTTP1Proxy starts the HTTP/1.1 proxy
func StartHTTP1Proxy() {
    addr := ":11095"
    fmt.Println("Starting HTTP/1 to HTTP/1 (with TLS) proxy on", addr)
    startHTTP1Proxy(addr)
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

func handleHTTP1Request(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Proxy received HTTP/1 request: %s %s\n", r.Method, r.URL)
    response, err := forwardRequestToHTTP1WithTLS(r)
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

func forwardRequestToHTTP1WithTLS(r *http.Request) (*http.Response, error) {
    // Create an HTTP client that supports HTTPS
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For development only
        },
    }

    // Construct the full URL for forwarding the request
    targetURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
    log.Printf("Forwarding HTTP/1 request to URL: %s with method: %s\n", targetURL, r.Method)

    req, err := http.NewRequest(r.Method, targetURL, r.Body)
    if err != nil {
        log.Printf("Error creating HTTPS request: %v\n", err)
        return nil, err
    }
    req.Header = r.Header

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending HTTPS request: %v\n", err)
        return nil, err
    }
    return resp, nil
}

