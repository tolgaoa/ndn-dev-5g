package main

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "net"

    quic "github.com/quic-go/quic-go/http3"
)

func main() {
    go startHTTP1toHTTP3Proxy()  // Start listening for HTTP/1 and forwarding as HTTP/3
    startHTTP3toHTTP1Proxy()     // Start listening for HTTP/3 and forwarding as HTTP/1
}

func startHTTP1toHTTP3Proxy() {
    fmt.Println("Starting HTTP/1 to HTTP/3 proxy on :11095")
    http.HandleFunc("/", handleHTTP1toHTTP3Request)
    log.Fatal(http.ListenAndServe(":11095", nil))
}

func handleHTTP1toHTTP3Request(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Proxy received HTTP/1 request: %s %s\n", r.Method, r.URL)
    body, statusCode, err := sendHTTP3Request(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(statusCode)
    w.Write(body)
}

func sendHTTP3Request(http1Req *http.Request) ([]byte, int, error) {
    // Create an HTTP client that supports HTTP/3 via QUIC
    client := &http.Client{
        Transport: &quic.RoundTripper{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For development only
        },
    }

    // Construct the URL for the HTTP/3 request
    targetURL := constructTargetURL(http1Req, true)

    log.Printf("Sending HTTP/3 request to URL: %s with method: %s\n", targetURL, http1Req.Method)

    // Read the request body
    reqBody, err := ioutil.ReadAll(http1Req.Body)
    if err != nil {
        log.Printf("Error reading request body: %v\n", err)
        return nil, 0, err
    }

    // Create a new request with the same parameters, ensuring we use a full URL
    req, err := http.NewRequest(http1Req.Method, targetURL, bytes.NewReader(reqBody))
    if err != nil {
        log.Printf("Error creating new HTTP/3 request: %v\n", err)
        return nil, 0, err
    }
    req.Header = http1Req.Header

    log.Println("Forwarding headers:")
    for name, headers := range req.Header {
        for _, h := range headers {
            log.Printf("%v: %v\n", name, h)
        }
    }

    // Perform the request using HTTP/3
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending HTTP/3 request: %v\n", err)
        return nil, 0, err
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v\n", err)
        return nil, 0, err
    }

    log.Printf("Received HTTP/3 response with status code: %d, body length: %d\n", resp.StatusCode, len(body))
    return body, resp.StatusCode, nil
}

func constructTargetURL(req *http.Request, toHTTP3 bool) string {
    // Ensure the URL has the appropriate scheme and port
    u := *req.URL
    if toHTTP3 {
        u.Scheme = "https"
        host, port, err := net.SplitHostPort(req.Host)
        if err != nil {
            host = req.Host
            port = "80"
        }
        if port == "80" {
            port = "8443"
        }
        u.Host = fmt.Sprintf("%s:%s", host, port)
    } else {
        u.Scheme = "http"
        host, port, err := net.SplitHostPort(req.Host)
        if err != nil {
            host = req.Host
            port = "8443"
        }
        if port == "8443" {
            port = "80"
        }
        u.Host = fmt.Sprintf("%s:%s", host, port)
    }
    return u.String()
}

func startHTTP3toHTTP1Proxy() {
    fmt.Println("Starting HTTP/3 to HTTP/1 proxy on :11096")
    server := &quic.Server{
        Addr:      ":11096",
        Handler:   http.HandlerFunc(handleHTTP3toHTTP1Request),
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{mustLoadCertificate()},
            NextProtos:   []string{"h3"},
        },
    }
    log.Fatal(server.ListenAndServeTLS("certs/server.crt", "certs/server.key"))
}

func handleHTTP3toHTTP1Request(w http.ResponseWriter, r *http.Request) {
    log.Printf("Proxy received HTTP/3 request for: %s %s\n", r.Method, r.URL)
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

    // Construct a fully qualified URL for forwarding the request
    forwardedURL := constructTargetURL(r, false)
    log.Printf("Forwarding HTTP/3 request to HTTP/1: %s %s\n", r.Method, forwardedURL)

    reqBody, _ := ioutil.ReadAll(r.Body)
    req, err := http.NewRequest(r.Method, forwardedURL, bytes.NewReader(reqBody))
    if err != nil {
        log.Printf("Error creating HTTP/1 request: %v\n", err)
        return nil, err
    }
    req.Header = r.Header

    log.Println("Forwarding HTTP/1 request headers:")
    for name, headers := range req.Header {
        for _, h := range headers {
            log.Printf("%v: %v\n", name, h)
        }
    }

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

