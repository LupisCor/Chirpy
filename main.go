package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/LupisCor/Chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbconn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbconn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))) // Serve static files from the current directory

	// HandleFunc format: mux.HandleFunc("/path", handlerFunction)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)             // Handle readiness check
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp) // Handle chirp validation

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics) // Handle metrics endpoint
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)    // Handle reset endpoint

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
