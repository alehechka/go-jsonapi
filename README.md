# go-jsonapi

This go module provides a useful SDK to create [JSON:API][jsonapi] HTTP interfaces. The primary usage of this library is to facilitate transformation from a flattened Go struct into the standardized JSON:API [resource object][jsonapi-resource-object].

With this, there are also struct interface options to facilitate relationships, metadata, and links, further information [below].

To help with standard error messaging, there also some available middleware functions available for use out-of-the-box if using [gin][gin].

## Installation

Currently, this module is built with Go 1.18, but has yet to implement anything with Generics so it should be fairly backwards compatible.

Install with:

```bash
go get github.com/alehechka/go-jsonapi
```

> gin will be listed as an indirect dependency, but the core of the library only uses core Go libraries.

## Usage

This library does not create any unnecessary marshalling assumptions, it will simply take in `Response` or `CollectionResponse` objects and transform them into JSON:API structs that keep the same provided object as the `Attributes` value. This allows an overlaying marshalling implementation to take over after transformation.

<details><summary>Example Implementation</summary>

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/alehechka/go-jsonapi"
)

func main() {
    engine := gin.Default()

    engine.GET("/record", getRecords)
    engine.GET("/record/:id", getRecord)

    engine.Run()
}

type Record struct {
    RecordID string `json:"-"`
    // variables
}

func (record Record) ID() string {
    return record.RecordID
}

func (record Record) Type() string {
    return "records"
}

func getRecords(ctx *gin.Context) {
    records := getRecordsFromDatabase()

    ctx.JSON(200, jsonapi.CreateCollectionResponse(ctx.Request)(jsonapi.CollectionResponse{
        Nodes: records,
    }))
}

func getRecord(ctx *gin.Context) {
    record := getRecordFromDatabase(ctx.Param("id"))

    ctx.JSON(200, jsonapi.CreateResponse(ctx.Request)(jsonapi.Response{
        Node: record,
    }))
}
```

</details>

### Interface Methods

This library is built with struct interfaces at its core, so any object that you'd like to transform into JSON:API must implement both the `ID() string` and `Type() string` methods. However, there are more optional methods that can be implemented to further extend the JSON:API transformed object.

<details><summary><code>ID() string</code></summary>

The `ID()` interface is always required and is used to select a data member from the struct to use as the `id` field in the JSON:API object.

> It is recommended to omit the selected variable with the following tag:
>
> ```go
> RecordID string `json:"-"`
> ```
>
> Although this is not required.

```go
func (record Record) ID() string {
    return record.RecordID
}
```

</details>

<details><summary><code>Links() jsonapi.Links</code></summary>

```go
func (record Record) Links() jsonapi.Links {
    return jsonapi.Links{
        jsonapi.SelfKey: {
            Href: "/records/:id",
            Params: jsonapi.Params{
                "id": record.ID(),
            },
            Queries: jsonapi.Queries{
                "page[size]": 10,
            },
            // Meta: jsonapi.Meta{
            //     "page": 10,
            // },
        }
    }
}
```

<details>
<summary>Resulting JSON</summary>
A `Record` object with `RecordID=1234` would have a resulting `links` object as follows:

```json
{
	"links": {
		"self": "http://example.com/records/1234?page[size]=10"
	}
}
```

If the commented `Meta` option were used the resulting `self` link would be an object as follows:

```json
{
	"links": {
		"self": {
			"href": "http://example.com/records/1234?page[size]=10",
			"meta": {
				"page": 10
			}
		}
	}
}
```

</details>

The `Href` variable is recommended to be provided a relative path and will be populated with a hostname from the provided `*http.Request` object. When no `Meta` variable is provided, the resulting `Link` will be a string of the generated URL.

The `Href` variable can also be created with path params in the prefixed with a colon `:` and will be substituted in with a provided matching key in the `Params` variable.

The `Params` variable represents a `map[string]any` and all matching keys will be substituted into the `Href`. `Params` will always be omitted when marshalling json.

The `Queries` variable represents a `map[string]any` and will generate query parameters to append to the `Href`. `Queries` will always be omitted when marshalling json.

The `Meta` variable represents a `map[string]any` and when provided will generate that `Link` as an object with the `meta` object included.

</details>

<!--- Links -->

[jsonapi]: (https://jsonapi.org/)
[jsonapi-resource-object]: (https://jsonapi.org/format/#document-resource-objects)
[gin]: (https://github.com/gin-gonic/gin)
