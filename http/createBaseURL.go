package http

import (
	"fmt"
	"net/http"
)

const (
	// ForwardedPrefix represents the prefix that is dropped when proxied through rest-api
	ForwardedPrefix string = "X-Forwarded-Prefix"
	// ForwardedProto represents the protocol that is received prior to being forwarded (http | https)
	ForwardedProto string = "X-Forwarded-Proto"
	// ForwardedHost represents the original host header received prior to being forwarded (example.com)
	ForwardedHost string = "X-Forwarded-Host"
)

func createBaseURL(request *http.Request) (baseURL string, path string) {
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
