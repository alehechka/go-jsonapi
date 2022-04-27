package jsonapi

import (
	"github.com/dmsi-io/go-logger"
	"github.com/dmsi-io/go-utils/errorshared"
)

//ErrorSource is the standard JSONAPI Error Source struct
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

//Error is the standard JSONAPI Error struct
type Error struct {
	ID     string      `json:"id,omitempty"`
	Links  Links       `json:"links,omitempty"`
	Status int         `json:"status"`
	Code   int         `json:"code"`
	Title  string      `json:"title"`
	Detail string      `json:"detail,omitempty"`
	Source interface{} `json:"source,omitempty"` // ErrorSource
	Meta   interface{} `json:"meta,omitempty"`
}

type internalError struct {
	ID     string      `json:"id,omitempty"`
	Links  LinkMap     `json:"links,omitempty"`
	Status int         `json:"status"`
	Code   int         `json:"code"`
	Title  string      `json:"title"`
	Detail string      `json:"detail,omitempty"`
	Source interface{} `json:"source,omitempty"` // ErrorSource
	Meta   interface{} `json:"meta,omitempty"`
}

func transformToInternalErrorStructs(errs []Error, baseURL string) []internalError {
	var internalErrors []internalError

	for _, err := range errs {
		internalErrors = append(internalErrors, transformToInternalErrorStruct(err, baseURL))
	}

	return internalErrors
}

func transformToInternalErrorStruct(err Error, baseURL string) internalError {
	return internalError{
		ID:     err.ID,
		Links:  transformLinks(err.Links, baseURL),
		Status: err.Status,
		Code:   err.Code,
		Title:  err.Title,
		Detail: err.Detail,
		Source: err.Source,
		Meta:   err.Meta,
	}
}

// CreateErrorsFromResponse creates an array of Error objects from provided RootResponse
func CreateErrorsFromResponse(statusCode int, returnCode int, messageNum int, messageText string, log logger.LogInterface) (errs []Error) {
	err, isError := CreateErrorFromResponse(statusCode, returnCode, messageNum, messageText, log)

	if isError {
		return []Error{err}
	}

	return make([]Error, 0)
}

// CreateErrorFromResponse creates an Error object from provided RootResponse
func CreateErrorFromResponse(statusCode int, returnCode int, messageNum int, messageText string, log logger.LogInterface) (err Error, isError bool) {

	if returnCode != 0 {
		if log != nil {
			log.Error(messageText)
		}
		return Error{
			Status: statusCode,
			Code:   messageNum,
			Title:  errorshared.GetErrorTitleFromInt(messageNum),
			Detail: messageText,
		}, true
	}

	return Error{}, false
}
