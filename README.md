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
├── server/           # Serveur Go (API REST)
├── client/           # Application React
├── docs/             # Documentation du projet
├── .gitignore
└── README.md
```

## Lancement

> Instructions à venir dans les prochaines phases.

## Licence

Projet universitaire — PC3R, Sorbonne Université.
