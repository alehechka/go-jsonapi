package jsonapi

import "reflect"

func transformCollectionResponseNodes(response CollectionResponse, baseURL string) (data interface{}, included []Node) {
	if response.Errors.HasErrors() {
		return nil, nil
	}
	return transformNodes(response.Nodes, baseURL)
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

	case reflect.Struct:
		appendNode(vals.Interface())

	case reflect.Ptr:
		return transformNodes(reflect.Indirect(vals).Interface(), baseURL)
	}

	return internalNodes, included
}
