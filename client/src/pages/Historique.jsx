import { useState, useEffect } from 'react'
import { historique } from '../api'
import { useNavigate } from 'react-router-dom'
import Navigation from '../components/Navigation'

export default function Historique() {
  const [votes, setVotes] = useState([])
  const [chargement, setChargement] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    historique()
      .then(setVotes)
      .catch(err => console.error('Erreur historique:', err))
      .finally(() => setChargement(false))
  }, [])

  // On sépare likes et unlikes en deux sections pour la lisibilité plutôt
  // que de tout mélanger chronologiquement.
  const likes = votes.filter(v => v.vote === 'like')
  const unlikes = votes.filter(v => v.vote === 'unlike')

  return (
    <div className="page-historique">
      <Navigation />

      <main className="historique-contenu">
        <h2>Mon historique</h2>

        {chargement && <p>Chargement...</p>}

        {!chargement && votes.length === 0 && (
          <p className="historique-vide">
            Vous n'avez pas encore voté. <button onClick={() => navigate('/')}>Découvrir des événements</button>
          </p>
        )}

        {likes.length > 0 && (
          <section>
            <h3>❤️ Événements likés ({likes.length})</h3>
            <ul className="historique-liste">
              {likes.map(v => (
                <li key={v.id} className="historique-item historique-like">
                  <span className="historique-titre">{v.titre_evenement || v.evenement_id}</span>
                  <span className="historique-date">
                    {new Date(v.cree_le).toLocaleDateString('fr-FR')}
                  </span>
                </li>
              ))}
            </ul>
          </section>
        )}

        {unlikes.length > 0 && (
          <section>
            <h3>❌ Événements passés ({unlikes.length})</h3>
            <ul className="historique-liste">
              {unlikes.map(v => (
                <li key={v.id} className="historique-item historique-unlike">
                  <span className="historique-titre">{v.titre_evenement || v.evenement_id}</span>
                  <span className="historique-date">
                    {new Date(v.cree_le).toLocaleDateString('fr-FR')}
                  </span>
                </li>
              ))}
            </ul>
          </section>
        )}
      </main>
    </div>
  )
}
