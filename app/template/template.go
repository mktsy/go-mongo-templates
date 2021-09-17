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
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, selectedField interface{}) (model.TemplateItem, error)
	PushTemplateItem(templateItem []model.TemplateItem, info map[string]interface{}) error
	UpdateTemplateItem(updatedItems model.TemplateItem, info map[string]interface{}) error
	Count(ctx context.Context, filter interface{}, skip int64) (int64, error)
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

func (t *Template) UpdateTemplateItem(updatedItems model.TemplateItem, info map[string]interface{}) error {
	var pageId = info["page_id"].(string)
	var templateId = info["template_id"].(primitive.ObjectID)
	session, err := t.mongo.StartSession()
	defer session.EndSession(context.Background())
	if err != nil {
		errStr := `unexpected error`
		log.Println(errStr + `:: ` + err.Error())
		return errors.New(`{"success": false, "code": 20006, "error": "` + errStr + `"}`)
	}

	updatedItems.Id = templateId
	query := bson.M{
		"page_id":               pageId,
		"templates.template_id": templateId,
	}
	update := bson.M{
		"$set": bson.M{
			"templates.$": updatedItems,
		},
	}
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err = t.mongo.Update(sessionContext, DB, Collection, query, update); err != nil {
			errStr := `unexpected database error occurred`
			log.Println(errStr + `:: ` + err.Error())
			return errors.New(`{"success": false, "code": 10509, "error": "` + errStr + `"}`)
		}
		return nil
	})
	return nil
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

func (t *Template) PushTemplateItem(templateItem []model.TemplateItem, info map[string]interface{}) error {
	var pageId = info["page_id"].(string)

	query := bson.M{
		"page_id": pageId,
	}

	session, err := t.mongo.StartSession()
	defer session.EndSession(context.Background())
	if err != nil {
		errStr := `unexpected error`
		log.Println(errStr + `:: ` + err.Error())
		return errors.New(`{"success": false, "code": 20006, "error": "` + errStr + `"}`)
	}

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		pushItem := bson.M{
			"$push": bson.M{
				"templates": bson.M{
					"$each":     templateItem,
					"$position": 0,
				},
			},
		}

		if err = t.mongo.Update(sessionContext, DB, Collection, query, pushItem); err != nil {
			errStr := `unexpected database error occurred`
			log.Println(errStr + `:: ` + err.Error())
			return errors.New(`{"success": false, "code": 10509, "error": "` + errStr + `"}`)
		}
		return nil
	})

	return nil
}

func (t *Template) Count(ctx context.Context, filter interface{}, skip int64) (int64, error) {
	return t.mongo.Count(ctx, DB, Collection, filter, skip)
}

func (t Template) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, selectedField interface{}) (model.TemplateItem, error) {
	var templateItem model.TemplateItem
	result, err := t.mongo.FindOneAndUpdate(ctx, DB, Collection, filter, update, selectedField, false)
	if err == nil {
		str, _ := bson.Marshal(result)
		err = bson.Unmarshal(str, *&templateItem)
	}

	return templateItem, err
}
