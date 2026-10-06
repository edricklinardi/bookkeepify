// Package config loads process configuration from the environment.
package config

import (
	"errors"
	"os"
)

type Config struct {
	Addr        string // API listen address
	DatabaseURL string

	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURL  string

	LastFMAPIKey string
}

func Load() (Config, error) {
	cfg := Config{
		Addr:                getenv("ADDR", ":8080"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURL:  getenv("SPOTIFY_REDIRECT_URL", "http://127.0.0.1:8080/auth/callback"),
		LastFMAPIKey:        os.Getenv("LASTFM_API_KEY"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
