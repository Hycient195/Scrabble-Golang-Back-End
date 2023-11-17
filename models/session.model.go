package models;

import (
	"scrabble_backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionModel struct{
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Initiator types.TPlayer `json:"initiator" bson:"initiator"`
	Players []types.TPlayer `json:"players" bson:"players"`
	Parameters primitive.ObjectID `json:"parameters" bson:"parameters"`
	RoomId string `json:"roomId" bson:"roomId" required:"true"`
	Chat primitive.ObjectID `json:"chat" bson:"chat"`
}