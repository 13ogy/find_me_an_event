package main

import "os"

// Config regroupe les paramètres de l'application
type Config struct {
	Port        string
	DatabaseURL string
}

// chargerConfig lit les variables d'environnement avec des valeurs par défaut
func chargerConfig() Config {
	return Config{
		Port:        envOuDefaut("PORT", "8080"),
		DatabaseURL: envOuDefaut("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/find_me_an_event?sslmode=disable"),
	}
}

func envOuDefaut(cle, defaut string) string {
	if v := os.Getenv(cle); v != "" {
		return v
	}
	return defaut
}
