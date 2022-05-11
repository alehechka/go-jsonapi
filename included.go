package jsonapi

func transformIncluded(includedNode []Node, node interface{}, baseURL string) (included []internalNode) {
	// included cannot exist if node does not exist: https://jsonapi.org/format/#document-top-level
	if node == nil {
		return
	}

	for _, node := range includedNode {
		included = append(included, transformIncludedNode(node, baseURL))
	}

	return
}

func transformIncludedNode(node Node, baseURL string) internalNode {
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
		ID:         node.ID(),
		Type:       node.Type(),
		Attributes: attributes,
		Links:      links,
		Meta:       meta,
	}
}
