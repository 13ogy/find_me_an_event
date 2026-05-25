package main

import (
	"database/sql"
	"time"
)

type Avis struct {
	ID             int       `json:"id"`
	UtilisateurID  int       `json:"utilisateur_id"`
	NomUtilisateur string    `json:"nom_utilisateur"`
	EvenementID    string    `json:"evenement_id"`
	Commentaire    string    `json:"commentaire"`
	Note           int       `json:"note"`
	CreeLe         time.Time `json:"cree_le"`
}

// creerAvis insère l'avis et renvoie l'objet complet (nom auteur inclus) en
// un seul aller-retour, pour que le client puisse l'afficher tout de suite
// sans relancer un GET /reviews.
func creerAvis(db *sql.DB, utilisateurID int, evenementID, commentaire string, note int) (*Avis, error) {
	a := &Avis{}
	err := db.QueryRow(
		`WITH nouveau AS (
		   INSERT INTO avis (utilisateur_id, evenement_id, commentaire, note)
		   VALUES ($1, $2, $3, $4)
		   RETURNING id, utilisateur_id, evenement_id, commentaire, note, cree_le
		 )
		 SELECT n.id, n.utilisateur_id, u.nom, n.evenement_id, n.commentaire, n.note, n.cree_le
		 FROM nouveau n
		 JOIN utilisateurs u ON n.utilisateur_id = u.id`,
		utilisateurID, evenementID, commentaire, note,
	).Scan(&a.ID, &a.UtilisateurID, &a.NomUtilisateur, &a.EvenementID, &a.Commentaire, &a.Note, &a.CreeLe)
	return a, err
}

// listerAvis fait la jointure avec utilisateurs pour récupérer le nom de
// l'auteur dans la même requête : pas de N+1 côté handler.
func listerAvis(db *sql.DB, evenementID string) ([]Avis, error) {
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return liste, nil
}
