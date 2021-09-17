package process

import (
	"context"
	"log"
	"net/http"
	"test-mongo/app/model"
	"test-mongo/app/process/reqreader"
	"test-mongo/app/process/respsender"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p *Process) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var templateItems []model.TemplateItem
	pageId := reqreader.ReadPathParam(r, "page-id")
	reqreader.ReadBody(r, &templateItems)

	if len(templateItems) == 0 {
		errStr := "no item in request body"
		log.Println(errStr)
		respsender.ResponseString(w, `{"success": false, "code": 20000, "error": "`+errStr+`"}`, http.StatusInternalServerError)
		return
	}

	for i := range templateItems {
		templateItems[i].Id = primitive.NewObjectID()
	}

	query := bson.M{
		"page_id": pageId,
	}
	count, err := p.Template.Count(context.TODO(), query, int64(0))
	if err != nil {
		errStr := `unexpected database error occurred`
		log.Println(errStr + `:: ` + err.Error())
		respsender.ResponseString(w, `{"success": false, "code": 20401, "error": "`+errStr+`"}`, http.StatusInternalServerError)
		return
	}

	info := map[string]interface{}{
		"page_id": pageId,
	}

	if count != 0 {
		if err = p.Template.PushTemplateItem(templateItems, info); err != nil {
			respsender.ResponseString(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("create template successfully")
		data := map[string]interface{}{
			"success": true,
		}

		respsender.ResponseMap(w, data, http.StatusCreated)
		return
	}

	if _, err := p.Template.InsertTemplate(templateItems, info); err != nil {
		respsender.ResponseString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("create template successfully")
	data := map[string]interface{}{
		"success": true,
	}

	respsender.ResponseMap(w, data, http.StatusCreated)
}
