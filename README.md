# Find Me An Event 🎉

Application web de découverte d'événements à Paris, développée dans le cadre du cours PC3R (Sorbonne Université).

L'utilisateur se connecte, l'application détecte sa position via son adresse IP, et lui propose des événements à proximité. Il peut parcourir les événements un par un (style Tinder), voter like ou unlike, consulter les avis des autres utilisateurs, et laisser des commentaires.

## Stack technique

| Composant | Technologie |
|-----------|-------------|
| Serveur | Go + net/http (sans framework) |
| Base de données | PostgreSQL |
| Client | React + Vite |
| API événements | [OpenData Paris — Que faire à Paris](https://opendata.paris.fr/explore/dataset/que-faire-a-paris-/) |
| API géolocalisation | [ip-api.com](https://ip-api.com/) |

## Fonctionnalités

- [x] Squelette du projet
- [ ] Inscription / connexion / déconnexion
- [ ] Géolocalisation automatique via l'IP
- [ ] Recherche d'événements dans un rayon choisi
- [ ] Parcours des événements (like / unlike)
- [ ] Compteur de likes et unlikes par événement
- [ ] Avis et commentaires des utilisateurs
- [ ] Page de détail d'un événement liké
- [ ] Historique des votes

## Structure du projet

```
find_me_an_event/
├── server/
│   ├── main.go           # Point d'entrée, routeur
│   ├── config.go         # Variables d'environnement
│   ├── db.go             # Connexion PostgreSQL + migrations
│   ├── cors.go           # Middleware CORS
│   ├── utilisateur.go    # Struct Utilisateur, requêtes SQL
│   ├── session.go        # Gestion des sessions (tokens)
│   ├── auth.go           # Handlers inscription/connexion/déconnexion
│   ├── middleware.go      # Vérification de session
│   ├── geolocalisation.go # Client ip-api.com
│   ├── evenements.go     # Client OpenData Paris
│   ├── vote.go           # Modèle Vote, requêtes SQL
│   ├── avis.go           # Modèle Avis, requêtes SQL
│   ├── handler_api.go    # Handlers localisation + événements
│   ├── handler_votes.go  # Handlers votes, avis, historique
│   └── go.mod / go.sum
├── client/
│   ├── src/
│   │   ├── main.jsx      # Point d'entrée React
│   │   ├── App.jsx       # Router + routes protégées
│   │   ├── api.js        # Appels fetch vers le serveur
│   │   ├── AuthContext.jsx # État d'authentification
│   │   ├── app.css       # Styles globaux
│   │   └── pages/
│   │       ├── Connexion.jsx
│   │       ├── Inscription.jsx
│   │       └── Accueil.jsx
│   ├── vite.config.js    # Proxy vers le serveur Go
│   └── package.json
├── docs/
│   └── dossier.md        # Documentation du projet
├── .gitignore
└── README.md
```

## Lancement

### Prérequis
- Go 1.22+
- PostgreSQL 16+

### Base de données
```bash
createdb find_me_an_event
```

### Serveur
```bash
cd server
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/find_me_an_event?sslmode=disable"
go run .
```

### Client React
```bash
cd client
npm install
npm run dev
```

Le client démarre sur `http://localhost:5173` et redirige les appels `/api` vers le serveur Go.

### API disponibles

| Méthode | Route | Description |
|---------|-------|-------------|
| POST | /api/register | Inscription (nom, email, mot_de_passe) |
| POST | /api/login | Connexion (nom, mot_de_passe) |
| POST | /api/logout | Déconnexion |
| GET | /api/me | Profil (authentifié) |
| GET | /api/location | Géolocalisation IP (authentifié) |
| GET | /api/events?lat=X&lon=Y&radius=Z | Événements à proximité (authentifié) |
| GET | /api/events/next?lat=X&lon=Y&radius=Z | Prochain événement non voté (authentifié) |
| POST | /api/events/{id}/vote | Voter like/unlike (authentifié) |
| DELETE | /api/events/{id}/vote | Retirer son vote (authentifié) |
| GET | /api/events/{id}/stats | Compteur likes/unlikes |
| GET | /api/events/{id}/reviews | Avis sur un événement |
| POST | /api/events/{id}/reviews | Poster un avis (authentifié) |
| GET | /api/me/history | Historique des votes (authentifié) |

## Licence

Projet universitaire — PC3R, Sorbonne Université.
