package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func beproxy() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	log.Println("Starting proxy server on :8080")
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
