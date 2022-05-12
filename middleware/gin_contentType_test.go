package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alehechka/go-jsonapi/jsonapi"
	"github.com/alehechka/go-jsonapi/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_ContentType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/example", nil)

	middleware.ContentType(c)

	assert.Equal(t, jsonapi.MediaType, c.Writer.Header().Get(jsonapi.ContentType))
}
