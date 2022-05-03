package jsonapi

import (
	"errors"
	"net/http"
	"strings"
)

// Pagination Query Parameters
const (
	// PageOffset represents the offset from previous pagination. Use in conjunction with PageLimit.
	PageOffset string = "page[offset]"
	// PageLimit represents the size limit of page offset. Use in conjunction with PageOffset.
	PageLimit string = "page[limit]"
	// PageNumber represents the page number. Use in conjunction with PageSize.
	PageNumber string = "page[number]"
	// PageSize represents the size limit of page. Use in conjunction with PageNumber.
	PageSize string = "page[size]"
	// PageCursor represents the cursor of the page. Use by itself.
	PageCursor string = "page[cursor]"
	// PageBefore used with cursor-based pagination. Use in conjunction with PageSize and PageAfter
	PageBefore string = "page[before]"
	// PageAfter used with cursor-based pagination. Use in conjunction with PageSize and PageBefore
	PageAfter string = "page[after]"
	// Include used to request extra resources to include
	Include string = "include"
)

// Included string array representing the included query parameter resources
type Included []string

// GetIncluded extracts included object names from request query parameters
func GetIncluded(request *http.Request) (included Included) {
	includedQuery := request.URL.Query().Get(Include)

	if len(includedQuery) == 0 {
		return
	}

	included = strings.Split(request.URL.Query().Get(Include), ",")
	return
}

// HasResource will check the provided included array for the requested resource name
func (included Included) HasResource(resource string) bool {
	for _, include := range included {
		if include == resource {
			return true
		}
	}

	return false
}

var (
	TooManyIncludedError error = errors.New("included query has too many resources")
	ResourceNotAvailable error = errors.New("resource from included query not available")
)

// VerifyResources verifies that all requested included members exist in available resources
func (included Included) VerifyResources(resources ...string) error {
	if len(included) > len(resources) {
		return TooManyIncludedError
	}

	resourceMap := make(map[string]bool)
	for _, resource := range resources {
		resourceMap[resource] = true
	}

	for _, include := range included {
		if ok, exists := resourceMap[include]; !ok || !exists {
			return ResourceNotAvailable
		}
	}

	return nil
}
