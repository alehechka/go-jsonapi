package jsonapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CreateResponse is a wrapper to TransformResponse that will create the baseURL parameters from *http.Request
func CreateResponse(request *http.Request) func(r Response) TransformedResponse {
	return func(r Response) TransformedResponse {
		baseURL, path := CreateBaseURL(request)
		r.Links = AppendGeneratedSelfLink(request)(r.Links, baseURL, path)

		return TransformResponse(r, baseURL)
	}
}

// CreateCollectionResponse is a wrapper to TransformCollectionResponse that will create the baseURL parameters from *http.Request
func CreateCollectionResponse(request *http.Request) func(r CollectionResponse) TransformedResponse {
	return func(r CollectionResponse) TransformedResponse {
		baseURL, path := CreateBaseURL(request)
		r.Links = AppendGeneratedSelfLink(request)(r.Links, baseURL, path)

		return TransformCollectionResponse(r, baseURL)
	}
}

// AppendGeneratedSelfLink will generate a self link object based on provided *http.Request
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

		links[SelfKey] = Link{
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

// CreateBaseURL will generate a baseURL, prioritizing proxied http headers
func CreateBaseURL(request *http.Request) (baseURL string, path string) {
	host := request.Header.Get(ForwardedHost) // host from primary source
	if len(host) == 0 {
		host = request.Host
	}

	protocol := request.Header.Get(ForwardedProto) // http | https
	if len(protocol) == 0 {
		if len(request.URL.Scheme) > 0 {
			protocol = request.URL.Scheme
		} else {
			protocol = "http"
		}
	}

	prefix := request.Header.Get(ForwardedPrefix) // proxied prefix
	path = request.URL.Path                       // self-path

	return fmt.Sprintf("%s://%s%s", protocol, host, prefix), path
}
