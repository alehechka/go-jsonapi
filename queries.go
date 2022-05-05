package jsonapi

import (
	"net/http"
	"strconv"
)

// PaginationOption represents the available JSON:API query parameters
type PaginationOption string

func (option PaginationOption) String() string {
	return string(option)
}

func (option PaginationOption) QueryExists(request *http.Request) bool {
	return request.URL.Query().Has(option.String())
}

// Pagination Query Parameters
const (
	// PageOffset represents the offset from previous pagination. Use in conjunction with PageLimit.
	PageOffset PaginationOption = "page[offset]"
	// PageLimit represents the size limit of page offset. Use in conjunction with PageOffset.
	PageLimit PaginationOption = "page[limit]"
	// PageNumber represents the page number. Use in conjunction with PageSize.
	PageNumber PaginationOption = "page[number]"
	// PageSize represents the size limit of page. Use in conjunction with PageNumber.
	PageSize PaginationOption = "page[size]"
	// PageCursor represents the cursor of the page. Use by itself.
	PageCursor PaginationOption = "page[cursor]"
	// PageBefore used with cursor-based pagination. Use in conjunction with PageSize and PageAfter
	PageBefore PaginationOption = "page[before]"
	// PageAfter used with cursor-based pagination. Use in conjunction with PageSize and PageBefore
	PageAfter PaginationOption = "page[after]"
	// Include used to request extra resources to include
	Include string = "include"
)

// GetPageOffset retrieves integer parsed PageOffset from query parameters
func GetPageOffset(request *http.Request) (int, error) {
	return getQueryInteger(request, PageOffset)
}

// GetPageLimit retrieves integer parsed PageOffset from query parameters
func GetPageLimit(request *http.Request) (int, error) {
	return getQueryInteger(request, PageLimit)
}

// GetPageNumber retrieves integer parsed PageOffset from query parameters
func GetPageNumber(request *http.Request) (int, error) {
	return getQueryInteger(request, PageNumber)
}

// GetPageSize retrieves integer parsed PageOffset from query parameters
func GetPageSize(request *http.Request) (int, error) {
	return getQueryInteger(request, PageSize)
}

// GetPageCursor retrieves integer parsed PageOffset from query parameters
func GetPageCursor(request *http.Request) (int, error) {
	return getQueryInteger(request, PageCursor)
}

// GetPageBefore retrieves integer parsed PageOffset from query parameters
func GetPageBefore(request *http.Request) (int, error) {
	return getQueryInteger(request, PageBefore)
}

// GetPageAfter retrieves integer parsed PageOffset from query parameters
func GetPageAfter(request *http.Request) (int, error) {
	return getQueryInteger(request, PageAfter)
}

func getQueryInteger(request *http.Request, option PaginationOption) (int, error) {
	return strconv.Atoi(request.URL.Query().Get(option.String()))
}
