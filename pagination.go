package jsonapi

import (
	"fmt"
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
)

var (
	// PaginationOptions is the array of all available PaginationOption items
	PaginationOptions []PaginationOption = []PaginationOption{
		PageOffset,
		PageLimit,
		PageNumber,
		PageSize,
		PageCursor,
		PageBefore,
		PageAfter,
	}
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

// CheckUnsupportedPagination will return with an array of Errors if any unsupported pagination options are found in query parameters
func CheckUnsupportedPagination(request *http.Request) func(unsupportedOptions ...PaginationOption) Errors {
	return func(unsupportedOptions ...PaginationOption) (errs Errors) {

		for _, option := range unsupportedOptions {
			if option.QueryExists(request) {
				errs = append(errs, CreateRangePaginationError(option))
			}
		}

		return
	}
}

func CheckSupportedPagination(request *http.Request) func(supportedOptions ...PaginationOption) Errors {
	return func(supportedOptions ...PaginationOption) (errs Errors) {
		var unsupportedOptions []PaginationOption
	Outer:
		for _, option := range PaginationOptions {
			for _, supportedOption := range supportedOptions {
				if option == supportedOption {
					continue Outer
				}
			}
			unsupportedOptions = append(unsupportedOptions, option)
		}

		return CheckUnsupportedPagination(request)(unsupportedOptions...)
	}
}

func CreateRangePaginationError(option PaginationOption) Error {
	return Error{
		Title:  "Range Pagination Not Supported.",
		Detail: fmt.Sprintf("%s is not a supported pagination option", option),
		Source: ErrorSource{
			Parameter: option.String(),
		},
		Status: http.StatusBadRequest,
		Links: Links{
			"type": {
				Href: "https://jsonapi.org/profiles/ethanresnick/cursor-pagination/#auto-id--range-pagination-not-supported-error",
			},
		},
	}
}

// CheckExceedsMaximumPaginationSize checks the provided request to see if any provided pagination options exceed the provided maximum
func CheckExceedsMaximumPaginationSize(request *http.Request) func(maxSize int) Errors {
	return func(maxSize int) (errs Errors) {

		if PageSize.QueryExists(request) {
			if pageSize, _ := GetPageSize(request); pageSize > maxSize {
				errs = append(errs, Error{
					Title:  "Page size requested is too large.",
					Detail: fmt.Sprintf("You requested a size of %d, but %d is the maximum.", pageSize, maxSize),
					Source: ErrorSource{
						Parameter: PageSize.String(),
					},
					Links: Links{
						"type": {
							Href: "https://jsonapi.org/profiles/ethanresnick/cursor-pagination/#auto-id--max-page-size-exceeded-error",
						},
					},
					Meta: Meta{
						"page": Meta{
							"maxSize": maxSize,
						},
					},
				})
			}
		}

		if PageLimit.QueryExists(request) {
			if pageLimit, _ := GetPageLimit(request); pageLimit > maxSize {
				errs = append(errs, Error{
					Title:  "Page limit requested is too large.",
					Detail: fmt.Sprintf("You requested a limit of %d, but %d is the maximum.", pageLimit, maxSize),
					Source: ErrorSource{
						Parameter: PageLimit.String(),
					},
					Links: Links{
						"type": {
							Href: "https://jsonapi.org/profiles/ethanresnick/cursor-pagination/#auto-id--max-page-size-exceeded-error",
						},
					},
					Meta: Meta{
						"page": Meta{
							"maxLimit": maxSize,
						},
					},
				})
			}
		}

		return
	}
}
