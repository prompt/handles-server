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
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New(explanation))
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
		_ = c.AbortWithError(http.StatusBadRequest, err)
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
			_ = c.AbortWithError(http.StatusInternalServerError, err)
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

		did, err := provider.GetDecentralizedIDForHandle(c, handle)

		if errors.Is(err, (*CannotGetHandelsFromDomainError)(nil)) {
			c.String(
				http.StatusBadRequest,
				fmt.Sprintf("Decentralized IDs for the domain %s are not provided by this server.", handle.Domain),
			)
			c.Abort()
			return
		}

		if err != nil {
			_ = c.AbortWithError(http.StatusBadGateway, err)
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

func RedirectUnmatchedRoute(redirectDid URLTemplate, redirectHandle URLTemplate) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := c.MustGet("result").(Result)

		if result.HasDecentralizedID {
			c.Redirect(http.StatusTemporaryRedirect, URLFromTemplate(
				redirectDid,
				c.Request,
				c.MustGet("handle").(Handle),
				result.DecentralizedID,
			))
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, URLFromTemplate(
			redirectHandle,
			c.Request,
			c.MustGet("handle").(Handle),
			DecentralizedID(""),
		))
	}
}
