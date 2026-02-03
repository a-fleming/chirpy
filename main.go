package main

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	godotenv.Load()

	const port = "8080"
	cfg := apiConfig{}
	cfg.platform = os.Getenv("PLATFORM")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	cfg.db = dbQueries

	fmt.Println("Hello, from chirpy!")

	serveMux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	serveMux.HandleFunc("GET /api/healthz", handlerHealth)
	serveMux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirp_id}", cfg.handlerGetChirpByID)
	serveMux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	serveMux.HandleFunc("POST /api/users", cfg.handlerUsers)

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
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("error encountered: %v\n", err)
	}

}
