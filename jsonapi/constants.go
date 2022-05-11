package jsonapi

import "errors"

// Standard JSON:API Link mapping keys
const (
	// NextKey represents the key for the next link
	NextKey string = "next"
	// PreviousKey represents the key for the previous link
	PreviousKey string = "prev"
	// SelfKey represents the key for the self link
	SelfKey string = "self"
	// FirstKey represents the key for the first link
	FirstKey string = "first"
	// LastKey represents the key for the last link
	LastKey string = "last"
	// RelatedKey represents the key for the related link
	RelatedKey string = "related"
)

// Include query parameter used to request extra resources to include in response
const Include string = "include"

// Standard HTTP Headers
const (
	// ForwardedPrefix represents the prefix that is dropped when proxied through rest-api
	ForwardedPrefix string = "X-Forwarded-Prefix"
	// ForwardedProto represents the protocol that is received prior to being forwarded (http | https)
	ForwardedProto string = "X-Forwarded-Proto"
	// ForwardedHost represents the original host header received prior to being forwarded (example.com)
	ForwardedHost string = "X-Forwarded-Host"
)

// Standard JSON:API Pagination Query Parameters
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

// PaginationOptions is the array of all available PaginationOption items
var PaginationOptions []PaginationOption = []PaginationOption{
	PageOffset,
	PageLimit,
	PageNumber,
	PageSize,
	PageCursor,
	PageBefore,
	PageAfter,
}

// Error messages
var (
	// ErrTooMinterface{}Included number of included is greater than number of available resources
	ErrTooManyIncluded error = errors.New("included query has too many resources")
	// ErrResourceNotAvailable member of included is not an available resource
	ErrResourceNotAvailable error = errors.New("resource from included query not available")
)
