package main

import (
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

var testProvider = NewInMemoryProvider(map[Hostname]DecentralizedID{
	"alice.example.com": "did:plc:example001",
	"bob.example.com":   "did:plc:example002",
}, map[Domain]bool{
	"example.com": true,
})

func TestHandleResultIsAddedToRequestContext(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	ctx.Set("handle", Handle{
		Domain:   "example.com",
		Username: "alice",
	})

	WithHandleResult(testProvider)(ctx)

	assert.Equal(t, ctx.MustGet("result"), Result{
		HasDecentralizedID: true,
		DecentralizedID:    "did:plc:example001",
	})
}

func TestRequestForUnsupportedDomainIsRejectedAsBadRequest(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.unsupported.domain"
	ctx.Request = req

	ctx.Set("handle", Handle{
		Domain:   "unsupported.domain",
		Username: "alice",
	})

	WithHandleResult(testProvider)(ctx)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestRedirectsUnmatchedRouteWithDecentralizedID(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	ctx.Set("handle", Handle{
		Domain:   "example.com",
		Username: "alice",
	})

	ctx.Set("result", Result{
		HasDecentralizedID: true,
		DecentralizedID:    DecentralizedID("did:plc:example"),
	})

	RedirectUnmatchedRoute(Config{
		RedirectDIDTemplate: "https://example.com/{did}",
	})(ctx)

	url, _ := res.Result().Location()

	assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	assert.Equal(t, "https://example.com/did:plc:example", url.String())
}

func TestRedirectsUnmatchedRouteWithoutDecentralizedID(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	ctx.Set("handle", Handle{
		Domain:   "example.com",
		Username: "alice",
	})

	ctx.Set("result", Result{
		HasDecentralizedID: false,
	})

	RedirectUnmatchedRoute(Config{
		RedirectHandleTemplate: "https://example.com/?from={handle}",
	})(ctx)

	url, _ := res.Result().Location()

	assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	assert.Equal(t, "https://example.com/?from=alice.example.com", url.String())
}

func TestServerIsHealthyWhenProviderIsHealthy(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	testProvider.SetHealthy(true)

	CheckServerIsHealthy(testProvider)(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestServerIsUnhealthyWhenProviderIsUnhealthy(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	testProvider.SetHealthy(false)

	CheckServerIsHealthy(testProvider)(ctx)

	assert.Equal(t, http.StatusBadGateway, res.Code)
}
