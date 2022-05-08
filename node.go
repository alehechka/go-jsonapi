package jsonapi

import (
	"reflect"
)

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

func transformNodes(payload interface{}, baseURL string) ([]internalNode, []Node) {
	internalNodes := make([]internalNode, 0)
	included := make([]Node, 0)

	appendNode := func(obj interface{}) {
		if node, isNodeable := obj.(Node); isNodeable {
			internalNode, inc := transformNode(node, baseURL)
			internalNodes = append(internalNodes, internalNode)
			included = append(included, inc...)
		}
	}

	switch vals := reflect.ValueOf(payload); vals.Kind() {
	case reflect.Slice:
		for x := 0; x < vals.Len(); x++ {
			appendNode(vals.Index(x).Interface())
		}
	case reflect.Ptr:
		switch ptrVals := reflect.Indirect(vals); ptrVals.Kind() {
		case reflect.Struct:
			appendNode(ptrVals.Interface())
		case reflect.Slice:
			for x := 0; x < ptrVals.Len(); x++ {
				appendNode(ptrVals.Index(x).Interface())
			}
		}

	case reflect.Struct:
		appendNode(vals.Interface())
	}

	return internalNodes, included
}

func transformNode(node Node, baseURL string) (internalNode, []Node) {
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

	relationships, included := transformRelationships(node, baseURL)

	return internalNode{
		ID:            node.ID(),
		Type:          node.Type(),
		Attributes:    attributes,
		Links:         links,
		Meta:          meta,
		Relationships: relationships,
	}, included
}

func transformResponseNode(response Response, baseURL string) (interface{}, []Node) {
	if response.Errors.HasErrors() {
		return nil, nil
	}
	return transformNode(response.Node, baseURL)
}

func transformCollectionResponseNodes(response CollectionResponse, baseURL string) (interface{}, []Node) {
	if response.Errors.HasErrors() {
		return nil, nil
	}
	return transformNodes(response.Nodes, baseURL)
}
