package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PageId    string             `json:"page_id" bson:"page_id"`
	Templates []TemplateItem     `json:"templates" bson:"templates"`
}

type TemplateItem struct {
	Id       primitive.ObjectID `json:"template_id" bson:"template_id"`
	Title    string             `json:"title" bson:"title"`
	Text     string             `json:"text,omitempty" bson:"text,omitempty"`
	ImageURL string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
}

func (t Template) GetItem(itemId primitive.ObjectID) (TemplateItem, bool) {
	for _, item := range t.Templates {
		if item.Id == itemId {
			return item, true
		}
	}
	return TemplateItem{}, false
}
