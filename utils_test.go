package main

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidHostnameIsParsedToHandle(t *testing.T) {
	tests := []struct {
		hostname       string
		expectedHandle Handle
	}{
		{
			hostname: "alice.example.com",
			expectedHandle: Handle{
				Hostname: "alice.example.com",
				Domain:   "example.com",
				Username: "alice",
			},
		},
		{
			hostname: "alice.at.people.example.com",
			expectedHandle: Handle{
				Hostname: "alice.at.people.example.com",
				Domain:   "at.people.example.com",
				Username: "alice",
			},
		},
		{
			hostname: "ALICE.at.people.example.com",
			expectedHandle: Handle{
				Hostname: "alice.at.people.example.com",
				Domain:   "at.people.example.com",
				Username: "alice",
			},
		},
	}

	for _, test := range tests {
		handle, err := HostnameToHandle(test.hostname)
		assert.Nil(t, err)
		assert.Equal(
			t,
			test.expectedHandle,
			handle,
			"Hostname %s was not parsed to expected Handle",
			test.hostname,
		)
	}
}

func TestInvalidHostnameReturnsError(t *testing.T) {
	handle, err := HostnameToHandle("not-a-valid-handle")
	assert.NotNil(t, err)
	assert.Equal(t, Handle{}, handle, "Handle returned for invalid hostname")
}

func TestTemplateUrlIsFormatted(t *testing.T) {
	tests := []struct {
		template    string
		expectedUrl string
	}{
		{
			template:    "https://example.com/?handle={handle}",
			expectedUrl: "https://example.com/?handle=alice.example.com",
		},
		{
			template:    "https://example.com/?handle={handle.hostname}",
			expectedUrl: "https://example.com/?handle=alice.example.com",
		},
		{
			template:    "https://{handle.domain}/?username={handle.username}",
			expectedUrl: "https://example.com/?username=alice",
		},
		{
			template:    "https://bsky.app/profile/{did}",
			expectedUrl: "https://bsky.app/profile/did:plc:example",
		},
		{
			template:    "https://example.com",
			expectedUrl: "https://example.com",
		},
		{
			template:    "https://example.com?{request.query}",
			expectedUrl: "https://example.com?a=b",
		},
	}

	request, _ := http.NewRequest("GET", "https://alice.example.com?a=b", bytes.NewReader([]byte{}))
	handle := Handle{Hostname: "alice.example.com", Domain: "example.com", Username: "alice"}
	did := DecentralizedID("did:plc:example")

	for _, test := range tests {
		url := FormatTemplateUrl(test.template, request, handle, did)
		assert.Equal(
			t,
			test.expectedUrl,
			url,
			"Template %s was not formatted correctly",
			url,
		)
	}
}
