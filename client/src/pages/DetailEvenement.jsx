import { useState, useEffect } from 'react'
import { useParams, useLocation, useNavigate } from 'react-router-dom'
import { statsEvenement, listerAvis } from '../api'
import Navigation from '../components/Navigation'
import FormulaireAvis from '../components/FormulaireAvis'

// Page détail — affichée après un like, avec infos complètes + avis
export default function DetailEvenement() {
  const { id } = useParams()
  const location = useLocation()
  const navigate = useNavigate()

  // Récupère l'événement depuis le state de navigation (passé par Decouverte)
  const evenement = location.state?.evenement

  const [stats, setStats] = useState(location.state?.stats || null)
  const [avis, setAvis] = useState([])
  const [chargement, setChargement] = useState(true)

  useEffect(() => {
    async function charger() {
      try {
        const [s, a] = await Promise.all([
          statsEvenement(id),
          listerAvis(id)
        ])
        setStats(s)
        setAvis(a)
      } catch (err) {
        console.error('Erreur chargement détail:', err)
      } finally {
        setChargement(false)
      }
    }
    charger()
  }, [id])

  // Lien Google Maps vers la localisation exacte
  const lienGoogleMaps = evenement?.lat && evenement?.lon
    ? `https://www.google.com/maps?q=${evenement.lat},${evenement.lon}`
    : null

  const formatDate = (iso) => {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString('fr-FR', {
      day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit'
    })
  }

  // Quand aucun événement dans le state (accès direct par URL)
  if (!evenement) {
    return (
      <div className="page-detail">
        <Navigation />
        <main className="detail-contenu">
          <p>Événement non trouvé. <button onClick={() => navigate('/')}>Retour</button></p>
        </main>
      </div>
    )
  }

  return (
    <div className="page-detail">
      <Navigation />

      <main className="detail-contenu">
        {evenement.cover_url && (
          <img
            className="detail-image"
            src={evenement.cover_url}
            alt={evenement.titre}
            onError={(e) => { e.target.style.display = 'none' }}
          />
        )}

        <h2>{evenement.titre}</h2>

        {stats && (
          <div className="carte-stats">
            <span className="stat-like">❤️ {stats.likes} likes</span>
            <span className="stat-unlike">❌ {stats.unlikes} unlikes</span>
          </div>
        )}

        {evenement.description && (
          <div
            className="detail-description"
            dangerouslySetInnerHTML={{ __html: evenement.description }}
          />
        )}

        <section className="detail-infos">
          <h3>Informations pratiques</h3>
          {evenement.nom_lieu && <p>📍 <strong>{evenement.nom_lieu}</strong></p>}
          {evenement.adresse && <p>{evenement.adresse}, {evenement.code_postal} {evenement.ville}</p>}
          <p>📅 Du {formatDate(evenement.date_debut)} au {formatDate(evenement.date_fin)}</p>
          {evenement.prix && <p>💰 {evenement.prix}</p>}
          {evenement.audience && <p>👥 {evenement.audience}</p>}
          {evenement.url_contact && (
            <p><a href={evenement.url_contact} target="_blank" rel="noopener noreferrer">🔗 Site officiel</a></p>
          )}
          {lienGoogleMaps && (
            <p><a href={lienGoogleMaps} target="_blank" rel="noopener noreferrer">🗺️ Voir sur Google Maps</a></p>
          )}
        </section>

        <section className="detail-avis">
          <h3>Avis ({avis.length})</h3>

          <FormulaireAvis
            evenementId={id}
            onAvisCree={(nouvel) => setAvis([nouvel, ...avis])}
          />

          {chargement && <p>Chargement des avis...</p>}

          {avis.length === 0 && !chargement && (
            <p className="avis-vide">Aucun avis pour le moment. Soyez le premier !</p>
          )}

          <div className="liste-avis">
            {avis.map(a => (
              <div key={a.id} className="avis-item">
                <div className="avis-entete">
                  <strong>{a.nom_utilisateur}</strong>
                  <span>{'⭐'.repeat(a.note)}</span>
                  <span className="avis-date">
                    {new Date(a.cree_le).toLocaleDateString('fr-FR')}
                  </span>
                </div>
                <p>{a.commentaire}</p>
              </div>
            ))}
          </div>
        </section>

        <div className="detail-actions">
          <button className="btn-continuer" onClick={() => navigate('/')}>
            ← Continuer à découvrir
          </button>
        </div>
      </main>
    </div>
  )
}
