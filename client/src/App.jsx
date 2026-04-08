import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './AuthContext'
import Connexion from './pages/Connexion'
import Inscription from './pages/Inscription'
import Accueil from './pages/Accueil'

// Redirige vers /connexion si l'utilisateur n'est pas connecté
function RouteProtegee({ children }) {
  const { utilisateur, chargement } = useAuth()

  if (chargement) return <p>Chargement...</p>
  if (!utilisateur) return <Navigate to="/connexion" />
  return children
}

// Redirige vers / si l'utilisateur est déjà connecté
function RoutePublique({ children }) {
  const { utilisateur, chargement } = useAuth()

  if (chargement) return <p>Chargement...</p>
  if (utilisateur) return <Navigate to="/" />
  return children
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/connexion" element={<RoutePublique><Connexion /></RoutePublique>} />
          <Route path="/inscription" element={<RoutePublique><Inscription /></RoutePublique>} />
          <Route path="/" element={<RouteProtegee><Accueil /></RouteProtegee>} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
