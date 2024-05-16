package main

import (
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    response, err := http.Get("http://localhost:8080/path")
    if err != nil {
        log.Fatalf("Error making request: %v", err)
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatalf("Error reading response: %v", err)
    }

    log.Printf("Response from final HTTP/1 server: %s", body)
}

