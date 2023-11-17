package models

import (
	"time"
	"scrabble_backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatModel struct{
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Session types.TCreateChat `json:"session" bson:"session"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}