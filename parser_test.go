package jsonapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

type dataExample struct {
	StringValue  string  `json:"stringValue"`
	IntValue     int     `json:"intValue"`
	DecimalValue float64 `json:"decimalValue"`
	BoolValue    bool    `json:"boolValue"`
}

type includedExampleOne struct {
	ValueOne   string  `json:"valueOne"`
	ValueTwo   int     `json:"valueTwo"`
	ValueThree float64 `json:"valueThree"`
	ValueFour  bool    `json:"valueFour"`
}

type includedExampleTwo struct {
	ThingOne   string  `json:"thingOne"`
	ThingTwo   int     `json:"thingTwo"`
	ThingThree float64 `json:"thingThree"`
	ThingFour  bool    `json:"thingFour"`
}

var arrObjs = map[string]interface{}{
	"dataExample":        dataExample{},
	"includedExampleOne": includedExampleOne{},
	"includedExampleTwo": includedExampleTwo{},
}

func TestParseReadCloser(t *testing.T) {

	reader, err := os.Open("parser.json")
	assert.NilError(t, err)
	defer reader.Close()

	included, err := ParseReadCloser(reader, arrObjs)
	assert.NilError(t, err)

	testParserResponse(t, included)
}

func TestParseGinContext(t *testing.T) {
	reader, err := os.Open("parser.json")
	assert.NilError(t, err)
	defer reader.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := &http.Request{
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	req.Body = reader
	c.Request = req

	included, err := ParseGinContext(c, arrObjs)
	assert.NilError(t, err)

	testParserResponse(t, included)
}

func TestParseJSONAPIArrays(t *testing.T) {

	m, err := openJSONFile("parser.json")
	assert.NilError(t, err)

	included := parseJSONAPIArrays(m, arrObjs)

	testParserResponse(t, included)
}

func testParserResponse(t *testing.T, resp map[string][]interface{}) {
	assert.Equal(t, len(resp), 4)

	assert.Equal(t, len(resp["dataExample"]), 1)
	_, ok := resp["dataExample"][0].(dataExample)
	assert.Equal(t, ok, true)

	assert.Equal(t, len(resp["includedExampleOne"]), 2)
	_, ok = resp["includedExampleOne"][0].(includedExampleOne)
	assert.Equal(t, ok, true)

	assert.Equal(t, len(resp["includedExampleTwo"]), 1)
	_, ok = resp["includedExampleTwo"][0].(includedExampleTwo)
	assert.Equal(t, ok, true)

	assert.Equal(t, len(resp["non-included"]), 1)
	_, ok = resp["non-included"][0].(map[string]interface{})
	assert.Equal(t, ok, true)
}

func openJSONFile(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	// Creating the maps for JSON
	m := map[string]interface{}{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &m)

	return m, nil
}
