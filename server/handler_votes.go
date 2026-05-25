package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// handlerVoter gère POST /api/events/{id}/vote
func handlerVoter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilisateurID := r.Context().Value(cleUtilisateurID).(int)
		evenementID := r.PathValue("id")

		var body struct {
			Vote  string `json:"vote"`
			Titre string `json:"titre"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			reponseErreur(w, http.StatusBadRequest, "JSON invalide")
			return
		}

		if body.Vote != "like" && body.Vote != "unlike" {
			reponseErreur(w, http.StatusBadRequest, "vote doit être 'like' ou 'unlike'")
			return
		}

		if err := voterPourEvenement(db, utilisateurID, evenementID, body.Titre, body.Vote); err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur lors du vote")
			return
		}

		reponseJSON(w, http.StatusOK, map[string]string{"message": "vote enregistré"})
	}
}

// handlerSupprimerVote gère DELETE /api/events/{id}/vote
func handlerSupprimerVote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilisateurID := r.Context().Value(cleUtilisateurID).(int)
		evenementID := r.PathValue("id")

		if err := supprimerVote(db, utilisateurID, evenementID); err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur suppression")
			return
		}

		reponseJSON(w, http.StatusOK, map[string]string{"message": "vote supprimé"})
	}
}

// handlerStats gère GET /api/events/{id}/stats
func handlerStats(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		evenementID := r.PathValue("id")

		stats, err := obtenirStatsVotes(db, evenementID)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur stats")
			return
		}

		reponseJSON(w, http.StatusOK, stats)
	}
}

// handlerListerAvis gère GET /api/events/{id}/reviews
func handlerListerAvis(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		evenementID := r.PathValue("id")

		avis, err := listerAvis(db, evenementID)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur avis")
			return
		}
		if avis == nil {
			avis = []Avis{} // Toujours renvoyer un tableau, jamais null
		}

		reponseJSON(w, http.StatusOK, avis)
	}
}

// handlerPosterAvis gère POST /api/events/{id}/reviews
func handlerPosterAvis(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilisateurID := r.Context().Value(cleUtilisateurID).(int)
		evenementID := r.PathValue("id")

		var body struct {
			Commentaire string `json:"commentaire"`
			Note        int    `json:"note"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			reponseErreur(w, http.StatusBadRequest, "JSON invalide")
			return
		}

		body.Commentaire = strings.TrimSpace(body.Commentaire)
		if body.Commentaire == "" {
			reponseErreur(w, http.StatusBadRequest, "commentaire vide")
			return
		}
		if body.Note < 1 || body.Note > 5 {
			reponseErreur(w, http.StatusBadRequest, "note entre 1 et 5")
			return
		}

		avis, err := creerAvis(db, utilisateurID, evenementID, body.Commentaire, body.Note)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur création avis")
			return
		}

		reponseJSON(w, http.StatusCreated, avis)
	}
}

// handlerHistorique gère GET /api/me/history
func handlerHistorique(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilisateurID := r.Context().Value(cleUtilisateurID).(int)

		votes, err := obtenirHistorique(db, utilisateurID)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur historique")
			return
		}
		if votes == nil {
			votes = []Vote{}
		}

		reponseJSON(w, http.StatusOK, votes)
	}
}

// handlerProchainEvenement gère GET /api/events/next?lat=X&lon=Y&radius=Z
// Renvoie le prochain événement que l'utilisateur n'a pas encore voté
func handlerProchainEvenement(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilisateurID := r.Context().Value(cleUtilisateurID).(int)

		lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			reponseErreur(w, http.StatusBadRequest, "paramètre lat invalide")
			return
		}
		lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		if err != nil {
			reponseErreur(w, http.StatusBadRequest, "paramètre lon invalide")
			return
		}
		rayon, _ := strconv.Atoi(r.URL.Query().Get("radius"))
		if rayon < 1 || rayon > 50 {
			rayon = 5
		}

		// Récupère un lot d'événements
		categorie := r.URL.Query().Get("categorie")
		evenements, err := rechercherEvenements(lat, lon, rayon, 50, categorie)
		if err != nil {
			reponseErreur(w, http.StatusBadGateway, "erreur API événements")
			return
		}

		// Filtre ceux déjà votés
		dejaVotes, err := evenementsDejaVotes(db, utilisateurID)
		if err != nil {
			reponseErreur(w, http.StatusInternalServerError, "erreur base de données")
			return
		}

		for _, ev := range evenements {
			if !dejaVotes[ev.ID] {
				// Ajoute les stats de votes à la réponse
				stats, _ := obtenirStatsVotes(db, ev.ID)
				reponseJSON(w, http.StatusOK, map[string]any{
					"evenement": ev,
					"stats":     stats,
				})
				return
			}
		}

		// Tous les événements ont été votés
		reponseJSON(w, http.StatusOK, map[string]any{
			"evenement": nil,
			"message":   "plus d'événements à découvrir dans ce rayon",
		})
	}
}
