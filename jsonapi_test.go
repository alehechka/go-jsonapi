package jsonapi

import (
	"encoding/json"
	"reflect"
	"testing"
)

type SomeRelatedAgilityData struct {
	CustomerID string `json:"customerId"`
}

func (d SomeRelatedAgilityData) ID() string {
	return d.CustomerID
}
func (d SomeRelatedAgilityData) Type() string {
	return "relatedData"
}
func (d SomeRelatedAgilityData) Meta() interface{} {
	return nil
}

type AgilityDataRelationship struct {
	ID string `json:"id"`
}

func (d AgilityDataRelationship) Links() Links {
	return Links{
		"resource": {
			Href: "/path/to/resource/:id",
			Params: map[string]interface{}{
				"id": d.ID,
			},
		},
	}
}

func (d AgilityDataRelationship) Data() ([]ResourceIdentifier, bool) {
	arr := make([]ResourceIdentifier, 1)

	arr[0] = SomeRelatedAgilityData{
		CustomerID: "cust1234",
	}

	return arr, false
}

func (d AgilityDataRelationship) Meta() interface{} {
	return nil
}

type SomeAgilityData struct {
	Name     string `json:"name"`
	TranID   string `json:"tranId"`
	ShipTo   string `json:"shipTo"`
	ItemName string `json:"itemName"`
}

func (d SomeAgilityData) ID() string {
	return d.TranID
}

func (d SomeAgilityData) Type() string {
	return "agilityData"
}

func (d SomeAgilityData) Attributes() interface{} {
	return d
}

// TODO Add an example here
func (d SomeAgilityData) Links() Links {
	return nil
}

func (d SomeAgilityData) Relationships() map[string]Relationship {
	if d.TranID == "1111" {
		return map[string]Relationship{
			"relatedData": AgilityDataRelationship{
				ID: d.TranID,
			},
		}
	}

	return nil
}

func (d SomeAgilityData) Meta() interface{} {
	return nil
}

func TestCreateResponse(t *testing.T) {

	type args struct {
		data     []Data
		included []Data
		errors   []Error
		links    Links
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
				data: []Data{
					SomeAgilityData{
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
		"type": "agilityData",
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
				data: []Data{
					SomeAgilityData{
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
		"type": "agilityData",
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
				data: []Data{
					SomeAgilityData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeAgilityData{
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
			"type": "agilityData",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "agilityData",
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
				data: []Data{
					SomeAgilityData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeAgilityData{
						Name:     "Testing data 2",
						TranID:   "12346",
						ShipTo:   "Location 2",
						ItemName: "Lot-o-nails",
					},
				},
				included: []Data{
					SomeAgilityData{
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
			"type": "agilityData",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "agilityData",
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
			"type": "agilityData",
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
				data: []Data{
					SomeAgilityData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeAgilityData{
						Name:     "Testing data 2",
						TranID:   "12346",
						ShipTo:   "Location 2",
						ItemName: "Lot-o-nails",
					},
				},
				included: nil,
				links: Links{
					"self": {
						Href: "/api/rest/someAgilityData",
					},
					"other": {
						Href: "/api/rest/someOtherAgilityData",
						Meta: map[string]interface{}{
							"count": 10,
						},
						Params: map[string]interface{}{
							"random": 11,
						},
					},
					"else": {
						Href: "/api/rest/elseAgilityData/:id",
						Params: map[string]interface{}{
							"id": 2,
						},
					},
					"with-protocol": {
						Href: "http://www.dmsi.com/api/rest/with-protocol",
					},
				},
				meta: nil,
			},
			want: `{
	"data": [
		{
			"id": "12345",
			"type": "agilityData",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "agilityData",
			"attributes": {
				"name": "Testing data 2",
				"tranId": "12346",
				"shipTo": "Location 2",
				"itemName": "Lot-o-nails"
			}
		}
	],
	"links": {
		"else": "https://example.com/api/rest/elseAgilityData/2",
		"other": {
			"href": "https://example.com/api/rest/someOtherAgilityData",
			"meta": {
				"count": 10
			}
		},
		"self": "https://example.com/api/rest/someAgilityData",
		"with-protocol": "http://www.dmsi.com/api/rest/with-protocol"
	}
}`,
		},
		{
			name: "Response test with multiple data values and links",
			args: args{
				data: []Data{
					SomeAgilityData{
						Name:     "Testing data 1",
						TranID:   "12345",
						ShipTo:   "Location 1",
						ItemName: "Box-o-shingles",
					},
					SomeAgilityData{
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
			"type": "agilityData",
			"attributes": {
				"name": "Testing data 1",
				"tranId": "12345",
				"shipTo": "Location 1",
				"itemName": "Box-o-shingles"
			}
		},
		{
			"id": "12346",
			"type": "agilityData",
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
				errors: []Error{
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
			got, err := json.MarshalIndent(CreateResponse(tt.args.data, tt.args.included, tt.args.errors, tt.args.links, tt.args.meta, "https://example.com"), "", "\t")
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
