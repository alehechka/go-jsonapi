package jsonapi_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi/jsonapi"
	"github.com/stretchr/testify/assert"
)

func Test_GetIncluded(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?include=resource1,resource2", nil)

	included := jsonapi.GetIncluded(req)

	assert.Equal(t, 2, len(included))
	assert.Equal(t, "resource1", included[0])
	assert.Equal(t, "resource2", included[1])
}

func Test_GetIncluded_Empty(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?include=", nil)

	included := jsonapi.GetIncluded(req)

	assert.Equal(t, 0, len(included))
}

func Test_GetIncluded_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	included := jsonapi.GetIncluded(req)

	assert.Equal(t, 0, len(included))
}

func Test_HasResource_True(t *testing.T) {
	included := jsonapi.Included{"resource1", "resource2"}

	hasResource := included.HasResource("resource1")

	assert.Equal(t, true, hasResource)
}

func Test_HasResource_False(t *testing.T) {
	included := jsonapi.Included{"resource1", "resource2"}

	hasResource := included.HasResource("resource13")

	assert.Equal(t, false, hasResource)
}

func Test_VerifyResources_OK(t *testing.T) {
	included := jsonapi.Included{"resource1", "resource2"}

	err := included.VerifyResources("resource1", "resource2")

	assert.Nil(t, err)
}

func Test_VerifyResources_TooManyIncluded(t *testing.T) {
	included := jsonapi.Included{"resource1", "resource2"}

	err := included.VerifyResources("resource1")

	assert.NotNil(t, err)
	assert.Equal(t, jsonapi.ErrTooManyIncluded, err)
}

func Test_VerifyResources_ResourceNotAvailable(t *testing.T) {
	included := jsonapi.Included{"resource1", "resource2"}

	err := included.VerifyResources("resource1", "resource3")

	assert.NotNil(t, err)
	assert.Equal(t, jsonapi.ErrResourceNotAvailable, err)
}
