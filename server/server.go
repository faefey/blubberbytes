package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/", helloHandler) // the route
    fmt.Println("server is running on port 3000...")

    http.ListenAndServe(":3000", nil)
}
