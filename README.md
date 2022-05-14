# go-jsonapi

This Go module provides a useful API to create [JSON:API][jsonapi] HTTP servers. The primary usage of this library is to facilitate transformation from a flattened Go structs into the standardized JSON:API [resource object][jsonapi-resource-object].

Additionally, there are optional methods that can be implemented with structs to add further standardized JSON:API structures such as links, relationships, included data, and metadata.

## Installation

```bash
go get github.com/alehechka/go-jsonapi
```

## Usage

### Defining a JSON:API struct

The primary resource object in JSON:API is of the following type:

```json
{
	"data": {
		"id": "1234",
		"type": "people",
		"attributes": {
			"firstName": "John",
			"lastName": "Doe",
			"age": 30
		}
	}
}
```

- The `attributes` object will be generated from the struct itself.
- The `id` field will be populated by the `ID()` interface method.
- The `type` field will be populated by the `Type()` interface method.

```go
type Person struct {
    // It is recommended to omit the primary ID from json marshalling, but not required
    PersonID    string  `json:"-"`
    FirstName   string  `json:"firstName"`
    LastName    string  `json:"lastName"`
    Age         int     `json:"age"`
}

func (person Person) ID() string {
    return person.PersonID
}

func (person Person) Type() string {
    return "people"
}
```

### Prepare for JSON marshalling

To prepare the struct for json marshalling it is required to use the provided `TransformResponse` or `TransformCollectionResponse` functions:

```go
response := jsonapi.TransformResponse(jsonapi.Response{
    Node: Person{},
    "http://example.com",
})

response := jsonapi.TransformCollectionResponse(jsonapi.CollectionResponse{
    Nodes: []Person{},
    "http://example.com",
})
```

The second parameter to these functions is for `baseURL`, this is used to dynamically populate relative URLs in `links` objects. More on this [here]().

### Recommended Usage

The above functions are effectively the top-level transformation tools, however, the dynamic link creation can be made easy by supplying an `*http.Request` object to the following functions instead:

```go
req := httptest.NewRequest("GET", "http://example.com/example", nil)

response := jsonapi.CreateResponse(req)(jsonapi.Response{
    Node: Person{},
})

response := jsonapi.CreateCollectionResponse(req)(jsonapi.CollectionResponse{
    Nodes: []Person{},
})
```

These versions will automatically extract the baseURL from the request and supply it to the respective `Transform` functions outlined above. This allows all generated links to display the same scheme and hostname as the server domain that the request was originally made to.

Additionally, using the `Create` functions will automatically generate a `self` link at the top-level object for every response.

### Extending the top-level resource

The JSON:API spec also allows for `links`, `errors`, and `meta` objects at the top-level of the document. Both `jsonapi.Response` and `jsonapi.CollectionResponse` have values available for these.

#### Links

