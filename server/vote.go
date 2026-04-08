package main

import (
	"database/sql"
	"time"
)

// Vote représente un like ou unlike d'un utilisateur sur un événement
type Vote struct {
	ID              int       `json:"id"`
	UtilisateurID   int       `json:"utilisateur_id"`
	EvenementID     string    `json:"evenement_id"`
	TitreEvenement  string    `json:"titre_evenement"`
	Vote            string    `json:"vote"`
	CreeLe          time.Time `json:"cree_le"`
}

// StatsVotes contient le décompte des likes et unlikes pour un événement
type StatsVotes struct {
	EvenementID string `json:"evenement_id"`
	Likes       int    `json:"likes"`
	Unlikes     int    `json:"unlikes"`
}

// voterPourEvenement crée ou met à jour un vote (UPSERT)
func voterPourEvenement(db *sql.DB, utilisateurID int, evenementID, titre, typeVote string) error {
	_, err := db.Exec(
		`INSERT INTO votes (utilisateur_id, evenement_id, titre_evenement, vote)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (utilisateur_id, evenement_id)
		 DO UPDATE SET vote = $4, titre_evenement = $3`,
		utilisateurID, evenementID, titre, typeVote,
	)
	return err
}

// supprimerVote retire le vote d'un utilisateur sur un événement
func supprimerVote(db *sql.DB, utilisateurID int, evenementID string) error {
	_, err := db.Exec(
		`DELETE FROM votes WHERE utilisateur_id = $1 AND evenement_id = $2`,
		utilisateurID, evenementID,
	)
	return err
}

// obtenirStatsVotes compte les likes et unlikes d'un événement
func obtenirStatsVotes(db *sql.DB, evenementID string) (*StatsVotes, error) {
	stats := &StatsVotes{EvenementID: evenementID}

	err := db.QueryRow(
		`SELECT COALESCE(SUM(CASE WHEN vote = 'like' THEN 1 ELSE 0 END), 0),
		        COALESCE(SUM(CASE WHEN vote = 'unlike' THEN 1 ELSE 0 END), 0)
		 FROM votes WHERE evenement_id = $1`,
		evenementID,
	).Scan(&stats.Likes, &stats.Unlikes)

	return stats, err
}

// obtenirHistorique renvoie les votes d'un utilisateur, du plus récent au plus ancien
func obtenirHistorique(db *sql.DB, utilisateurID int) ([]Vote, error) {
	rows, err := db.Query(
		`SELECT id, utilisateur_id, evenement_id, titre_evenement, vote, cree_le
		 FROM votes WHERE utilisateur_id = $1
		 ORDER BY cree_le DESC`,
		utilisateurID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []Vote
	for rows.Next() {
		var v Vote
		if err := rows.Scan(&v.ID, &v.UtilisateurID, &v.EvenementID, &v.TitreEvenement, &v.Vote, &v.CreeLe); err != nil {
			return nil, err
		}
		votes = append(votes, v)
	}
	return votes, nil
}

// evenementsDejaVotes renvoie les IDs des événements déjà votés par un utilisateur
func evenementsDejaVotes(db *sql.DB, utilisateurID int) (map[string]bool, error) {
	rows, err := db.Query(
		`SELECT evenement_id FROM votes WHERE utilisateur_id = $1`,
		utilisateurID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = true
	}
	return ids, nil
}
