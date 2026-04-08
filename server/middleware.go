package main

import (
	"context"
	"database/sql"
	"net/http"
)

// Type privé pour la clé de contexte — évite les collisions
type cleContexte string

const cleUtilisateurID cleContexte = "utilisateur_id"

// proteger vérifie la session avant d'appeler le handler suivant
func proteger(db *sql.DB, suivant http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			reponseErreur(w, http.StatusUnauthorized, "non connecté")
			return
		}

		utilisateurID, err := validerSession(db, cookie.Value)
		if err != nil {
			reponseErreur(w, http.StatusUnauthorized, "session invalide ou expirée")
			return
		}

		// Transmet l'ID utilisateur au handler via le contexte
		ctx := context.WithValue(r.Context(), cleUtilisateurID, utilisateurID)
		suivant(w, r.WithContext(ctx))
	}
}
