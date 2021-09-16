package mongo

import (
	"context"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongoer interface {
	GetCollectionNames(ctx context.Context, database string, filter interface{}) ([]string, error)
	InsertOne(ctx context.Context, database string, collection string, docs interface{}) (interface{}, error)
	InsertMany(ctx context.Context, database, collection string, docs []interface{}) ([]interface{}, error)
	FindOne(ctx context.Context, database, collection string, query interface{}, selectedField interface{}) (interface{}, error)
	FindAll(ctx context.Context, database, collection string, query interface{}, selectedField interface{}, limit, offset int64, sort interface{}) ([]bson.M, error)
	Update(ctx context.Context, database, collection string, filter interface{}, update interface{}) error
	UpdateMany(ctx context.Context, database, collection string, filter interface{}, update interface{}) error
	Delete(ctx context.Context, database, collection string, filter interface{}) (err error)
	GetClient() *mongo.Client
	Close(ctx context.Context)
}

type Mongo struct {
	client *mongo.Client
}

func New(url string) (*Mongo, error) {
	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{}))
	clientOptions := options.Client().ApplyURI(url).SetRegistry(rb.Build())

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &Mongo{client: client}, err
}

func (m *Mongo) GetCollectionNames(ctx context.Context, database string, filter interface{}) ([]string, error) {
	return m.client.Database(database).ListCollectionNames(ctx, filter)
}

func (m *Mongo) InsertOne(ctx context.Context, database string, collection string, docs interface{}) (interface{}, error) {
	result, err := m.client.Database(database).Collection(collection).InsertOne(ctx, docs)
	if err != nil {
		return nil, err
	}

	return result.InsertedID.(primitive.ObjectID), err
}

func (m *Mongo) InsertMany(ctx context.Context, database, collection string, docs []interface{}) ([]interface{}, error) {
	result, err := m.client.Database(database).Collection(collection).InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	return result.InsertedIDs, err
}

func (m *Mongo) FindOne(ctx context.Context, database, collection string, query interface{}, selectedField interface{}) (interface{}, error) {
	var result bson.M
	err := m.client.Database(database).Collection(collection).FindOne(ctx, query, options.FindOne().SetProjection(selectedField)).Decode(&result)
	return result, err
}

func (m *Mongo) FindAll(ctx context.Context, database, collection string, query interface{}, selectedField interface{}, limit, offset int64, sort interface{}) ([]bson.M, error) {
	var result []bson.M
	cursor, err := m.client.Database(database).Collection(collection).Find(ctx, query,
		options.Find().SetProjection(selectedField),
		options.Find().SetLimit(limit),
		options.Find().SetSkip(offset),
		options.Find().SetSort(sort))
	if err != nil {
		return make([]bson.M, 0), err
	}

	err = cursor.All(context.TODO(), &result)
	return result, err
}

func (m *Mongo) Update(ctx context.Context, database, collection string, filter interface{}, update interface{}) error {
	updateresult, err := m.client.Database(database).Collection(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updateresult.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	return err
}

func (m *Mongo) UpdateMany(ctx context.Context, database, collection string, filter interface{}, update interface{}) error {
	_, err := m.client.Database(database).Collection(collection).UpdateMany(ctx, filter, update)
	return err
}

func (m *Mongo) Delete(ctx context.Context, database, collection string, filter interface{}) (err error) {
	_, err = m.client.Database(database).Collection(collection).DeleteOne(ctx, filter)
	return
}

func (m *Mongo) GetClient() *mongo.Client {
	return m.client
}

func (m *Mongo) Close(ctx context.Context) {
	err := m.client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Mongo) StartSession() (mongo.Session, error) {
	return m.client.StartSession()
}
