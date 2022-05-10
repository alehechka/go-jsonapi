package jsonapi

// Node is the standard JSONAPI Data struct
type Node interface {
	ID() string
	Type() string
}

// Attributeable method for custom Attributes structs
type Attributeable interface {
	Attributes() any
}

// Metable method for meta data objects
type Metable interface {
	Meta() any
}

type internalNode struct {
	ID            string                          `json:"id"`
	Type          string                          `json:"type"`
	Attributes    any                             `json:"attributes,omitempty"`
	Links         LinkMap                         `json:"links,omitempty"`
	Relationships map[string]internalRelationship `json:"relationships,omitempty"`
	Meta          any                             `json:"meta,omitempty"`
}

func transformResponseNode(response Response, baseURL string) (node any, included []Node) {
	if response.Errors.HasErrors() {
		return nil, nil
	}
	return transformNode(response.Node, baseURL)
}

func transformNode(node Node, baseURL string) (internalNode, []Node) {
	if node == nil {
		return internalNode{}, nil
	}

	var links LinkMap
	if linkableNode, isLinkable := node.(Linkable); isLinkable {
		links = TransformLinks(linkableNode.Links(), baseURL)
	}

	var meta any
	if metaNode, isMetable := node.(Metable); isMetable {
		meta = metaNode.Meta()
	}

	var attributes any = node
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
