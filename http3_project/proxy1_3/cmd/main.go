// HTTP1 to HTTP3 proxy

package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"fmt"
	"net/http"

	quic "github.com/quic-go/quic-go/http3"
)

// sendHTTP3Request sends an HTTP/1 request as HTTP/3 to the specified destination server.
func sendHTTP3Request(http1Req *http.Request) ([]byte, int, error) {

	fmt.Printf("Proxy forwarding request: %s %s\n", http1Req.Method, http1Req.URL.String())
    client := &http.Client{
        Transport: &quic.RoundTripper{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For development only
        },
    }

    reqBody, _ := ioutil.ReadAll(http1Req.Body)
    req, err := http.NewRequest(http1Req.Method, "https://localhost:8443"+http1Req.URL.String(), bytes.NewReader(reqBody))
    if err != nil {
        return nil, 0, err
    }
    req.Header = http1Req.Header

    resp, err := client.Do(req)
    if err != nil {
        return nil, 0, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, 0, err
    }

	fmt.Printf("Proxy received status: %d with body length: %d\n", resp.StatusCode, len(body))
    return body, resp.StatusCode, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Proxy received request: %s %s\n", r.Method, r.URL.Path)
    body, statusCode, err := sendHTTP3Request(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Printf("Proxy encountered an error: %v\n", err)
        return
    }
    w.WriteHeader(statusCode)
    w.Write(body)
    fmt.Printf("Proxy responded with status: %d\n", statusCode)
}



func main() {
	fmt.Printf("Starting the HTTP1 server")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":11095", nil))
}

