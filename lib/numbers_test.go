package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumbersClient_GetMathFact(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/random/math", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "json", r.URL.RawQuery)

		// Return a sample response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"text":"42 is the meaning of life", "number":42, "found":true, "type":"math"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewNumbersClient()
	client.baseURL = server.URL

	// Call the method
	fact, err := client.GetMathFact()

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, "42 is the meaning of life", fact.Text)
	assert.Equal(t, 42, fact.Number)
	assert.True(t, fact.Found)
	assert.Equal(t, "math", fact.Type)
}
