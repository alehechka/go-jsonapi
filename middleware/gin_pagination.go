package middleware

import (
	"net/http"

	"github.com/alehechka/go-jsonapi/jsonapi"
	"github.com/gin-gonic/gin"
)

// SupportedPagination will short-circuit if a pagination query that is not in the provided, supported options.
func SupportedPagination(supportedOptions ...jsonapi.PaginationOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		errs := jsonapi.CheckSupportedPagination(c.Request)(supportedOptions...)

		if errs.HasErrors() {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonapi.CreateResponse(c.Request)(jsonapi.Response{Errors: errs}))
			return
		}

		c.Next()
	}
}

// UnsupportedPagination will short-circuit if one of the provided, unsupported query options is provided in the request.
func UnsupportedPagination(unsupportedOptions ...jsonapi.PaginationOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		errs := jsonapi.CheckUnsupportedPagination(c.Request)(unsupportedOptions...)

		if errs.HasErrors() {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonapi.CreateResponse(c.Request)(jsonapi.Response{Errors: errs}))
			return
		}

		c.Next()
	}
}

// MaximumPaginationSize will short-circuit if one of the provided pagination query options exceeds the provided maximum.
func MaximumPaginationSize(maxSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		errs := jsonapi.CheckExceedsMaximumPaginationSize(c.Request)(maxSize)

		if errs.HasErrors() {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonapi.CreateResponse(c.Request)(jsonapi.Response{Errors: errs}))
			return
		}

		c.Next()
	}
}
