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

	resolver := NewInMemoryResolver(map[Hostname]DecentralizedID{
		"alice.example.com": "did:plc:example001",
		"bob.example.com":   "did:plc:example002",
	}, map[Domain]bool{
		"example.com": true,
	})

	resolver.SetHealthy(false)

	router := gin.New()

	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())

	router.GET("/healthz", CheckServerIsHealthy(resolver))
	router.GET("/domainz", CheckDomainIsResolvedByServer(resolver))

	router.Use(ParseHandleFromHostname)
	router.Use(ResolveHandle(resolver))

	router.GET("/.well-known/atproto-did", VerifyHandle)

	router.NoRoute(RedirectUnmatchedRoute(config))

	router.Run(net.JoinHostPort(config.Host, config.Port))
}
