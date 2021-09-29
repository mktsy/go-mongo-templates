package process

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"test-mongo/app/process/reqreader"
	"test-mongo/app/process/respsender"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (p *Process) GetTemplates(w http.ResponseWriter, r *http.Request) {
	pageId := reqreader.ReadPathParam(r, "page-id")
	limitStr := reqreader.ReadQueryParam(r, "limit")
	offsetStr := reqreader.ReadQueryParam(r, "offset")

	limitInt, _ := strconv.Atoi(limitStr)
	offsetInt, _ := strconv.Atoi(offsetStr)

	query := bson.M{
		"page_id": pageId,
	}
	projectionOpt := bson.M{
		"_id":       1,
		"page_id":   1,
		"title":     1,
		"text":      1,
		"image_url": 1,
	}
	sortOpt := bson.M{}

	limitOpt := int64(limitInt)
	offsetOpt := int64(offsetInt)

	template, err := p.Template.FindAll(context.TODO(), &query, projectionOpt, limitOpt, offsetOpt, sortOpt)
	if err != nil {
		var errStr string
		var errResp string
		switch err {
		case mongo.ErrNoDocuments:
			errStr = `not found this page_id ` + pageId + ``
			errResp = `{"success": false, "code": 20300, "error": "` + errStr + `"}`
		default:
			errStr = `unexpected database error occurred`
			errResp = `{"success": false, "code": 20301, "error": "` + errStr + `"}`
		}

		log.Println(errStr + `:: ` + err.Error())
		respsender.ResponseString(w, errResp, http.StatusInternalServerError)
		return
	}

	log.Println("found that template")
	data := map[string]interface{}{
		"success": true,
		"result":  template,
	}

	respsender.ResponseMap(w, data, http.StatusOK)
}
