package process

import (
	"log"
	"net/http"
	"test-mongo/app/model"
	"test-mongo/app/process/reqreader"
	"test-mongo/app/process/respsender"
)

func (p *Process) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var template model.Template
	pageId := reqreader.ReadPathParam(r, "page-id")
	reqreader.ReadBody(r, &template)

	if template.Title == "" {
		errStr := "no item in request body"
		log.Println(errStr)
		respsender.ResponseString(w, `{"success": false, "code": 20000, "error": "`+errStr+`"}`, http.StatusInternalServerError)
		return
	}

	info := map[string]interface{}{
		"page_id": pageId,
	}

	if _, err := p.Template.InsertTemplate(template, info); err != nil {
		respsender.ResponseString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("create template successfully")
	data := map[string]interface{}{
		"success": true,
	}

	respsender.ResponseMap(w, data, http.StatusCreated)
}
