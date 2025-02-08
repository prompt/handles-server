package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckServerIsHealthy(resolver ResolvesHandlesToDecentralizedIDs) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthy, explanation := resolver.IsResolverHealthy()
		if !healthy {
			c.AbortWithError(http.StatusBadGateway, errors.New(explanation))
		}
		c.String(http.StatusOK, explanation)
	}
}

type Resolution struct {
	HasDecentralizedID bool
	DecentralizedID    DecentralizedID
	Err                error
}

func ParseHandleFromHostname(c *gin.Context) {
	handle, err := HostnameToHandle(c.Request.Host)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Set("handle", handle)

	c.Next()
}

func CheckDomainIsResolvedByServer(resolver ResolvesHandlesToDecentralizedIDs) gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := Domain(strings.ToLower(c.Query("domain")))

		canResolve := resolver.CanResolveHandlesAtDomain(domain)

		if !canResolve {
			c.String(http.StatusNotFound, "Handles for domain %s are not resolved by this server.", domain)
			return
		}

		c.String(http.StatusOK, "Handles for domain %s are resolved by this server.", domain)
	}
}

func ResolveHandle(resolver ResolvesHandlesToDecentralizedIDs) gin.HandlerFunc {
	return func(c *gin.Context) {
		handle := c.MustGet("handle").(Handle)

		did, err := resolver.ResolveHandleToDID(handle)

		c.Set("resolution", Resolution{
			HasDecentralizedID: err == nil,
			DecentralizedID:    did,
			Err:                err,
		})
	}
}

func VerifyHandle(c *gin.Context) {
	resolution := c.MustGet("resolution").(Resolution)

	if resolution.Err != nil {
		c.String(http.StatusNotFound, resolution.Err.Error())
		return
	}

	c.String(http.StatusOK, string(resolution.DecentralizedID))
}

func RedirectUnmatchedRoute(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		resolution := c.MustGet("resolution").(Resolution)

		if resolution.HasDecentralizedID {
			c.Redirect(http.StatusTemporaryRedirect, FormatTemplateUrl(
				config.HandleRedirect,
				c.Request,
				c.MustGet("handle").(Handle),
				resolution.DecentralizedID,
			))
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, FormatTemplateUrl(
			config.UnresolvableRedirect,
			c.Request,
			c.MustGet("handle").(Handle),
			DecentralizedID(""),
		))
	}
}
