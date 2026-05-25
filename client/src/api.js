// Petit wrapper fetch : injecte le content-type JSON et `credentials: 'include'`
// pour que le cookie de session traverse le proxy Vite vers le serveur Go.
// Lève une Error contenant le message renvoyé par le serveur si le statut HTTP
// n'est pas 2xx, ce qui permet aux composants de faire un simple try/catch.

const BASE = '/api'

async function appel(methode, chemin, body = null) {
  const options = {
    method: methode,
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  }
  if (body) {
    options.body = JSON.stringify(body)
  }

  const reponse = await fetch(BASE + chemin, options)
  const data = await reponse.json()

  if (!reponse.ok) {
    throw new Error(data.erreur || 'Erreur serveur')
  }
  return data
}

export const inscription = (nom, email, mot_de_passe) =>
  appel('POST', '/register', { nom, email, mot_de_passe })

export const connexion = (nom, mot_de_passe) =>
  appel('POST', '/login', { nom, mot_de_passe })

export const deconnexion = () =>
  appel('POST', '/logout')

export const profil = () =>
  appel('GET', '/me')

export const localisation = () =>
  appel('GET', '/location')

export const evenements = (lat, lon, radius, categorie = '') =>
  appel(
    'GET',
    `/events?lat=${lat}&lon=${lon}&radius=${radius}&categorie=${encodeURIComponent(categorie)}`
  )

export const prochainEvenement = (lat, lon, radius, categorie = '') =>
  appel(
    'GET',
    `/events/next?lat=${lat}&lon=${lon}&radius=${radius}&categorie=${encodeURIComponent(categorie)}`
  )

export const voter = (id, vote, titre) =>
  appel('POST', `/events/${id}/vote`, { vote, titre })

export const supprimerVote = (id) =>
  appel('DELETE', `/events/${id}/vote`)

export const statsEvenement = (id) =>
  appel('GET', `/events/${id}/stats`)

export const listerAvis = (id) =>
  appel('GET', `/events/${id}/reviews`)

export const posterAvis = (id, commentaire, note) =>
  appel('POST', `/events/${id}/reviews`, { commentaire, note })

export const historique = () =>
  appel('GET', '/me/history')
