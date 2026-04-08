# Dossier du projet — Find Me An Event

## Sujet

Application web de découverte d'événements à Paris. L'utilisateur se connecte, l'application détecte sa localisation via son adresse IP, et lui propose des événements proches. Il peut parcourir les événements un par un, voter (like/unlike), consulter les avis d'autres utilisateurs et laisser des commentaires sur les événements qui l'intéressent.

## API externe

### OpenData Paris — "Que faire à Paris"

- **URL** : https://opendata.paris.fr/explore/dataset/que-faire-a-paris-/
- **Contenu** : événements culturels, sportifs, festifs à Paris (concerts, expos, spectacles, ateliers…)
- **Mise à jour** : quotidienne
- **Utilisation** : recherche d'événements dans un rayon autour de la position de l'utilisateur, filtrés par date
- **Exemple de requête** :
  ```
  GET https://opendata.paris.fr/api/explore/v2.1/catalog/datasets/que-faire-a-paris-/records
    ?where=within_distance(lat_lon, GEOM'POINT(2.3488 48.8534)', 2km)
      AND date_start IS NOT NULL
      AND date_end >= "2026-04-08"
    &limit=10
  ```

### ip-api.com

- **URL** : https://ip-api.com/
- **Contenu** : géolocalisation à partir d'une adresse IP (latitude, longitude, ville, pays)
- **Utilisation** : détecter la position de l'utilisateur à la connexion
- **Exemple** :
  ```
  GET http://ip-api.com/json/
  → { "lat": 48.8566, "lon": 2.3522, "city": "Paris", ... }
  ```

## Fonctionnalités

1. **Inscription et connexion** : création de compte avec nom d'utilisateur, email et mot de passe. Authentification par session (cookies).
2. **Géolocalisation automatique** : à la connexion, l'IP de l'utilisateur est utilisée pour déterminer sa position.
3. **Choix du rayon** : l'utilisateur choisit un rayon de recherche (1km à 20km).
4. **Découverte d'événements** : les événements sont présentés un par un. L'utilisateur vote like ou unlike.
5. **Statistiques sociales** : chaque événement affiche le nombre total de likes et unlikes des autres utilisateurs.
6. **Page de détail** : quand l'utilisateur like un événement, il accède à une page avec les informations complètes (adresse, horaires, prix) et les avis.
7. **Avis et commentaires** : l'utilisateur peut laisser un commentaire et une note sur un événement.
8. **Historique** : l'utilisateur peut consulter la liste de ses votes passés.

## Cas d'utilisation

### Alice découvre un concert

Alice se connecte à l'application. Sa position est automatiquement détectée (Paris, 11e). Elle choisit un rayon de 5km. L'application lui présente un événement : "Festival Jazz à la Villette". Elle voit que 23 utilisateurs ont liké et 4 ont unliké. Elle clique sur "Like". L'application l'emmène sur la page de détail où elle voit l'adresse exacte, les horaires, et deux avis positifs d'autres utilisateurs. Elle laisse un commentaire : "J'y vais avec des amis, hâte !" puis clique sur "Continuer" pour voir d'autres événements.

### Bob cherche une sortie rapide

Bob se connecte et choisit un rayon de 1km. Il parcourt trois événements : il unlike les deux premiers (une expo qu'il a déjà vue et un atelier cuisine trop cher), puis like le troisième (un marché nocturne). Sur la page de détail, il consulte l'adresse sur Google Maps.

## Données stockées

### Table `utilisateurs`
| Champ | Type | Description |
|-------|------|-------------|
| id | SERIAL | Identifiant unique |
| nom | VARCHAR(50) | Nom d'utilisateur |
| email | VARCHAR(100) | Email |
| mot_de_passe_hash | VARCHAR(255) | Mot de passe hashé (bcrypt) |
| cree_le | TIMESTAMP | Date de création |

### Table `sessions`
| Champ | Type | Description |
|-------|------|-------------|
| id | VARCHAR(64) | Token de session |
| utilisateur_id | INTEGER | Référence vers utilisateurs |
| expire_le | TIMESTAMP | Date d'expiration |

### Table `votes`
| Champ | Type | Description |
|-------|------|-------------|
| id | SERIAL | Identifiant unique |
| utilisateur_id | INTEGER | Référence vers utilisateurs |
| evenement_id | VARCHAR(255) | ID OpenData de l'événement |
| titre_evenement | VARCHAR(500) | Titre (cache local) |
| vote | VARCHAR(10) | 'like' ou 'unlike' |
| cree_le | TIMESTAMP | Date du vote |

### Table `avis`
| Champ | Type | Description |
|-------|------|-------------|
| id | SERIAL | Identifiant unique |
| utilisateur_id | INTEGER | Référence vers utilisateurs |
| evenement_id | VARCHAR(255) | ID OpenData de l'événement |
| commentaire | TEXT | Texte de l'avis |
| note | INTEGER | Note de 1 à 5 |
| cree_le | TIMESTAMP | Date de l'avis |

## Architecture serveur

Approche **Ressources (REST)**. Le serveur ne génère pas de HTML — il sert uniquement du JSON. Le client React effectue des appels AJAX (fetch) et construit l'interface.

### Endpoints

| Méthode | Route | Description |
|---------|-------|-------------|
| POST | /api/register | Inscription |
| POST | /api/login | Connexion |
| POST | /api/logout | Déconnexion |
| GET | /api/me | Profil de l'utilisateur connecté |
| GET | /api/location | Géolocalisation IP |
| GET | /api/events?lat=X&lon=Y&radius=Z | Événements dans le rayon |
| GET | /api/events/next?lat=X&lon=Y&radius=Z | Prochain événement non voté |
| POST | /api/events/{id}/vote | Voter like/unlike |
| DELETE | /api/events/{id}/vote | Retirer son vote |
| GET | /api/events/{id}/stats | Statistiques de votes |
| GET | /api/events/{id}/reviews | Avis sur un événement |
| POST | /api/events/{id}/reviews | Poster un avis |
| GET | /api/me/history | Historique des votes |

## Architecture client

Application **monopage (SPA)** en React.

### Écrans
1. **Connexion** — formulaire login/mot de passe
2. **Inscription** — formulaire nom/email/mot de passe
3. **Découverte** — carte d'événement + boutons like/unlike + choix du rayon
4. **Détail événement** — infos complètes, avis, formulaire de commentaire
5. **Historique** — liste des événements likés/unlikés

### Appels AJAX
Tous les appels utilisent `fetch()` avec `credentials: 'include'` pour envoyer le cookie de session. Les réponses sont en JSON.
