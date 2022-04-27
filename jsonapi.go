// Package jsonapi is a collection of useful wrappers for creating responses that
// adhere to the JSON API spec
package jsonapi

/*
	Usage:
	Create a data struct that implements the Data interface and pass that data into the CreateResponse function

	For relationships you will have to create two additional structs that implement the Relationship and ResourceIdentifier respectively
*/

// Response is the standard JSONAPI Response struct
type Response struct {
	Data     []Data
	Included []Data
	Errors   []Error
	Links    Links
	Meta     interface{}
}

type TransformedResponse struct {
	Data     interface{}     `json:"data,omitempty"` // Data | []Data
	Errors   []internalError `json:"errors,omitempty"`
	Included []internalData  `json:"included,omitempty"`
	Links    LinkMap         `json:"links,omitempty"`
	Meta     interface{}     `json:"meta,omitempty"`
}

// CreateResponse transforms provided parameters into standardized JSONAPI format
func CreateResponse(r Response, baseURL string) TransformedResponse {

	return TransformedResponse{
		Data:     transformToInternalDataStructs(r.Data, baseURL),
		Included: transformToInternalDataStructArray(r.Included, baseURL),
		Errors:   transformToInternalErrorStructs(r.Errors, baseURL),
		Links:    TransformLinks(r.Links, baseURL),
		Meta:     r.Meta,
	}
}
