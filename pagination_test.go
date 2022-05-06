package jsonapi_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi"
	"github.com/stretchr/testify/assert"
)

func Test_GetPageOffset(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[offset]=10", nil)

	offset, err := jsonapi.GetPageOffset(req)

	assert.Nil(t, err)
	assert.Equal(t, 10, offset)
}

func Test_GetPageOffset_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	offset, err := jsonapi.GetPageOffset(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, offset)
}

func Test_GetPageOffset_ParsingError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[offset]=ten", nil)

	offset, err := jsonapi.GetPageOffset(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, offset)
}

func Test_GetPageLimit(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[limit]=10", nil)

	limit, err := jsonapi.GetPageLimit(req)

	assert.Nil(t, err)
	assert.Equal(t, 10, limit)
}

func Test_GetPageLimit_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	limit, err := jsonapi.GetPageLimit(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, limit)
}

func Test_GetPageLimit_ParsingError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[limit]=ten", nil)

	limit, err := jsonapi.GetPageLimit(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, limit)
}

func Test_GetPageNumber(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[number]=10", nil)

	number, err := jsonapi.GetPageNumber(req)

	assert.Nil(t, err)
	assert.Equal(t, 10, number)
}

func Test_GetPageNumber_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	number, err := jsonapi.GetPageNumber(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, number)
}

func Test_GetPageNumber_ParsingError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[number]=ten", nil)

	number, err := jsonapi.GetPageNumber(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, number)
}

func Test_GetPageSize(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[size]=10", nil)

	size, err := jsonapi.GetPageSize(req)

	assert.Nil(t, err)
	assert.Equal(t, 10, size)
}

func Test_GetPageSize_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	size, err := jsonapi.GetPageSize(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, size)
}

func Test_GetPageSize_ParsingError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[size]=ten", nil)

	size, err := jsonapi.GetPageSize(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, size)
}

func Test_GetPageCursor(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[cursor]=10", nil)

	cursor, err := jsonapi.GetPageCursor(req)

	assert.Nil(t, err)
	assert.Equal(t, 10, cursor)
}

func Test_GetPageCursor_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	cursor, err := jsonapi.GetPageCursor(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, cursor)
}

func Test_GetPageCursor_ParsingError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[cursor]=ten", nil)

	cursor, err := jsonapi.GetPageCursor(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, cursor)
}

func Test_GetPageBefore(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[before]=10", nil)

	before, err := jsonapi.GetPageBefore(req)

	assert.Nil(t, err)
	assert.Equal(t, 10, before)
}

func Test_GetPageBefore_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	before, err := jsonapi.GetPageBefore(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, before)
}

func Test_GetPageBefore_ParsingError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[before]=ten", nil)

	before, err := jsonapi.GetPageBefore(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, before)
}

func Test_GetPageAfter(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[after]=10", nil)

	after, err := jsonapi.GetPageAfter(req)

	assert.Nil(t, err)
	assert.Equal(t, 10, after)
}

func Test_GetPageAfter_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example", nil)

	after, err := jsonapi.GetPageAfter(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, after)
}

func Test_GetPageAfter_ParsingError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[after]=ten", nil)

	after, err := jsonapi.GetPageAfter(req)

	assert.NotNil(t, err)
	assert.Equal(t, 0, after)
}

func Test_FindUnsupportedPagination(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/example?page[after]=ten&page[number]=10", nil)

	errs := jsonapi.FindUnsupportedPagination(req)(jsonapi.PageAfter)

	assert.NotNil(t, errs)
	assert.Equal(t, 1, len(errs))

	err := errs[0]
	assert.Equal(t, jsonapi.PageAfter.String(), err.Source.(jsonapi.ErrorSource).Parameter)
}
