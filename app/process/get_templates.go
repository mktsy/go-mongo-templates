package process

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"test-mongo/app/process/reqreader"
	"test-mongo/app/process/respsender"
	"test-mongo/app/utils"

	"go.mongodb.org/mongo-driver/bson"
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

	projectionOpt := bson.M{}
	sortOpt := bson.M{}
	limitOpt := int64(limitInt)
	offsetOpt := int64(offsetInt)

	template, err := p.Template.FindAll(context.TODO(), &query, projectionOpt, limitOpt, offsetOpt, sortOpt)
	if err != nil {
		errStr := `unexpected database error occurred`
		log.Println(errStr + `:: ` + err.Error())
		respsender.ResponseString(w, `{"success": false, "code": 20400, "error": "`+errStr+`"}`, http.StatusInternalServerError)
		return
	}

	count, err := p.Template.Count(context.TODO(), query, int64(0))
	if err != nil {
		errStr := `unexpected database error occurred`
		log.Println(errStr + `:: ` + err.Error())
		respsender.ResponseString(w, `{"success": false, "code": 20401, "error": "`+errStr+`"}`, http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"total":  count,
		"limit":  limitInt,
		"offset": offsetInt,
		"result": template,
	}

	if count > 0 && limitInt != 0 {
		prev, next := utils.Pagination(offsetInt, limitInt, int(count))
		if next != -999 {
			result["next"] = `` + r.URL.Path + `?limit=` + limitStr + `&offset=` + strconv.Itoa(next) + ``
		} else {
			result["next"] = ""
		}

		if prev != -999 {
			result["prev"] = `` + r.URL.Path + `?limit=` + limitStr + `&offset=` + strconv.Itoa(prev) + ``
		} else {
			result["prev"] = ""
		}

		log.Println("found template in this page_id")
	} else {
		result["next"] = ""
		result["prev"] = ""
		log.Println("no template in this page_id")
	}

	data := map[string]interface{}{
		"success": true,
		"result":  result,
	}

	respsender.ResponseMap(w, data, http.StatusOK)
}
