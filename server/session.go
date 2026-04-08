package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

const dureeSession = 24 * time.Hour // Session valide 24h

// creerSession génère un token aléatoire et l'insère en base
func creerSession(db *sql.DB, utilisateurID int) (string, error) {
	token, err := genererToken()
	if err != nil {
		return "", err
	}

	expiration := time.Now().Add(dureeSession)
	_, err = db.Exec(
		`INSERT INTO sessions (id, utilisateur_id, expire_le) VALUES ($1, $2, $3)`,
		token, utilisateurID, expiration,
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

// validerSession vérifie qu'un token de session existe et n'est pas expiré
// Retourne l'ID de l'utilisateur associé
func validerSession(db *sql.DB, token string) (int, error) {
	var utilisateurID int
	var expireLe time.Time

	err := db.QueryRow(
		`SELECT utilisateur_id, expire_le FROM sessions WHERE id = $1`,
		token,
	).Scan(&utilisateurID, &expireLe)
	if err != nil {
		return 0, err
	}

	// Suppression automatique si expirée
	if time.Now().After(expireLe) {
		supprimerSession(db, token)
		return 0, sql.ErrNoRows
	}

	return utilisateurID, nil
}

// supprimerSession détruit une session (logout)
func supprimerSession(db *sql.DB, token string) {
	db.Exec(`DELETE FROM sessions WHERE id = $1`, token)
}

// genererToken produit 32 octets aléatoires en hexadécimal
func genererToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
