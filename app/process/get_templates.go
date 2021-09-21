package process

import (
	"context"
	"log"
	"net/http"
	"test-mongo/app/model"
	"test-mongo/app/process/reqreader"
	"test-mongo/app/process/respsender"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (p *Process) GetAllTemplate(w http.ResponseWriter, r *http.Request) {
	pageId := reqreader.ReadPathParam(r, "page-id")
	var template model.Template
	query := bson.M{
		"page_id": pageId,
	}
	template, err := p.Template.FindOne(context.TODO(), query, bson.M{})
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
