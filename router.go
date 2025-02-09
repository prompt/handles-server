package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckServerIsHealthy(provider ProvidesDecentralizedIDs) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthy, explanation := provider.IsHealthy(c)
		if !healthy {
			c.AbortWithError(http.StatusBadGateway, errors.New(explanation))
		}
		c.String(http.StatusOK, explanation)
	}
}

type Result struct {
	HasDecentralizedID bool
	DecentralizedID    DecentralizedID
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

func CheckServerProvidesForDomain(provider ProvidesDecentralizedIDs) gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := Domain(strings.ToLower(c.Query("domain")))

		canProvide, err := provider.CanProvideForDomain(c, domain)

		if err != nil {
			c.AbortWithError(http.StatusBadGateway, err)
			return
		}

		if !canProvide {
			c.String(http.StatusNotFound, "Decentralized IDs are not provided for %s by this server.", domain)
			return
		}

		c.String(http.StatusOK, "Decentralized IDs are not provided for %s by this server.", domain)
	}
}

func WithHandleResult(provider ProvidesDecentralizedIDs) gin.HandlerFunc {
	return func(c *gin.Context) {
		handle := c.MustGet("handle").(Handle)

		// TODO: Check can provide for domain

		did, err := provider.GetDecentralizedIDForHandle(c, handle)

		if err != nil {
			c.AbortWithError(http.StatusBadGateway, err)
			return
		}

		c.Set("result", Result{
			HasDecentralizedID: did != "",
			DecentralizedID:    did,
		})
	}
}

func VerifyHandle(c *gin.Context) {
	result := c.MustGet("result").(Result)

	if !result.HasDecentralizedID {
		c.String(
			http.StatusNotFound,
			fmt.Sprintf("Decentralized ID not found for %s", c.MustGet("handle").(Handle).String()),
		)
		return
	}

	c.String(http.StatusOK, string(result.DecentralizedID))
}

func RedirectUnmatchedRoute(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := c.MustGet("result").(Result)

		if result.HasDecentralizedID {
			c.Redirect(http.StatusTemporaryRedirect, URLFromTemplate(
				config.RedirectDIDTemplate,
				c.Request,
				c.MustGet("handle").(Handle),
				result.DecentralizedID,
			))
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, URLFromTemplate(
			config.RedirectHandleTemplate,
			c.Request,
			c.MustGet("handle").(Handle),
			DecentralizedID(""),
		))
	}
}
