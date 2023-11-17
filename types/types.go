package types;

import (
	"time"
	// "scrabble_backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreEntry struct{
	Score int32 `json:"score" bson:"score"`
	Word string `json:"word" bson:"word"`
}

type TDroppedTiles struct{
	BoardIndex int32 `json:"boardIndex"`
	RackIndex int32 `json:"RackIndex"`
	Character string `json:"character"`
}

type TPlayer struct {
	PlayerName string `json:"playerName" bson:"playerName" required:"true"`
	PlayerEmail string `json:"playerEmail" bson:"playerEmail" required:"true"`
	PlayerId primitive.ObjectID `json:"playerId" bson:"playerId"`
}

type TOwnPlayer struct {
	PlayerName string `json:"playerName" bson:"playerName" required:"true"`
	PlayerEmail string `json:"playerEmail" bson:"playerEmail" required:"true"`
	PlayerPhotoUrl string `json:"playerPhotoUrl" bson:"playerPhotoUrl"`
}

type TParameter struct{
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ScoreSheet [][][]ScoreEntry `json:"scoresheet" bson:"scoresheet"`
	TileRack [][]string `json:"tileRack" bson:"tileRack"`
	TileBoard interface{} `json:"tileBoard" bson:"tileBoard"`
	PlayerTurn int32 `json:"playerTurn" bson:"playerTurn"`
	TileBag []string `json:"tileBag" bson:"tileBag"`
	NumberOfPlayers int32 `json:"numberOfPlayers" bson:"numberOfPlayers"`
	PointerDirectionOffset int32 `json:"pointerDirectionOffset" bson:"pointerDirectionOffset"`
	LastDroppedTiles []TDroppedTiles
	Timed bool `json:"timed" bson:"timed"`
	Time int32 `json:"time" bson:"time"`
	Spectators bool `json:"spectators" bson:"spectators"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TCheckRoomParameters struct{
	RoomId string `json:"roomId,omitempty"`
	Name string `json:"name"`
	Email string `json:"email"`
	PlayerPhotoUrl string `json:"playerPhotoUrl"`
	Parameters TParameter `json:"parameters,omitempty"`
}

type TCheckRoomResponse struct{
	ScoreSheet [][][]ScoreEntry `json:"scoresheet"`
	TileRack [][] string `json:"tileRack"`
	TileBoard interface{} `json:"tileBoard"`
	PlayerTurn int32 `json:"playerTurn"`
	TileBag []string `json:"tileBag"`
	NumberOfPlayers int32 `json:"numberOfPlayers"`
	Timed bool `json:"timed"`
	Time int32 `json:"time"`
	Spectators bool `json:"spectators"`
	Initiator TPlayer `json:"initiator"`
	Players []TPlayer `json:"players"`
	Parameters interface{} `json:"parameters,omitempty"`
	RoomId string `json:"roomId"`
	OwnPlayer TPlayer `json:"ownPlayer"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TSubmitResponse struct{

}

type TCreateChat struct{
	RoomId string `json:"roomId" bson:"roomId"`
	Messages []TMessage `json:"messages" bson:"messages"`
}


type TMessage struct{
	Player TPlayer `json:"player"`
	Text string `json:"text"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type TQuitResponse struct{
	Status bool `json:"status"`
	Action string `json:"action"`
	Message string `json:"message"`
	ShouldTimerBeRestarted bool `json:"shouldTimerBeRestarted"`
	Parameters TCheckRoomResponse `json:"parameters"`
}