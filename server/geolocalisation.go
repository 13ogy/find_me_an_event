package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Localisation représente la position géographique d'un utilisateur
type Localisation struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Ville string  `json:"ville"`
	Pays  string  `json:"pays"`
}

// reponseIPAPI correspond au format de ip-api.com
type reponseIPAPI struct {
	Status  string  `json:"status"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	City    string  `json:"city"`
	Country string  `json:"country"`
	Message string  `json:"message"`
}

// localiserIP interroge ip-api.com pour géolocaliser une adresse IP
func localiserIP(ip string) (*Localisation, error) {
	// ip-api.com accepte l'IP dans l'URL, ou renvoie celle de l'appelant si vide
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

// extraireIP récupère l'IP du client depuis les en-têtes ou RemoteAddr
func extraireIP(r *http.Request) string {
	// X-Forwarded-For est rempli par les proxys/load balancers
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// RemoteAddr contient ip:port, on garde juste l'IP
	ip := r.RemoteAddr
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] == ':' {
			return ip[:i]
		}
	}
	return ip
}
