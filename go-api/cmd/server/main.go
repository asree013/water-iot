package server

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go IoT Ingestion API is running!")
	})
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
