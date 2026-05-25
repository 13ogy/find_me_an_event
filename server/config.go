package main

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

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
