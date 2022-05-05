package jsonapi_test

import (
	"testing"

	"github.com/alehechka/go-jsonapi"
	"github.com/stretchr/testify/assert"
)

func Test_TransformLinks(t *testing.T) {

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

func Test_TransformLink(t *testing.T) {
	link := jsonapi.Link{
		Href: "/api/objects",
		Queries: jsonapi.Queries{
			"offset": 25,
		},
	}

	transformed := jsonapi.TransformLink(link, "")

	assert.Equal(t, "/api/objects?offset=25", transformed)
}

func Test_NumberNextLinks(t *testing.T) {
	path := "/example"
	num := 10
	links := jsonapi.NumberNextLinks(path, jsonapi.Params{"id": num}, true, num, num)

	assert.Equal(t, 1, len(links))

	link := links["next"]

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[number]=10&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageNumber])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
}

func Test_NumberNextLink_IsLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.NumberNextLink(path, jsonapi.Params{"id": num}, true, num, num)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[number]=10&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageNumber])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func Test_NumberNextLink_IsNotLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.NumberNextLink(path, jsonapi.Params{"id": num}, false, num, num)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.Equal(t, false, isLink)
}

func Test_OffsetNextLinks(t *testing.T) {
	path := "/example"
	num := 10
	links := jsonapi.OffsetNextLinks(path, jsonapi.Params{"id": num}, true, num, num)

	assert.Equal(t, 1, len(links))

	link := links["next"]

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[limit]=10&page[offset]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageOffset])
	assert.Equal(t, num, link.Queries[jsonapi.PageLimit])
}

func Test_OffsetNextLink_IsLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.OffsetNextLink(path, jsonapi.Params{"id": num}, true, num, num)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[limit]=10&page[offset]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageOffset])
	assert.Equal(t, num, link.Queries[jsonapi.PageLimit])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func Test_OffsetNextLink_IsNotLink(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.OffsetNextLink(path, jsonapi.Params{"id": num}, false, num, num)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.Equal(t, false, isLink)
}

func Test_CursorNextPrevLinks(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	after := "4321"
	links := jsonapi.CursorNextPrevLinks(path, jsonapi.Params{"id": num}, num, &before, &after)

	assert.Equal(t, 2, len(links))

	nextLink := links["next"]
	transformedNext := jsonapi.TransformLink(nextLink, "https://example.com")
	assert.Equal(t, path, nextLink.Href)
	assert.Equal(t, "https://example.com/example?page[after]=4321&page[size]=10", transformedNext)
	assert.Equal(t, num, nextLink.Params["id"])
	assert.Equal(t, num, nextLink.Queries[jsonapi.PageSize])
	assert.Equal(t, after, nextLink.Queries[jsonapi.PageAfter])

	prevLink := links["prev"]
	transformedPrev := jsonapi.TransformLink(prevLink, "https://example.com")
	assert.Equal(t, path, prevLink.Href)
	assert.Equal(t, "https://example.com/example?page[before]=1234&page[size]=10", transformedPrev)
	assert.Equal(t, num, prevLink.Params["id"])
	assert.Equal(t, num, prevLink.Queries[jsonapi.PageSize])
	assert.Equal(t, before, prevLink.Queries[jsonapi.PageBefore])
}

func Test_CursorNextLink(t *testing.T) {
	path := "/example"
	num := 10
	after := "4321"
	link, key, isLink := jsonapi.CursorNextLink(path, jsonapi.Params{"id": num}, num, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[after]=4321&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func Test_CursorNextLink_NilAfter(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CursorNextLink(path, jsonapi.Params{"id": num}, num, nil)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, nil, link.Queries[jsonapi.PageAfter])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func Test_CursorNextLink_EmptyString(t *testing.T) {
	path := "/example"
	num := 10
	after := ""
	link, key, isLink := jsonapi.CursorNextLink(path, jsonapi.Params{"id": num}, num, &after)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter])
	assert.Equal(t, "next", key)
	assert.Equal(t, true, isLink)
}

func Test_CursorNextLink_NoSize_NilCursor(t *testing.T) {
	path := "/example"
	num := 0
	link, key, isLink := jsonapi.CursorNextLink(path, jsonapi.Params{"id": num}, num, nil)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.Equal(t, false, isLink)
}

func Test_CursorPrevLink(t *testing.T) {
	path := "/example"
	num := 10
	before := "4321"
	link, key, isLink := jsonapi.CursorPrevLink(path, jsonapi.Params{"id": num}, num, &before)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[before]=4321&page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, "prev", key)
	assert.Equal(t, true, isLink)
}

func Test_CursorPrevLink_NilBefore(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CursorPrevLink(path, jsonapi.Params{"id": num}, num, nil)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, "", link.Href)
	assert.Equal(t, "https://example.com", transformed)
	assert.Equal(t, nil, link.Params["id"])
	assert.Equal(t, nil, link.Queries[jsonapi.PageSize])
	assert.Equal(t, nil, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, "", key)
	assert.Equal(t, false, isLink)
}

func Test_CursorPrevLink_EmptyString(t *testing.T) {
	path := "/example"
	num := 10
	before := ""
	link, key, isLink := jsonapi.CursorPrevLink(path, jsonapi.Params{"id": num}, num, &before)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore])
	assert.Equal(t, "prev", key)
	assert.Equal(t, true, isLink)
}

func Test_CursorPrevLink_NoSize_NilCursor(t *testing.T) {
	path := "/example"
	num := 0
	link, key, isLink := jsonapi.CursorPrevLink(path, jsonapi.Params{"id": num}, num, nil)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.Equal(t, false, isLink)
}
