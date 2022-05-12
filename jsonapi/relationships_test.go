package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_transformRelationships_Relationshipable(t *testing.T) {
	relationships, included := transformRelationships(testObject, baseURL)

	assert.NotNil(t, relationships)

	assert.NotNil(t, relationships["tests"])
	tests := relationships["tests"]
	assert.Nil(t, tests.Meta)
	assert.Nil(t, tests.Links)
	assert.NotNil(t, tests.Data)
	testsIdentifiers, ok := tests.Data.([]internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, testObject.TestData[0].ID(), testsIdentifiers[0].ID)
	assert.Equal(t, testObject.TestData[0].Type(), testsIdentifiers[0].Type)
	assert.Nil(t, testsIdentifiers[0].Meta)

	assert.NotNil(t, relationships["test"])
	test := relationships["test"]
	testIdentifiers, ok := test.Data.(internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, testObject.SingleTest.ID(), testIdentifiers.ID)
	assert.Equal(t, testObject.SingleTest.Type(), testIdentifiers.Type)
	assert.Nil(t, testIdentifiers.Meta)

	assert.NotNil(t, included)
	assert.Equal(t, 2, len(included))
	assert.Equal(t, testObject.TestData[0], included[0])
	assert.Equal(t, testObject.SingleTest, included[1])
}

func Test_transformRelationships_Nil(t *testing.T) {
	test := testStruct{
		TestID: "1234",
		Number: 123,
		IsTest: true,
	}

	relationships, included := transformRelationships(test, baseURL)

	assert.Nil(t, relationships)
	assert.Nil(t, included)
}

func Test_transformRelationship(t *testing.T) {
	test := testStruct{
		TestID: "1234",
		Number: 123,
		IsTest: true,
	}

	relationship, included := transformRelationship(test, "4321", baseURL)

	assert.NotNil(t, relationship)
	assert.Nil(t, relationship.Links)
	assert.Nil(t, relationship.Meta)

	identifier, ok := relationship.Data.(internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, test.ID(), identifier.ID)
	assert.Equal(t, test.Type(), identifier.Type)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, test, included[0])
}

func Test_transformRelationship_Methods(t *testing.T) {

	relationship, included := transformRelationship(testObject, "4321", baseURL)

	assert.NotNil(t, relationship)

	identifier, ok := relationship.Data.(internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, testObject.ID(), identifier.ID)
	assert.Equal(t, testObject.Type(), identifier.Type)

	assert.NotNil(t, relationship.Meta)
	assert.Equal(t, testObject.Meta(), relationship.Meta)

	assert.NotNil(t, relationship.Links)
	assert.NotNil(t, relationship.Links[SelfKey])
	assert.Equal(t, baseURL+"/path/to/resource/4321/child", relationship.Links[SelfKey])

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, testObject, included[0])
}

func Test_transformRelationshipData_Nodeable(t *testing.T) {
	resource, included := transformRelationshipData(testObject)

	assert.NotNil(t, resource)
	identifier, ok := resource.(internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, testObject.ID(), identifier.ID)
	assert.Equal(t, testObject.Type(), identifier.Type)
	assert.NotNil(t, identifier.Meta)
	assert.Equal(t, testObject.Meta(), identifier.Meta)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, testObject, included[0])
}

func Test_transformRelationshipData(t *testing.T) {
	test := testStruct{
		TestID: "1234",
		Number: 123,
		IsTest: true,
	}

	resource, included := transformRelationshipData(test)

	assert.NotNil(t, resource)
	identifier, ok := resource.(internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, test.ID(), identifier.ID)
	assert.Equal(t, test.Type(), identifier.Type)
	assert.Nil(t, identifier.Meta)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, test, included[0])
}

func Test_transformRelationNodes_PtrSlice(t *testing.T) {
	resource, included := transformRelationNodes(&testObjects)

	assert.NotNil(t, resource)
	identifiers, ok := resource.([]internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, 1, len(identifiers))
	assert.Equal(t, testObject.ID(), identifiers[0].ID)
	assert.Equal(t, testObject.Type(), identifiers[0].Type)
	assert.NotNil(t, identifiers[0].Meta)
	assert.Equal(t, testObject.Meta(), identifiers[0].Meta)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, testObject, included[0])
}

func Test_transformRelationNodes_PtrStruct(t *testing.T) {
	resource, included := transformRelationNodes(&testObject)

	assert.NotNil(t, resource)
	identifier, ok := resource.(internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, testObject.ID(), identifier.ID)
	assert.Equal(t, testObject.Type(), identifier.Type)
	assert.NotNil(t, identifier.Meta)
	assert.Equal(t, testObject.Meta(), identifier.Meta)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, testObject, included[0])
}

func Test_transformRelationNodes_Slice(t *testing.T) {
	resource, included := transformRelationNodes(testObjects)

	assert.NotNil(t, resource)
	identifiers, ok := resource.([]internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, 1, len(identifiers))
	assert.Equal(t, testObject.ID(), identifiers[0].ID)
	assert.Equal(t, testObject.Type(), identifiers[0].Type)
	assert.NotNil(t, identifiers[0].Meta)
	assert.Equal(t, testObject.Meta(), identifiers[0].Meta)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, testObject, included[0])
}

func Test_transformRelationNodes_Struct(t *testing.T) {
	resource, included := transformRelationNodes(testObject)

	assert.NotNil(t, resource)
	identifier, ok := resource.(internalResourceIdentifier)
	assert.True(t, ok)
	assert.Equal(t, testObject.ID(), identifier.ID)
	assert.Equal(t, testObject.Type(), identifier.Type)
	assert.NotNil(t, identifier.Meta)
	assert.Equal(t, testObject.Meta(), identifier.Meta)

	assert.NotNil(t, included)
	assert.Equal(t, 1, len(included))
	assert.Equal(t, testObject, included[0])
}

func Test_transformRelationNodes_Nil(t *testing.T) {
	resource, included := transformRelationNodes("testObject")

	assert.Nil(t, resource)
	assert.Nil(t, included)
}

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
