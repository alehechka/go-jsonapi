package jsonapi

import "reflect"

type internalResourceIdentifier struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Meta interface{} `json:"meta,omitempty"`
}

type internalRelationship struct {
	Links LinkMap     `json:"links,omitempty"`
	Data  interface{} `json:"data"` // ResourceIdentifier | []ResourceIdentifier
	Meta  interface{} `json:"meta,omitempty"`
}

type Relationship interface {
	Data() interface{} // Node | []Node
}

type Relationable interface {
	Relationships() map[string]interface{} // Node | []Node
}

func transformRelationships(node Node, baseURL string) (map[string]internalRelationship, []Node) {
	if relationshipNode, isRelationable := node.(Relationable); isRelationable {

		relationships := relationshipNode.Relationships()

		internalRelationships := make(map[string]internalRelationship)
		included := make([]Node, 0)

		for k, v := range relationships {
			relationship, inc := transformRelationship(v, node.ID(), baseURL)
			internalRelationships[k] = relationship
			included = append(included, inc...)
		}

		return internalRelationships, included
	}

	return nil, nil
}

func transformRelationship(relationship interface{}, parentID string, baseURL string) (internalRelationship, []Node) {

	var links LinkMap
	if linkableNode, isLinkable := relationship.(RelationshipLinkable); isLinkable {
		links = TransformLinks(linkableNode.RelationshipLinks(parentID), baseURL)
	}

	var meta interface{}
	if metaNode, isMetable := relationship.(Metable); isMetable {
		meta = metaNode.Meta()
	}

	data, included := transformRelationshipData(relationship)

	return internalRelationship{
		Links: links,
		Meta:  meta,
		Data:  data,
	}, included
}

func transformRelationshipData(r interface{}) (interface{}, []Node) {
	if relationship, isRelationship := r.(Relationship); isRelationship {
		return transformRelationNodes(relationship.Data())
	}
	return transformRelationNodes(r)
}

func transformRelationNodes(r interface{}) (interface{}, []Node) {
	switch vals := reflect.ValueOf(r); vals.Kind() {
	case reflect.Slice:
		internalResources := make([]internalResourceIdentifier, 0)
		included := make([]Node, 0)
		for x := 0; x < vals.Len(); x++ {
			if node, isNodeable := vals.Index(x).Interface().(Node); isNodeable {
				internalResources = append(internalResources, createResourceIdentifier(node))
				included = append(included, node)
			}
		}
		return internalResources, included
	case reflect.Ptr:
		if reflect.Indirect(vals).Kind() == reflect.Struct {
			if node, isNodeable := vals.Interface().(Node); isNodeable {
				return createResourceIdentifier(node), []Node{node}
			}
		}

	case reflect.Struct:
		if node, isNodeable := vals.Interface().(Node); isNodeable {
			return createResourceIdentifier(node), []Node{node}
		}
	}
	return nil, nil
}

func createResourceIdentifier(resource Node) internalResourceIdentifier {
	var meta interface{}
	if metaNode, isMetable := resource.(Metable); isMetable {
		meta = metaNode.Meta()
	}

	return internalResourceIdentifier{
		ID:   resource.ID(),
		Type: resource.Type(),
		Meta: meta,
	}
}
