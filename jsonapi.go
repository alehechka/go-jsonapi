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
	Node     Node
	Included interface{} // []Node
	Errors   Errors
	Links    Links
	Meta     interface{}
}

// CollectionResponse is the standard JSONAPI collection Response struct
type CollectionResponse struct {
	Nodes    interface{} // Node | []Node
	Included interface{} // []Node
	Errors   Errors
	Links    Links
	Meta     interface{}
}

// TransformedResponse is the resulting Data struct after transforming via TransformResponse/TransformCollectionResponse
type TransformedResponse struct {
	Data     interface{}     `json:"data,omitempty"` // Node | []Node
	Errors   []internalError `json:"errors,omitempty"`
	Included []internalNode  `json:"included,omitempty"`
	Links    LinkMap         `json:"links,omitempty"`
	Meta     interface{}     `json:"meta,omitempty"`
}

// TransformResponse transforms provided parameters into standardized JSONAPI format
func TransformResponse(r Response, baseURL string) TransformedResponse {
	data := transformResponseNode(r, baseURL)

	return TransformedResponse{
		Data:     data,
		Included: transformIncluded(r.Included, data, baseURL),
		Errors:   transformToInternalErrorStructs(r.Errors, baseURL),
		Links:    TransformLinks(r.Links, baseURL),
		Meta:     r.Meta,
	}
}

// TransformCollectionResponse transforms provided parameters into standardized collection JSONAPI format
func TransformCollectionResponse(r CollectionResponse, baseURL string) TransformedResponse {
	nodes := transformCollectionResponseNode(r, baseURL)

	return TransformedResponse{
		Data:     nodes,
		Included: transformIncluded(r.Included, nodes, baseURL),
		Errors:   transformToInternalErrorStructs(r.Errors, baseURL),
		Links:    TransformLinks(r.Links, baseURL),
		Meta:     r.Meta,
	}
}
