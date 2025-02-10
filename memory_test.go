package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var memory = NewInMemoryProvider(map[Hostname]DecentralizedID{
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

	did, err := memory.GetDecentralizedIDForHandle(context.Background(), handle)

	assert.Nil(t, err)
	assert.Equal(t, did, DecentralizedID("did:plc:example001"))
}

func TestMemoryProviderCanProvideForDomain(t *testing.T) {
	expectedToProvide, err := memory.CanProvideForDomain(context.Background(), "example.com")

	assert.Nil(t, err)
	assert.True(t, expectedToProvide)

	expectedNotToProvide, err := memory.CanProvideForDomain(context.Background(), "example.net")

	assert.Nil(t, err)
	assert.False(t, expectedNotToProvide)
}

func TestMemoryIsHealthyByDefault(t *testing.T) {
	healthy, _ := memory.IsHealthy(context.Background())

	assert.True(t, healthy)
}

func TestMemoryCanBeSetUnhealthy(t *testing.T) {
	memory.SetHealthy(false)

	healthy, _ := memory.IsHealthy(context.Background())

	assert.False(t, healthy)
}
