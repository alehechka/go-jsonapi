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
func (d SomeRelatedData) Meta() interface{} {
	return nil
}

type DataRelationship struct {
	ID string `json:"id"`
}

func (d DataRelationship) Links() jsonapi.Links {
	return jsonapi.Links{
		"resource": {
			Href: "/path/to/resource/:id",
			Params: map[string]interface{}{
				"id": d.ID,
			},
		},
	}
}

func (d DataRelationship) Data() ([]jsonapi.ResourceIdentifier, bool) {
	arr := make([]jsonapi.ResourceIdentifier, 1)

	arr[0] = SomeRelatedData{
		CustomerID: "cust1234",
	}

	return arr, false
}

func (d DataRelationship) Meta() interface{} {
	return nil
}

type SomeData struct {
	Name     string `json:"name"`
	TranID   string `json:"tranId"`
	ShipTo   string `json:"shipTo"`
	ItemName string `json:"itemName"`
}

func (d SomeData) ID() string {
	return d.TranID
}

func (d SomeData) Type() string {
	return "Data"
}

func (d SomeData) Attributes() interface{} {
	return d
}

// TODO Add an example here
func (d SomeData) Links() jsonapi.Links {
	return nil
}

func (d SomeData) Relationships() map[string]jsonapi.Relationship {
	if d.TranID == "1111" {
		return map[string]jsonapi.Relationship{
			"relatedData": DataRelationship{
				ID: d.TranID,
			},
		}
	}

	return nil
}

func (d SomeData) Meta() interface{} {
	return nil
}

func TestCreateResponse(t *testing.T) {

	type args struct {
		data     []jsonapi.Data
		included []jsonapi.Data
		errors   []jsonapi.Error
		links    jsonapi.Links
		meta     interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Response test with a single data item",
			args: args{
				data: []jsonapi.Data{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
				},
				included: nil,
				links:    nil,
				meta:     nil,
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
				data: []jsonapi.Data{
					SomeData{
						Name:     "Testing data 1",
						TranID:   "1111",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
				},
				included: nil,
				links:    nil,
				meta:     nil,
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
					"resource": "https://example.com/path/to/resource/1111"
				},
				"data": {
					"id": "cust1234",
					"type": "relatedData"
				}
			}
		}
	}
}`,
		},
		{
			name: "Response test with multiple data values",
			args: args{
				data: []jsonapi.Data{
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
				included: nil,
				links:    nil,
				meta:     nil,
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
				data: []jsonapi.Data{
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
				included: []jsonapi.Data{
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
	"included": [
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
			name: "Response test with multiple data values and links",
			args: args{
				data: []jsonapi.Data{
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
				included: nil,
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
				data: []jsonapi.Data{
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
				included: nil,
				links:    nil,
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
			"status": 500,
			"code": 0,
			"title": ""
		}
	]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.MarshalIndent(jsonapi.CreateResponse(jsonapi.Response{
				Data:     tt.args.data,
				Included: tt.args.included,
				Errors:   tt.args.errors,
				Links:    tt.args.links,
				Meta:     tt.args.meta,
			}, "https://example.com"), "", "\t")
			if err != nil {
				t.Errorf("CreateResponse() error %v", err)
				return
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("CreateResponse() = \n%v, want \n%v", string(got), tt.want)
			}

		})
	}
}
