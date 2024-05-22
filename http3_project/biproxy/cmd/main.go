package main

import (
	h3 "biproxy/pkg/http3proxy"
)

func main() {
    go h3.StartHTTP1toHTTP3Proxy()  // Start listening for HTTP/1 and forwarding as HTTP/3
    h3.StartHTTP3toHTTP1Proxy()     // Start listening for HTTP/3 and forwarding as HTTP/1
}

