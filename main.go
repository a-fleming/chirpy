package main

import (
	"fmt"
	"net/http"
	"time"
)

func handlerHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	msg := "OK"
	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Printf("error writing response: %v\n", err)
	}
}

func main() {
	const port = "8080"
	fmt.Println("Hello, from chirpy!")
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./"))))
	serveMux.HandleFunc("/healthz", handlerHealth)

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
