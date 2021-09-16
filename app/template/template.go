package template

import (
	"context"
	"errors"
	"log"
	"test-mongo/app/config"
	"test-mongo/app/model"
	db "test-mongo/app/storage/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Templater interface {
	FindOne(ctx context.Context, query, selectedField interface{}) (model.Template, error)
	InsertTemplate(templateItems []model.TemplateItem, info map[string]interface{}) (interface{}, error)
}

var DB string

const Collection = "template"

type Template struct {
	mongo *db.Mongo
}

func New(dbMongo *db.Mongo) *Template {
	DB = config.MongoDBName
	template := &Template{
		mongo: dbMongo,
	}

	return template
}

func (t *Template) FindOne(ctx context.Context, query, selectedField interface{}) (model.Template, error) {
	var template model.Template
	result, err := t.mongo.FindOne(ctx, DB, Collection, query, selectedField)
	if err == nil {
		str, _ := bson.Marshal(result)
		err = bson.Unmarshal(str, &template)
	}

	return template, err
}

func (t *Template) InsertTemplate(templateItems []model.TemplateItem, info map[string]interface{}) (interface{}, error) {
	var templateId interface{}
	var pageId = info["page_id"].(string)
	session, err := t.mongo.StartSession()
	defer session.EndSession(context.Background())
	if err != nil {
		errStr := `unexpected error`
		log.Println(errStr + `:: ` + err.Error())
		return primitive.NilObjectID, errors.New(`{"success": false, "code": 20006, "error": "` + errStr + `"}`)
	}

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		var template = model.Template{
			PageId:    pageId,
			Templates: templateItems,
		}

		templateId, err = t.mongo.InsertOne(sessionContext, DB, Collection, template)
		if err != nil {
			errStr := `unexpected database error occurred`
			log.Println(errStr + ":: " + err.Error())
			return errors.New(`{"success": false, "code": 20012, "error": "` + errStr + `"}`)
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return primitive.NilObjectID, err
	}

	return templateId, nil
}
