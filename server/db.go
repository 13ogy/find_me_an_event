package main

import (
	"database/sql"
	"fmt"

	// Import en blanc : le pilote s'enregistre auprès de database/sql via son init().
	_ "github.com/lib/pq"
)

func ouvrirDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	// sql.Open ne contacte pas la base : on force un round-trip pour
	// échouer tôt si la chaîne de connexion est mauvaise.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	return db, nil
}

// IF NOT EXISTS rend la migration rejouable à chaque démarrage du serveur.
func creerTables(db *sql.DB) error {
	requetes := []string{
		`CREATE TABLE IF NOT EXISTS utilisateurs (
			id SERIAL PRIMARY KEY,
			nom VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			mot_de_passe_hash VARCHAR(255) NOT NULL,
			cree_le TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id VARCHAR(64) PRIMARY KEY,
			utilisateur_id INTEGER REFERENCES utilisateurs(id) ON DELETE CASCADE,
			expire_le TIMESTAMP NOT NULL,
			cree_le TIMESTAMP DEFAULT NOW()
		)`,
		// UNIQUE(utilisateur_id, evenement_id) : un seul vote par utilisateur
		// et par événement — l'UPSERT côté Go s'appuie là-dessus.
		`CREATE TABLE IF NOT EXISTS votes (
			id SERIAL PRIMARY KEY,
			utilisateur_id INTEGER REFERENCES utilisateurs(id) ON DELETE CASCADE,
			evenement_id VARCHAR(255) NOT NULL,
			titre_evenement VARCHAR(500),
			vote VARCHAR(10) NOT NULL CHECK (vote IN ('like', 'unlike')),
			cree_le TIMESTAMP DEFAULT NOW(),
			UNIQUE(utilisateur_id, evenement_id)
		)`,
		`CREATE TABLE IF NOT EXISTS avis (
			id SERIAL PRIMARY KEY,
			utilisateur_id INTEGER REFERENCES utilisateurs(id) ON DELETE CASCADE,
			evenement_id VARCHAR(255) NOT NULL,
			commentaire TEXT NOT NULL,
			note INTEGER CHECK (note BETWEEN 1 AND 5),
			cree_le TIMESTAMP DEFAULT NOW()
		)`,
	}

	for _, req := range requetes {
		if _, err := db.Exec(req); err != nil {
			return fmt.Errorf("migration: %w", err)
		}
	}
	return nil
}