A top-level `links` object can be provided to both `Response` and `CollectionResponse`. See [Link](#link) below for further details.

```go
res := jsonapi.Response{
    Links: jsonapi.Links{
        jsonapi.NextKey: jsonapi.Link{
            Href: "/path/to/next/resource",
        },
    },
}
```

> When using either `CreateResponse` or `CreateCollectionResponse` the `self` link will be automatically generated and always override an existing `self` link.

#### Meta

A top-level `meta` object can be provided to both `Response` and `CollectionResponse` in the form of any interface or key-value map.

```go
res := jsonapi.Response{
    Meta: jsonapi.Meta{
        "page": jsonapi.Meta{
            "size": 10,
            "number": 2,
        },
    },
}
```

> The `Meta` struct is simply an alias for `map[string]interface{}`

#### Errors

A top-level `errors` array can be provided to both `Response` and `CollectionResponse` in the form of an array of `Error` objects. See [Error](#error) below for further detail.

```go
res := jsonapi.Response{
    Errors: jsonapi.Errors{
        {
            Status: http.StatusBadRequest,
            Title: "Error Occurred",
            Detail: "Failed to retrieve resource",
        },
    },
}
```

> It is important to note that if at least 1 error is present in this array than the top-level `data` object/array and `included` array will not be available as per the JSON:API spec for [Top Level][jsonapi-top-level]

### Extending `Node` interface

By default, to be considered a JSON:API resource, a struct must include the `ID()` and `Type()` methods.

However, this functionality can be extended further with other methods as follows:

#### `Links()`

The `Links()` method allows an individual resource to generate the `links` object for itself using data from the object. See [Link](#link) below for further details.

```go
func (person Person) Links() jsonapi.Links {
    return jsonapi.Links{
        jsonapi.SelfKey: jsonapi.Link{
            Href: "/people/:id",
            Params: jsonapi.Params{
                "id": person.ID(),
            }
        },
    }
}
```

The above scenario makes use of the `Params` field which will not be included in the resulting json, but will use the key-value pairs to substitute the values into the `href` based on keys that it finds. (Ex. `:id` in the href will be substituted with the value of `person.ID()`)

#### `Relationships()`

[Relationships][jsonapi-relationships] are a key object within a resource to provide linkage and information about related resources. To facilitate the mapping, the `Relationships()` method gives access to the parent struct and allows definition of the `relationships` map as follows:

```go
type Company struct {
    CompanyID string `json:"-"`
    Name string `json:"name"`
    Address string `json:"address"`
    Employees []Person `json:"-"` // recommended to omit children resources
    Owner Person `json:"-"`
}

func (company Company) Relationships() map[string]interface{
    return map[string]interface{}{
        "employees": company.Employees,
        "owner": company.Owner,
    }
}
```

> In the above example it is crucial that the children relationship objects adhere to the JSON:API methods, i.e. initialize their own `ID()` and `Type()` methods.

#### `RelationshipLinks(parentID string)`

Typically in the `relationships` object, there will be included `links` object with links to the [related resources][jsonapi-related-links]. This can be facilitated by included the `RelationshipLinks(parentID string`) on children structs. The `parentID` parameter will automatically be supplied when generated as part of a relationship by the parent struct, it is recommended to use this in generating path params for the href variable.

```go
func (person Person) RelationshipLinks(companyID string) jsonapi.Links {
    return jsonapi.Links{
        jsonapi.SelfKey: jsonapi.Link{
            Href: "/companies/:companyID/relationships/employees",
            Params: jsonapi.Params{
                "companyID": companyID,
            },
        },
        jsonapi.RelatedKey: jsonapi.Link{
            Href: "/companies/:companyID/employees",
            Params: jsonapi.Params{
                "companyID": companyID,
            },
        },
    }
}
```

If the relationship will point to an array of resources, it is recommended to instead create a unique type for that array of structs as follows:

```go
type People []Person

func (people People) RelationshipLinks(companyID string) jsonapi.Links {
    return jsonapi.Links{
        jsonapi.SelfKey: jsonapi.Link{
            Href: "/companies/:companyID/relationships/employees",
            Params: jsonapi.Params{
                "companyID": companyID,
            },
        },
    }
}
```

#### `Meta()`

The `Meta()` method is simply a means to generate a `meta` object for an individual resource by using the object as an input.

```go
func (person Person) Meta() interface{} {
    return jsonapi.Meta{
        "fullName": fmt.Sprintf("%s %s", person.FirstName, person.LastName),
    }
}
```

### Structs Explained

#### `Link`

[/jsonapi/links.go](/jsonapi/links.go#L35-L44)

#### `Error`

[/jsonapi/errors.go](/jsonapi/errors.go#L10-L19)

<!--- Links -->

[jsonapi]: (https://jsonapi.org/)
[jsonapi-resource-object]: (https://jsonapi.org/format/#document-resource-objects)
[jsonapi-top-level]: (https://jsonapi.org/format/#document-top-level)
[jsonapi-relationships]: (https://jsonapi.org/format/#document-resource-object-relationships)
[jsonapi-related-links]: (https://jsonapi.org/format/#document-resource-object-related-resource-links)
[gin]: (https://github.com/gin-gonic/gin)
