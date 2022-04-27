package jsonapi

// Data is the standard JSONAPI Data struct
type Data interface {
	ID() string
	Type() string
	Attributes() interface{}
	Links() Links
	Relationships() map[string]Relationship
	Meta() interface{}
}

type internalData struct {
	ID            string                          `json:"id"`
	Datatype      string                          `json:"type"`
	Attributes    interface{}                     `json:"attributes,omitempty"`
	Links         LinkMap                         `json:"links,omitempty"`
	Relationships map[string]internalRelationship `json:"relationships,omitempty"`
	Meta          interface{}                     `json:"meta,omitempty"`
}

func transformToInternalDataStructs(data []Data, baseURL string) interface{} {
	if len(data) == 0 {
		return nil
	}

	if len(data) == 1 {
		return transformToInternalDataStruct(data[0], baseURL)
	}

	return transformToInternalDataStructArray(data, baseURL)

}

func transformToInternalDataStructArray(data []Data, baseURL string) []internalData {
	var internalData []internalData

	for _, d := range data {
		internalData = append(internalData, transformToInternalDataStruct(d, baseURL))
	}

	return internalData
}

func transformToInternalDataStruct(d Data, baseURL string) internalData {

	return internalData{
		ID:            d.ID(),
		Datatype:      d.Type(),
		Attributes:    d.Attributes(),
		Links:         transformLinks(d.Links(), baseURL),
		Meta:          d.Meta(),
		Relationships: transformToInternalRelationships(d, baseURL),
	}
}
