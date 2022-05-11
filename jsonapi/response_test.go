package jsonapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi/jsonapi"
	"github.com/stretchr/testify/assert"
)

func Test_CreateResponse(t *testing.T) {
	url := "http://localhost:8080/example?id=123"
	req := httptest.NewRequest("GET", url, nil)

	response := jsonapi.CreateResponse(req)(jsonapi.Response{
		Links: jsonapi.Links{
			jsonapi.SelfKey: jsonapi.Link{
				Href: "https://example.com/resource",
			},
		},
	})

	assert.NotNil(t, response)
	assert.NotNil(t, response.Links)
	assert.NotNil(t, response.Links[jsonapi.SelfKey])
	assert.Equal(t, url, response.Links[jsonapi.SelfKey])
}

func Test_CreateCollectionResponse(t *testing.T) {
	url := "http://localhost:8080/example?id=123"
	req := httptest.NewRequest("GET", url, nil)

	response := jsonapi.CreateCollectionResponse(req)(jsonapi.CollectionResponse{
		Links: jsonapi.Links{
			jsonapi.SelfKey: jsonapi.Link{
				Href: "https://example.com/resource",
			},
		},
	})

	assert.NotNil(t, response)
	assert.NotNil(t, response.Links)
	assert.NotNil(t, response.Links[jsonapi.SelfKey])
	assert.Equal(t, url, response.Links[jsonapi.SelfKey])
}

func Test_AppendGeneratedSelfLink(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?id=123", nil)
	baseURL, path := jsonapi.CreateBaseURL(req)

	links := jsonapi.AppendGeneratedSelfLink(req)(nil, baseURL, path)
	expectedLink := jsonapi.Link{
		Href: "http://localhost:8080/example?id=123",
	}

	assert.Equal(t, 1, len(links))
	assert.NotNil(t, links["self"])
	assert.Equal(t, expectedLink, links["self"])
}

func Test_CreateBaseURL_Proxied(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)
	req.Header = http.Header{
		jsonapi.ForwardedHost:   {"example.com"},
		jsonapi.ForwardedProto:  {"https"},
		jsonapi.ForwardedPrefix: {"/rest"},
	}

	baseURL, path := jsonapi.CreateBaseURL(req)

	assert.Equal(t, "https://example.com/rest", baseURL)
	assert.Equal(t, "/example", path)
}

func Test_createBaseURL_Raw(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/example", nil)

	baseURL, path := jsonapi.CreateBaseURL(req)

	assert.Equal(t, "https://example.com", baseURL)
	assert.Equal(t, "/example", path)
}

func Test_createBaseURL_Localhost(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	baseURL, path := jsonapi.CreateBaseURL(req)

	assert.Equal(t, "http://localhost:8080", baseURL)
	assert.Equal(t, "/example", path)
}

func Test_createBaseURL_NoScheme(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)
	req.URL.Scheme = ""

	baseURL, path := jsonapi.CreateBaseURL(req)

	assert.Equal(t, "http://localhost:8080", baseURL)
	assert.Equal(t, "/example", path)
}
