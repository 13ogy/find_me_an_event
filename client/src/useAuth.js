import { createContext, useContext } from 'react'

// Contexte d'authentification partagé entre AuthProvider et useAuth ; isolé pour
// rester compatible avec le HMR Vite (le composant et le hook vivent dans des
// fichiers séparés).
export const AuthContext = createContext(null)

export function useAuth() {
  return useContext(AuthContext)
}
