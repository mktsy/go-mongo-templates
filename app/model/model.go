package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PageId   string             `json:"page_id" bson:"page_id"`
	Title    string             `json:"title" bson:"title"`
	Text     string             `json:"text" bson:"text"`
	ImageURL string             `json:"image_url" bson:"image_url"`
}
