import { useState, useEffect } from 'react'
import { AuthContext } from './useAuth'
import {
  profil,
  connexion as apiConnexion,
  deconnexion as apiDeconnexion,
  inscription as apiInscription,
} from './api'

// Fournit l'état d'auth à l'arbre React et tente une reprise de session au
// montage (le cookie HttpOnly n'est pas lisible depuis JS, seul un appel
// serveur peut nous dire si on est connecté).
export function AuthProvider({ children }) {
  const [utilisateur, setUtilisateur] = useState(null)
  const [chargement, setChargement] = useState(true)

  useEffect(() => {
    let actif = true
    profil()
      .then((u) => { if (actif) setUtilisateur(u) })
      .catch(() => { if (actif) setUtilisateur(null) })
      .finally(() => { if (actif) setChargement(false) })
    return () => { actif = false }
  }, [])

  async function connecter(nom, mot_de_passe) {
    const u = await apiConnexion(nom, mot_de_passe)
    setUtilisateur(u)
    return u
  }

  // Le serveur ne renvoie pas de Set-Cookie sur /register : on enchaîne donc
  // un /login pour démarrer immédiatement la session.
  async function inscrire(nom, email, mot_de_passe) {
    const u = await apiInscription(nom, email, mot_de_passe)
    await connecter(nom, mot_de_passe)
    return u
  }

  async function deconnecter() {
    await apiDeconnexion()
    setUtilisateur(null)
  }

  return (
    <AuthContext.Provider value={{ utilisateur, chargement, connecter, inscrire, deconnecter }}>
      {children}
    </AuthContext.Provider>
  )
}
