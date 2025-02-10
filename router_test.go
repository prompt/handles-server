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

var testProviderForRouter = NewInMemoryProvider(map[Hostname]DecentralizedID{
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

	WithHandleResult(testProviderForRouter)(ctx)

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

	WithHandleResult(testProviderForRouter)(ctx)

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

	RedirectUnmatchedRoute(
		URLTemplate("https://example.com/{did}"),
		URLTemplate("https://example.com/{handle}"),
	)(ctx)

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

	RedirectUnmatchedRoute(
		URLTemplate("https://example.com/{did}"),
		URLTemplate("https://example.com/?from={handle}"),
	)(ctx)

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

	testProviderForRouter.SetHealthy(true)

	CheckServerIsHealthy(testProviderForRouter)(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestServerIsUnhealthyWhenProviderIsUnhealthy(t *testing.T) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alice.example.com"
	ctx.Request = req

	testProviderForRouter.SetHealthy(false)

	CheckServerIsHealthy(testProviderForRouter)(ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
