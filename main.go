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

func (cfg *apiConfig) middlewareHandlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Printf("error writing response: %v\n", err)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) middlewareHandlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	msg := "metrics reset"
	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Printf("error writing response: %v\n", err)
	}
}

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
	cfg := apiConfig{}

	fmt.Println("Hello, from chirpy!")

	serveMux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	serveMux.HandleFunc("GET /api/healthz", handlerHealth)
	serveMux.HandleFunc("GET /api/metrics", cfg.middlewareHandlerMetrics)
	serveMux.HandleFunc("POST /api/reset", cfg.middlewareHandlerReset)

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
