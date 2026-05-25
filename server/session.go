package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

const dureeSession = 24 * time.Hour

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

// validerSession renvoie l'ID utilisateur si la session existe et n'a pas
// expiré. Une session expirée est aussitôt purgée pour ne pas s'accumuler.
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

	if time.Now().After(expireLe) {
		supprimerSession(db, token)
		return 0, sql.ErrNoRows
	}

	return utilisateurID, nil
}

func supprimerSession(db *sql.DB, token string) {
	db.Exec(`DELETE FROM sessions WHERE id = $1`, token)
}

// genererToken produit 32 octets aléatoires (256 bits d'entropie) encodés en
// hex : suffisamment imprévisible pour ne pas être deviné par bruteforce.
func genererToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
