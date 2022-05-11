package jsonapi_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/alehechka/go-jsonapi"
)

type SomeRelatedData struct {
	CustomerID string `json:"customerId"`
}

func (d SomeRelatedData) ID() string {
	return d.CustomerID
}
func (d SomeRelatedData) Type() string {
	return "relatedData"
}

type DataRelationship struct {
	UUID string `json:"UUID"`
}

func (d DataRelationship) ID() string {
	return d.UUID
}

func (d DataRelationship) Type() string {
	return "dataRelationship"
}

func (d DataRelationship) RelationshipLinks(parentID string) jsonapi.Links {
	return jsonapi.Links{
		"resource": {
			Href: "/path/to/resource/:id/data",
			Params: map[string]interface{}{
				"id": parentID,
			},
		},
	}
}

type SomeData struct {
	Name             string           `json:"name"`
	TranID           string           `json:"tranId"`
	ShipTo           string           `json:"shipTo"`
	ItemName         string           `json:"itemName"`
	DataRelationship DataRelationship `json:"-"`
}

func (d SomeData) ID() string {
	return d.TranID
}

func (d SomeData) Type() string {
	return "Data"
}

// TODO Add an example here
func (d SomeData) Links() jsonapi.Links {
	return nil
}

func (d SomeData) Relationships() map[string]interface{} {
	if d.TranID == "1111" {
		return map[string]interface{}{
			"relatedData": d.DataRelationship,
		}
	}

	return nil
}

