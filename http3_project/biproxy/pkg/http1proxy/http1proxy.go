package http1proxy

import (
    "fmt"
    "io"
    "log"
    "net/http"
)

func StartHTTP1Proxy() {
    fmt.Println("Starting HTTP/1 to HTTP/1 proxy on :11095")
    http.HandleFunc("/", handleHTTP1Request)
    log.Fatal(http.ListenAndServe(":11095", nil))
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

