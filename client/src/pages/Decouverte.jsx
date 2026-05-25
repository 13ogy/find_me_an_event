import { useState, useEffect, useCallback } from 'react'
import { localisation, prochainEvenement, voter } from '../api'
import { useNavigate } from 'react-router-dom'
import Navigation from '../components/Navigation'
import CarteEvenement from '../components/CarteEvenement'

export default function Decouverte() {
  const navigate = useNavigate()
  const [position, setPosition] = useState(null)
  const [rayon, setRayon] = useState(5)
  const [categorie, setCategorie] = useState('')
  const [evenement, setEvenement] = useState(null)
  const [stats, setStats] = useState(null)
  const [chargement, setChargement] = useState(true)
  const [message, setMessage] = useState('')
  const [voteEnCours, setVoteEnCours] = useState(false)

  // Géoloc IP une seule fois à l'ouverture de la page : le rayon peut
  // changer sans avoir à redemander la position.
  useEffect(() => {
    localisation()
      .then(loc => setPosition(loc))
      .catch(() => setMessage('Impossible de détecter votre position'))
      .finally(() => setChargement(false))
  }, [])

  const chargerSuivant = useCallback(async () => {
    if (!position) return

    setChargement(true)
    setMessage('')
    try {
      const data = await prochainEvenement(position.lat, position.lon, rayon, categorie)
      if (data.evenement) {
        setEvenement(data.evenement)
        setStats(data.stats)
      } else {
        setEvenement(null)
        setStats(null)
        setMessage(data.message || 'Plus d\'événements à découvrir')
      }
    } catch (err) {
      setMessage('Erreur : ' + err.message)
    } finally {
      setChargement(false)
    }
  }, [position, rayon, categorie])

  useEffect(() => {
    chargerSuivant()
  }, [chargerSuivant])

  async function handleVote(typeVote) {
    if (!evenement || voteEnCours) return

    setVoteEnCours(true)
    try {
      await voter(evenement.id, typeVote, evenement.titre)

      if (typeVote === 'like') {
        // On passe evenement et stats dans le state du router : la page de
        // détail évite ainsi un re-fetch des infos déjà connues.
        navigate(`/evenement/${evenement.id}`, {
          state: { evenement, stats }
        })
      } else {
        await chargerSuivant()
      }
    } catch (err) {
      setMessage('Erreur : ' + err.message)
    } finally {
      setVoteEnCours(false)
    }
  }

  return (
    <div className="page-decouverte">
      <Navigation />

      <main className="decouverte-contenu">
        <div className="decouverte-controles">
          <label>
            Rayon :
            <select value={rayon} onChange={(e) => setRayon(Number(e.target.value))}>
              <option value={1}>1 km</option>
              <option value={2}>2 km</option>
              <option value={5}>5 km</option>
              <option value={10}>10 km</option>
              <option value={20}>20 km</option>
            </select>
          </label>

          <label>
            Type :
            <select value={categorie} onChange={(e) => setCategorie(e.target.value)}>
              <option value="">Tous</option>
              <option value="cinéma">Cinéma</option>
              <option value="concert">Concert</option>
              <option value="exposition">Exposition</option>
              <option value="théâtre">Théâtre</option>
              <option value="art">Art</option>
              <option value="sport">Sport</option>
              <option value="famille">Famille</option>
            </select>
          </label>

          {position && (
            <span className="position-info">
              📍 {position.ville || 'Position détectée'}
            </span>
          )}
        </div>

        {chargement && <p className="chargement">Chargement...</p>}

        {!chargement && message && (
          <div className="decouverte-message">
            <p>{message}</p>
            <button onClick={chargerSuivant}>Réessayer</button>
          </div>
        )}

        {!chargement && evenement && (
          <>
            <CarteEvenement evenement={evenement} stats={stats} />

            <div className="boutons-vote">
              <button
                className="btn-unlike"
                onClick={() => handleVote('unlike')}
                disabled={voteEnCours}
              >
                ❌ Passer
              </button>
              <button
                className="btn-like"
                onClick={() => handleVote('like')}
                disabled={voteEnCours}
              >
                ❤️ J'y vais !
              </button>
            </div>
          </>
        )}
      </main>
    </div>
  )
}
