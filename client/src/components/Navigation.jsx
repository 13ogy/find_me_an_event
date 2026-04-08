import { useAuth } from '../AuthContext'
import { NavLink, useNavigate } from 'react-router-dom'

// Barre de navigation commune à toutes les pages protégées
export default function Navigation() {
  const { utilisateur, deconnecter } = useAuth()
  const navigate = useNavigate()

  async function handleDeconnexion() {
    await deconnecter()
    navigate('/connexion')
  }

  return (
    <header className="nav-bar">
      <h1 className="nav-logo">🎉 Find Me An Event</h1>
      <nav className="nav-liens">
        <NavLink to="/" end>Découvrir</NavLink>
        <NavLink to="/historique">Historique</NavLink>
      </nav>
      <div className="nav-user">
        <span>{utilisateur?.nom}</span>
        <button onClick={handleDeconnexion}>Déconnexion</button>
      </div>
    </header>
  )
}
