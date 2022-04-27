package jsonapi

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// ParseReadCloser will decode JSON:API from an io.ReadCloser object and return a map of expected JSON:API objects.
func ParseReadCloser(body io.ReadCloser, objs map[string]interface{}) (map[string][]interface{}, error) {
	var bodyMap map[string]interface{}

	if err := json.NewDecoder(body).Decode(&bodyMap); err != nil {
		return nil, err
	}

	return parseJSONAPIArrays(bodyMap, objs), nil
}

// ParseGinContext will decode the body from a gin.Context and return a map of expected JSON:API objects.
func ParseGinContext(c *gin.Context, objs map[string]interface{}) (map[string][]interface{}, error) {
	var body map[string]interface{}

	if err := c.ShouldBind(&body); err != nil {
		return nil, err
	}

	return parseJSONAPIArrays(body, objs), nil
}

// ParseJSONAPIArrays parses all objects from "included" array
func parseJSONAPIArrays(body map[string]interface{}, objs map[string]interface{}) map[string][]interface{} {
	arrObjs := make(map[string][]interface{}, 0)

	arrObjs = parseJSONAPI(objs, arrObjs, body["data"])
	arrObjs = parseJSONAPI(objs, arrObjs, body["included"])

	return arrObjs
}

func parseJSONAPI(objs map[string]interface{}, arrObjs map[string][]interface{}, jsonapi interface{}) map[string][]interface{} {

	switch jsonapi.(type) {
	case map[string]interface{}:
		valMap := jsonapi.(map[string]interface{})
		if valMap["type"] != nil {
			jsonapiType := valMap["type"].(string)
			temp := objs[jsonapiType]
			if err := parseJSONAPIObject(jsonapiType, &temp, valMap); err == nil {
				arrObjs[jsonapiType] = append(arrObjs[jsonapiType], temp)
			}
		}
	case []interface{}:
		for _, val := range jsonapi.([]interface{}) {
			arrObjs = parseJSONAPI(objs, arrObjs, val)
		}
	}

	return arrObjs
}

func parseJSONAPIObject(jsonapiType string, obj interface{}, val map[string]interface{}) error {
	if val["type"] != nil && val["type"].(string) == jsonapiType {
		return mapstructure.Decode(val["attributes"].(map[string]interface{}), obj)
	}

	return errors.New("provided type does not match object type")
}
