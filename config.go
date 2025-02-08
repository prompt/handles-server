package main

import (
	"log/slog"
	"reflect"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host     string     `env:"HOST" envDefault:"localhost"`
	Port     string     `env:"PORT" envDefault:"80"`
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"error"`

	HandleRedirect       string `env:"HANDLE_REDIRECT" envDefault:"https://bsky.app/profile/{did}"`
	UnresolvableRedirect string `env:"UNRESOLVABLE_REDIRECT" envDefault:"https://{handle.domain}?handle={handle}"`
}

func ConfigFromEnvironment() (Config, error) {
	config := Config{}

	err := env.ParseWithOptions(&config, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[slog.Level](): func(v string) (interface{}, error) {
				var level slog.Level
				return level, level.UnmarshalText([]byte(v))
			},
		},
	})

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
