package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "9001", "Port to run the backend server on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request at backend server on port %s\n", *port)
		fmt.Fprintf(w, "Hello from backend server on port %s!", *port)
	})

	fmt.Printf("Backend server is running on port %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
