# Dossier du projet — Find Me An Event

## 1. Sujet de l'application

Find Me An Event est une application web qui aide un utilisateur connecté à
découvrir des événements culturels et festifs qui ont lieu près de chez lui à
Paris.

Le principe est inspiré des applications de type *swipe* : l'application
détecte la position de l'utilisateur via son adresse IP, lui propose des
événements un par un dans le rayon choisi, et l'utilisateur exprime son
intérêt en likant ou en passant. Quand il like, il accède à la fiche
détaillée de l'événement avec les informations pratiques, les avis laissés
par les autres utilisateurs, et peut publier le sien.

## 2. API web externes

### 2.1 OpenData Paris — *Que faire à Paris*

- **URL** : <https://opendata.paris.fr/explore/dataset/que-faire-a-paris-/>
- **Contenu** : agenda officiel des événements parisiens (concerts,
  expositions, ateliers, spectacles, festivals…). Chaque enregistrement
  contient un titre, une description, des dates de début et de fin, un lieu
  (nom, adresse, coordonnées géographiques), un prix, une URL de contact, une
  image de couverture, etc.
- **Mise à jour** : quotidienne, par la Mairie de Paris.
- **Accès** : libre, pas de clé API requise. Une instance ODSQL (variante de
  SQL pour OpenData Soft) permet de filtrer côté serveur.
- **Comment l'application l'utilise** : à chaque demande d'événement, le
  serveur Go construit une requête ODSQL qui filtre par distance géographique
  (`within_distance`) et par date (`date_end >= today`), trie par date de
  début croissante, et limite le nombre de résultats. Les enregistrements
  retournés sont mappés vers la structure `Evenement` que consomme le client.
- **Exemple de requête HTTP** :
  ```
  GET https://opendata.paris.fr/api/explore/v2.1/catalog/datasets/que-faire-a-paris-/records
    ?where=within_distance(lat_lon, GEOM'POINT(2.3488 48.8534)', 2km)
       AND date_start IS NOT NULL
       AND date_end >= '2026-04-08'
    &order_by=date_start ASC
    &limit=50
  ```

### 2.2 ip-api.com

- **URL** : <https://ip-api.com/>
- **Contenu** : géolocalisation à partir d'une adresse IP (latitude,
  longitude, ville, pays).
- **Accès** : libre en HTTP, 45 requêtes/minute pour le tier gratuit (très
  largement suffisant pour ce projet).
- **Utilisation** : à l'arrivée sur la page Découverte, le serveur extrait
  l'IP du client (en-têtes `X-Forwarded-For`/`X-Real-IP` puis `RemoteAddr`)
  et la transmet à ip-api pour obtenir la position. En local, l'IP est
  privée et ip-api renvoie une erreur : on retombe alors sur Paris
  (48.8566 / 2.3522) pour que l'application reste utilisable en
  développement.
- **Exemple de requête HTTP** :
  ```
  GET http://ip-api.com/json/8.8.8.8
  → { "status": "success", "lat": 37.4056, "lon": -122.0775,
      "city": "Mountain View", "country": "United States" }
  ```

## 3. Fonctionnalités

1. **Inscription / connexion / déconnexion** — création d'un compte avec
   pseudonyme, email et mot de passe (validé à 6 caractères minimum),
   authentification par session côté serveur, cookie HttpOnly.
2. **Géolocalisation IP automatique** — position détectée sans intervention
   de l'utilisateur.
3. **Choix du rayon de recherche** — sélecteur 1 / 2 / 5 / 10 / 20 km.
4. **Filtre par catégorie** — menu déroulant (cinéma, concert, exposition,
   théâtre, art, sport, famille) qui restreint le flux de découverte aux
   événements dont le mot-clé apparaît dans le titre, la description ou le
   chapeau.
5. **Découverte d'événements façon swipe** — un événement à la fois, deux
   boutons « ❌ Passer » et « ❤️ J'y vais ! ».
5. **Filtrage côté serveur des événements déjà votés** — l'utilisateur ne
   revoit jamais deux fois le même événement.
6. **Statistiques sociales** — pour chaque événement on affiche le nombre
   total de likes et de unlikes de la communauté.
7. **Page de détail** — accessible après un like : informations pratiques
   complètes, lien Google Maps, lien vers le site officiel, avis et
   formulaire d'avis.
