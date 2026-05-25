import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    // En dev, le client (vite) tourne sur :5173 et le serveur Go sur :8080.
    // Le proxy évite d'embarquer une URL absolue dans le code client.
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
