package middleware

import (
	"net/http"

	"github.com/alehechka/go-jsonapi"
	"github.com/gin-gonic/gin"
)

// UnsupportedPagination will short-circuit if one of the provided, unsupported query options is provided in the request.
func UnsupportedPagination(paginationOptions ...jsonapi.PaginationOption) gin.HandlerFunc {

	return func(c *gin.Context) {

		errs := jsonapi.FindUnsupportedPagination(c.Request)(paginationOptions...)

		if errs.HasErrors() {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonapi.CreateResponse(c.Request)(jsonapi.Response{Errors: errs}))
			return
		}

		c.Next()
	}
}
