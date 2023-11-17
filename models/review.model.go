package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewModel struct {
	ObjectId  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	Title     string             `json:"title"`
	Text      string             `json:"text"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updateAt"`
}