8. **Avis et commentaires** — un utilisateur peut publier un commentaire et
   une note de 1 à 5 sur un événement ; tous les avis sont listés du plus
   récent au plus ancien.
9. **Historique personnel** — liste des événements likés et passés, classés
   du plus récent au plus ancien.
10. **Design responsive** — interface utilisable sur mobile et desktop.

## 4. Cas d'utilisation

### 4.1 Alice découvre un concert

Alice se connecte à l'application. Sa position est détectée
automatiquement (Paris 11e). Elle laisse le rayon par défaut (5 km).
L'application lui propose un événement : *Festival Jazz à la Villette*.
La carte affiche le titre, le chapeau, le lieu, les dates et le compteur
de votes (23 likes, 4 unlikes). Alice clique sur ❤️ J'y vais. Le serveur
enregistre son like, et le client la redirige vers la page de détail. Elle
y lit la description complète, consulte l'adresse, suit le lien Google
Maps, et publie un avis : « J'y vais avec des amis, hâte ! » avec une note
de 5/5. Elle revient ensuite sur la page Découverte pour voir l'événement
suivant.

### 4.2 Bob enchaîne les passes

Bob se connecte et resserre le rayon à 1 km. Il parcourt trois événements :
il passe les deux premiers (une expo qu'il a déjà vue, un atelier cuisine
trop cher), puis like le troisième (un marché nocturne). Sur la page de
détail il lit les avis : un utilisateur a noté 4/5 en disant « ambiance
sympa, prévoir liquide ». Il revient à la découverte mais il a épuisé les
événements de son rayon ; l'application affiche un message d'invite à
élargir la zone ou à revenir plus tard.

### 4.3 Carla revoit son historique

Carla, déjà inscrite depuis quelques semaines, va dans l'onglet Historique.
Elle y voit ses likes et ses unlikes regroupés par catégorie, avec les
dates de vote. Elle se souvient d'un concert liké il y a deux semaines et
clique pour rafraîchir sa mémoire.

## 5. Données stockées

Quatre tables PostgreSQL (créées automatiquement au démarrage du serveur via
`CREATE TABLE IF NOT EXISTS`).

### `utilisateurs`

| Champ | Type | Description |
|-------|------|-------------|
| id | SERIAL PK | Identifiant numérique |
| nom | VARCHAR(50) UNIQUE | Pseudonyme |
| email | VARCHAR(100) UNIQUE | Email |
| mot_de_passe_hash | VARCHAR(255) | Hash bcrypt (jamais le mot de passe en clair) |
| cree_le | TIMESTAMP | Date d'inscription |

### `sessions`

| Champ | Type | Description |
|-------|------|-------------|
| id | VARCHAR(64) PK | Token de session (256 bits aléatoires, hex) |
| utilisateur_id | INTEGER FK utilisateurs | Propriétaire de la session |
| expire_le | TIMESTAMP | Expiration (24 h après création) |
| cree_le | TIMESTAMP | Date de création |

### `votes`

| Champ | Type | Description |
|-------|------|-------------|
| id | SERIAL PK | Identifiant |
| utilisateur_id | INTEGER FK utilisateurs | Auteur du vote |
| evenement_id | VARCHAR(255) | ID OpenData de l'événement |
| titre_evenement | VARCHAR(500) | Titre dénormalisé pour l'historique |
| vote | VARCHAR(10) CHECK in ('like','unlike') | Type de vote |
| cree_le | TIMESTAMP | Date du vote |
| **Contrainte** | UNIQUE(utilisateur_id, evenement_id) | Un seul vote par couple |

### `avis`

| Champ | Type | Description |
|-------|------|-------------|
| id | SERIAL PK | Identifiant |
| utilisateur_id | INTEGER FK utilisateurs | Auteur de l'avis |
| evenement_id | VARCHAR(255) | ID OpenData |
| commentaire | TEXT NOT NULL | Texte de l'avis |
| note | INTEGER CHECK BETWEEN 1 AND 5 | Note sur 5 |
| cree_le | TIMESTAMP | Date de publication |

