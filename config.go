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

	RedirectDIDTemplate    URLTemplate `env:"REDIRECT_DID_TEMPLATE" envDefault:"https://bsky.app/profile/{did}"`
	RedirectHandleTemplate URLTemplate `env:"REDIRECT_HANDLE_TEMPLATE" envDefault:"https://{handle.domain}?handle={handle}"`

	Postgres             *pgxpool.Config `env:"DATABASE_URL"`
	PostgresDidsTable    string          `env:"DATABASE_TABLE_DIDS" envDefault:"dids"`
	PostgresDomainsTable string          `env:"DATABASE_TABLE_DOMAINS" envDefault:"domains"`

	MemoryDids    map[string]string `env:"MEMORY_DIDS" envKeyValSeparator:"@"`
	MemoryDomains []string          `env:"MEMORY_DOMAINS"`

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
					if config.Postgres == nil {
						return &PostgresHandles{}, errors.New("a database connection (`DATABASE_URL`) is required to use the postgres provider")
					}

					return NewPostgresHandlesProvider(
						config.Postgres,
						config.PostgresDidsTable,
						config.PostgresDomainsTable,
					)
				case "memory":
					if config.MemoryDids == nil || config.MemoryDomains == nil {
						return nil, errors.New("a map of Decentralized IDs (`MEMORY_DIDS`) and domains (`MEMORY_DOMAINS`) is required to use the memory provider")
					}

					dids := make(MapOfDids)

					for handle, did := range config.MemoryDids {
						dids[Hostname(handle)] = DecentralizedID(did)
					}

					domains := make(MapOfDomains)

					for _, domain := range config.MemoryDomains {
						domains[Domain(domain)] = true
					}

					provider := NewInMemoryProvider(dids, domains)
					return provider, nil
				default:
					return nil, errors.New("no valid provider of decentralized IDs specified")
				}
			},
		},
	})

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
