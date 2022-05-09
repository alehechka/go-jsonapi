package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://example.com"

type testStruct struct {
	TestID string
	Number int
	IsTest bool
}

func (d testStruct) ID() string {
	return d.TestID
}

func (d testStruct) Type() string {
	return "test"
}

func Test_transformNode(t *testing.T) {
	node := testStruct{
		TestID: "1234",
		Number: 1234,
		IsTest: true,
	}

	transformed, _ := transformNode(node, baseURL)

	assert.Equal(t, node, transformed.Attributes)
	assert.Equal(t, node.TestID, transformed.ID)
	assert.Equal(t, node.Type(), transformed.Type)
}

type testStructMethods struct {
	TestID     string
	Number     int
	IsTest     bool
	TestData   []testStruct
	SingleTest testStruct
}

func (d testStructMethods) ID() string {
	return d.TestID
}

func (d testStructMethods) Type() string {
	return "test"
}

func (d testStructMethods) Attributes() interface{} {
	return d
}

func (d testStructMethods) Meta() interface{} {
	return Meta{
		"something": "interesting",
	}
}

func (d testStructMethods) Links() Links {
	return Links{
		SelfKey: Link{
			Href: "/path/to/resource",
		},
	}
}

func (d testStructMethods) Relationships() map[string]interface{} {
	return map[string]interface{}{
		"tests": d.TestData,
		"test":  d.SingleTest,
	}
}

func Test_transformNode_Methods(t *testing.T) {
	node := testStructMethods{
		TestID: "1234",
		Number: 1234,
		IsTest: true,
		TestData: []testStruct{{
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

	transformed, included := transformNode(node, baseURL)

	assertOutput(t, []internalNode{transformed}, included, "Methods")
}

func Test_transformResponseNode_WithErrors(t *testing.T) {
	response := Response{
		Node: testStruct{},
		Errors: Errors{
			{
				Status: 400,
				Detail: "this has an error",
			},
		},
	}

	node, included := transformResponseNode(response, baseURL)

	assert.Nil(t, node)
	assert.Nil(t, included)
}
