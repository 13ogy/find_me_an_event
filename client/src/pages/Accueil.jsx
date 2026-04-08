import { useAuth } from '../AuthContext'

// Page d'accueil — placeholder pour la phase 5
export default function Accueil() {
  const { utilisateur, deconnecter } = useAuth()

  return (
    <div className="page-accueil">
      <header>
        <h1>🎉 Find Me An Event</h1>
        <div className="header-actions">
          <span>Bonjour, {utilisateur?.nom}</span>
          <button onClick={deconnecter}>Déconnexion</button>
        </div>
      </header>
      <main>
        <p>Bienvenue ! La découverte d'événements arrive bientôt.</p>
      </main>
    </div>
  )
}
