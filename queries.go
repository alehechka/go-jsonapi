package jsonapi

import "github.com/gin-gonic/gin"

// QueryParam represents available standard JSONAPI query parameters.
type QueryParam string

// Pagination Query Parameters
const (
	// PageOffset represents the offset from previous pagination. Use in conjunction with PageLimit.
	PageOffset QueryParam = "page[offset]"
	// PageLimit represents the size limit of page offset. Use in conjunction with PageOffset.
	PageLimit QueryParam = "page[limit]"
	// PageNumber represents the page number. Use in conjunction with PageSize.
	PageNumber QueryParam = "page[number]"
	// PageSize represents the size limit of page. Use in conjunction with PageNumber.
	PageSize QueryParam = "page[size]"
	// PageCursor represents the cursor of the page. Use by itself.
	PageCursor QueryParam = "page[cursor]"
)

func (q QueryParam) String() string {
	return string(q)
}

// GetQueryParameter retrieves standard JSONAPI query parameters.
func GetQueryParameter(c *gin.Context, query QueryParam) string {
	return c.Query(string(query))
}

// Agility Pagination Parameters
const (
	// ApplyDataChunking is a bool flag that applied data chunking
	ApplyDataChunking string = "applyDataChunking"
	// ChunkStartPointer is an int value that represents chunk start pointer
	ChunkStartPointer string = "chunkStartPointer"
	// RecordFetchLimit is an int value that represents record fetch limit
	RecordFetchLimit string = "recordFetchLimit"
	// MoreResultsAvailable is a bool flag that represents more results available
	MoreResultsAvailable string = "moreResultsAvailable"
	// NextChunkStartPointer is an int value that represents next chunk start pointer
	NextChunkStartPointer string = "nextChunkStartPointer"
)
