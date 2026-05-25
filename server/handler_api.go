package main

import (
	"net/http"
	"strconv"
)

// handlerLocalisation renvoie la position géographique correspondant à l'IP
// du client. En local (IP privée 127.x / 192.168.x) ip-api.com échoue : on
// retombe alors sur Paris pour que l'app reste utilisable en dev.
func handlerLocalisation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := extraireIP(r)

		loc, err := localiserIP(ip)
		if err != nil {
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

		// Borne le rayon : éviter qu'un client farceur ne demande 100 000 km
		// et fasse retourner toute la base OpenData d'un coup.
		rayon, err := strconv.Atoi(r.URL.Query().Get("radius"))
		if err != nil || rayon < 1 || rayon > 50 {
			rayon = 5
		}

		categorie := r.URL.Query().Get("categorie")

		evenements, err := rechercherEvenements(lat, lon, rayon, 20, categorie)
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
