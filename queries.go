package jsonapi

// Pagination Query Parameters
const (
	// PageOffset represents the offset from previous pagination. Use in conjunction with PageLimit.
	PageOffset string = "page[offset]"
	// PageLimit represents the size limit of page offset. Use in conjunction with PageOffset.
	PageLimit string = "page[limit]"
	// PageNumber represents the page number. Use in conjunction with PageSize.
	PageNumber string = "page[number]"
	// PageSize represents the size limit of page. Use in conjunction with PageNumber.
	PageSize string = "page[size]"
	// PageCursor represents the cursor of the page. Use by itself.
	PageCursor string = "page[cursor]"
	// PageBefore used with cursor-based pagination. Use in conjunction with PageSize and PageAfter
	PageBefore string = "page[before]"
	// PageAfter used with cursor-based pagination. Use in conjunction with PageSize and PageBefore
	PageAfter string = "page[after]"
)
