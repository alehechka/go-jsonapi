package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HasErrors_True(t *testing.T) {
	errs := Errors{
		{
			Status: 400,
			Detail: "has an error",
		},
	}

	assert.Equal(t, true, errs.HasErrors())
}

func Test_HasErrors_False(t *testing.T) {

	assert.Equal(t, false, Errors{}.HasErrors())
}

func Test_transformErrors(t *testing.T) {
	errs := Errors{
		{
			Status: 400,
			Detail: "has an error",
			Links: Links{
				RelatedKey: Link{
					Href: "/example/error",
				},
			},
		},
	}

	transformed := transformErrors(errs, baseURL)

	assert.NotNil(t, transformed)
	assert.Equal(t, 1, len(transformed))

	err := transformed[0]
	assert.Equal(t, errs[0].Status, err.Status)
	assert.Equal(t, errs[0].Detail, err.Detail)
	assert.Equal(t, 1, len(err.Links))
	assert.NotNil(t, err.Links[RelatedKey])
	assert.Equal(t, baseURL+errs[0].Links[RelatedKey].Href, err.Links[RelatedKey])
}

func Test_transformError(t *testing.T) {
	err := Error{
		Status: 400,
		Detail: "has an error",
		Links: Links{
			RelatedKey: Link{
				Href: "/example/error",
			},
		},
	}

	transformed := transformError(err, baseURL)

	assert.NotNil(t, transformed)

	assert.Equal(t, err.Status, transformed.Status)
	assert.Equal(t, err.Detail, transformed.Detail)
	assert.Equal(t, 1, len(transformed.Links))
	assert.NotNil(t, transformed.Links[RelatedKey])
	assert.Equal(t, baseURL+err.Links[RelatedKey].Href, transformed.Links[RelatedKey])
}
