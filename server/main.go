package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := chargerConfig()

	db, err := ouvrirDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("impossible de se connecter à la base : ", err)
	}
	defer db.Close()

	if err := creerTables(db); err != nil {
		log.Fatal("erreur migration : ", err)
	}

	// On utilise le ServeMux 1.22+ pour ne pas avoir à coder le routage par
	// méthode HTTP à la main — c'est encore du net/http, donc compatible
	// avec la contrainte du sujet.
	mux := http.NewServeMux()

	// Routes publiques d'authentification.
	mux.HandleFunc("POST /api/register", handlerInscription(db))
	mux.HandleFunc("POST /api/login", handlerConnexion(db))
	mux.HandleFunc("POST /api/logout", handlerDeconnexion(db))

	// Routes protégées : profil, géoloc, événements, votes, avis, historique.
	mux.HandleFunc("GET /api/me", proteger(db, handlerProfil(db)))
	mux.HandleFunc("GET /api/location", proteger(db, handlerLocalisation()))
	mux.HandleFunc("GET /api/events", proteger(db, handlerEvenements()))
	mux.HandleFunc("GET /api/events/next", proteger(db, handlerProchainEvenement(db)))
	mux.HandleFunc("POST /api/events/{id}/vote", proteger(db, handlerVoter(db)))
	mux.HandleFunc("DELETE /api/events/{id}/vote", proteger(db, handlerSupprimerVote(db)))
	mux.HandleFunc("GET /api/events/{id}/stats", proteger(db, handlerStats(db)))
	mux.HandleFunc("GET /api/events/{id}/reviews", proteger(db, handlerListerAvis(db)))
	mux.HandleFunc("POST /api/events/{id}/reviews", proteger(db, handlerPosterAvis(db)))
	mux.HandleFunc("GET /api/me/history", proteger(db, handlerHistorique(db)))

	log.Printf("serveur démarré sur :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, cors(mux)))
}
