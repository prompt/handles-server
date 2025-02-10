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
		Level: config.LogLevel,
	}))

	router := gin.New()

	AddApplicationRoutes(router, logger, config)

	if err := router.Run(net.JoinHostPort(config.Host, config.Port)); err != nil {
		log.Fatal(err)
	}
}

func AddApplicationRoutes(router *gin.Engine, logger *slog.Logger, config Config) {
	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())

	router.GET("/healthz", CheckServerIsHealthy(config.Provider))
	router.GET("/domainz", CheckServerProvidesForDomain(config.Provider))

	router.Use(ParseHandleFromHostname)
	router.Use(WithHandleResult(config.Provider))

	router.GET("/.well-known/atproto-did", VerifyHandle)

	router.NoRoute(RedirectUnmatchedRoute(config.RedirectDIDTemplate, config.RedirectHandleTemplate))
}
