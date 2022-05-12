package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_transformRelationships_Relationshipable(t *testing.T) {}
func Test_transformRelationships(t *testing.T)                  {}
func Test_transformRelationship(t *testing.T)                   {}
func Test_transformRelationship_Methods(t *testing.T)           {}
func Test_transformRelationshipData_Nodeable(t *testing.T)      {}
func Test_transformRelationshipData(t *testing.T)               {}
func Test_transformRelationNodes_PtrSlice(t *testing.T)         {}
func Test_transformRelationNodes_PtrStruct(t *testing.T)        {}
func Test_transformRelationNodes_Slice(t *testing.T)            {}
func Test_transformRelationNodes_Struct(t *testing.T)           {}

func Test_createResourceIdentifier(t *testing.T) {
	test := testStruct{
		TestID: "1234",
		Number: 123,
		IsTest: true,
	}

	resource := createResourceIdentifier(test)

	assert.NotNil(t, resource)
	assert.Equal(t, testObject.ID(), resource.ID)
	assert.Equal(t, testObject.Type(), resource.Type)
	assert.Nil(t, resource.Meta)
}

func Test_createResourceIdentifier_Meta(t *testing.T) {
	resource := createResourceIdentifier(testObject)

	assert.NotNil(t, resource)
	assert.Equal(t, testObject.ID(), resource.ID)
	assert.Equal(t, testObject.Type(), resource.Type)
	assert.NotNil(t, resource.Meta)
	assert.Equal(t, testObject.Meta(), resource.Meta)
}
