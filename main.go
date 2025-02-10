package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	config, err := ConfigFromEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()

	AddApplicationRoutes(router, config)

	if err := router.Run(net.JoinHostPort(config.Host, config.Port)); err != nil {
		log.Fatal(err)
	}
}

func AddApplicationRoutes(router *gin.Engine, config Config) {
	router.Use(sloggin.New(config.Logger))
	router.Use(gin.Recovery())

	router.GET("/healthz", CheckServerIsHealthy(config.Provider))
	router.GET("/domainz", CheckServerProvidesForDomain(config.Provider, config.CheckDomainParameter))

	router.Use(ParseHandleFromHostname)
	router.Use(WithHandleResult(config.Provider))

	router.GET("/.well-known/atproto-did", VerifyHandle)

	router.NoRoute(RedirectUnmatchedRoute(config.RedirectDIDTemplate, config.RedirectHandleTemplate))
}
