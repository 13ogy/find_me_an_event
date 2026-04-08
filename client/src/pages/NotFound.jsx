import { Link } from 'react-router-dom'

// Page 404 — affichée quand aucune route ne correspond
export default function NotFound() {
  return (
    <div className="page-404">
      <h2>404</h2>
      <p>Cette page n'existe pas.</p>
      <Link to="/">← Retour à l'accueil</Link>
    </div>
  )
}
