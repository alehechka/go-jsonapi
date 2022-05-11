package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testObject = testStructMethods{
	TestID: "1234",
	Number: 1234,
	IsTest: true,
	TestData: []testStruct{
		{
			TestID: "4321",
			Number: 4321,
			IsTest: false,
		},
	},
	SingleTest: testStruct{
		TestID: "9876",
		Number: 9876,
		IsTest: true,
	},
}

var testObjects = []testStructMethods{testObject}

func Test_transformNodes_NodeSlice(t *testing.T) {
	nodes, included := transformNodes(testObjects, baseURL)

	assertOutput(t, nodes, included, "NodeSlice")
}

func Test_transformNodes_NodePtrSlice(t *testing.T) {
	nodes, included := transformNodes(&testObjects, baseURL)

	assertOutput(t, nodes, included, "NodePtrSlice")
}

func Test_transformNodes_NodeStruct(t *testing.T) {
	nodes, included := transformNodes(testObjects[0], baseURL)

	assertOutput(t, nodes, included, "NodeStruct")
}

func Test_transformNodes_NodePtrStruct(t *testing.T) {
	nodes, included := transformNodes(&testObjects[0], baseURL)

	assertOutput(t, nodes, included, "NodePtrStruct")
}

func assertOutput(t *testing.T, nodes []internalNode, included []Node, message string) {
	assert.NotNil(t, nodes, message)
	assert.Equal(t, 1, len(nodes), message)

	node := nodes[0]

	assert.Equal(t, testObject.ID(), node.ID, message)
	assert.Equal(t, testObject.Type(), node.Type, message)
	assert.Equal(t, testObject, node.Attributes, message)

	assert.NotNil(t, node.Meta)
	assert.Equal(t, testObject.Meta(), node.Meta)

	assert.NotNil(t, node.Links)
	assert.NotNil(t, node.Links[SelfKey])
	assert.Equal(t, baseURL+"/path/to/resource", node.Links[SelfKey])

	assert.NotNil(t, node.Relationships)
	assert.NotNil(t, node.Relationships["tests"])
	assert.Equal(t, testObject.TestData[0].ID(), node.Relationships["tests"].Data.([]internalResourceIdentifier)[0].ID)
	assert.Equal(t, testObject.TestData[0].Type(), node.Relationships["tests"].Data.([]internalResourceIdentifier)[0].Type)
	assert.NotNil(t, node.Relationships["test"])
	assert.Equal(t, testObject.SingleTest.ID(), node.Relationships["test"].Data.(internalResourceIdentifier).ID)
	assert.Equal(t, testObject.SingleTest.Type(), node.Relationships["test"].Data.(internalResourceIdentifier).Type)

	assert.NotNil(t, included)
	assert.Equal(t, 2, len(included))

	for _, inc := range included {
		switch inc.ID() {
		case testObject.TestData[0].ID():
			assert.Equal(t, testObject.TestData[0].ID(), inc.ID())
			assert.Equal(t, testObject.TestData[0].Type(), inc.Type())
			assert.Equal(t, testObject.TestData[0], inc)
		case testObject.SingleTest.ID():
			assert.Equal(t, testObject.SingleTest.ID(), inc.ID())
			assert.Equal(t, testObject.SingleTest.Type(), inc.Type())
			assert.Equal(t, testObject.SingleTest, inc)
		}
	}
}

func Test_transformCollectionResponseNodes_WithErrors(t *testing.T) {
	response := CollectionResponse{
		Nodes: []testStruct{},
		Errors: Errors{
			{
				Status: 400,
				Detail: "this has an error",
			},
		},
	}

	node, included := transformCollectionResponseNodes(response, baseURL)

	assert.Nil(t, node)
	assert.Nil(t, included)
}
