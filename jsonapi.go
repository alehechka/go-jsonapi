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
	Data     interface{}     `json:"data,omitempty"` // Data | []Data
	Errors   []internalError `json:"errors,omitempty"`
	Included []internalData  `json:"included,omitempty"`
	Links    LinkMap         `json:"links,omitempty"`
	Meta     interface{}     `json:"meta,omitempty"`
}

// CreateResponse transforms provided parameters into standardized JSONAPI format
func CreateResponse(data, included []Data, errors []Error, links Links, meta interface{}, baseURL string) Response {

	return Response{
		Data:     transformToInternalDataStructs(data, baseURL),
		Included: transformToInternalDataStructArray(included, baseURL),
		Errors:   transformToInternalErrorStructs(errors, baseURL),
		Links:    transformLinks(links, baseURL),
		Meta:     meta,
	}
}
