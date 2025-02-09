package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleIsRenderedToString(t *testing.T) {
	tests := []struct {
		handle         Handle
		expectedString string
	}{
		{
			handle: Handle{
				Domain:   "example.com",
				Username: "alice",
			},
			expectedString: "alice.example.com",
		},
		{
			handle: Handle{
				Domain:   "example.com",
				Username: "ALICE",
			},
			expectedString: "alice.example.com",
		},
		{
			handle: Handle{
				Domain:   "AT.HANDLES.EXAMPLE.COM",
				Username: "alice",
			},
			expectedString: "alice.at.handles.example.com",
		},
	}

	for _, test := range tests {
		assert.Equal(
			t,
			test.expectedString,
			test.handle.String(),
		)
	}
}
