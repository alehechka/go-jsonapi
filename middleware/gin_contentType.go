package middleware

import (
	"github.com/alehechka/go-jsonapi/jsonapi"
	"github.com/gin-gonic/gin"
)

// ContentType is a middleware function to set the Content-Type header to the official JSON:API header.
func ContentType(c *gin.Context) {
	c.Header(jsonapi.ContentType, jsonapi.MediaType)
	c.Next()
}
