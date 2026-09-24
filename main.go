package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/PPaterka/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db             *database.Queries
	fileserverHits atomic.Int32
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM is not set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)
	apiCfg := apiConfig{
		db:             dbQueries,
		fileserverHits: atomic.Int32{},
		platform:       platform,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)
	mux.HandleFunc("GET /admin/metrics", apiCfg.returnHits)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	mux.HandleFunc("POST /api/chirps", apiCfg.newChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
