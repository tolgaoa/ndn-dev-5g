package http2proxy

import (
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "net/http"

    "golang.org/x/net/http2"
)

func StartHTTP1toHTTP2Proxy() {
    fmt.Println("Starting HTTP/1 to HTTP/2 proxy on :11095")
    http.HandleFunc("/", handleHTTP1toHTTP2Request)
    log.Fatal(http.ListenAndServe(":11095", nil))
}

func handleHTTP1toHTTP2Request(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Proxy received HTTP/1 request: %s %s\n", r.Method, r.URL)
    body, statusCode, err := sendHTTP2Request(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(statusCode)
    w.Write(body)
}

func sendHTTP2Request(http1Req *http.Request) ([]byte, int, error) {
    // Create an HTTP client that supports HTTP/2
    client := &http.Client{
        Transport: &http2.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For development only
        },
    }

    targetURL := fmt.Sprintf("https://%s%s", http1Req.Host, http1Req.RequestURI)
    req, err := http.NewRequest(http1Req.Method, targetURL, http1Req.Body)
    if err != nil {
        log.Printf("Error creating new HTTP/2 request: %v\n", err)
        return nil, 0, err
    }
    req.Header = http1Req.Header

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending HTTP/2 request: %v\n", err)
        return nil, 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v\n", err)
        return nil, 0, err
    }

    return body, resp.StatusCode, nil
}

func StartHTTP2toHTTP1Proxy() {
    fmt.Println("Starting HTTP/2 to HTTP/1 proxy on :11096")
    
    // Create an HTTPS server with HTTP/2 support
    srv := &http.Server{
        Addr:    ":11096",
        Handler: http.HandlerFunc(handleHTTP2toHTTP1Request),
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{mustLoadCertificate()},
            NextProtos:   []string{http2.NextProtoTLS}, // Enables HTTP/2
        },
    }

    // Enable HTTP/2 support
    http2.ConfigureServer(srv, &http2.Server{})

    log.Fatal(srv.ListenAndServeTLS("certs/server.crt", "certs/server.key"))
}

func handleHTTP2toHTTP1Request(w http.ResponseWriter, r *http.Request) {
    log.Printf("Proxy received HTTP/2 request for: %s %s\n", r.Method, r.URL)
    response, err := forwardRequestToHTTP1(r)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to forward request: %v", err), http.StatusInternalServerError)
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
    targetURL := fmt.Sprintf("http://%s%s", r.Host, r.RequestURI)

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

