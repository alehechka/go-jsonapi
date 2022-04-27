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

	assert.Equal(t, transformed["1"], "https://example.com/api/objects")
	assert.Equal(t, transformed["2"], "https://example.com/api/objects/123")
	assert.Equal(t, transformed["3"], "https://example.com/api/objects?offset=25")
	assert.Equal(t, transformed["4"], "https://example.com/api/objects?page[limit]=25")
	assert.Equal(t, transformed["5"].(jsonapi.Link).Href, "https://example.com/api/objects")
	assert.Equal(t, transformed["5"].(jsonapi.Link).Meta["random"], links["5"].Meta["random"])
}

func TestTransformLink(t *testing.T) {
	link := jsonapi.Link{
		Href: "/api/objects",
		Queries: jsonapi.Queries{
			"offset": 25,
		},
	}

	transformed := jsonapi.TransformLink(link, "")

	assert.Equal(t, transformed, "/api/objects?offset=25")
}

func TestCreateNextLinksFromOffsetPaginationResponse(t *testing.T) {
	path := "/example"
	num := 10
	links := jsonapi.CreateNextLinksFromOffsetPaginationResponse(path, jsonapi.Params{"id": num}, true, num, num)

	assert.Equal(t, len(links), 1)

	link := links["next"]

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, link.Href, path)
	assert.Equal(t, transformed, "https://example.com/example?page[limit]=10&page[offset]=10")
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[jsonapi.PageOffset], num)
	assert.Equal(t, link.Queries[jsonapi.PageLimit], num)
}

func TestCreateNextLinkFromOffsetPaginationResponse_IsLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CreateNextLinkFromOffsetPaginationResponse(path, jsonapi.Params{"id": num}, true, num, num)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, link.Href, path)
	assert.Equal(t, transformed, "https://example.com/example?page[limit]=10&page[offset]=10")
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[jsonapi.PageOffset], num)
	assert.Equal(t, link.Queries[jsonapi.PageLimit], num)
	assert.Equal(t, key, "next")
	assert.Equal(t, isLink, true)
}

func TestCreateNextLinkFromOffsetPaginationResponse_IsNotLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CreateNextLinkFromOffsetPaginationResponse(path, jsonapi.Params{"id": num}, false, num, num)

	assert.Equal(t, link.Href, "")
	assert.Equal(t, len(link.Params), 0)
	assert.Equal(t, len(link.Queries), 0)
	assert.Equal(t, len(link.Queries), 0)
	assert.Equal(t, key, "")
	assert.Equal(t, isLink, false)
}

func TestCreateNextLinksFromCursorPaginationResponse(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	after := "4321"
	links := jsonapi.CreateNextLinksFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, &after)

	assert.Equal(t, len(links), 1)

	link := links["next"]

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, link.Href, path)
	assert.Equal(t, transformed, "https://example.com/example?page[after]=4321&page[before]=1234&page[size]=10")
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[jsonapi.PageSize], num)
	assert.Equal(t, link.Queries[jsonapi.PageBefore], before)
	assert.Equal(t, link.Queries[jsonapi.PageAfter], after)
}

func TestCreateNextLinkFromCursorPaginationResponse(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	after := "4321"
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, link.Href, path)
	assert.Equal(t, transformed, "https://example.com/example?page[after]=4321&page[before]=1234&page[size]=10")
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[jsonapi.PageSize], num)
	assert.Equal(t, link.Queries[jsonapi.PageBefore], before)
	assert.Equal(t, link.Queries[jsonapi.PageAfter], after)
	assert.Equal(t, key, "next")
	assert.Equal(t, isLink, true)
}

func TestCreateNextLinkFromCursorPaginationResponse_NilBefore(t *testing.T) {
	path := "/example"
	num := 10
	after := "4321"
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, nil, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, link.Href, path)
	assert.Equal(t, transformed, "https://example.com/example?page[after]=4321&page[size]=10")
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[jsonapi.PageSize], num)
	assert.Equal(t, link.Queries[jsonapi.PageBefore], nil)
	assert.Equal(t, link.Queries[jsonapi.PageAfter], after)
	assert.Equal(t, key, "next")
	assert.Equal(t, isLink, true)
}

func TestCreateNextLinkFromCursorPaginationResponse_NilAfter(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, nil)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, link.Href, path)
	assert.Equal(t, transformed, "https://example.com/example?page[before]=1234&page[size]=10")
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[jsonapi.PageSize], num)
	assert.Equal(t, link.Queries[jsonapi.PageBefore], before)
	assert.Equal(t, link.Queries[jsonapi.PageAfter], nil)
	assert.Equal(t, key, "next")
	assert.Equal(t, isLink, true)
}

func TestCreateNextLinkFromCursorPaginationResponse_EmptyStrings(t *testing.T) {
	path := "/example"
	num := 10
	before := ""
	after := ""
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, &before, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, link.Href, path)
	assert.Equal(t, transformed, "https://example.com/example?page[size]=10")
	assert.Equal(t, link.Params["id"], num)
	assert.Equal(t, link.Queries[jsonapi.PageSize], num)
	assert.Equal(t, link.Queries[jsonapi.PageBefore], before)
	assert.Equal(t, link.Queries[jsonapi.PageAfter], after)
	assert.Equal(t, key, "next")
	assert.Equal(t, isLink, true)
}

func TestCreateNextLinkFromCursorPaginationResponse_NilPointers(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CreateNextLinkFromCursorPaginationResponse(path, jsonapi.Params{"id": num}, num, nil, nil)

	assert.Equal(t, link.Href, "")
	assert.Equal(t, len(link.Params), 0)
	assert.Equal(t, len(link.Queries), 0)
	assert.Equal(t, len(link.Queries), 0)
	assert.Equal(t, key, "")
	assert.Equal(t, isLink, false)
}
