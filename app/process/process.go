package process

import (
	db "test-mongo/app/storage/mongo"
	"test-mongo/app/template"

	"go.mongodb.org/mongo-driver/mongo"
)

type Process struct {
	mongo    *mongo.Client
	Template template.Templater
}

func New(dbMongo *db.Mongo) (*Process, error) {
	p := new(Process)
	p.Template = template.New(dbMongo)
	p.mongo = dbMongo.GetClient()
	return p, nil
}
