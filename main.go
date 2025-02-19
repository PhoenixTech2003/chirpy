package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/phoenixTech2003/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareWritesMetrics(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		serverHits := strconv.FormatInt(int64(cfg.fileServerHits.Load()), 10)
		w.Write([]byte(serverHits))
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) middlewareResetServerHits(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(200)
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err!=nil {
		log.Printf("Failed to open database connection %s", err)
	}
	apiCfg := apiConfig{
		dbQueries: database.New(db),
	}
	
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets", http.FileServer(http.Dir("./assets/logo.png")))
	mux.HandleFunc("GET /api/healthz", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(200)
		response.Write([]byte("OK"))

	})
	mux.HandleFunc("POST /api/users", apiCfg.postUsers)
	mux.HandleFunc("POST /api/validate_chirp", validator)
	mux.Handle("/admin/metrics", apiCfg.middlewareWritesMetrics(http.FileServer(http.Dir("./admin.html"))))
	mux.HandleFunc("POST /api/reset", apiCfg.middlewareResetServerHits)
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()

}
