package examples

import (
	"github.com/alehechka/go-jsonapi/jsonapi"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	engine.GET("/record", getRecords)
	engine.GET("/record/:id", getRecord)

	engine.Run()
}

type Record struct {
	RecordID  string `json:"-"`
	FirstName string `json:"firstName"`
	Age       int    `json:"age"`
	IsAdmin   bool   `json:"isAdmin"`
}

func (record Record) ID() string {
	return record.RecordID
}

func (record Record) Type() string {
	return "records"
}

func getRecords(ctx *gin.Context) {

	ctx.JSON(200, jsonapi.CreateCollectionResponse(ctx.Request)(jsonapi.CollectionResponse{
		Nodes: records(),
	}))
}

func getRecord(ctx *gin.Context) {
	records := records()
	recordID := ctx.Param("id")

	var record Record
	for _, r := range records {
		if r.RecordID == recordID {
			record = r
		}
	}

	ctx.JSON(200, jsonapi.CreateResponse(ctx.Request)(jsonapi.Response{
		Node: record,
	}))
}

func records() []Record {
	return []Record{{
		RecordID:  "1234",
		FirstName: "Joe",
		Age:       25,
		IsAdmin:   true,
	},
		{
			RecordID:  "4321",
			FirstName: "Sally",
			Age:       30,
			IsAdmin:   true,
		}}
}