**Choix de modélisation.** Les événements en eux-mêmes ne sont pas stockés
en base : on garde uniquement leur identifiant OpenData et leur titre (pour
l'historique). Cela évite d'avoir à synchroniser une copie locale du
catalogue, qui change tous les jours.

## 6. Mise à jour des données et appels à l'API externe

Le projet n'utilise pas de tâche planifiée. Tous les appels aux APIs
externes sont **déclenchés à la demande** par les actions de l'utilisateur :

- **OpenData Paris** est appelée :
  - À chaque chargement de la page Découverte (via
    `GET /api/events/next`).
  - À chaque clic sur ❌ Passer ou Réessayer, qui recharge un événement.
  - Côté serveur on demande systématiquement un lot de 50 événements, puis
    on filtre en mémoire ceux que l'utilisateur a déjà votés.
- **ip-api.com** est appelée une fois par session de découverte, à
  l'ouverture de la page Découverte (via `GET /api/location`). Le résultat
  est gardé dans le state React pour la durée de la navigation.

**Pourquoi pas un cache local ?** Le catalogue OpenData évolue tous les
jours (nouveaux événements, événements terminés). Cacher en base
exposerait l'utilisateur à des résultats périmés. Le filtrage géographique
côté API est efficace et la latence nous a semblé acceptable.

## 7. Architecture serveur

### 7.1 Approche

Approche **Ressources / REST**, sans framework web (contrainte du sujet) :
on utilise uniquement le package `net/http` de la bibliothèque standard,
avec le routeur `ServeMux` enrichi de Go 1.22 (méthodes HTTP et paramètres
de chemin directement dans les patterns).

Le serveur n'émet jamais de HTML : il ne sert que du JSON. Le client React,
servi par Vite en développement et statiquement en production, consomme
l'API par appels `fetch` asynchrones.

### 7.2 Choix techniques

- **Langage** : Go 1.25 — typage statique, compilation rapide, bibliothèque
  standard riche pour HTTP.
- **Base de données** : PostgreSQL via `database/sql` + pilote `lib/pq`,
  toutes les requêtes utilisent des paramètres préparés (`$1`, `$2`…), ce
  qui élimine les injections SQL.
- **Mot de passe** : `golang.org/x/crypto/bcrypt` avec coût 10.
- **Sessions** : tokens 256 bits générés par `crypto/rand`, stockés en base
  avec date d'expiration. Pas de JWT : un token opaque côté client suffit
  et permet une révocation immédiate (DELETE en base).
- **CORS** : middleware maison qui reflète l'`Origin` du client et autorise
  les credentials, pour que le cookie de session traverse le proxy Vite en
  développement.

### 7.3 Composants (handlers)

Chaque handler est une fonction `http.HandlerFunc` ou un constructeur qui
en renvoie une (closure capturant `*sql.DB`). Tous les handlers protégés
passent par le middleware `proteger`, qui vérifie le cookie de session et
injecte l'ID utilisateur dans le `context.Context` de la requête.

| Composant | Fichier | Rôle |
|-----------|---------|------|
| Routeur | `main.go` | Câblage des routes, démarrage du serveur, application de CORS. |
| Config | `config.go` | Lecture des variables d'environnement. |
| DB | `db.go` | Ouverture de la connexion et migration des tables. |
| CORS | `cors.go` | Middleware d'autorisation cross-origin. |
| Modèle Utilisateur | `utilisateur.go` | Requêtes CRUD utilisateurs. |
| Sessions | `session.go` | Génération, validation, suppression de tokens. |
| Auth | `auth.go` | Handlers `/register`, `/login`, `/logout`, `/me`. |
| Middleware | `middleware.go` | Garde `proteger` pour les routes authentifiées. |
| Géoloc | `geolocalisation.go` | Client ip-api et extraction d'IP. |
| Événements | `evenements.go` | Client OpenData Paris et conversion des records. |
| Votes | `vote.go` | Requêtes SQL pour les votes (UPSERT, stats, historique). |
| Avis | `avis.go` | Requêtes SQL pour les avis. |
| Handlers API | `handler_api.go` | `/location`, `/events`. |
| Handlers Votes | `handler_votes.go` | `/events/{id}/vote`, `/stats`, `/reviews`, `/me/history`, `/events/next`. |

### 7.4 Endpoints

| Méthode | Route | Auth | Description |
|---------|-------|------|-------------|
| POST | `/api/register` | non | Crée un utilisateur |
| POST | `/api/login` | non | Ouvre une session (cookie) |
| POST | `/api/logout` | non | Ferme la session |
| GET | `/api/me` | oui | Renvoie le profil |
| GET | `/api/location` | oui | Géolocalise l'IP du client |
| GET | `/api/events?lat=&lon=&radius=&categorie=` | oui | Liste les événements dans le rayon, filtre par mot-clé optionnel |
| GET | `/api/events/next?lat=&lon=&radius=&categorie=` | oui | Renvoie le prochain événement non voté |
| POST | `/api/events/{id}/vote` | oui | Enregistre un like ou un unlike |
| DELETE | `/api/events/{id}/vote` | oui | Retire son vote |
| GET | `/api/events/{id}/stats` | oui | Compte des likes / unlikes |
| GET | `/api/events/{id}/reviews` | oui | Liste les avis (avec nom auteur) |
| POST | `/api/events/{id}/reviews` | oui | Publie un avis |
| GET | `/api/me/history` | oui | Historique des votes |

## 8. Architecture client

### 8.1 Technologies utilisées

- **React 19** (encouragé par le sujet) — composants fonctionnels, hooks,
  état local + contexte pour l'authentification.
- **Vite** — bundler de développement et de production. Le serveur dev
  proxifie `/api` vers le serveur Go (port 8080), ce qui évite les problèmes
  CORS en local et permet de garder dans le code client des chemins
  relatifs (`/api/...`) sans URL absolue.
- **React Router 7** — application monopage avec routes protégées :
  `RouteProtegee` et `RoutePublique` redirigent selon l'état
  d'authentification (et attendent la fin de la reprise de session avant
  de décider).
