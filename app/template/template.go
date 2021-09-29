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
	FindAll(ctx context.Context, query, selectedField interface{}, limit, skip int64, sort interface{}) (model.Templates, error)
	InsertTemplate(templateItems model.Template, info map[string]interface{}) (interface{}, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, selectedField interface{}) (model.Template, error)
	UpdateTemplateItem(updatedItems model.Template, info map[string]interface{}) error
	Count(ctx context.Context, filter interface{}, skip int64) (int64, error)
	DeleteTemplate(info map[string]interface{}) error
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

func (i *Template) FindAll(ctx context.Context, query, selectedField interface{}, limit, skip int64, sort interface{}) (model.Templates, error) {
	var templates = make(model.Templates, 0)
	result, err := i.mongo.FindAll(ctx, DB, Collection, query, selectedField, limit, skip, sort)
	if err == nil {
		for _, r := range result {
			var item model.Template
			str, _ := bson.Marshal(r)
			err = bson.Unmarshal(str, &item)
			if err != nil {
				return model.Templates{}, err
			}
			templates = append(templates, item)
		}
	}
	return templates, err
}

func (t *Template) UpdateTemplateItem(updatedItem model.Template, info map[string]interface{}) error {
	var templateId = info["template_id"].(primitive.ObjectID)

	session, err := t.mongo.StartSession()
	defer session.EndSession(context.Background())
	if err != nil {
		errStr := `unexpected error`
		log.Println(errStr + `:: ` + err.Error())
		return errors.New(`{"success": false, "code": 20006, "error": "` + errStr + `"}`)
	}

	updatedItem.Id = templateId
	query := bson.M{
		"_id": templateId,
	}
	update := bson.M{
		"$set": updatedItem,
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

func (t *Template) InsertTemplate(template model.Template, info map[string]interface{}) (interface{}, error) {
	var templateId interface{}
	template.PageId = info["page_id"].(string)

	session, err := t.mongo.StartSession()
	defer session.EndSession(context.Background())
	if err != nil {
		errStr := `unexpected error`
		log.Println(errStr + `:: ` + err.Error())
		return primitive.NilObjectID, errors.New(`{"success": false, "code": 20006, "error": "` + errStr + `"}`)
	}

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
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

func (t *Template) DeleteTemplate(info map[string]interface{}) error {
	var templateId = info["template_id"].(primitive.ObjectID)

	session, err := t.mongo.StartSession()
	defer session.EndSession(context.Background())
	if err != nil {
		errStr := `unexpected error`
		log.Println(errStr + `:: ` + err.Error())
		return errors.New(`{"success": false, "code": 20006, "error": "` + errStr + `"}`)
	}

	delete := bson.M{
		"_id": templateId,
	}

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err = t.mongo.Delete(sessionContext, DB, Collection, delete); err != nil {
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

func (t Template) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, selectedField interface{}) (model.Template, error) {
	var template model.Template
	result, err := t.mongo.FindOneAndUpdate(ctx, DB, Collection, filter, update, selectedField, false)
	if err == nil {
		str, _ := bson.Marshal(result)
		err = bson.Unmarshal(str, *&template)
	}

	return template, err
}
