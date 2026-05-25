package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Les tests qui suivent ne touchent pas PostgreSQL : on ne couvre que les
// helpers et les middlewares purs. C'est suffisant pour vérifier les
// invariants importants (génération de tokens, parsing OpenData, CORS,
// extraction d'IP) sans imposer une base de test à la chaîne d'intégration.

func TestGenererToken(t *testing.T) {
	tok, err := genererToken()
	if err != nil {
		t.Fatalf("genererToken: %v", err)
	}
	// 32 octets en hex = 64 caractères ; on s'assure aussi que c'est bien
	// du hex décodable, sinon le stockage en base et le cookie casseraient.
	if len(tok) != 64 {
		t.Errorf("longueur attendue 64, obtenu %d", len(tok))
	}
	if _, err := hex.DecodeString(tok); err != nil {
		t.Errorf("token non hex: %v", err)
	}

	tok2, _ := genererToken()
	if tok == tok2 {
		t.Error("deux tokens consécutifs identiques — RNG cassé")
	}
}

func TestExtraireIP(t *testing.T) {
	mkReq := func(remote string, h map[string]string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = remote
		for k, v := range h {
			req.Header.Set(k, v)
		}
		return req
	}

	cas := []struct {
		nom     string
		req     *http.Request
		attendu string
	}{
		{
			nom:     "X-Forwarded-For simple",
			req:     mkReq("10.0.0.1:1234", map[string]string{"X-Forwarded-For": "203.0.113.7"}),
			attendu: "203.0.113.7",
		},
		{
			nom:     "X-Forwarded-For chaîné garde l'IP d'origine",
			req:     mkReq("10.0.0.1:1234", map[string]string{"X-Forwarded-For": "203.0.113.7, 10.0.0.1"}),
			attendu: "203.0.113.7",
		},
		{
			nom:     "X-Real-IP en repli",
			req:     mkReq("10.0.0.1:1234", map[string]string{"X-Real-IP": "198.51.100.42"}),
			attendu: "198.51.100.42",
		},
		{
			nom:     "RemoteAddr ipv4 sans en-tête",
			req:     mkReq("192.0.2.1:5678", nil),
			attendu: "192.0.2.1",
		},
	}

	for _, c := range cas {
		t.Run(c.nom, func(t *testing.T) {
			got := extraireIP(c.req)
			if got != c.attendu {
				t.Errorf("attendu %q, obtenu %q", c.attendu, got)
			}
		})
	}
}

func TestCORSPreflight(t *testing.T) {
	// OPTIONS doit court-circuiter le handler aval avec un 204 et poser
	// les bons en-têtes.
	suiteAppelee := false
	h := cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suiteAppelee = true
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if suiteAppelee {
		t.Error("le handler aval ne doit pas être appelé sur un preflight")
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("statut attendu 204, obtenu %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Error("Access-Control-Allow-Origin doit refléter l'Origin")
	}
	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("credentials doivent être autorisés pour que la session passe")
	}
}

func TestCORSGet(t *testing.T) {
	// Sur une vraie requête, CORS doit poser les en-têtes puis déléguer.
	h := cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Errorf("le handler aval doit être appelé, statut attendu 418, obtenu %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("en-tête CORS manquant sur la réponse")
	}
}

func TestReponseJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	reponseJSON(rec, http.StatusCreated, map[string]string{"ok": "yes"})

	if rec.Code != http.StatusCreated {
		t.Errorf("statut attendu 201, obtenu %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type attendu application/json, obtenu %q", ct)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("corps non JSON: %v", err)
	}
	if payload["ok"] != "yes" {
		t.Errorf("payload mal sérialisé: %v", payload)
	}
}

func TestReponseErreur(t *testing.T) {
	rec := httptest.NewRecorder()
	reponseErreur(rec, http.StatusBadRequest, "champ requis")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("statut attendu 400, obtenu %d", rec.Code)
	}
	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("corps non JSON: %v", err)
	}
	if payload["erreur"] != "champ requis" {
		t.Errorf("le champ 'erreur' doit être renvoyé tel quel, obtenu %v", payload)
	}
}

