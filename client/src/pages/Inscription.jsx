import { useState } from 'react'
import { useAuth } from '../AuthContext'
import { useNavigate, Link } from 'react-router-dom'

export default function Inscription() {
  const { inscrire } = useAuth()
  const navigate = useNavigate()
  const [nom, setNom] = useState('')
  const [email, setEmail] = useState('')
  const [motDePasse, setMotDePasse] = useState('')
  const [erreur, setErreur] = useState('')

  async function soumettre(e) {
    e.preventDefault()
    setErreur('')
    try {
      await inscrire(nom, email, motDePasse)
      navigate('/')
    } catch (err) {
      setErreur(err.message)
    }
  }

  return (
    <div className="page-auth">
      <h1>Inscription</h1>
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
          Email
          <input
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
        </label>
        <label>
          Mot de passe (6 caractères min.)
          <input
            type="password"
            value={motDePasse}
            onChange={e => setMotDePasse(e.target.value)}
            required
            minLength={6}
          />
        </label>
        {erreur && <p className="erreur">{erreur}</p>}
        <button type="submit">Créer mon compte</button>
      </form>
      <p>
        Déjà un compte ? <Link to="/connexion">Se connecter</Link>
      </p>
    </div>
  )
}
