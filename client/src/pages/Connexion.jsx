import { useState } from 'react'
import { useAuth } from '../AuthContext'
import { useNavigate, Link } from 'react-router-dom'

export default function Connexion() {
  const { connecter } = useAuth()
  const navigate = useNavigate()
  const [nom, setNom] = useState('')
  const [motDePasse, setMotDePasse] = useState('')
  const [erreur, setErreur] = useState('')

  async function soumettre(e) {
    e.preventDefault()
    setErreur('')
    try {
      await connecter(nom, motDePasse)
      navigate('/')
    } catch (err) {
      setErreur(err.message)
    }
  }

  return (
    <div className="page-auth">
      <h1>Connexion</h1>
      <form onSubmit={soumettre}>
        <label>
          Nom d'utilisateur
          <input
            type="text"
            value={nom}
            onChange={e => setNom(e.target.value)}
            required
          />
        </label>
        <label>
          Mot de passe
          <input
            type="password"
            value={motDePasse}
            onChange={e => setMotDePasse(e.target.value)}
            required
          />
        </label>
        {erreur && <p className="erreur">{erreur}</p>}
        <button type="submit">Se connecter</button>
      </form>
      <p>
        Pas encore de compte ? <Link to="/inscription">Créer un compte</Link>
      </p>
    </div>
  )
}
