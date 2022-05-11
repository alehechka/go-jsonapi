package jsonapi

// ErrorSource is the standard JSONAPI Error Source struct
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// Error is the standard JSONAPI Error struct
type Error struct {
	ID     string      `json:"id,omitempty"`
	Links  Links       `json:"links,omitempty"`
	Status int         `json:"status,omitempty"`
	Code   int         `json:"code,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Source interface{} `json:"source,omitempty"` // ErrorSource
	Meta   interface{} `json:"meta,omitempty"`
}

// Errors is a standard array of JSONAPI Error structs
type Errors []Error

// HasErrors checks if Errors array has one more errors.
// data key cannot coexist in response with errors: https://jsonapi.org/format/#document-top-level
func (errs Errors) HasErrors() bool {
	return len(errs) > 0
}

type internalError struct {
	ID     string      `json:"id,omitempty"`
	Links  LinkMap     `json:"links,omitempty"`
	Status int         `json:"status,omitempty"`
	Code   int         `json:"code,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Source interface{} `json:"source,omitempty"` // ErrorSource
	Meta   interface{} `json:"meta,omitempty"`
}

func transformErrors(errs []Error, baseURL string) []internalError {
	var internalErrors []internalError

	for _, err := range errs {
		internalErrors = append(internalErrors, transformError(err, baseURL))
	}

	return internalErrors
}

func transformError(err Error, baseURL string) internalError {
	return internalError{
		ID:     err.ID,
		Links:  TransformLinks(err.Links, baseURL),
		Status: err.Status,
		Code:   err.Code,
		Title:  err.Title,
		Detail: err.Detail,
		Source: err.Source,
		Meta:   err.Meta,
	}
}
