package jsonapi

// Node is the standard JSONAPI Data struct
type Node interface {
	ID() string
	Type() string
}

// Attributeable method for custom Attributes structs
type Attributeable interface {
	Attributes() interface{}
}

// Metable method for meta data objects
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

func transformResponseNode(response Response, baseURL string) (node interface{}, included []Node) {
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
