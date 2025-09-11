package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealth_Ok checks that the HTTP client handles a successful call to the Health endpoint as expected
func TestHealth_Ok(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/health", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	err := client.Health(context.Background())

	// Assert
	assert.NoError(t, err)
}

// TestHealth_HTTPError checks that the HTTP client handles an unsuccessful call to the Health endpoint as expected
func TestHealth_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	err := client.Health(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected http status code")
}
