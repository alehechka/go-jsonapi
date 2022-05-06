package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi"
	"github.com/alehechka/go-jsonapi/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_UnsupportedPagination_Abort(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/?page[offset]=10&page[size]=10", nil)

	middleware.UnsupportedPagination(jsonapi.PageOffset)(c)

	assert.Equal(t, true, c.IsAborted())
}

func Test_UnsupportedPagination_Next(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/?page[offset]=10&page[size]=10", nil)

	middleware.UnsupportedPagination(jsonapi.PageAfter)(c)

	assert.Equal(t, false, c.IsAborted())
}

func Test_ExceedsMaximumPaginationSize_Abort(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/?page[limit]=1000&page[size]=10", nil)

	middleware.ExceedsMaximumPaginationSize(100)(c)

	assert.Equal(t, true, c.IsAborted())
}

func Test_ExceedsMaximumPaginationSize_Next(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/?page[offset]=10&page[size]=10", nil)

	middleware.ExceedsMaximumPaginationSize(100)(c)

	assert.Equal(t, false, c.IsAborted())
}
