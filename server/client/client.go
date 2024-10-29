package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func sendRequest(proxy string, requestUrl string) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Fatalf("Error parsing proxy URL: %v", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get(requestUrl)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(body))

	defer resp.Body.Close()

}

func main() {
	sendRequest(os.Args[1], os.Args[2])
}
