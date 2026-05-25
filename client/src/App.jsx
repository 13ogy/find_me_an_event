import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './AuthContext'
import { useAuth } from './useAuth'
import Connexion from './pages/Connexion'
import Inscription from './pages/Inscription'
import Decouverte from './pages/Decouverte'
import DetailEvenement from './pages/DetailEvenement'
import Historique from './pages/Historique'
import NotFound from './pages/NotFound'

// Les deux gardes attendent que l'AuthProvider ait fini sa reprise de session
// avant de décider d'une redirection : sinon, un utilisateur connecté serait
// brièvement renvoyé vers /connexion à chaque refresh.
function RouteProtegee({ children }) {
  const { utilisateur, chargement } = useAuth()

  if (chargement) return <p>Chargement...</p>
  if (!utilisateur) return <Navigate to="/connexion" />
  return children
}

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
