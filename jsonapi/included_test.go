package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_transformIncluded_Nil(t *testing.T) {
	included := transformIncluded([]Node{}, nil, baseURL)

	assert.Nil(t, included)
}

func Test_transformIncluded(t *testing.T) {
	included := transformIncluded([]Node{testObject}, []Node{}, baseURL)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))

	include := included[0]

	assert.Equal(t, testObject.ID(), include.ID)
	assert.Equal(t, testObject.Type(), include.Type)
	assert.Equal(t, testObject, include.Attributes)
	assert.Equal(t, testObject.Meta(), include.Meta)

	assert.NotNil(t, include.Links)
	assert.NotNil(t, include.Links[SelfKey])
	assert.Equal(t, baseURL+testObject.Links()[SelfKey].Href, include.Links[SelfKey])
}

func Test_transformIncludedNode(t *testing.T) {
	include := transformIncludedNode(testObject, baseURL)

	assert.NotNil(t, include)
	assert.Equal(t, testObject.ID(), include.ID)
	assert.Equal(t, testObject.Type(), include.Type)
	assert.Equal(t, testObject, include.Attributes)
	assert.Equal(t, testObject.Meta(), include.Meta)

	assert.NotNil(t, include.Links)
	assert.NotNil(t, include.Links[SelfKey])
	assert.Equal(t, baseURL+testObject.Links()[SelfKey].Href, include.Links[SelfKey])
}