- **CSS vanille** — un seul fichier `app.css` avec des classes
  parlantes ; pas de framework UI, responsive grâce à des media queries ciblées.

### 8.2 Plan du site (SPA)

```
/connexion        Connexion (publique)
/inscription      Inscription (publique)
/                 Découverte d'événements (protégée)
/evenement/:id    Détail d'un événement liké (protégée)
/historique       Historique des votes (protégée)
*                 404
```

### 8.3 Écrans et enchaînement

| Écran | Composants | Appels serveur déclenchés |
|-------|------------|---------------------------|
| Connexion | formulaire login | `POST /api/login` |
| Inscription | formulaire | `POST /api/register` + `POST /api/login` |
| Découverte | `CarteEvenement`, boutons vote | `GET /api/location` (au montage), `GET /api/events/next` (au montage et à chaque vote unlike), `POST /api/events/{id}/vote` |
| Détail événement | `FormulaireAvis`, liste d'avis | `GET /api/events/{id}/stats` + `GET /api/events/{id}/reviews` (parallèles au montage), `POST /api/events/{id}/reviews` (à la publication d'un avis) |
| Historique | liste likes / unlikes | `GET /api/me/history` (au montage) |

L'enchaînement classique est : Inscription → Découverte → like → Détail →
publication d'un avis → retour à Découverte → unlike → … → Historique.

### 8.4 Liste des appels AJAX

Tous les appels passent par `appel(methode, chemin, body)` défini dans
`src/api.js`, qui :
- ajoute systématiquement `credentials: 'include'` pour envoyer le cookie
  de session,
- pose `Content-Type: application/json` et sérialise le corps si présent,
- lève une `Error` contenant le champ `erreur` de la réponse si le statut
  HTTP n'est pas 2xx.

Fonctions exportées : `inscription`, `connexion`, `deconnexion`, `profil`,
`localisation`, `evenements`, `prochainEvenement`, `voter`,
`supprimerVote`, `statsEvenement`, `listerAvis`, `posterAvis`,
`historique`.

## 9. Exemples de requêtes / réponses

### Inscription

```
POST /api/register
Content-Type: application/json

{ "nom": "alice", "email": "alice@example.com", "mot_de_passe": "secret123" }

→ 201 Created
{ "id": 1, "nom": "alice", "email": "alice@example.com", "cree_le": "2026-04-09T18:21:03Z" }
```

### Connexion

```
POST /api/login
Content-Type: application/json

{ "nom": "alice", "mot_de_passe": "secret123" }

→ 200 OK
Set-Cookie: session=<token>; HttpOnly; SameSite=Lax; Max-Age=86400; Path=/
{ "id": 1, "nom": "alice", "email": "alice@example.com", "cree_le": "..." }
```

### Géolocalisation

```
GET /api/location
Cookie: session=<token>

→ 200 OK
{ "lat": 48.8566, "lon": 2.3522, "ville": "Paris", "pays": "France" }
```

### Prochain événement non voté

