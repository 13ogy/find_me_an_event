import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './AuthContext'
import Connexion from './pages/Connexion'
import Inscription from './pages/Inscription'
import Decouverte from './pages/Decouverte'
import DetailEvenement from './pages/DetailEvenement'
import Historique from './pages/Historique'
import NotFound from './pages/NotFound'

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
          <Route path="/" element={<RouteProtegee><Decouverte /></RouteProtegee>} />
          <Route path="/evenement/:id" element={<RouteProtegee><DetailEvenement /></RouteProtegee>} />
          <Route path="/historique" element={<RouteProtegee><Historique /></RouteProtegee>} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
