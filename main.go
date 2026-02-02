package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type Chirp struct {
	Body string `json:"body"`
}

type ValidationResponse struct {
	Error string `json:"error"`
	Valid bool   `json:"valid"`
}

func (cfg *apiConfig) middlewareHandlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, cfg.fileserverHits.Load())
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

func handlerHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	msg := "OK"
	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Printf("error writing response: %v\n", err)
	}
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	chirp := Chirp{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	if len(chirp.Body) > 140 {
		msg := "Chirp is too long"
		respondWithError(w, http.StatusBadRequest, msg, nil)
		return
	}
	resp := ValidationResponse{
		Valid: true,
	}
	respondWithJSON(w, http.StatusOK, resp)
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

	serveMux.HandleFunc("GET /admin/metrics", cfg.middlewareHandlerMetrics)
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
