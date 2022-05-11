package jsonapi_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi/jsonapi"
	"github.com/stretchr/testify/assert"
)

func Test_Params_Initialize(t *testing.T) {
	var params jsonapi.Params
	assert.Nil(t, params)

	params.Initialize()
	assert.NotNil(t, params)
}

func Test_Queries_Initialize(t *testing.T) {
	var queries jsonapi.Queries
	assert.Nil(t, queries)

	queries.Initialize()
	assert.NotNil(t, queries)
}

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
				jsonapi.PageLimit.String(): 25,
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

func Test_PageSizeNextLinks(t *testing.T) {
	path := "/example"
	num := 10
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[number]=10&page[size]=10", nil)
	link := jsonapi.Link{Href: path, Params: jsonapi.Params{"id": num}, Queries: jsonapi.Queries{"something": "else"}}
	links := jsonapi.PageSizeNextLinks(req)(link, true)

	assert.NotNil(t, links)
	assert.Equal(t, 1, len(links))

	nextLink := links[jsonapi.NextKey]

	transformed := jsonapi.TransformLink(nextLink, "https://example.com")

	assert.Equal(t, path, nextLink.Href)
	assert.Equal(t, "https://example.com/example?page[number]=11&page[size]=10&something=else", transformed)
	assert.Equal(t, num, nextLink.Params["id"])
	assert.Equal(t, num+1, nextLink.Queries[jsonapi.PageNumber.String()])
	assert.Equal(t, num, nextLink.Queries[jsonapi.PageSize.String()])
	assert.Equal(t, link.Queries["something"], nextLink.Queries["something"])
}

func Test_PageSizeNextLink(t *testing.T) {
	path := "/example"
	num := 10
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[number]=10&page[size]=10", nil)
	link := jsonapi.Link{Href: path, Params: jsonapi.Params{"id": num}}
	nextLink := jsonapi.PageSizeNextLink(req)(link)

	transformed := jsonapi.TransformLink(nextLink, "https://example.com")

	assert.Equal(t, path, nextLink.Href)
	assert.Equal(t, "https://example.com/example?page[number]=11&page[size]=10", transformed)
	assert.Equal(t, num, nextLink.Params["id"])
	assert.Equal(t, num+1, nextLink.Queries[jsonapi.PageNumber.String()])
	assert.Equal(t, num, nextLink.Queries[jsonapi.PageSize.String()])
}

func Test_PageLimitNextLinks(t *testing.T) {
	path := "/example"
	num := 10
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[offset]=10&page[limit]=10", nil)
	link := jsonapi.Link{Href: path, Params: jsonapi.Params{"id": num}}
	links := jsonapi.PageLimitNextLinks(req)(link, true, num)

	assert.Equal(t, 1, len(links))

	nextLink := links[jsonapi.NextKey]

	transformed := jsonapi.TransformLink(nextLink, "https://example.com")

	assert.Equal(t, path, nextLink.Href)
	assert.Equal(t, "https://example.com/example?page[limit]=10&page[offset]=20", transformed)
	assert.Equal(t, num, nextLink.Params["id"])
	assert.Equal(t, num+num, nextLink.Queries[jsonapi.PageOffset.String()])
	assert.Equal(t, num, nextLink.Queries[jsonapi.PageLimit.String()])
}

func Test_PageLimitNextLink(t *testing.T) {
	path := "/example"
	num := 10
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[offset]=10&page[limit]=10", nil)
	link := jsonapi.Link{Href: path, Params: jsonapi.Params{"id": num}}
	nextLink := jsonapi.PageLimitNextLink(req)(link, num)

	transformed := jsonapi.TransformLink(nextLink, "https://example.com")

	assert.Equal(t, path, nextLink.Href)
	assert.Equal(t, "https://example.com/example?page[limit]=10&page[offset]=20", transformed)
	assert.Equal(t, num, nextLink.Params["id"])
	assert.Equal(t, num+num, nextLink.Queries[jsonapi.PageOffset.String()])
	assert.Equal(t, num, nextLink.Queries[jsonapi.PageLimit.String()])
}

