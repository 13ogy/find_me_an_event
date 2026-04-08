// Carte d'événement — affiche les infos principales pour le swipe
export default function CarteEvenement({ evenement, stats }) {
  if (!evenement) return null

  // Formatage lisible des dates ISO
  const formatDate = (iso) => {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString('fr-FR', {
      day: 'numeric', month: 'long', year: 'numeric'
    })
  }

  return (
    <div className="carte-evenement">
      {evenement.cover_url && (
        <img
          className="carte-image"
          src={evenement.cover_url}
          alt={evenement.titre}
          onError={(e) => { e.target.style.display = 'none' }}
        />
      )}

      <div className="carte-contenu">
        <h2>{evenement.titre}</h2>

        {evenement.chapeau && (
          <p className="carte-chapeau">{evenement.chapeau}</p>
        )}

        <div className="carte-details">
          {evenement.nom_lieu && <p>📍 {evenement.nom_lieu}</p>}
          {evenement.adresse && <p>{evenement.adresse}, {evenement.code_postal} {evenement.ville}</p>}
          <p>📅 {formatDate(evenement.date_debut)} — {formatDate(evenement.date_fin)}</p>
          {evenement.prix && <p>💰 {evenement.prix}</p>}
        </div>

        {stats && (
          <div className="carte-stats">
            <span className="stat-like">❤️ {stats.likes}</span>
            <span className="stat-unlike">❌ {stats.unlikes}</span>
          </div>
        )}
      </div>
    </div>
  )
}
