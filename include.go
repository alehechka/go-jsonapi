package jsonapi

import (
	"errors"
	"net/http"
	"strings"
)

const (
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
	// ErrTooManyIncluded number of included is greater than number of available resources
	ErrTooManyIncluded error = errors.New("included query has too many resources")
	// ErrResourceNotAvailable member of included is not an available resource
	ErrResourceNotAvailable error = errors.New("resource from included query not available")
)

// VerifyResources verifies that all requested included members exist in available resources
func (included Included) VerifyResources(resources ...string) error {
	if len(included) > len(resources) {
		return ErrTooManyIncluded
	}

	resourceMap := make(map[string]bool)
	for _, resource := range resources {
		resourceMap[resource] = true
	}

	for _, include := range included {
		if ok, exists := resourceMap[include]; !ok || !exists {
			return ErrResourceNotAvailable
		}
	}

	return nil
}
