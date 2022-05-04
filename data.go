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

func transformToInternalDataStructArray(data []Data, baseURL string) []internalData {
	internalData := make([]internalData, 0)

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
		Links:         TransformLinks(d.Links(), baseURL),
		Meta:          d.Meta(),
		Relationships: transformToInternalRelationships(d, baseURL),
	}
}

func transformData(r interface{}, baseURL string) interface{} {

	switch response := r.(type) {
	case Response:
		if response.Errors.HasErrors() {
			return nil
		}
		return transformToInternalDataStruct(response.Data, baseURL)
	case CollectionResponse:
		if response.Errors.HasErrors() {
			return nil
		}
		return transformToInternalDataStructArray(response.Data, baseURL)
	}

	return nil
}

func transformIncluded(includedData []Data, data interface{}, baseURL string) (included []internalData) {
	// included cannot exist if data does not exist: https://jsonapi.org/format/#document-top-level
	if data == nil {
		return
	}

	return transformToInternalDataStructArray(includedData, baseURL)
}
