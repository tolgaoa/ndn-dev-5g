package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	quic "github.com/quic-go/quic-go/http3"
)

// sendHTTP3Request sends an HTTP/1 request as HTTP/3 to the specified destination server.
func sendHTTP3Request(http1Req *http.Request) ([]byte, int, error) {
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

    return body, resp.StatusCode, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    body, statusCode, err := sendHTTP3Request(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(statusCode)
    _, err = w.Write(body)
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}



func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