func Test_CursorNextPrevLinks(t *testing.T) {
	path := "/example"
	num := 10
	before := "1234"
	after := "4321"
	links := jsonapi.CursorNextPrevLinks(path, jsonapi.Params{"id": num}, num, &before, &after)

	assert.Equal(t, 2, len(links))

	nextLink := links[jsonapi.NextKey]
	transformedNext := jsonapi.TransformLink(nextLink, "https://example.com")
	assert.Equal(t, path, nextLink.Href)
	assert.Equal(t, "https://example.com/example?page[after]=4321&page[size]=10", transformedNext)
	assert.Equal(t, num, nextLink.Params["id"])
	assert.Equal(t, num, nextLink.Queries[jsonapi.PageSize.String()])
	assert.Equal(t, after, nextLink.Queries[jsonapi.PageAfter.String()])

	prevLink := links["prev"]
	transformedPrev := jsonapi.TransformLink(prevLink, "https://example.com")
	assert.Equal(t, path, prevLink.Href)
	assert.Equal(t, "https://example.com/example?page[before]=1234&page[size]=10", transformedPrev)
	assert.Equal(t, num, prevLink.Params["id"])
	assert.Equal(t, num, prevLink.Queries[jsonapi.PageSize.String()])
	assert.Equal(t, before, prevLink.Queries[jsonapi.PageBefore.String()])
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
	assert.Equal(t, num, link.Queries[jsonapi.PageSize.String()])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter.String()])
	assert.Equal(t, jsonapi.NextKey, key)
	assert.True(t, isLink)
}

func Test_CursorNextLink_NilAfter(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CursorNextLink(path, jsonapi.Params{"id": num}, num, nil)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, path, link.Href)
	assert.Equal(t, "https://example.com/example?page[size]=10", transformed)
	assert.Equal(t, num, link.Params["id"])
	assert.Equal(t, num, link.Queries[jsonapi.PageSize.String()])
	assert.Nil(t, link.Queries[jsonapi.PageAfter.String()])
	assert.Equal(t, jsonapi.NextKey, key)
	assert.True(t, isLink)
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
	assert.Equal(t, num, link.Queries[jsonapi.PageSize.String()])
	assert.Equal(t, after, link.Queries[jsonapi.PageAfter.String()])
	assert.Equal(t, jsonapi.NextKey, key)
	assert.True(t, isLink)
}

func Test_CursorNextLink_NoSize_NilCursor(t *testing.T) {
	path := "/example"
	num := 0
	link, key, isLink := jsonapi.CursorNextLink(path, jsonapi.Params{"id": num}, num, nil)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.False(t, isLink)
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
	assert.Equal(t, num, link.Queries[jsonapi.PageSize.String()])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore.String()])
	assert.Equal(t, jsonapi.PreviousKey, key)
	assert.True(t, isLink)
}

func Test_CursorPrevLink_NilBefore(t *testing.T) {
	path := "/example"
	num := 10
	link, key, isLink := jsonapi.CursorPrevLink(path, jsonapi.Params{"id": num}, num, nil)

	transformed := jsonapi.TransformLink(link, "https://example.com")

	assert.Equal(t, "", link.Href)
	assert.Equal(t, "https://example.com", transformed)
	assert.Nil(t, link.Params["id"])
	assert.Nil(t, link.Queries[jsonapi.PageSize.String()])
	assert.Nil(t, link.Queries[jsonapi.PageBefore.String()])
	assert.Equal(t, "", key)
	assert.False(t, isLink)
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
	assert.Equal(t, num, link.Queries[jsonapi.PageSize.String()])
	assert.Equal(t, before, link.Queries[jsonapi.PageBefore.String()])
	assert.Equal(t, jsonapi.PreviousKey, key)
	assert.True(t, isLink)
}

func Test_CursorPrevLink_NoSize_NilCursor(t *testing.T) {
	path := "/example"
	num := 0
	link, key, isLink := jsonapi.CursorPrevLink(path, jsonapi.Params{"id": num}, num, nil)

	assert.Equal(t, "", link.Href)
	assert.Equal(t, 0, len(link.Params))
	assert.Equal(t, 0, len(link.Queries))
	assert.Equal(t, "", key)
	assert.False(t, isLink)
}

func Test_IsAbsoluteURL_True(t *testing.T) {
	isAbsolute := jsonapi.IsAbsoluteURL("https://example.com")
	assert.True(t, isAbsolute)
}

func Test_IsAbsoluteURL_False(t *testing.T) {
	isAbsolute := jsonapi.IsAbsoluteURL("example.com")
	assert.False(t, isAbsolute)
}

func Test_IsRelativeURL_True(t *testing.T) {
	isRelative := jsonapi.IsRelativeURL("example.com")
	assert.True(t, isRelative)
}

func Test_IsRelativeURL_TruePath(t *testing.T) {
	isRelative := jsonapi.IsRelativeURL("/path/to/resource")
	assert.True(t, isRelative)
}

func Test_IsARelativeURL_False(t *testing.T) {
	isRelative := jsonapi.IsRelativeURL("https://example.com")
	assert.False(t, isRelative)
}
