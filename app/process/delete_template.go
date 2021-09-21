package process

import (
	"context"
	"log"
	"net/http"
	"test-mongo/app/process/reqreader"
	"test-mongo/app/process/respsender"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (p *Process) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	pageId := reqreader.ReadPathParam(r, "page-id")
	templateId, _ := primitive.ObjectIDFromHex(reqreader.ReadPathParam(r, "template-id"))

	query := bson.M{
		"page_id":               pageId,
		"templates.template_id": templateId,
	}
	_, err := p.Template.FindOne(context.TODO(), query, bson.M{})
	if err != nil {
		var errStr string
		var errResp string
		switch err {
		case mongo.ErrNoDocuments:
			errStr = `not found this template_id ` + templateId.Hex() + ``
			errResp = `{"success": false, "code": 20300, "error": "` + errStr + `"}`
		default:
			errStr = `unexpected database error occurred`
			errResp = `{"success": false, "code": 20301, "error": "` + errStr + `"}`
		}

		log.Println(errStr + `:: ` + err.Error())
		respsender.ResponseString(w, errResp, http.StatusInternalServerError)
		return
	}

	info := map[string]interface{}{
		"page_id":     pageId,
		"template_id": templateId,
	}

	if err := p.Template.PullTemplateItem(info); err != nil {
		respsender.ResponseString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"success": true,
	}

	respsender.ResponseMap(w, data, http.StatusOK)
}
