package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	if err != nil {
		fmt.Printf("error displaying number of hits: %v", err)
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error writing response: %s", err)
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpBody struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := chirpBody{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleanSplit := strings.Split(chirp.Body, " ")
	chirpString := strings.ToLower(chirp.Body)
	chirpSplit := strings.Split(chirpString, " ")
	for i, word := range chirpSplit {
		switch word {
		case "kerfuffle", "sharbert", "fornax":
			chirpSplit[i] = "****"
		default:
			chirpSplit[i] = cleanSplit[i]
		}
	}
	type cleanChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}
	cleanedChirp := cleanChirp{}
	cleanedChirp.CleanedBody = strings.Join(chirpSplit, " ")
	respondWithJSON(w, 200, cleanedChirp)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("could not load env: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("could not open postres: %v", err)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		db: dbQueries,
	}

	port := "8080"
	mux := http.NewServeMux()
	filepathRoot := http.Dir(".")
	handler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(filepathRoot)))
	mux.Handle("/app/", handler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			log.Printf("could not write to healthz: %v", err)
		}
	})
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
