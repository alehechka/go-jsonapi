package jsonapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CreateJSONAPIResponse is a wrapper to jsonapi.CreateResponse that will create the baseURL parameters from gin.Context
func CreateJSONAPIResponse(request *http.Request) func(r Response) TransformedResponse {
	return func(r Response) TransformedResponse {
		baseURL, path := CreateBaseURL(request)
		r.Links = AppendGeneratedSelfLink(request)(r.Links, baseURL, path)

		return CreateResponse(r, baseURL)
	}
}

func AppendGeneratedSelfLink(request *http.Request) func(links Links, baseURL string, path string) Links {
	return func(links Links, baseURL string, path string) Links {

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
			links = make(Links)
		}

		links["self"] = Link{
			Href: href,
		}

		return links
	}
}

const (
	// ForwardedPrefix represents the prefix that is dropped when proxied through rest-api
	ForwardedPrefix string = "X-Forwarded-Prefix"
	// ForwardedProto represents the protocol that is received prior to being forwarded (http | https)
	ForwardedProto string = "X-Forwarded-Proto"
	// ForwardedHost represents the original host header received prior to being forwarded (example.com)
	ForwardedHost string = "X-Forwarded-Host"
)

func CreateBaseURL(request *http.Request) (baseURL string, path string) {
	host := request.Header.Get(ForwardedHost) // host from primary source
	if len(host) == 0 {
		host = request.Host
	}

	protocol := request.Header.Get(ForwardedProto) // http | https
	if len(protocol) == 0 {
		protocol = request.URL.Scheme
	}

	prefix := request.Header.Get(ForwardedPrefix) // proxied prefix
	path = request.URL.Path                       // self-path

	return fmt.Sprintf("%s://%s%s", protocol, host, prefix), path
}
