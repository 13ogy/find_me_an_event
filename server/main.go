package main

import (
	"log"
	"net/http"
)

func main() {
	// Charge la config depuis les variables d'environnement
	cfg := chargerConfig()

	// Connexion PostgreSQL
	db, err := ouvrirDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("impossible de se connecter à la base :", err)
	}
	defer db.Close()

	// Création des tables si elles n'existent pas
	if err := creerTables(db); err != nil {
		log.Fatal("erreur migration :", err)
	}

	// Routeur principal — net/http standard, pas de framework
	mux := http.NewServeMux()

	// Routes d'authentification (publiques)
	mux.HandleFunc("POST /api/register", handlerInscription(db))
	mux.HandleFunc("POST /api/login", handlerConnexion(db))
	mux.HandleFunc("POST /api/logout", handlerDeconnexion(db))

	// Route protégée : profil utilisateur
	mux.HandleFunc("GET /api/me", proteger(db, handlerProfil(db)))

	// Routes événements et géolocalisation (protégées)
	mux.HandleFunc("GET /api/location", proteger(db, handlerLocalisation()))
	mux.HandleFunc("GET /api/events", proteger(db, handlerEvenements()))

	log.Printf("serveur démarré sur :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
