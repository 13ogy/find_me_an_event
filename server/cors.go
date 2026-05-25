package main

import "net/http"

// Liste blanche d'origines autorisées : on évite de refléter aveuglément
// l'Origin de l'appelant, ce qui combiné avec Allow-Credentials laisserait
// n'importe quel site faire des requêtes authentifiées au nom de
// l'utilisateur connecté. Pour la démo on n'autorise que le client Vite
// local.
var originsAutorisees = map[string]bool{
	"http://localhost:5173": true,
	"http://127.0.0.1:5173": true,
}

// En dev le client React tourne sur :5173 et le serveur sur :8080 :
// sans ce middleware, le navigateur bloque toutes les requêtes /api.
// On autorise les cookies pour que la session traverse le proxy.
func cors(suivant http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if originsAutorisees[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		suivant.ServeHTTP(w, r)
	})
}
