package jsonapi

import "reflect"

// Node is the standard JSONAPI Data struct
type Node interface {
	ID() string
	Type() string
	Attributes() interface{}
}

type Metable interface {
	Meta() interface{}
}

type internalNode struct {
	ID            string                          `json:"id"`
	Type          string                          `json:"type"`
	Attributes    interface{}                     `json:"attributes,omitempty"`
	Links         LinkMap                         `json:"links,omitempty"`
	Relationships map[string]internalRelationship `json:"relationships,omitempty"`
	Meta          interface{}                     `json:"meta,omitempty"`
}

func transformToInternalDataStructArray(payload interface{}, baseURL string) []internalNode {
	internalNodes := make([]internalNode, 0)

	var nodes []Node
	switch vals := reflect.ValueOf(payload); vals.Kind() {
	case reflect.Slice:
		for x := 0; x < vals.Len(); x++ {
			if node, isNodeable := vals.Index(x).Interface().(Node); isNodeable {
				nodes = append(nodes, node)
			}
		}
	case reflect.Ptr:
		if reflect.Indirect(vals).Kind() == reflect.Struct {
			if node, isNodeable := vals.Interface().(Node); isNodeable {
				nodes = append(nodes, node)
			}
		}
	}

	for _, node := range nodes {
		internalNodes = append(internalNodes, transformToInternalDataStruct(node, baseURL))
	}

	return internalNodes
}

func transformToInternalDataStruct(node Node, baseURL string) internalNode {
	var links LinkMap
	if linkableNode, isLinkable := node.(Linkable); isLinkable {
		links = TransformLinks(linkableNode.Links(), baseURL)
	}

	var meta interface{}
	if metaNode, isMetable := node.(Metable); isMetable {
		meta = metaNode.Meta()
	}

	return internalNode{
		ID:            node.ID(),
		Type:          node.Type(),
		Attributes:    node.Attributes(),
		Links:         links,
		Meta:          meta,
		Relationships: transformToInternalRelationships(node, baseURL),
	}
}

func transformResponseData(response Response, baseURL string) interface{} {
	if response.Errors.HasErrors() {
		return nil
	}
	return transformToInternalDataStruct(response.Node, baseURL)
}

func transformCollectionResponseData(response CollectionResponse, baseURL string) interface{} {
	if response.Errors.HasErrors() {
		return nil
	}
	return transformToInternalDataStructArray(response.Nodes, baseURL)
}

func transformIncluded(includedNode interface{}, node interface{}, baseURL string) (included []internalNode) {
	// included cannot exist if node does not exist: https://jsonapi.org/format/#document-top-level
	if node == nil {
		return
	}

	return transformToInternalDataStructArray(includedNode, baseURL)
}
