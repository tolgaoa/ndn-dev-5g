package main

import (
	// built-in
	"strings"
	"os"

	// external
	log "github.com/sirupsen/logrus"

	// internal
	h1 "biproxy/pkg/http1proxy"
	h2 "biproxy/pkg/http2proxy"
	h3 "biproxy/pkg/http3proxy"
	ev "biproxy/utils/envProc"
)

func setLogLevel() {


    log.SetFormatter(&log.TextFormatter{
        ForceColors: true,
        TimestampFormat: "2006-01-02 15:04:05.000000",
        FullTimestamp:   true,
    })

    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" {
        logLevel = "info"
    }

    switch strings.ToLower(logLevel) {
    case "trace":
        log.SetLevel(log.TraceLevel)
    case "debug":
        log.SetLevel(log.DebugLevel)
    case "info":
        log.SetLevel(log.InfoLevel)
    case "warn":
        log.SetLevel(log.WarnLevel)
    case "error":
        log.SetLevel(log.ErrorLevel)
    case "fatal":
        log.SetLevel(log.FatalLevel)
    case "panic":
        log.SetLevel(log.PanicLevel)
    default:
        log.Warnf("Unknown log level: %s, defaulting to info", logLevel)
        log.SetLevel(log.InfoLevel)
    }
}


func main() {

	setLogLevel()

	opmode := ev.GetEnv("OPERATION_MODE","HTTP1") // select the operation mode from ENV: HTTP1, HTTP2, HTTP3 -- default to http1


	switch opmode {
	case "HTTP1":
		log.Info("Starting Proxy in HTTP1 <--> HTTP1 Forwarding Mode")
		h1.StartHTTP1Proxy()
	case "HTTP2":
		log.Info("Starting Proxy in HTTP1 <--> HTTP2 Forwarding Mode")
        go h2.StartHTTP1toHTTP2Proxy()
        h2.StartHTTP2toHTTP1Proxy()
	case "HTTP3":
		log.Info("Starting Proxy in HTTP1 <--> HTTP3 Forwarding Mode")
		go h3.StartHTTP1toHTTP3Proxy()  // Start listening for HTTP/1 and forwarding as HTTP/3
		h3.StartHTTP3toHTTP1Proxy()     // Start listening for HTTP/3 and forwarding as HTTP/1
	}
}





