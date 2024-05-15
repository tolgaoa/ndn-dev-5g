package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	quic "github.com/lucas-clemente/quic-go/http3"
)

// sendHTTP3Request sends an HTTP/1 request as HTTP/3 to the specified destination server.
func sendHTTP3Request(http1Req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Transport: &quic.RoundTripper{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For development only
		},
	}

	reqBody, _ := ioutil.ReadAll(http1Req.Body)
	req, err := http.NewRequest(http1Req.Method, "https://localhost:8443"+http1Req.URL.String(), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header = http1Req.Header

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: resp.StatusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	response, err := sendHTTP3Request(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	_, err = w.Write(response.Body)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

