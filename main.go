package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	cfg := apiConfig{}

	fmt.Println("Hello, from chirpy!")

	serveMux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	serveMux.HandleFunc("GET /api/healthz", handlerHealth)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        serveMux,
		ReadTimeout:    10 * time.Second, // from net/http example
		WriteTimeout:   10 * time.Second, // from net/http example
		MaxHeaderBytes: 1 << 20,          // from net/http example
	}
	fmt.Printf("Serving on port: %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("error encountered: %v\n", err)
	}

}
