package controllers;

import (
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/mux"
	"scrabble_backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"scrabble_backend/helpers"
	"fmt"
	"encoding/json"
	"scrabble_backend/types"
	"net/http"
)

/* ====================================== */
/* Handler To Send A Message In A Session */
/* ====================================== */
func SendChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	type TBody struct{
		RoomId string `json:"roomId"`
		PlayerId string `json:"playerId"`
		Text string `json:"text"`
	}

	var rBodyInst TBody;
	json.NewDecoder(r.Body).Decode(&rBodyInst);
	r.Body.Close();

	playerObjectId, err := primitive.ObjectIDFromHex(rBodyInst.PlayerId);
	helpers.HandlePanicError(err);

	var playerInst types.TPlayer;
	filter := bson.M{"_id": playerObjectId};
	helpers.PlayerCollection.FindOne(context.Background(), filter).Decode(&playerInst);

	fmt.Println(rBodyInst.RoomId);

	type TUpdate struct{
		Player types.TPlayer `json:"player"`
		Text string `json:"text"`
		CreatedAt time.Time `json:"createdAt,omitempty"`
		UpdatedAt time.Time `json:"updatedAt,omitempty"`
	}

	chatFilter := bson.M{"session.roomId": rBodyInst.RoomId};

	chatUpdate := bson.M{
    "$push": bson.M{
			"session.messages": TUpdate{
				Player: types.TPlayer{ PlayerName: playerInst.PlayerName, PlayerEmail: playerInst.PlayerEmail, PlayerId: playerObjectId },
				Text: rBodyInst.Text,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	upsert := true;
	after := options.After;
	chatOptions := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert: &upsert,
	}

	var respInst models.ChatModel;
	findErr := helpers.ChatCollection.FindOneAndUpdate(context.Background(), chatFilter, chatUpdate, &chatOptions).Decode(&respInst);
	fmt.Println(findErr);

	json.NewEncoder(w).Encode(respInst);
}

/* ===================================== */
/* Handler To Get All Chats In A Session */
/* ===================================== */
func GetChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r);

	var chatInst models.ChatModel;
	filter := bson.M{"session.roomId": params["roomId"]}
	helpers.ChatCollection.FindOne(context.Background(), filter).Decode(&chatInst);

	json.NewEncoder(w).Encode(chatInst);
}