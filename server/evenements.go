package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Evenement représente un événement parisien, mappé depuis l'API OpenData
type Evenement struct {
	ID          string  `json:"id"`
	Titre       string  `json:"titre"`
	Description string  `json:"description"`
	Chapeau     string  `json:"chapeau"` // Texte d'accroche
	DateDebut   string  `json:"date_debut"`
	DateFin     string  `json:"date_fin"`
	CoverURL    string  `json:"cover_url"`
	Adresse     string  `json:"adresse"`
	Ville       string  `json:"ville"`
	CodePostal  string  `json:"code_postal"`
	NomLieu     string  `json:"nom_lieu"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Prix        string  `json:"prix"`
	URLContact  string  `json:"url_contact"`
	Audience    string  `json:"audience"`
}

const baseURLOpenData = "https://opendata.paris.fr/api/explore/v2.1/catalog/datasets/que-faire-a-paris-/records"

// rechercherEvenements appelle l'API OpenData Paris avec un filtre géographique
func rechercherEvenements(lat, lon float64, rayonKm int, limite int) ([]Evenement, error) {
	// Filtre ODSQL : événements dans le rayon, avec des dates définies
	filtre := fmt.Sprintf(
		"within_distance(lat_lon, GEOM'POINT(%f %f)', %dkm) AND date_start IS NOT NULL AND date_end >= '%s'",
		lon, lat, rayonKm, time.Now().Format("2006-01-02"),
	)

	params := url.Values{}
	params.Set("where", filtre)
	params.Set("limit", fmt.Sprintf("%d", limite))
	params.Set("order_by", "date_start ASC")

	reqURL := baseURLOpenData + "?" + params.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("appel opendata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opendata status %d", resp.StatusCode)
	}

	// Parsing de la réponse — structure spécifique à l'API OpenData v2.1
	var resultat struct {
		TotalCount int              `json:"total_count"`
		Results    []json.RawMessage `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resultat); err != nil {
		return nil, fmt.Errorf("décodage opendata: %w", err)
	}

	evenements := make([]Evenement, 0, len(resultat.Results))
	for _, raw := range resultat.Results {
		ev, err := convertirRecord(raw)
		if err != nil {
			continue // On ignore les records mal formés
		}
		evenements = append(evenements, ev)
	}

	return evenements, nil
}

// convertirRecord transforme un record brut OpenData en Evenement
func convertirRecord(raw json.RawMessage) (Evenement, error) {
	var record struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		LeadText    string `json:"lead_text"`
		DateStart   string `json:"date_start"`
		DateEnd     string `json:"date_end"`
		CoverURL    string `json:"cover_url"`
		AddrStreet  string `json:"address_street"`
		AddrCity    string `json:"address_city"`
		AddrZip     string `json:"address_zipcode"`
		AddrName    string `json:"address_name"`
		LatLon      *struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"lat_lon"`
		PriceType  string `json:"price_type"`
		ContactURL string `json:"contact_url"`
		Audience   string `json:"audience"`
	}

	if err := json.Unmarshal(raw, &record); err != nil {
		return Evenement{}, err
	}

	ev := Evenement{
		ID:          record.ID,
		Titre:       record.Title,
		Description: record.Description,
		Chapeau:     record.LeadText,
		DateDebut:   record.DateStart,
		DateFin:     record.DateEnd,
		CoverURL:    record.CoverURL,
		Adresse:     record.AddrStreet,
		Ville:       record.AddrCity,
		CodePostal:  record.AddrZip,
		NomLieu:     record.AddrName,
		Prix:        record.PriceType,
		URLContact:  record.ContactURL,
		Audience:    record.Audience,
	}

	if record.LatLon != nil {
		ev.Lat = record.LatLon.Lat
		ev.Lon = record.LatLon.Lon
	}

	return ev, nil
}
