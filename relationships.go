package jsonapi

// ResourceIdentifier is the standard Resource Identifier struct
type ResourceIdentifier interface {
	ID() string
	Type() string
	Meta() interface{}
}

type internalResourceIdentifier struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Meta interface{} `json:"meta,omitempty"`
}

// Relationship is the standard JSONAPI Relationship struct
type Relationship interface {
	Links() Links
	Data() ([]ResourceIdentifier, bool) // identifiers, is relationship to many
	Meta() interface{}
}

type internalRelationship struct {
	Links LinkMap     `json:"links,omitempty"`
	Data  interface{} `json:"data"` // ResourceIdentifier | []ResourceIdentifier
	Meta  interface{} `json:"meta,omitempty"`
}

type Relationable interface {
	Relationships() map[string]Relationship
}

func transformToInternalRelationships(node Node, baseURL string) map[string]internalRelationship {
	if relationshipNode, isRelationable := node.(Relationable); isRelationable {

		relationships := relationshipNode.Relationships()

		internalRelationships := make(map[string]internalRelationship)

		for k, v := range relationships {
			internalRelationships[k] = transformToInternalRelationship(v, baseURL)
		}

		return internalRelationships
	}

	return nil
}

func transformToInternalRelationship(r Relationship, baseURL string) internalRelationship {

	return internalRelationship{
		Links: TransformLinks(r.Links(), baseURL),
		Meta:  r.Meta(),
		Data:  transformToInternalRelationshipData(r, baseURL),
	}
}

func transformToInternalRelationshipData(r Relationship, baseURL string) interface{} {
	relationshipData, isToMany := r.Data()

	if relationshipData == nil {
		return nil
	}

	if len(relationshipData) == 0 && isToMany {
		return [0]internalResourceIdentifier{}
	}

	if len(relationshipData) == 1 && !isToMany {
		return internalResourceIdentifier{
			ID:   relationshipData[0].ID(),
			Type: relationshipData[0].Type(),
			Meta: relationshipData[0].Meta(),
		}
	}

	data := make([]internalResourceIdentifier, 0)
	for _, relationship := range relationshipData {
		data = append(data, internalResourceIdentifier{
			ID:   relationship.ID(),
			Type: relationship.Type(),
			Meta: relationship.Meta(),
		})
	}

	return data
}
