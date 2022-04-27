package jsonapi

import (
	"testing"

	"gotest.tools/assert"
)

func TestTransformLinks(t *testing.T) {

	baseURL := "https://example.com"

	links := Links{
		"1": {
			Href: "/api/objects",
		},
		"2": {
			Href: "/api/objects/:id",
			Params: Params{
				"id": 123,
			},
		},
		"3": {
			Href: "/api/objects",
			Queries: Queries{
				"offset": 25,
			},
		},
		"4": {
			Href: "/api/objects",
			Queries: Queries{
				PageLimit.String(): 25,
			},
		},
		"5": {
			Href: "/api/objects",
			Meta: Meta{
				"random": "value",
			},
		},
	}

	transformed := transformLinks(links, baseURL)

	assert.Equal(t, transformed["1"], "https://example.com/api/objects")
	assert.Equal(t, transformed["2"], "https://example.com/api/objects/123")
	assert.Equal(t, transformed["3"], "https://example.com/api/objects?offset=25")
	assert.Equal(t, transformed["4"], "https://example.com/api/objects?page[limit]=25")
	assert.Equal(t, transformed["5"].(Link).Href, "https://example.com/api/objects")
	assert.Equal(t, transformed["5"].(Link).Meta["random"], links["5"].Meta["random"])
}

func TestTransformLink(t *testing.T) {
	link := Link{
		Href: "/api/objects",
		Queries: Queries{
			"offset": 25,
		},
	}

	transformed := transformLink(link, "")

	assert.Equal(t, transformed, "/api/objects?offset=25")
}

func TestCreateNextLinksFromPaginationResponse(t *testing.T) {
	path := "/example"
	num := 10
	links := CreateNextLinksFromPaginationResponse(path, Params{"id": num}, true, num, num)

	assert.Equal(t, len(links), 1)

	link := links["next"]

	assert.Equal(t, link.Href, path)
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[PageOffset.String()], num)
	assert.Equal(t, link.Queries[PageLimit.String()], num)
}

func TestCreateNextLinkFromPaginationResponse_IsLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := CreateNextLinkFromPaginationResponse(path, Params{"id": num}, true, num, num)

	assert.Equal(t, link.Href, path)
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[PageOffset.String()], num)
	assert.Equal(t, link.Queries[PageLimit.String()], num)
	assert.Equal(t, key, "next")
	assert.Equal(t, isLink, true)
}

func TestCreateNextLinkFromPaginationResponse_IsNotLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := CreateNextLinkFromPaginationResponse(path, Params{"id": num}, false, num, num)

	assert.Equal(t, link.Href, "")
	assert.Equal(t, len(link.Params), 0)
	assert.Equal(t, len(link.Queries), 0)
	assert.Equal(t, len(link.Queries), 0)
	assert.Equal(t, key, "next")
	assert.Equal(t, isLink, false)
}
