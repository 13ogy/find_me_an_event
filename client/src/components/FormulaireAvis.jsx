import { useState } from 'react'
import { posterAvis } from '../api'

export default function FormulaireAvis({ evenementId, onAvisCree }) {
  const [commentaire, setCommentaire] = useState('')
  const [note, setNote] = useState(3)
  const [envoi, setEnvoi] = useState(false)
  const [erreur, setErreur] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    if (!commentaire.trim()) return

    setEnvoi(true)
    setErreur('')
    try {
      const avis = await posterAvis(evenementId, commentaire.trim(), note)
      setCommentaire('')
      setNote(3)
      // Le parent insère l'avis en tête de liste sans avoir à refaire un
      // GET /reviews : affichage instantané.
      if (onAvisCree) onAvisCree(avis)
    } catch (err) {
      setErreur(err.message)
    } finally {
      setEnvoi(false)
    }
  }

  return (
    <form className="formulaire-avis" onSubmit={handleSubmit}>
      <h3>Laisser un avis</h3>

      <label>
        Note
        <select value={note} onChange={(e) => setNote(Number(e.target.value))}>
          {[1, 2, 3, 4, 5].map(n => (
            <option key={n} value={n}>{'⭐'.repeat(n)}</option>
          ))}
        </select>
      </label>

      <label>
        Commentaire
        <textarea
          value={commentaire}
          onChange={(e) => setCommentaire(e.target.value)}
          placeholder="Votre avis sur cet événement..."
          rows={3}
          required
        />
      </label>

      {erreur && <p className="erreur">{erreur}</p>}

      <button type="submit" disabled={envoi}>
        {envoi ? 'Envoi...' : 'Publier'}
      </button>
    </form>
  )
}