func TestConvertirRecord(t *testing.T) {
	// Record OpenData réduit aux champs qui nous intéressent : on vérifie le
	// renommage anglais → français et l'aplatissement de lat_lon.
	raw := json.RawMessage(`{
		"id": "evt-1",
		"title": "Concert au parc",
		"description": "<p>Concert gratuit</p>",
		"lead_text": "Concert en plein air",
		"date_start": "2026-06-01T18:00:00+00:00",
		"date_end":   "2026-06-01T22:00:00+00:00",
		"cover_url":  "https://example.org/img.jpg",
		"address_street": "1 rue X",
		"address_city":   "Paris",
		"address_zipcode":"75011",
		"address_name":   "Parc Y",
		"lat_lon": {"lat": 48.86, "lon": 2.37},
		"price_type": "gratuit",
		"contact_url": "https://example.org/info",
		"audience": "tout public"
	}`)

	ev, err := convertirRecord(raw)
	if err != nil {
		t.Fatalf("convertirRecord: %v", err)
	}

	if ev.ID != "evt-1" || ev.Titre != "Concert au parc" {
		t.Errorf("ID/Titre mal mappés: %+v", ev)
	}
	if ev.Chapeau != "Concert en plein air" {
		t.Errorf("lead_text → Chapeau mal mappé: %q", ev.Chapeau)
	}
	if ev.Lat != 48.86 || ev.Lon != 2.37 {
		t.Errorf("lat_lon mal aplati: lat=%v lon=%v", ev.Lat, ev.Lon)
	}
	if ev.Adresse != "1 rue X" || ev.CodePostal != "75011" || ev.NomLieu != "Parc Y" {
		t.Errorf("adresse mal mappée: %+v", ev)
	}
	if ev.Prix != "gratuit" || ev.URLContact != "https://example.org/info" {
		t.Errorf("prix/contact mal mappés: %+v", ev)
	}
}

func TestConvertirRecordSansLatLon(t *testing.T) {
	// Certains records n'ont pas de coordonnées : on doit accepter sans
	// paniquer et laisser Lat/Lon à zéro.
	raw := json.RawMessage(`{"id": "evt-2", "title": "Sans lieu"}`)
	ev, err := convertirRecord(raw)
	if err != nil {
		t.Fatalf("convertirRecord: %v", err)
	}
	if ev.Lat != 0 || ev.Lon != 0 {
		t.Errorf("Lat/Lon devraient être 0 sans lat_lon, obtenu %v / %v", ev.Lat, ev.Lon)
	}
}

func TestConvertirRecordJSONInvalide(t *testing.T) {
	// JSON syntaxiquement cassé → erreur, pas de panique.
	_, err := convertirRecord(json.RawMessage(`{invalide`))
	if err == nil {
		t.Error("un JSON malformé doit remonter une erreur")
	}
}

func TestChargerConfigParDefaut(t *testing.T) {
	// Si on n'a pas posé PORT/DATABASE_URL dans l'environnement de test,
	// les valeurs par défaut doivent répondre.
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	cfg := chargerConfig()
	if cfg.Port != "8080" {
		t.Errorf("port par défaut attendu 8080, obtenu %q", cfg.Port)
	}
	if !strings.Contains(cfg.DatabaseURL, "find_me_an_event") {
		t.Errorf("URL par défaut suspecte: %q", cfg.DatabaseURL)
	}
}

func TestChargerConfigSurchargee(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://x:y@host/db")
	cfg := chargerConfig()
	if cfg.Port != "9090" || cfg.DatabaseURL != "postgres://x:y@host/db" {
		t.Errorf("la config ne reflète pas l'environnement: %+v", cfg)
	}
}
