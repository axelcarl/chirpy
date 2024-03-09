package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port string = "8080"

	cfg := apiConfig{}

	r := chi.NewRouter()
	corsMux := middlewareCors(r)

	fsHandler := cfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	r.Route("/api", func(r chi.Router) {
		r.Get("/healthz", handlerHealth)
		r.Get("/reset", cfg.handlerReset)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Get("/metrics", cfg.handlerMetrics)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Starting server on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
