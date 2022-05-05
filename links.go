package jsonapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	// NextKey represents the key for the next link
	NextKey string = "next"
	// PreviousKey represents the key for the previous link
	PreviousKey string = "prev"
	// SelfKey represents the key for the self link
	SelfKey string = "self"
)

// Meta is a map of meta data to attach with a response.
type Meta map[string]interface{}

// Params is a map of path parameters to substitute into a url.
type Params map[string]interface{}

// Initialize will initialize the Params map if it is nil
func (p *Params) Initialize() {
	if *p == nil {
		*p = make(Params)
	}
}

// Queries is a map of query parameters to append to a url.
type Queries map[string]interface{}

// Initialize will initialize the Queries map if it is nil
func (q *Queries) Initialize() {
	if *q == nil {
		*q = make(Queries)
	}
}

// Link is the standard JSONAPI Link object.
// For more info: https://jsonapi.org/format/#document-links
type Link struct {
	// Href is a string container the link's URL. Using relative URLs will allow a baseURL to be prefixed.
	Href string `json:"href,omitempty"`
	// Meta is an object containing non-standard meta information about the link.
	Meta Meta `json:"meta,omitempty"`
	// Params represent path parameters that will be substituted into Href. Params will be omitted from the response.
	Params Params `json:"-"`
	// Queries represent query parameters that will be appended to Href. Queries will be omitted from the response.
	Queries Queries `json:"-"`
}

// Links is a map of JsonAPILink objects
type Links map[string]Link

// LinkMap should have values of type JsonAPILink or string
type LinkMap map[string]interface{} // JsonAPILink | string

// TransformLinks transforms provided Links map into a JSON:API LinkMap
func TransformLinks(jsonLinks Links, baseURL string) LinkMap {
	links := make(LinkMap)

	for key, jsonLink := range jsonLinks {
		links[key] = TransformLink(jsonLink, baseURL)
	}

	return links
}

// TransformLink transforms an individual Link into a JSON:API link object
func TransformLink(jsonLink Link, baseURL string) (link interface{}) {

	jsonLink = appendBaseURL(jsonLink, baseURL)

	jsonLink = substitutePathParams(jsonLink)

	jsonLink = appendQueryParams(jsonLink)

	return stringOrLinkObject(jsonLink)
}

func appendBaseURL(link Link, baseURL string) Link {
	// only append baseURL if href is a relative URL
	if IsRelativeURL(link.Href) {
		link.Href = fmt.Sprintf("%s%s", baseURL, link.Href)
	}

	return link
}

func substitutePathParams(link Link) Link {
	if link.Params != nil && len(link.Params) > 0 {
		pathParts := strings.Split(link.Href, "/")

		for index, pathPart := range pathParts {
			if strings.HasPrefix(pathPart, ":") {
				paramString := strings.TrimPrefix(pathPart, ":")
				if param, exists := link.Params[paramString]; exists {
					pathParts[index] = fmt.Sprintf("%v", param)
				}
			}
		}

		link.Href = strings.Join(pathParts, "/")
	}

	return link
}

func appendQueryParams(link Link) Link {
	if link.Queries == nil || len(link.Queries) == 0 {
		return link
	}

	u, err := url.Parse(link.Href)
	if err != nil {
		return link
	}

	q := u.Query()
	for key, value := range link.Queries {
		if v, ok := value.(string); ok && len(v) == 0 {
			continue
		}
		q.Set(key, fmt.Sprintf("%v", value))
	}
	u.RawQuery, _ = url.QueryUnescape(q.Encode())

	link.Href = u.String()

	return link
}

func stringOrLinkObject(jsonLink Link) (link interface{}) {
	if jsonLink.Meta == nil || len(jsonLink.Meta) == 0 {
		return jsonLink.Href
	}

	return jsonLink
}

// IsURL parses string and returns boolean if string is valid URL
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsRelativeURL parses string and returns boolean is string is a relative URL
func IsRelativeURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme == ""
}

// NumberNextLinks creates a Links map for next pagination step (using PageNumber/PageSize)
func NumberNextLinks(request *http.Request) func(link Link, moreResultsAvailable bool) Links {
	return func(link Link, moreResultsAvailable bool) Links {
		links := make(Links)

		if moreResultsAvailable {
			links[NextKey] = NumberNextLink(request)(link)
		}

		return links
	}
}

// NumberNextLink creates a Link object for next pagination step (using PageNumber/PageSize)
func NumberNextLink(request *http.Request) func(link Link) (nextLink Link) {
	return func(link Link) Link {

		link.Queries.Initialize()
		pageNumber, _ := GetPageNumber(request)
		link.Queries[PageNumber.String()] = pageNumber + 1
		link.Queries[PageSize.String()], _ = GetPageSize(request)

		return link
	}
}

// OffsetNextLinks creates a Links map for next pagination step (using PageOffset/PageLimit)
func OffsetNextLinks(request *http.Request) func(link Link, moreResultsAvailable bool) Links {
	return func(link Link, moreResultsAvailable bool) Links {
		links := make(Links)

		if moreResultsAvailable {
			links[NextKey] = OffsetNextLink(request)(link)
		}

		return links
	}
}

// OffsetNextLink creates a Link object for next pagination step (using PageOffset/PageLimit)
func OffsetNextLink(request *http.Request) func(link Link) (nextLink Link) {
	return func(link Link) (nextLink Link) {

		link.Queries.Initialize()
		pageOffset, _ := GetPageOffset(request)
		pageLimit, _ := GetPageLimit(request)
		link.Queries[PageOffset.String()] = pageOffset + pageLimit
		link.Queries[PageLimit.String()] = pageLimit

		return link
	}
}

// CursorNextPrevLinks creates a Links map for next pagination step (using PageSize/PageBefore/PageAfter)
func CursorNextPrevLinks(href string, params Params, size int, before *string, after *string) Links {
	links := make(Links)

	// Next Link
	if link, key, isLink := CursorNextLink(href, params, size, after); isLink {
		links[key] = link
	}

	// Previous Link
	if link, key, isLink := CursorPrevLink(href, params, size, before); isLink {
		links[key] = link
	}

	return links
}

// CursorNextLink creates a Link object for next pagination step (using PageSize/PageAfter)
func CursorNextLink(href string, params Params, size int, after *string) (link Link, key string, isLink bool) {
	if size == 0 && after == nil {
		return
	}

	queries := Queries{}

	if size > 0 {
		queries[PageSize.String()] = size
	}

	if after != nil {
		queries[PageAfter.String()] = *after
	}

	return Link{
		Href:    href,
		Params:  params,
		Queries: queries,
	}, "next", true
}

// CursorPrevLink creates a Link object for previous pagination step (using PageSize/PageBefore)
func CursorPrevLink(href string, params Params, size int, before *string) (link Link, key string, isLink bool) {
	if before == nil {
		return
	}

	queries := Queries{}

	if size > 0 {
		queries[PageSize.String()] = size
	}

	if before != nil {
		queries[PageBefore.String()] = *before
	}

	return Link{
		Href:    href,
		Params:  params,
		Queries: queries,
	}, "prev", true
}
