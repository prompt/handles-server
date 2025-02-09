package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	config, err := ConfigFromEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	provider := config.Provider

	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()

	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())

	router.GET("/healthz", CheckServerIsHealthy(provider))
	router.GET("/domainz", CheckServerProvidesForDomain(provider))

	router.Use(ParseHandleFromHostname)
	router.Use(WithHandleResult(provider))

	router.GET("/.well-known/atproto-did", VerifyHandle)

	router.NoRoute(RedirectUnmatchedRoute(config))

	router.Run(net.JoinHostPort(config.Host, config.Port))
}
