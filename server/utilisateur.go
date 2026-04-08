package main

import (
	"database/sql"
	"time"
)

// Utilisateur correspond à une ligne de la table utilisateurs
type Utilisateur struct {
	ID              int       `json:"id"`
	Nom             string    `json:"nom"`
	Email           string    `json:"email"`
	MotDePasseHash  string    `json:"-"` // Jamais renvoyé en JSON
	CreeLe          time.Time `json:"cree_le"`
}

// creerUtilisateur insère un nouvel utilisateur en base
func creerUtilisateur(db *sql.DB, nom, email, hash string) (*Utilisateur, error) {
	u := &Utilisateur{}
	err := db.QueryRow(
		`INSERT INTO utilisateurs (nom, email, mot_de_passe_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, nom, email, cree_le`,
		nom, email, hash,
	).Scan(&u.ID, &u.Nom, &u.Email, &u.CreeLe)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// trouverParNom cherche un utilisateur par son nom (pour le login)
func trouverParNom(db *sql.DB, nom string) (*Utilisateur, error) {
	u := &Utilisateur{}
	err := db.QueryRow(
		`SELECT id, nom, email, mot_de_passe_hash, cree_le
		 FROM utilisateurs WHERE nom = $1`,
		nom,
	).Scan(&u.ID, &u.Nom, &u.Email, &u.MotDePasseHash, &u.CreeLe)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// trouverParID cherche un utilisateur par son ID
func trouverParID(db *sql.DB, id int) (*Utilisateur, error) {
	u := &Utilisateur{}
	err := db.QueryRow(
		`SELECT id, nom, email, cree_le
		 FROM utilisateurs WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Nom, &u.Email, &u.CreeLe)
	if err != nil {
		return nil, err
	}
	return u, nil
}
