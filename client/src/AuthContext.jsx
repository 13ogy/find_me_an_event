import { createContext, useContext, useState, useEffect } from 'react'
import { profil, connexion as apiConnexion, deconnexion as apiDeconnexion, inscription as apiInscription } from './api'

const AuthContext = createContext(null)

// Fournit l'état d'authentification à toute l'application
export function AuthProvider({ children }) {
  const [utilisateur, setUtilisateur] = useState(null)
  const [chargement, setChargement] = useState(true)

  // Vérifie la session au chargement (cookie existant ?)
  useEffect(() => {
    profil()
      .then(setUtilisateur)
      .catch(() => setUtilisateur(null))
      .finally(() => setChargement(false))
  }, [])

  async function connecter(nom, mot_de_passe) {
    const u = await apiConnexion(nom, mot_de_passe)
    setUtilisateur(u)
    return u
  }

  async function inscrire(nom, email, mot_de_passe) {
    const u = await apiInscription(nom, email, mot_de_passe)
    // Connexion automatique après inscription
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

// Hook pour accéder au contexte d'auth depuis n'importe quel composant
export function useAuth() {
  return useContext(AuthContext)
}
