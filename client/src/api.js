// Fonctions d'appel au serveur — simples et directes
// Toutes les requêtes incluent les cookies (credentials) pour l'authentification

const BASE = '/api'

async function appel(methode, chemin, body = null) {
  const options = {
    method: methode,
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include', // Envoie le cookie de session
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

// Auth
export const inscription = (nom, email, mot_de_passe) =>
  appel('POST', '/register', { nom, email, mot_de_passe })

export const connexion = (nom, mot_de_passe) =>
  appel('POST', '/login', { nom, mot_de_passe })

export const deconnexion = () =>
  appel('POST', '/logout')

export const profil = () =>
  appel('GET', '/me')

// Géolocalisation
export const localisation = () =>
  appel('GET', '/location')

// Événements
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
// Votes
export const voter = (id, vote, titre) =>
  appel('POST', `/events/${id}/vote`, { vote, titre })

export const supprimerVote = (id) =>
  appel('DELETE', `/events/${id}/vote`)

export const statsEvenement = (id) =>
  appel('GET', `/events/${id}/stats`)

// Avis
export const listerAvis = (id) =>
  appel('GET', `/events/${id}/reviews`)

export const posterAvis = (id, commentaire, note) =>
  appel('POST', `/events/${id}/reviews`, { commentaire, note })

// Historique
export const historique = () =>
  appel('GET', '/me/history')