```
GET /api/events/next?lat=48.85&lon=2.35&radius=5
Cookie: session=<token>

→ 200 OK
{
  "evenement": {
    "id": "abc123",
    "titre": "Festival Jazz à la Villette",
    "chapeau": "Un festival de jazz en plein air",
    "date_debut": "2026-06-15T18:00:00+00:00",
    "date_fin":   "2026-06-15T23:00:00+00:00",
    "nom_lieu":   "Parc de la Villette",
    "adresse":    "211 Avenue Jean Jaurès",
    "ville":      "Paris",
    "code_postal":"75019",
    "lat": 48.8938, "lon": 2.3913,
    "prix":       "gratuit",
    "cover_url":  "https://..."
  },
  "stats": { "evenement_id": "abc123", "likes": 23, "unlikes": 4 }
}
```

Quand tous les événements ont déjà été vus :
```
→ 200 OK
{ "evenement": null, "message": "plus d'événements à découvrir dans ce rayon" }
```

### Vote

```
POST /api/events/abc123/vote
Cookie: session=<token>
Content-Type: application/json

{ "vote": "like", "titre": "Festival Jazz" }

→ 200 OK
{ "message": "vote enregistré" }
```

### Publier un avis

```
POST /api/events/abc123/reviews
Cookie: session=<token>
Content-Type: application/json

{ "commentaire": "Super ambiance !", "note": 5 }

→ 201 Created
{
  "id": 1, "utilisateur_id": 1, "evenement_id": "abc123",
  "commentaire": "Super ambiance !", "note": 5, "cree_le": "2026-04-09T19:00:00Z"
}
```

### Historique

```
GET /api/me/history
Cookie: session=<token>

→ 200 OK
[
  { "id": 7, "utilisateur_id": 1, "evenement_id": "abc123",
    "titre_evenement": "Festival Jazz", "vote": "like",
    "cree_le": "2026-04-09T18:30:00Z" },
  { "id": 6, "utilisateur_id": 1, "evenement_id": "xyz789",
    "titre_evenement": "Atelier cuisine", "vote": "unlike",
    "cree_le": "2026-04-09T18:21:00Z" }
]
```

## 10. Sécurité

- **Mots de passe** stockés avec `bcrypt`.
- **Sessions** : tokens cryptographiquement aléatoires (256 bits), expirés
  après 24 h, supprimés à la déconnexion. Validation et révocation
  immédiates en base.
- **Cookies** : `HttpOnly` (inaccessibles depuis JavaScript, ce qui les
  protège d'un XSS) et `SameSite=Lax` (qui bloque les CSRF basiques).
- **Injection SQL** : toutes les requêtes utilisent des paramètres préparés
  (`$1`, `$2`…), aucune concaténation de chaîne.
- **Validation d'entrée** : noms et emails trimés, mots de passe à 6 cars
  minimum, vote restreint à `like` / `unlike` (liste blanche), notes
  bornées à 1–5 côté serveur en plus de la contrainte `CHECK` SQL.
- **Bornes API** : rayon plafonné à 50 km côté serveur pour éviter qu'un
  client malveillant ne déclenche d'énormes appels OpenData.
- **CORS** : reflète l'`Origin` au lieu d'autoriser `*`, ce qui reste
  compatible avec `Access-Control-Allow-Credentials: true`.

## 11. Schéma global

```
   ┌────────────────┐   AJAX/JSON, cookie de session   ┌──────────────────────┐
   │  Client React  │ ◄──────────────────────────────► │  Serveur Go          │
   │  + Vite        │                                   │  net/http  :8080     │
   │  :5173 (dev)   │                                   │  (sans framework)    │
   └────────────────┘                                   └──────────┬───────────┘
                                                                   │
                                       ┌───────────────────────────┼──────────────────────────┐
                                       │                           │                          │
                                       ▼                           ▼                          ▼
                              ┌──────────────────┐      ┌──────────────────┐      ┌────────────────────┐
                              │  PostgreSQL      │      │  OpenData Paris  │      │  ip-api.com        │
                              │  utilisateurs    │      │  événements      │      │  géoloc IP         │
                              │  sessions        │      │  HTTPS, sans clé │      │  HTTP, sans clé    │
                              │  votes / avis    │      └──────────────────┘      └────────────────────┘
                              └──────────────────┘
```

Le serveur Go centralise tous les appels externes : le client React n'appelle
jamais directement OpenData ni ip-api, ce qui évite d'exposer la logique
métier au navigateur.
