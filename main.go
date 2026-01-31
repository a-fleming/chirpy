package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	const port = "8080"
	fmt.Println("Hello, from chirpy!")
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        serveMux,
		ReadTimeout:    10 * time.Second, // from net/http example
		WriteTimeout:   10 * time.Second, // from net/http example
		MaxHeaderBytes: 1 << 20,          // from net/http example
	}
	fmt.Printf("Serving on port: %s\n", port)
	server.ListenAndServe()

}
