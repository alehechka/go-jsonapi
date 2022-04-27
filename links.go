package jsonapi

import (
	"fmt"
	"net/url"
	"strings"
)

// Meta is a map of meta data to attach with a response.
type Meta map[string]interface{}

// Params is a map of path parameters to substitute into a url.
type Params map[string]interface{}

// Queries is a map of query parameters to append to a url.
type Queries map[string]interface{}

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

func TransformLinks(jsonLinks Links, baseURL string) LinkMap {
	links := make(LinkMap)

	for key, jsonLink := range jsonLinks {
		links[key] = TransformLink(jsonLink, baseURL)
	}

	return links
}

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

// CreateNextLinksFromOffsetPaginationResponse creates a Links map for next pagination step (using PageOffset/PageLimit)
func CreateNextLinksFromOffsetPaginationResponse(href string, params Params, moreResultsAvailable bool, pageOffset int, pageLimit int) Links {
	link, key, isLink := CreateNextLinkFromOffsetPaginationResponse(href, params, moreResultsAvailable, pageOffset, pageLimit)

	links := make(Links)

	if isLink {
		links[key] = link
	}

	return links
}

// CreateNextLinkFromOffsetPaginationResponse creates a Link object for next pagination step (using PageOffset/PageLimit)
func CreateNextLinkFromOffsetPaginationResponse(href string, params Params, moreResultsAvailable bool, pageOffset int, pageLimit int) (link Link, key string, isLink bool) {
	if moreResultsAvailable == false {
		return
	}

	return Link{
		Href:   href,
		Params: params,
		Queries: Queries{
			PageOffset: pageOffset,
			PageLimit:  pageLimit,
		},
	}, "next", true
}

// CreateNextLinksFromCursorPaginationResponse creates a Links map for next pagination step (using PageSize/PageBefore/PageAfter)
func CreateNextLinksFromCursorPaginationResponse(href string, params Params, size int, before *string, after *string) Links {
	link, key, isLink := CreateNextLinkFromCursorPaginationResponse(href, params, size, before, after)

	links := make(Links)

	if isLink {
		links[key] = link
	}

	return links
}

// CreateNextLinkFromCursorPaginationResponse creates a Link object for next pagination step (using PageSize/PageBefore/PageAfter)
func CreateNextLinkFromCursorPaginationResponse(href string, params Params, size int, before *string, after *string) (link Link, key string, isLink bool) {
	if before == nil && after == nil {
		return
	}

	queries := Queries{PageSize: size}

	if before != nil {
		queries[PageBefore] = *before
	}

	if after != nil {
		queries[PageAfter] = *after
	}

	return Link{
		Href:    href,
		Params:  params,
		Queries: queries,
	}, "next", true
}
