package main

import (
	"database/sql"
	"time"
)

// Avis représente un commentaire d'un utilisateur sur un événement
type Avis struct {
	ID            int       `json:"id"`
	UtilisateurID int       `json:"utilisateur_id"`
	NomUtilisateur string   `json:"nom_utilisateur"`
	EvenementID   string    `json:"evenement_id"`
	Commentaire   string    `json:"commentaire"`
	Note          int       `json:"note"`
	CreeLe        time.Time `json:"cree_le"`
}

// creerAvis insère un nouvel avis en base
func creerAvis(db *sql.DB, utilisateurID int, evenementID, commentaire string, note int) (*Avis, error) {
	a := &Avis{}
	err := db.QueryRow(
		`INSERT INTO avis (utilisateur_id, evenement_id, commentaire, note)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, utilisateur_id, evenement_id, commentaire, note, cree_le`,
		utilisateurID, evenementID, commentaire, note,
	).Scan(&a.ID, &a.UtilisateurID, &a.EvenementID, &a.Commentaire, &a.Note, &a.CreeLe)
	return a, err
}

// listerAvis renvoie tous les avis pour un événement donné
func listerAvis(db *sql.DB, evenementID string) ([]Avis, error) {
	// Jointure pour récupérer le nom de l'auteur sans requête supplémentaire
	rows, err := db.Query(
		`SELECT a.id, a.utilisateur_id, u.nom, a.evenement_id, a.commentaire, a.note, a.cree_le
		 FROM avis a
		 JOIN utilisateurs u ON a.utilisateur_id = u.id
		 WHERE a.evenement_id = $1
		 ORDER BY a.cree_le DESC`,
		evenementID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var liste []Avis
	for rows.Next() {
		var a Avis
		if err := rows.Scan(&a.ID, &a.UtilisateurID, &a.NomUtilisateur, &a.EvenementID, &a.Commentaire, &a.Note, &a.CreeLe); err != nil {
			return nil, err
		}
		liste = append(liste, a)
	}
	return liste, nil
}
