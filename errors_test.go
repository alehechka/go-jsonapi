package jsonapi

import (
	"testing"

	"github.com/dmsi-io/go-utils/errorshared"
	"gotest.tools/assert"
)

func TestCreateErrorsFromResponse(t *testing.T) {
	type response struct {
		ReturnCode    int
		MessageNum    int
		MessageText   string
		SomethingElse bool
	}

	resp := response{
		ReturnCode:    2,
		MessageNum:    errorshared.InvalidCredentials.Int(),
		MessageText:   "Error Message",
		SomethingElse: true,
	}

	statusCode := 400

	errs := CreateErrorsFromResponse(statusCode, resp.ReturnCode, resp.MessageNum, resp.MessageText, nil)
	assert.Equal(t, len(errs), 1)

	err := errs[0]

	assert.Equal(t, err.Status, statusCode)
	assert.Equal(t, err.Code, resp.MessageNum)
	assert.Equal(t, err.Title, errorshared.GetErrorTitleFromCode(errorshared.InvalidCredentials))
	assert.Equal(t, err.Detail, resp.MessageText)
}

func TestCreateErrorFromResponse(t *testing.T) {
	type response struct {
		ReturnCode    int
		MessageNum    int
		MessageText   string
		SomethingElse bool
	}

	resp := response{
		ReturnCode:    2,
		MessageNum:    errorshared.InvalidCredentials.Int(),
		MessageText:   "Error Message",
		SomethingElse: true,
	}

	statusCode := 400

	err, isError := CreateErrorFromResponse(statusCode, resp.ReturnCode, resp.MessageNum, resp.MessageText, nil)

	assert.Equal(t, err.Status, statusCode)
	assert.Equal(t, err.Code, resp.MessageNum)
	assert.Equal(t, err.Title, errorshared.GetErrorTitleFromCode(errorshared.InvalidCredentials))
	assert.Equal(t, err.Detail, resp.MessageText)
	assert.Equal(t, isError, true)
}
