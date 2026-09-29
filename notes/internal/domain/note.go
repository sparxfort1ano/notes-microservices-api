package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Note struct {
	ID       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty" binding:"required,min=1" bson:"name"`
	Content  *string       `json:"content,omitempty" bson:"content,omitempty"`
	AuthorID int           `json:"-" bson:"author_id"`
}
