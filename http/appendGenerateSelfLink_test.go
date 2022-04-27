package http

import (
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi"
	"github.com/stretchr/testify/assert"
)

func Test_appendGeneratedSelfLink(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?id=123", nil)
	baseURL, path := createBaseURL(req)

	links := appendGeneratedSelfLink(req)(nil, baseURL, path)
	expectedLink := jsonapi.Link{
		Href: "http://localhost:8080/example?id=123",
	}

	assert.Equal(t, len(links), 1)
	assert.NotNil(t, links["self"])
	assert.Equal(t, links["self"], expectedLink)
}
