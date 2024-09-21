// backend/cmd/main.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Summit Social!")
	})

	fmt.Println("Backend server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
