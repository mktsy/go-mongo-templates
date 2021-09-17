package process

import (
	"log"
	"net/http"
	"test-mongo/app/model"
	"test-mongo/app/process/reqreader"
	"test-mongo/app/process/respsender"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p *Process) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	// var ok bool
	// var err error

	var updateTemplateItem []model.TemplateItem
	pageId := reqreader.ReadPathParam(r, "page-id")
	templateId, _ := primitive.ObjectIDFromHex(reqreader.ReadPathParam(r, "template-id"))
	reqreader.ReadBody(r, &updateTemplateItem)
	if len(updateTemplateItem) == 0 {
		errStr := "no item in request body"
		log.Println(errStr)
		respsender.ResponseString(w, `{"success": false, "code": 20000, "error": "`+errStr+`"}`, http.StatusInternalServerError)
		return
	}

	if len(updateTemplateItem) > 1 {
		errStr := "too many item in request body"
		log.Println(errStr)
		respsender.ResponseString(w, `{"success": false, "code": 20000, "error": "`+errStr+`"}`, http.StatusInternalServerError)
		return
	}

	info := map[string]interface{}{
		"page_id":     pageId,
		"template_id": templateId,
	}
	if err := p.Template.UpdateTemplateItem(updateTemplateItem[0], info); err != nil {

	}
}
