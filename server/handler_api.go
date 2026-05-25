package main

import (
	"net/http"
	"strconv"
)

// handlerLocalisation gère GET /api/location
// Détecte l'IP du client et renvoie sa position géographique
func handlerLocalisation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := extraireIP(r)

		loc, err := localiserIP(ip)
		if err != nil {
			// En dev (localhost), ip-api refuse les IP privées
			// On renvoie Paris par défaut
			reponseJSON(w, http.StatusOK, Localisation{
				Lat:   48.8566,
				Lon:   2.3522,
				Ville: "Paris",
				Pays:  "France",
			})
			return
		}

		reponseJSON(w, http.StatusOK, loc)
	}
}

// handlerEvenements gère GET /api/events?lat=X&lon=Y&radius=Z
// Interroge l'API OpenData Paris et renvoie les événements à proximité
func handlerEvenements() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			reponseErreur(w, http.StatusBadRequest, "paramètre lat manquant ou invalide")
			return
		}
		lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		if err != nil {
			reponseErreur(w, http.StatusBadRequest, "paramètre lon manquant ou invalide")
			return
		}

		rayon, err := strconv.Atoi(r.URL.Query().Get("radius"))
		if err != nil || rayon < 1 || rayon > 50 {
			rayon = 5 // Rayon par défaut : 5km
		}

		categorie := r.URL.Query().Get("categorie")

		limite := 20 // Nombre d'événements par requête

		evenements, err := rechercherEvenements(lat, lon, rayon, limite, categorie)
		if err != nil {
			reponseErreur(w, http.StatusBadGateway, "erreur lors de la récupération des événements")
			return
		}

		reponseJSON(w, http.StatusOK, map[string]any{
			"evenements": evenements,
			"total":      len(evenements),
		})
	}
}
