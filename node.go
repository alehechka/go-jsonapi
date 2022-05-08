package jsonapi

import "reflect"

// Node is the standard JSONAPI Data struct
type Node interface {
	ID() string
	Type() string
}

type Attributeable interface {
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

func transformToInternalNodeStructArray(payload interface{}, baseURL string) []internalNode {
	internalNodes := make([]internalNode, 0)

	switch vals := reflect.ValueOf(payload); vals.Kind() {
	case reflect.Slice:
		for x := 0; x < vals.Len(); x++ {
			if node, isNodeable := vals.Index(x).Interface().(Node); isNodeable {
				internalNodes = append(internalNodes, transformToInternalNodeStruct(node, baseURL))
			}
		}
	case reflect.Ptr:
		if reflect.Indirect(vals).Kind() == reflect.Struct {
			if node, isNodeable := vals.Interface().(Node); isNodeable {
				internalNodes = append(internalNodes, transformToInternalNodeStruct(node, baseURL))
			}
		}
	}

	return internalNodes
}

func transformToInternalNodeStruct(node Node, baseURL string) internalNode {
	var links LinkMap
	if linkableNode, isLinkable := node.(Linkable); isLinkable {
		links = TransformLinks(linkableNode.Links(), baseURL)
	}

	var meta interface{}
	if metaNode, isMetable := node.(Metable); isMetable {
		meta = metaNode.Meta()
	}

	var attributes interface{} = node
	if attributeNode, isAttributeable := node.(Attributeable); isAttributeable {
		attributes = attributeNode.Attributes()
	}

	return internalNode{
		ID:            node.ID(),
		Type:          node.Type(),
		Attributes:    attributes,
		Links:         links,
		Meta:          meta,
		Relationships: transformToInternalRelationships(node, baseURL),
	}
}

func transformResponseNode(response Response, baseURL string) interface{} {
	if response.Errors.HasErrors() {
		return nil
	}
	return transformToInternalNodeStruct(response.Node, baseURL)
}

func transformCollectionResponseNode(response CollectionResponse, baseURL string) interface{} {
	if response.Errors.HasErrors() {
		return nil
	}
	return transformToInternalNodeStructArray(response.Nodes, baseURL)
}

func transformIncluded(includedNode interface{}, node interface{}, baseURL string) (included []internalNode) {
	// included cannot exist if node does not exist: https://jsonapi.org/format/#document-top-level
	if node == nil {
		return
	}

	return transformToInternalNodeStructArray(includedNode, baseURL)
}
