package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func proxyChecker(next http.Handler, expectedHost string, useProxy string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != expectedHost {
			if useProxy != "" {
				proxyHandler(w, r)
			} else {
				fmt.Fprintf(w, "Please do not use me as a proxy!")
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("helloHandler has been accessed")
	fmt.Fprintf(w, "Good morning")
}

// Code obtained from https://github.com/eliben/code-for-blog/blob/main/2022/go-and-proxies/forward-proxy-using-reverseproxy.go
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Proxy has been accessed")
	target, err := url.Parse(r.URL.Scheme + "://" + r.URL.Host)
	if err != nil {
		panic(err)
	}

	//reqb, err := httputil.DumpRequest(r, true)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(string(reqb))

	p := httputil.NewSingleHostReverseProxy(target)
	p.ServeHTTP(w, r)
}

func main() {
	//Port is 3000 if no arg else entered port [go run server.go [port]]
	port := "3000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	useProxy := ""
	if len(os.Args) > 2 {
		useProxy = os.Args[2]
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", proxyHandler)
	mux.HandleFunc("/hello", helloHandler)

	address := "localhost:" + port

	proxyCheckMux := proxyChecker(mux, address, useProxy)

	fmt.Println("server is running on port " + port + "...")

	if err := http.ListenAndServe(":"+port, proxyCheckMux); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
