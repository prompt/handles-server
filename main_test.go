package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func NewTestEnvironment() (*gin.Engine, *InMemoryProvider) {
	var testProvider = NewInMemoryProvider(map[Hostname]DecentralizedID{
		"alice.example.com": "did:plc:example001",
		"bob.example.com":   "did:plc:example002",
	}, map[Domain]bool{
		"example.com": true,
	})

	var testConfig = Config{
		Provider:               testProvider,
		RedirectDIDTemplate:    "https://example.com/profile/{did}",
		RedirectHandleTemplate: "https://example.com/register?handle={handle}",
	}

	var testRouter = gin.New()
	var testLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	AddApplicationRoutes(testRouter, testLogger, testConfig)

	return testRouter, testProvider
}

func TestHealthEndpointReturnsHealthyStatus(t *testing.T) {
	router, provider := NewTestEnvironment()

	provider.SetHealthy(true)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHealthEndpointReturnsUnhealthyStatus(t *testing.T) {
	router, provider := NewTestEnvironment()

	provider.SetHealthy(false)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestDomainEndpointReturnsOKForProvidedDomain(t *testing.T) {
	router, _ := NewTestEnvironment()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/domainz?domain=example.com", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestDomainEndpointReturnsErrorForUnprovidedDomain(t *testing.T) {
	router, _ := NewTestEnvironment()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/domainz?domain=unprovided.test", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestDidEndpointReturnsDidForKnownHandle(t *testing.T) {
	router, _ := NewTestEnvironment()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://alice.example.com/.well-known/atproto-did", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "did:plc:example001", res.Body.String())
}

func TestDidEndpointReturnsNotFoundForUnknownHandle(t *testing.T) {
	router, _ := NewTestEnvironment()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://unknown.example.com/.well-known/atproto-did", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestDidEndpointReturnsRequestErrorForUnknownDomain(t *testing.T) {
	router, _ := NewTestEnvironment()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://username.unprovided.test/.well-known/atproto-did", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestKnownHandleRedirectsToDestination(t *testing.T) {
	router, _ := NewTestEnvironment()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://alice.example.com", nil)
	router.ServeHTTP(res, req)

	url, _ := res.Result().Location()

	assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	assert.Equal(t, "https://example.com/profile/did:plc:example001", url.String())

}

func TestUnknownHandleRedirectsToDestination(t *testing.T) {
	router, _ := NewTestEnvironment()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://carol.example.com", nil)
	router.ServeHTTP(res, req)

	url, _ := res.Result().Location()

	assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	assert.Equal(t, "https://example.com/register?handle=carol.example.com", url.String())

}
