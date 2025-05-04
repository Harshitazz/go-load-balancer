package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🔥 Hello from Server 2")
	})
	http.ListenAndServe(":9002", nil)
}
