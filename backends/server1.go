package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🍀 Hello from Server 1")
	})
	http.ListenAndServe(":9001", nil)
}
