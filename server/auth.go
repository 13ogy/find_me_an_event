package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// reponseJSON écrit une réponse JSON avec le code HTTP donné
func reponseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// reponseErreur écrit une erreur JSON
func reponseErreur(w http.ResponseWriter, code int, message string) {
	reponseJSON(w, code, map[string]string{"erreur": message})
}

// handlerInscription gère POST /api/register
func handlerInscription(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Nom      string `json:"nom"`
			Email    string `json:"email"`
			MotDePasse string `json:"mot_de_passe"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			reponseErreur(w, http.StatusBadRequest, "JSON invalide")
			return
		}

		// Validation basique
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

		// Hash bcrypt — coût 10, bon compromis sécurité/performance
		hash, err := bcrypt.GenerateFromPassword([]byte(body.MotDePasse), 10)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		utilisateur, err := creerUtilisateur(db, body.Nom, body.Email, string(hash))
		if err != nil {
			// Contrainte UNIQUE violée → nom ou email déjà pris
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

// handlerConnexion gère POST /api/login
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
			reponseErreur(w, http.StatusUnauthorized, "identifiants incorrects")
			return
		}

		// Comparaison bcrypt — résiste aux attaques par timing
		if err := bcrypt.CompareHashAndPassword([]byte(utilisateur.MotDePasseHash), []byte(body.MotDePasse)); err != nil {
			reponseErreur(w, http.StatusUnauthorized, "identifiants incorrects")
			return
		}

		token, err := creerSession(db, utilisateur.ID)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		// Cookie HttpOnly — inaccessible depuis JS, protège contre XSS
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400, // 24h en secondes
			SameSite: http.SameSiteLaxMode,
		})

		reponseJSON(w, http.StatusOK, utilisateur)
	}
}

// handlerDeconnexion gère POST /api/logout
func handlerDeconnexion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == nil {
			supprimerSession(db, cookie.Value)
		}

		// Supprime le cookie côté client
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

// handlerProfil gère GET /api/me
func handlerProfil(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// L'ID utilisateur est placé dans le contexte par le middleware proteger
		utilisateurID := r.Context().Value(cleUtilisateurID).(int)

		utilisateur, err := trouverParID(db, utilisateurID)
		if err != nil {
			reponseErreur(w, http.StatusNotFound, "utilisateur introuvable")
			return
		}

		reponseJSON(w, http.StatusOK, utilisateur)
	}
}
