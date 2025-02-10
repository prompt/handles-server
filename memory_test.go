package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMemoryProvider = NewInMemoryProvider(map[Hostname]DecentralizedID{
	"alice.example.com": "did:plc:example001",
	"bob.example.com":   "did:plc:example002",
}, map[Domain]bool{
	"example.com": true,
})

func TestMemoryProviderHasDecentralizedIdForHandle(t *testing.T) {
	handle := Handle{
		Domain:   "example.com",
		Username: "alice",
	}

	did, err := testMemoryProvider.GetDecentralizedIDForHandle(context.Background(), handle)

	assert.Nil(t, err)
	assert.Equal(t, did, DecentralizedID("did:plc:example001"))
}

func TestMemoryProviderCanProvideForDomain(t *testing.T) {
	expectedToProvide, err := testMemoryProvider.CanProvideForDomain(context.Background(), "example.com")

	assert.Nil(t, err)
	assert.True(t, expectedToProvide)

	expectedNotToProvide, err := testMemoryProvider.CanProvideForDomain(context.Background(), "example.net")

	assert.Nil(t, err)
	assert.False(t, expectedNotToProvide)
}

func TestMemoryIsHealthyByDefault(t *testing.T) {
	healthy, _ := testMemoryProvider.IsHealthy(context.Background())

	assert.True(t, healthy)
}

func TestMemoryCanBeSetUnhealthy(t *testing.T) {
	testMemoryProvider.SetHealthy(false)

	healthy, _ := testMemoryProvider.IsHealthy(context.Background())

	assert.False(t, healthy)
}
