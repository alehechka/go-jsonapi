package jsonapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi"
	"github.com/stretchr/testify/assert"
)

func Test_AppendGeneratedSelfLink(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?id=123", nil)
	baseURL, path := jsonapi.CreateBaseURL(req)

	links := jsonapi.AppendGeneratedSelfLink(req)(nil, baseURL, path)
	expectedLink := jsonapi.Link{
		Href: "http://localhost:8080/example?id=123",
	}

	assert.Equal(t, len(links), 1)
	assert.NotNil(t, links["self"])
	assert.Equal(t, links["self"], expectedLink)
}

func Test_CreateBaseURL_Proxied(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)
	req.Header = http.Header{
		jsonapi.ForwardedHost:   {"example.com"},
		jsonapi.ForwardedProto:  {"https"},
		jsonapi.ForwardedPrefix: {"/rest"},
	}

	baseURL, path := jsonapi.CreateBaseURL(req)

	assert.Equal(t, baseURL, "https://example.com/rest")
	assert.Equal(t, path, "/example")
}

func Test_createBaseURL_Raw(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/example", nil)

	baseURL, path := jsonapi.CreateBaseURL(req)

	assert.Equal(t, baseURL, "https://example.com")
	assert.Equal(t, path, "/example")
}

func Test_createBaseURL_Localhost(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	baseURL, path := jsonapi.CreateBaseURL(req)

	assert.Equal(t, baseURL, "http://localhost:8080")
	assert.Equal(t, path, "/example")
}
