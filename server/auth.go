package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func reponseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func reponseErreur(w http.ResponseWriter, code int, message string) {
	reponseJSON(w, code, map[string]string{"erreur": message})
}

func handlerInscription(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Nom        string `json:"nom"`
			Email      string `json:"email"`
			MotDePasse string `json:"mot_de_passe"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			reponseErreur(w, http.StatusBadRequest, "JSON invalide")
			return
		}

		body.Nom = strings.TrimSpace(body.Nom)
		body.Email = strings.TrimSpace(body.Email)
		if body.Nom == "" || body.Email == "" || body.MotDePasse == "" {
			reponseErreur(w, http.StatusBadRequest, "tous les champs sont requis")
			return
		}
		if len(body.MotDePasse) < 6 {
			reponseErreur(w, http.StatusBadRequest, "mot de passe trop court (6 caractères minimum)")
			return
		}

		// Coût 10 : compromis sécurité/perf raisonnable, bcrypt prend ~70ms
		// par hash sur une machine moderne et ralentit donc le bruteforce.
		hash, err := bcrypt.GenerateFromPassword([]byte(body.MotDePasse), 10)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		utilisateur, err := creerUtilisateur(db, body.Nom, body.Email, string(hash))
		if err != nil {
			// Détection de la violation d'UNIQUE par le texte d'erreur : on
			// préfère un 409 propre à l'utilisateur plutôt qu'un 500 générique.
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				reponseErreur(w, http.StatusConflict, "nom ou email déjà utilisé")
				return
			}
			reponseErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		reponseJSON(w, http.StatusCreated, utilisateur)
	}
}

func handlerConnexion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Nom        string `json:"nom"`
			MotDePasse string `json:"mot_de_passe"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			reponseErreur(w, http.StatusBadRequest, "JSON invalide")
			return
		}

		utilisateur, err := trouverParNom(db, body.Nom)
		if err != nil {
			// On renvoie le même message qu'un mauvais mot de passe pour ne
			// pas révéler quels noms d'utilisateur existent.
			reponseErreur(w, http.StatusUnauthorized, "identifiants incorrects")
			return
		}

		// bcrypt.CompareHashAndPassword est en temps constant : pas de fuite
		// par timing sur la comparaison des hash.
		if err := bcrypt.CompareHashAndPassword([]byte(utilisateur.MotDePasseHash), []byte(body.MotDePasse)); err != nil {
			reponseErreur(w, http.StatusUnauthorized, "identifiants incorrects")
			return
		}

		token, err := creerSession(db, utilisateur.ID)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		// HttpOnly : inaccessible depuis document.cookie, donc inexploitable
		// par un script injecté (XSS). SameSite=Lax bloque les CSRF basiques
		// tout en laissant marcher la navigation depuis nos propres pages.
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400,
			SameSite: http.SameSiteLaxMode,
		})

		reponseJSON(w, http.StatusOK, utilisateur)
	}
}

func handlerDeconnexion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == nil {
			supprimerSession(db, cookie.Value)
		}

		// MaxAge négatif : force le navigateur à supprimer le cookie.
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		reponseJSON(w, http.StatusOK, map[string]string{"message": "déconnecté"})
	}
}

func handlerProfil(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilisateurID, ok := utilisateurDeContexte(r)
		if !ok {
			reponseErreur(w, http.StatusUnauthorized, "non connecté")
			return
		}

		utilisateur, err := trouverParID(db, utilisateurID)
		if err != nil {
			reponseErreur(w, http.StatusNotFound, "utilisateur introuvable")
			return
		}

		reponseJSON(w, http.StatusOK, utilisateur)
	}
}