func TestTransformResponse(t *testing.T) {

	type args struct {
		node   jsonapi.Node
		errors []jsonapi.Error
		links  jsonapi.Links
		meta   interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Response test with a single data item",
			args: args{
				node: SomeData{
					Name:     "Testing data 1",
					TranID:   "12345",
					ShipTo:   "Location 1",
					ItemName: "Box-o-shingles",
				},
				links: nil,
				meta:  nil,
			},
			want: `{
	"data": {
		"id": "12345",
		"type": "Data",
		"attributes": {
			"name": "Testing data 1",
			"tranId": "12345",
			"shipTo": "Location 1",
			"itemName": "Box-o-shingles"
		}
	}
}`,
		},
		{
			name: "Response test with a single data item relationship",
			args: args{
				node: SomeData{
					Name:     "Testing data 1",
					TranID:   "1111",
					ShipTo:   "Location 1",
					ItemName: "Box-o-shingles",
					DataRelationship: DataRelationship{
						UUID: "cust1234",
					},
				},
				links: nil,
				meta:  nil,
			},
			want: `{
	"data": {
		"id": "1111",
		"type": "Data",
		"attributes": {
			"name": "Testing data 1",
			"tranId": "1111",
			"shipTo": "Location 1",
			"itemName": "Box-o-shingles"
		},
		"relationships": {
			"relatedData": {
				"links": {
					"resource": "https://example.com/path/to/resource/1111/data"
				},
				"data": {
					"id": "cust1234",
					"type": "dataRelationship"
				}
			}
		}
	},
	"included": [
		{
			"id": "cust1234",
			"type": "dataRelationship",
			"attributes": {
				"UUID": "cust1234"
			}
		}
	]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.MarshalIndent(jsonapi.TransformResponse(jsonapi.Response{
				Node:   tt.args.node,
				Errors: tt.args.errors,
				Links:  tt.args.links,
				Meta:   tt.args.meta,
			}, "https://example.com"), "", "\t")
			if err != nil {
				t.Errorf("TransformResponse() error %v", err)
				return
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("TransformResponse() = \n%v, want \n%v", string(got), tt.want)
			}

		})
	}
}

func TestTransformCollectionResponse(t *testing.T) {

	type args struct {
		node   []jsonapi.Node
		errors []jsonapi.Error
		links  jsonapi.Links
		meta   interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Response test with a single data item",
			args: args{
				node: []jsonapi.Node{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
				},
				links: nil,
				meta:  nil,
			},
			want: `{
	"data": [
		{
			"id": "12345",
			"type": "Data",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		}
	]
}`,
		},
		{
			name: "Response test with a single data item relationship",
			args: args{
				node: []jsonapi.Node{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "1111",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
						DataRelationship: DataRelationship{
							UUID: "cust1234",
						},
					},
				},
				links: nil,
				meta:  nil,
			},
			want: `{
	"data": [
		{
			"id": "1111",
			"type": "Data",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "1111",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			},
			"relationships": {
				"relatedData": {
					"links": {
						"resource": "https://example.com/path/to/resource/1111/data"
					},
					"data": {
						"id": "cust1234",
						"type": "dataRelationship"
					}
				}
			}
		}
	],
	"included": [
		{
			"id": "cust1234",
			"type": "dataRelationship",
			"attributes": {
				"UUID": "cust1234"
			}
		}
	]
}`,
		},
		{
			name: "Response test with multiple data values",
			args: args{
				node: []jsonapi.Node{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeData{
						Name:     "Testing data 2",
						TranID:   "12346",
						ShipTo:   "Location 2",
						ItemName: "Lot-o-nails",
					},
				},
				links: nil,
				meta:  nil,
			},
			want: `{
	"data": [
		{
			"id": "12345",
			"type": "Data",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "Data",
			"attributes": {
				"name": "Testing data 2",
				"tranId": "12346",
				"shipTo": "Location 2",
				"itemName": "Lot-o-nails"
			}
		}
	]
}`,
		},
		{
			name: "Response test with multiple data values and included",
			args: args{
				node: []jsonapi.Node{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeData{
						Name:     "Testing data 2",
						TranID:   "12346",
						ShipTo:   "Location 2",
						ItemName: "Lot-o-nails",
					},
				},
				links: nil,
				meta:  nil,
			},
			want: `{
	"data": [
		{
			"id": "12345",
			"type": "Data",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "Data",
			"attributes": {
				"name": "Testing data 2",
				"tranId": "12346",
				"shipTo": "Location 2",
				"itemName": "Lot-o-nails"
			}
		}
	]
}`,
		},
		{
			name: "Response test with multiple data values and links",
			args: args{
				node: []jsonapi.Node{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeData{
						Name:     "Testing data 2",
						TranID:   "12346",
						ShipTo:   "Location 2",
						ItemName: "Lot-o-nails",
					},
				},
				links: jsonapi.Links{
					"self": {
						Href: "/api/rest/someData",
					},
					"other": {
						Href: "/api/rest/someOtherData",
						Meta: map[string]interface{}{
							"count": 10,
						},
						Params: map[string]interface{}{
							"random": 11,
						},
					},
					"else": {
						Href: "/api/rest/elseData/:id",
						Params: map[string]interface{}{
							"id": 2,
						},
					},
					"with-protocol": {
						Href: "http://www.example.com/api/rest/with-protocol",
					},
				},
				meta: nil,
			},
			want: `{
	"data": [
		{
			"id": "12345",
			"type": "Data",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "Data",
			"attributes": {
				"name": "Testing data 2",
				"tranId": "12346",
				"shipTo": "Location 2",
				"itemName": "Lot-o-nails"
			}
		}
	],
	"links": {
		"else": "https://example.com/api/rest/elseData/2",
		"other": {
			"href": "https://example.com/api/rest/someOtherData",
			"meta": {
				"count": 10
			}
		},
		"self": "https://example.com/api/rest/someData",
		"with-protocol": "http://www.example.com/api/rest/with-protocol"
	}
}`,
		},
		{
			name: "Response test with multiple data values and links",
			args: args{
				node: []jsonapi.Node{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeData{
						Name:     "Testing data 2",
						TranID:   "12346",
						ShipTo:   "Location 2",
						ItemName: "Lot-o-nails",
					},
				},
				links: nil,
				meta: struct {
					Test string `json:"test"`
				}{
					Test: "test",
				},
			},
			want: `{
	"data": [
		{
			"id": "12345",
			"type": "Data",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "Data",
			"attributes": {
				"name": "Testing data 2",
				"tranId": "12346",
				"shipTo": "Location 2",
				"itemName": "Lot-o-nails"
			}
		}
	],
	"meta": {
		"test": "test"
	}
}`,
		},
		{
			name: "Error Response test with a single error",
			args: args{
				errors: []jsonapi.Error{
					{
						ID:     "12345",
						Status: 500,
					},
				},
				links: nil,
				meta:  nil,
			},
			want: `{
	"errors": [
		{
			"id": "12345",
			"status": 500
		}
	]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.MarshalIndent(jsonapi.TransformCollectionResponse(jsonapi.CollectionResponse{
				Nodes:  tt.args.node,
				Errors: tt.args.errors,
				Links:  tt.args.links,
				Meta:   tt.args.meta,
			}, "https://example.com"), "", "\t")
			if err != nil {
				t.Errorf("TransformCollectionResponse() error %v", err)
				return
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("TransformCollectionResponse() = \n%v, want \n%v", string(got), tt.want)
			}

		})
	}
}
