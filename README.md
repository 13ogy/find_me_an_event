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
- [x] Inscription / connexion / déconnexion
- [x] Géolocalisation automatique via l'IP
- [x] Recherche d'événements dans un rayon choisi
- [x] Parcours des événements (like / unlike)
- [x] Compteur de likes et unlikes par événement
- [x] Avis et commentaires des utilisateurs
- [x] Page de détail d'un événement liké
- [x] Historique des votes
- [ ] Design responsive + finitions

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
│   │   ├── components/
│   │   │   ├── Navigation.jsx
│   │   │   ├── CarteEvenement.jsx
│   │   │   └── FormulaireAvis.jsx
│   │   └── pages/
│   │       ├── Connexion.jsx
│   │       ├── Inscription.jsx
│   │       ├── Decouverte.jsx    # Swipe like/unlike
│   │       ├── DetailEvenement.jsx
│   │       └── Historique.jsx
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
- PostgreSQL 14+
- Node.js 18+

### Première installation

1. **Démarrer PostgreSQL** (si ce n'est pas déjà fait) :
   ```bash
   brew services start postgresql@14
   ```

2. **Créer le rôle et la base de données** :
   ```bash
   psql -d postgres -c "CREATE ROLE postgres WITH LOGIN SUPERUSER PASSWORD 'postgres';"
   createdb -U postgres find_me_an_event
   ```

3. **Installer les dépendances du client** :
   ```bash
   cd client
   npm install
   ```

### Lancement (usage quotidien)

Dans un premier terminal, lancer le serveur Go :
```bash
cd server
go run .
```

Dans un second terminal, lancer le client React :
```bash
cd client
npm run dev
```

Le serveur démarre sur `http://localhost:8080` (tables créées automatiquement).
Le client démarre sur `http://localhost:5173` et redirige les appels `/api` vers le serveur Go.

> **Note :** si PostgreSQL n'est pas démarré (ex. après un redémarrage), lancez d'abord :
> ```bash
> brew services start postgresql@14
> ```

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
