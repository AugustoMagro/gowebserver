package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AugustoMagro/gowebserver/internal/chirpyapi"
	"github.com/AugustoMagro/gowebserver/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	const filepathRoot = "./templates"
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	plat := os.Getenv("PLATFORM")
	if plat == "" {
		log.Printf("PLATFORM must be set")
	}

	dbQueries := database.New(db)

	if err != nil {
		log.Printf("Connection to Database failed: %s", err)
	}

	apiCfg := chirpyapi.NewClient(dbQueries, plat, os.Getenv("SECRET_KEY"))

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serverMux.HandleFunc("GET /api/healthz/", chirpyapi.HandlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)
	serverMux.HandleFunc("POST /api/validate_chirp", chirpyapi.ValidateChirp)
	serverMux.HandleFunc("POST /api/users", apiCfg.HandleUser)
	serverMux.HandleFunc("GET /api/users", apiCfg.GetUsers)
	serverMux.HandleFunc("POST /api/login", apiCfg.HandleUserLogin)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.CreateChirpy)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.GetChirps)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpsByID)
	serverMux.HandleFunc("GET /api/refresh", apiCfg.GetChirpsByID)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
