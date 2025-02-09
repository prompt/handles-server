package main

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string     `env:"HOST" envDefault:"localhost"`
	Port     string     `env:"PORT" envDefault:"80"`
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"error"`

	RedirectDID    string `env:"REDIRECT_DID" envDefault:"https://bsky.app/profile/{did}"`
	RedirectHandle string `env:"REDIRECT_HANDLE" envDefault:"https://{handle.domain}?handle={handle}"`

	Database *pgxpool.Config `env:"DATABASE_URL"`

	Provider ProvidesDecentralizedIDs `env:"DID_PROVIDER,required"`
}

func ConfigFromEnvironment() (Config, error) {
	config := Config{}

	err := env.ParseWithOptions(&config, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[slog.Level](): func(v string) (interface{}, error) {
				var level slog.Level
				return level, level.UnmarshalText([]byte(v))
			},
			reflect.TypeFor[pgxpool.Config](): func(v string) (interface{}, error) {
				config, err := pgxpool.ParseConfig(v)
				return *config, err
			},
			reflect.TypeFor[ProvidesDecentralizedIDs](): func(v string) (interface{}, error) {
				switch v {
				case "postgres":
					if config.Database == nil {
						return &PostgresHandles{}, errors.New("a database connection (`DATABASE_URL`) is required to use the postgres provider")
					}

					return NewPostgresHandlesProvider(config.Database)
				case "memory":
					provider := NewInMemoryProvider(map[Hostname]DecentralizedID{
						"alice.example.com": "did:plc:example001",
						"bob.example.com":   "did:plc:example002",
					}, map[Domain]bool{
						"example.com": true,
					})
					return provider, nil
				default:
					return nil, errors.New("no provider of decentralized IDs specified")
				}
			},
		},
	})

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
