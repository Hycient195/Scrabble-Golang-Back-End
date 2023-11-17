package controllers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"scrabble_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"scrabble_backend/helpers"
	"encoding/json"
	"net/http"
)

/* =========================================================== */
/* Handler To Check The Current Number Of Players In A Session */
/* =========================================================== */
func CheckPlayerCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	type TrBody struct{
		RoomId string `json:"roomId"`
		Email string `json:"email"`
	}

	/* Parsing the request body into the "rBody" structure */
	var rBody TrBody;
	json.NewDecoder(r.Body).Decode(&rBody);
	defer r.Body.Close();

	/* Fetching details of the player using email */
	var playerInst models.PlayerModel;
	filter := bson.M{"playerEmail": rBody.Email};
	playerErr := helpers.PlayerCollection.FindOne(context.Background(), filter).Decode(&playerInst);

	/* Fetching the session with the provided roomId */
	var sessionInst models.SessionModel;
	sessInstFilter := bson.M{"roomId": rBody.RoomId};
	sessErr := helpers.SessionCollection.FindOne(context.Background(), sessInstFilter).Decode(&sessionInst);

	/* If the sesssion does exists */
	if !helpers.CheckIsEmpty(sessErr) {
		var sessionParamInst models.SessionParameterModel;
		filter := bson.M{"_id": sessionInst.Parameters};
		helpers.SeseionParameterCollection.FindOne(context.Background(), filter).Decode(&sessionParamInst);

		if int32(len(sessionInst.Players)) >= sessionParamInst.NumberOfPlayers && playerInst.PlayerEmail != "" {
			if !helpers.CheckIsEmpty(playerErr) && helpers.PlayerExistsInSession(sessionInst, playerInst.ObjectId) || sessionParamInst.Spectators {
				w.Write([]byte(`{ "status": true, "message": "More players can be accomodated yo" }`));
			} else {
				w.Write([]byte(`{ "status": false, "message": "Maximum number of players reached" }`));
			}
		} else {
			w.Write([]byte(`{ "status": true, "message": "More players can be accomodated fe" }`))
		}
	}
}


/* ============================================= */
/* Handler To Get The Thubnail Images Of Players */
/* ============================================= */
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	
	cursor, err := helpers.PlayerCollection.Find(context.Background(), bson.M{});
	helpers.HandlePanicError(err);

	var urls []primitive.M;
	for cursor.Next(context.Background()) {
		var url bson.M;
		err := cursor.Decode(&url)
		helpers.HandlePanicError(err);
		urls = append(urls, url);
	}

	json.NewEncoder(w).Encode(urls);
}