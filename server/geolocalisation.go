package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Localisation struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Ville string  `json:"ville"`
	Pays  string  `json:"pays"`
}

// reponseIPAPI ne mappe que les champs que l'on utilise réellement, le reste
// du JSON renvoyé par ip-api.com est ignoré silencieusement par le décodeur.
type reponseIPAPI struct {
	Status  string  `json:"status"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	City    string  `json:"city"`
	Country string  `json:"country"`
	Message string  `json:"message"`
}

func localiserIP(ip string) (*Localisation, error) {
	// Sans IP dans l'URL, ip-api.com géolocalise l'IP appelante — ce qui
	// donne celle du serveur en local et n'a pas de sens : on passe donc
	// systématiquement l'IP du client.
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("appel ip-api: %w", err)
	}
	defer resp.Body.Close()

	var data reponseIPAPI
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("décodage ip-api: %w", err)
	}

	if data.Status != "success" {
		return nil, fmt.Errorf("ip-api erreur: %s", data.Message)
	}

	return &Localisation{
		Lat:   data.Lat,
		Lon:   data.Lon,
		Ville: data.City,
		Pays:  data.Country,
	}, nil
}

// extraireIP regarde d'abord les en-têtes posés par un éventuel reverse
// proxy, puis retombe sur RemoteAddr (qui contient ip:port).
func extraireIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For peut chaîner plusieurs IPs séparées par ", " :
		// la première est celle du client d'origine.
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	ip := r.RemoteAddr
	if i := strings.LastIndexByte(ip, ':'); i >= 0 {
		return ip[:i]
	}
	return ip
}
