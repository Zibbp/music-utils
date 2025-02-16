package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Debug               bool   `env:"DEBUG, default=false"`
	JSONLog             bool   `env:"JSON_LOG, default=false"`
	TidalClientId       string `env:"TIDAL_CLIENT_ID"`
	TidalClientSecret   string `env:"TIDAL_CLIENT_SECRET"`
	SpotifyClientId     string `env:"SPOTIFY_CLIENT_ID"`
	SpotifyClientSecret string `env:"SPOTIFY_CLIENT_SECRET"`
	SpotifyRedirectUri  string `env:"SPOTIFY_CLIENT_REDIRECT_URI, default=http://localhost:28542/callback"`
}

func Init() (*Config, error) {
	ctx := context.Background()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &c, nil
}
