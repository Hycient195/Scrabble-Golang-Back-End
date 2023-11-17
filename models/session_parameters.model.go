package models;

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"scrabble_backend/types"
)

type SessionParameterModel struct{
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ScoreSheet [][][]types.ScoreEntry `json:"scoresheet" bson:"scoresheet"`
	TileRack [][]string `json:"tileRack" bson:"tileRack"`
	TileBoard interface{} `json:"tileBoard" bson:"tileBoard"`
	PlayerTurn int32 `json:"playerTurn" bson:"playerTurn"`
	TileBag []string `json:"tileBag" bson:"tileBag"`
	NumberOfPlayers int32 `json:"numberOfPlayers" bson:"numberOfPlayers"`
	PointerDirectionOffset int32 `json:"pointerDirectionOffset" bson:"pointerDirectionOffset"`
	LastDroppedTiles []types.TDroppedTiles
	Timed bool `json:"timed" bson:"timed"`
	Time int32 `json:"time" bson:"time"`
	Spectators bool `json:"spectators" bson:"spectators"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}