package http

import (
	"net/http"

	"github.com/alehechka/go-jsonapi"
)

// CreateJSONAPIResponse is a wrapper to jsonapi.CreateResponse that will create the baseURL parameters from gin.Context
func CreateJSONAPIResponse(request *http.Request) func(r jsonapi.Response) jsonapi.TransformedResponse {
	return func(r jsonapi.Response) jsonapi.TransformedResponse {
		baseURL, path := createBaseURL(request)
		r.Links = appendGeneratedSelfLink(request)(r.Links, baseURL, path)

		return jsonapi.CreateResponse(r, baseURL)
	}
}
