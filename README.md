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

The above functions are effectively the top-level transformation tools, however, the dynamic link creation can be made easy by supply an `*http.Request` object to the following functions instead:

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

A top-level `links` object can be provided to both `Response` and `CollectionResponse`. See [/jsonapi/links.go](/jsonapi/links.go#L35-L44) for all available variables in `Link`.

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

A top-level `errors` array can be provided to both `Response` and `CollectionResponse` in the form of an array of `Error` objects. See [/jsonapi/errors.go](/jsonapi/errors.go#L10-L19) for all available variables in `Error`.

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

<!--- Links -->

[jsonapi]: (https://jsonapi.org/)
[jsonapi-resource-object]: (https://jsonapi.org/format/#document-resource-objects)
[jsonapi-top-level]: (https://jsonapi.org/format/#document-top-level)
[gin]: (https://github.com/gin-gonic/gin)
