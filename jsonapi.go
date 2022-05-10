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
	Node   Node
	Errors Errors
	Links  Links
	Meta   any
}

// CollectionResponse is the standard JSONAPI collection Response struct
type CollectionResponse struct {
	Nodes  any // Node | []Node
	Errors Errors
	Links  Links
	Meta   any
}

// TransformedResponse is the resulting Data struct after transforming via TransformResponse/TransformCollectionResponse
type TransformedResponse struct {
	Data     any             `json:"data,omitempty"` // Node | []Node
	Errors   []internalError `json:"errors,omitempty"`
	Included []internalNode  `json:"included,omitempty"`
	Links    LinkMap         `json:"links,omitempty"`
	Meta     any             `json:"meta,omitempty"`
}

// TransformResponse transforms provided parameters into standardized JSONAPI format
func TransformResponse(r Response, baseURL string) TransformedResponse {
	data, included := transformResponseNode(r, baseURL)

	return TransformedResponse{
		Data:     data,
		Included: transformIncluded(included, data, baseURL),
		Errors:   transformErrors(r.Errors, baseURL),
		Links:    TransformLinks(r.Links, baseURL),
		Meta:     r.Meta,
	}
}

// TransformCollectionResponse transforms provided parameters into standardized collection JSONAPI format
func TransformCollectionResponse(r CollectionResponse, baseURL string) TransformedResponse {
	nodes, included := transformCollectionResponseNodes(r, baseURL)

	return TransformedResponse{
		Data:     nodes,
		Included: transformIncluded(included, nodes, baseURL),
		Errors:   transformErrors(r.Errors, baseURL),
		Links:    TransformLinks(r.Links, baseURL),
		Meta:     r.Meta,
	}
}
