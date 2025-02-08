package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidHandleIsAddedToRequestContext(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	ParseHandleFromHostname(ctx)

	assert.Equal(t, ctx.MustGet("handle"), Handle{
		Hostname: "alice.example.com",
		Domain:   "example.com",
		Username: "alice",
	})
}

func TestInvalidHostnameCausesBadRequest(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice"
	ctx.Request = req

	ParseHandleFromHostname(ctx)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

var testResolver = NewInMemoryResolver(map[Hostname]DecentralizedID{
	"alice.example.com": "did:plc:example001",
	"bob.example.com":   "did:plc:example002",
}, map[Domain]bool{
	"example.com": true,
})

func TestResolvedHandleIsAddedToRequestContext(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	ctx.Set("handle", Handle{
		Hostname: "alice.example.com",
		Domain:   "example.com",
		Username: "alice",
	})

	ResolveHandle(testResolver)(ctx)

	assert.Equal(t, ctx.MustGet("resolution"), Resolution{
		HasDecentralizedID: true,
		DecentralizedID:    "did:plc:example001",
		Err:                nil,
	})
}

func TestRedirectsUnmatchedRouteWithResolvedHandle(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	ctx.Set("handle", Handle{
		Hostname: "alice.example.com",
		Domain:   "example.com",
		Username: "alice",
	})

	ctx.Set("resolution", Resolution{
		HasDecentralizedID: true,
		DecentralizedID:    DecentralizedID("did:plc:example"),
		Err:                nil,
	})

	RedirectUnmatchedRoute(Config{
		HandleRedirect: "https://example.com/{did}",
	})(ctx)

	url, _ := res.Result().Location()

	assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	assert.Equal(t, "https://example.com/did:plc:example", url.String())
}

func TestRedirectsUnmatchedRouteWithUnresolvedHandle(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	ctx.Set("handle", Handle{
		Hostname: "alice.example.com",
		Domain:   "example.com",
		Username: "alice",
	})

	ctx.Set("resolution", Resolution{
		HasDecentralizedID: false,
		Err:                errors.New("Handle not resolved to a DID."),
	})

	RedirectUnmatchedRoute(Config{
		UnresolvableRedirect: "https://example.com/?from={handle}",
	})(ctx)

	url, _ := res.Result().Location()

	assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	assert.Equal(t, "https://example.com/?from=alice.example.com", url.String())
}

func TestServerIsHealthyWhenResolverIsHealthy(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	testResolver.SetHealthy(true)

	CheckServerIsHealthy(testResolver)(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestServerIsUnhealthyWhenResolverIsUnhealthy(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	testResolver.SetHealthy(false)

	CheckServerIsHealthy(testResolver)(ctx)

	assert.Equal(t, http.StatusBadGateway, res.Code)
}
