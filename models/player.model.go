package models;

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerModel struct{
	ObjectId primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PlayerName string `json:"playerName" bson:"playerName"`
	PlayerEmail string `json:"playerEmail" bson:"playerEmail"`
	PlayerPhotoUrl string `json:"playerPhotoUrl" bson:"playerPhotoUrl"` 
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}