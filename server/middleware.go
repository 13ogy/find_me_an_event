package main

import (
	"context"
	"database/sql"
	"net/http"
)

// Type privé pour éviter les collisions de clés dans context.Value :
// recommandation explicite de la documentation de context.
type cleContexte string

const cleUtilisateurID cleContexte = "utilisateur_id"

// utilisateurDeContexte récupère l'ID injecté par proteger() et renvoie un
// indicateur d'absence : évite un panic en cas de mauvais câblage de routes
// (handler protégé monté sans le middleware par erreur).
func utilisateurDeContexte(r *http.Request) (int, bool) {
	id, ok := r.Context().Value(cleUtilisateurID).(int)
	return id, ok
}

// proteger refuse l'accès aux requêtes sans cookie de session valide et
// injecte l'ID utilisateur dans le contexte pour le handler aval.
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

		ctx := context.WithValue(r.Context(), cleUtilisateurID, utilisateurID)
		suivant(w, r.WithContext(ctx))
	}
}
