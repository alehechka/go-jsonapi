package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/alehechka/go-jsonapi"
)

func appendGeneratedSelfLink(request *http.Request) func(links jsonapi.Links, baseURL string, path string) jsonapi.Links {
	return func(links jsonapi.Links, baseURL string, path string) jsonapi.Links {

		href := baseURL
		// only append path if it is not already the suffix of baseURL to avoid duplicating path
		if !strings.HasSuffix(baseURL, path) {
			href += path
		}

		queryValues := request.URL.Query()
		if queryValues != nil && len(queryValues) > 0 {
			query, _ := url.QueryUnescape(queryValues.Encode())
			href += fmt.Sprintf("?%s", query)
		}

		if links == nil {
			links = make(jsonapi.Links)
		}

		links["self"] = jsonapi.Link{
			Href: href,
		}

		return links
	}
}
