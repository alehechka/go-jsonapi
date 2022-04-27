package jsonapi_test

import (
	"testing"

	"github.com/alehechka/go-jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestTransformLinks(t *testing.T) {

	baseURL := "https://example.com"

	links := jsonapi.Links{
		"1": {
			Href: "/api/objects",
		},
		"2": {
			Href: "/api/objects/:id",
			Params: jsonapi.Params{
				"id": 123,
			},
		},
		"3": {
			Href: "/api/objects",
			Queries: jsonapi.Queries{
				"offset": 25,
			},
		},
		"4": {
			Href: "/api/objects",
			Queries: jsonapi.Queries{
				jsonapi.PageLimit: 25,
			},
		},
		"5": {
			Href: "/api/objects",
			Meta: jsonapi.Meta{
				"random": "value",
			},
		},
	}

	transformed := jsonapi.TransformLinks(links, baseURL)

	assert.Equal(t, "https://example.com/api/objects", transformed["1"])
	assert.Equal(t, "https://example.com/api/objects/123", transformed["2"])
	assert.Equal(t, "https://example.com/api/objects?offset=25", transformed["3"])
	assert.Equal(t, "https://example.com/api/objects?page[limit]=25", transformed["4"])
	assert.Equal(t, "https://example.com/api/objects", transformed["5"].(jsonapi.Link).Href)
	assert.Equal(t, links["5"].Meta["random"], transformed["5"].(jsonapi.Link).Meta["random"])
}

func TestTransformLink(t *testing.T) {
	link := jsonapi.Link{
		Href: "/api/objects",
		Queries: jsonapi.Queries{
			"offset": 25,
		},
	}

	transformed := jsonapi.TransformLink(link, "")

	assert.Equal(t, "/api/objects?offset=25", transformed)
}

func TestCreateNextLinksFromOffsetPaginationResponse(t *testing.T) {
	path := "/example"
	num := 10
	links := jsonapi.CreateNextLinksFromOffsetPaginationResponse(path, jsonapi.Params{"id": num}, true, num, num)

	assert.Equal(t, 1, len(links))

	link := links["next"]

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[limit]=10&page[offset]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageOffset])
	assert.Equal(t, num, link.Queries[jsonapi.PageLimit])
}

func TestCreateNextLinkFromOffsetPaginationResponse_IsLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CreateNextLinkFromOffsetPaginationResponse(path, jsonapi.Params{"id": num}, true, num, num)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[limit]=10&page[offset]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageOffset])
	assert.Equal(t, num, link.Queries[jsonapi.PageLimit])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func TestCreateNextLinkFromOffsetPaginationResponse_IsNotLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CreateNextLinkFromOffsetPaginationResponse(path, jsonapi.Params{"id": num}, false, num, num)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.Equal(t, false, isLink)
}

func TestCreateNextLinksFromCursorPaginationResponse(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	after := "4321"
	links := jsonapi.CreateNextLinksFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, &after)

	assert.Equal(t, 1, len(links))

	link := links["next"]

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[after]=4321&page[before]=1234&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter])
}

func TestCreateNextLinkFromCursorPaginationResponse(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	after := "4321"
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[after]=4321&page[before]=1234&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func TestCreateNextLinkFromCursorPaginationResponse_NilBefore(t *testing.T) {
	path := "/example"
	num := 10
	after := "4321"
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, nil, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[after]=4321&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, nil, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func TestCreateNextLinkFromCursorPaginationResponse_NilAfter(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, nil)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[before]=1234&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, nil, link.Queries[jsonapi.PageAfter])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func TestCreateNextLinkFromCursorPaginationResponse_EmptyStrings(t *testing.T) {
	path := "/example"
	num := 10
	before := ""
	after := ""
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func TestCreateNextLinkFromCursorPaginationResponse_NilPointers(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, nil, nil)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.Equal(t, false, isLink)
}
