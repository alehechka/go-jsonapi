package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createBaseURL_Proxied(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)
	req.Header = http.Header{
		ForwardedHost:   {"example.com"},
		ForwardedProto:  {"https"},
		ForwardedPrefix: {"/rest"},
	}

	baseURL, path := createBaseURL(req)

	assert.Equal(t, baseURL, "https://example.com/rest")
	assert.Equal(t, path, "/example")
}

func Test_createBaseURL_Raw(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/example", nil)

	baseURL, path := createBaseURL(req)

	assert.Equal(t, baseURL, "https://example.com")
	assert.Equal(t, path, "/example")
}

func Test_createBaseURL_Localhost(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	baseURL, path := createBaseURL(req)

	assert.Equal(t, baseURL, "http://localhost:8080")
	assert.Equal(t, path, "/example")
}
