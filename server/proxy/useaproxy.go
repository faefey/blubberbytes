package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

	"github.com/elazarl/goproxy"
)

func useproxy() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxyURL, err := url.Parse("http://23.239.12.179:8080")
	if err != nil {
		log.Fatal("Error parsing proxy URL:", err)
	}

	//HTTP
	proxy.Tr = &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//HTTPS
	proxy.ConnectDial = proxy.NewConnectDialToProxy(proxyURL.String())

	// Start the local proxy server on port 9000
	log.Println("Starting local proxy server on :9000")
	log.Fatal(http.ListenAndServe(":9000", proxy))
}
