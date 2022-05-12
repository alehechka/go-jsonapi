package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_appendBaseURL_RelativeURL(t *testing.T) {
	link := Link{
		Href: "/example",
	}

	transformed := appendBaseURL(link, baseURL)

	assert.NotEqual(t, link, transformed)
	assert.Equal(t, baseURL+link.Href, transformed.Href)
}

func Test_appendBaseURL_AbsoluteURL(t *testing.T) {
	link := Link{
		Href: "http://absolute.example.com/example",
	}

	transformed := appendBaseURL(link, baseURL)

	assert.Equal(t, link, transformed)
	assert.Equal(t, link.Href, transformed.Href)
}

func Test_substitutePathParams_NilParams(t *testing.T) {
	link := Link{
		Href: "/example",
	}

	transformed := substitutePathParams(link)

	assert.Equal(t, link, transformed)
}

func Test_substitutePathParams_EmptyParams(t *testing.T) {
	link := Link{
		Href:   "/example",
		Params: make(Params),
	}

	transformed := substitutePathParams(link)

	assert.Equal(t, link, transformed)
}

func Test_substitutePathParams_WithParams_NoPathParams(t *testing.T) {
	link := Link{
		Href: "/example",
		Params: Params{
			"id": 123,
		},
	}

	transformed := substitutePathParams(link)

	assert.Equal(t, link, transformed)
}

func Test_substitutePathParams_WithParams_WithPathParams(t *testing.T) {
	link := Link{
		Href: "/example/:id",
		Params: Params{
			"id": 123,
		},
	}

	transformed := substitutePathParams(link)

	assert.NotEqual(t, link, transformed)
	assert.Equal(t, "/example/123", transformed.Href)
}

func Test_substitutePathParams_WithParams_WithExtraPathParams(t *testing.T) {
	link := Link{
		Href: "/example/:id/children/:childID",
		Params: Params{
			"id": 123,
		},
	}

	transformed := substitutePathParams(link)

	assert.NotEqual(t, link, transformed)
	assert.Equal(t, "/example/123/children/:childID", transformed.Href)
}

func Test_substitutePathParams_WithParams_AbsoluteURL(t *testing.T) {
	link := Link{
		Href: "http://example.com/example/:id",
		Params: Params{
			"id": 123,
		},
	}

	transformed := substitutePathParams(link)

	assert.NotEqual(t, link, transformed)
	assert.Equal(t, "http://example.com/example/123", transformed.Href)
}

func Test_appendQueryParams_InvalidURL(t *testing.T) {
	link := Link{
		Href: "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
		Params: Params{
			"abc": 123,
		},
		Queries: Queries{
			PageLimit.String(): 10,
		},
	}

	transformed := appendQueryParams(link)

	assert.Equal(t, link, transformed)
}

func Test_appendQueryParams_NoQueries(t *testing.T) {
	link := Link{
		Href: "/example",
	}

	transformed := appendQueryParams(link)

	assert.Equal(t, link, transformed)
}

func Test_appendQueryParams_WithQueries(t *testing.T) {
	link := Link{
		Href: "/example",
		Queries: Queries{
			PageSize.String(): 123,
		},
	}

	transformed := appendQueryParams(link)

	assert.NotEqual(t, link, transformed)
	assert.Equal(t, link.Href+"?page[size]=123", transformed.Href)
}

func Test_appendQueryParams_WithExistingQueries(t *testing.T) {
	link := Link{
		Href: "/example?page[number]=5",
		Queries: Queries{
			PageSize.String(): 123,
		},
	}

	transformed := appendQueryParams(link)

	assert.NotEqual(t, link, transformed)
	assert.Equal(t, link.Href+"&page[size]=123", transformed.Href)
}

func Test_stringOrLinkObject_NilMeta(t *testing.T) {
	link := Link{
		Href: "/example",
	}

	transformed := stringOrLinkObject(link)

	href, ok := transformed.(string)
	assert.True(t, ok)
	assert.Equal(t, href, "/example")
}

func Test_stringOrLinkObject_EmptyMeta(t *testing.T) {
	link := Link{
		Href: "/example",
		Meta: make(Meta),
	}

	transformed := stringOrLinkObject(link)

	href, ok := transformed.(string)
	assert.True(t, ok)
	assert.Equal(t, href, "/example")
}

func Test_stringOrLinkObject_WithMeta(t *testing.T) {
	link := Link{
		Href: "/example",
		Meta: Meta{"something": "cool"},
	}

	transformed := stringOrLinkObject(link)

	transformedLink, ok := transformed.(Link)
	assert.True(t, ok)
	assert.Equal(t, link, transformedLink)
}
