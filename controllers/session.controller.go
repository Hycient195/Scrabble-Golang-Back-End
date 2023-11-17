package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scrabble_backend/helpers"
	"scrabble_backend/models"
	"scrabble_backend/types"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mu sync.Mutex

/* Option for returning updated document after update for the different collections */
var GlobalUpdateOption = options.FindOneAndUpdate().SetReturnDocument(options.After)

func CheckRoom(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	dataByte, err := io.ReadAll(r.Body)
	helpers.HandlePanicError(err)
	r.Body.Close()

	/* Unmarshalling (extracting) JSON from request body */
	var reqBody types.TCheckRoomParameters
	json.Unmarshal(dataByte, &reqBody)

	/* Checking if a session already exists in the database */
	var sessionInst models.SessionModel
	sessionFilter := bson.M{"roomId": reqBody.RoomId}
	errSes := helpers.SessionCollection.FindOne(context.Background(), sessionFilter).Decode(&sessionInst)

	/* Checking if a player already exists in the database */
	var playerInst models.PlayerModel
	playerFilter := bson.M{"playerEmail": reqBody.Email}
	errPl := helpers.PlayerCollection.FindOne(context.Background(), playerFilter).Decode(&playerInst)

	/* Obtaining The Game Parameters For The Session */
	var sessionParamInst models.SessionParameterModel
	sessParamFilter := bson.M{"_id": sessionInst.Parameters}
	helpers.SeseionParameterCollection.FindOne(context.Background(), sessParamFilter).Decode(&sessionParamInst)

	fmt.Println(sessionParamInst.Timed)

	if reqBody.Name != "" && reqBody.Email != "" && reqBody.RoomId != "" {
		mu.Lock()
		if !helpers.CheckIsEmpty(errSes) {
			if !helpers.CheckIsEmpty(errPl) {

				/* If the Player Does Not Already Exists In The Session */
				if !helpers.PlayerExistsInSession(sessionInst, playerInst.ObjectId) {
					fmt.Println("Case 1 A")

					plExistFilter := bson.M{"_id": sessionInst.ObjectID}
					plExistUpdate := bson.M{
						"$push": bson.M{
							"players": types.TPlayer{
								PlayerName:  playerInst.PlayerName,
								PlayerEmail: playerInst.PlayerEmail,
								PlayerId:    playerInst.ObjectId,
							},
						},
					}
					var updateResInst models.SessionModel
					errUp := helpers.SessionCollection.FindOneAndUpdate(context.Background(), plExistFilter, plExistUpdate, GlobalUpdateOption).Decode(&updateResInst)
					helpers.HandlePanicError(errUp)

					sendResp := types.TCheckRoomResponse{
						ScoreSheet:      sessionParamInst.ScoreSheet,
						TileRack:        sessionParamInst.TileRack,
						TileBoard:       sessionParamInst.TileBoard,
						PlayerTurn:      sessionParamInst.PlayerTurn,
						TileBag:         sessionParamInst.TileBag,
						NumberOfPlayers: sessionParamInst.NumberOfPlayers,
						Timed:           sessionParamInst.Timed,
						Time:            sessionParamInst.Time,
						Spectators:      sessionParamInst.Spectators,
						Initiator:       updateResInst.Initiator,
						Players:         updateResInst.Players,
						Parameters:      updateResInst.Parameters,
						RoomId:          updateResInst.RoomId,
						OwnPlayer: types.TPlayer{
							PlayerName:  playerInst.PlayerName,
							PlayerEmail: playerInst.PlayerEmail,
							PlayerId:    playerInst.ObjectId,
						},
					}

					json.NewEncoder(w).Encode(sendResp)
					mu.Unlock()
				} else { // if The Player Does Not Exist In the Session;
					fmt.Println("Case 1 B")
					sendResp := types.TCheckRoomResponse{
						ScoreSheet:      sessionParamInst.ScoreSheet,
						TileRack:        sessionParamInst.TileRack,
						TileBoard:       sessionParamInst.TileBoard,
						PlayerTurn:      sessionParamInst.PlayerTurn,
						TileBag:         sessionParamInst.TileBag,
						NumberOfPlayers: sessionParamInst.NumberOfPlayers,
						Timed:           sessionParamInst.Timed,
						Time:            sessionParamInst.Time,
						Spectators:      sessionParamInst.Spectators,
						Initiator:       sessionInst.Initiator,
						Players:         sessionInst.Players,
						Parameters:      sessionInst.Parameters,
						RoomId:          sessionInst.RoomId,
						OwnPlayer: types.TPlayer{
							PlayerName:  playerInst.PlayerName,
							PlayerEmail: playerInst.PlayerEmail,
							PlayerId:    playerInst.ObjectId,
						},
					}

					json.NewEncoder(w).Encode(sendResp)
					mu.Unlock()
				}
			} else {
				fmt.Println("Case 2")

				// var newPlayerInst models.PlayerModel;
				playerInsertResult, err := helpers.PlayerCollection.InsertOne(context.Background(), models.PlayerModel{
					PlayerName: reqBody.Name, PlayerEmail: reqBody.Email, PlayerPhotoUrl: reqBody.PlayerPhotoUrl,
					CreatedAt: time.Now(), UpdatedAt: time.Now(),
				})

				helpers.HandlePanicError(err)

				var newPlayerInst models.PlayerModel
				findPlayerFilter := bson.M{"_id": playerInsertResult.InsertedID}
				helpers.PlayerCollection.FindOne(context.Background(), findPlayerFilter).Decode(&newPlayerInst)

				// io.ReadAll(playerInsertResult.)
				// json.Unmarshal(playerInsertResult, &newPlayerInst);

				// fmt.Println(resy)
				// fmt.Println(playerInsertResult.InsertedID)
				fmt.Println(newPlayerInst)

				plExistFilter := bson.M{"_id": sessionInst.ObjectID}
				plExistUpdate := bson.M{
					"$push": bson.M{
						"players": types.TPlayer{
							PlayerName:  newPlayerInst.PlayerName,
							PlayerEmail: newPlayerInst.PlayerEmail,
							PlayerId:    newPlayerInst.ObjectId,
						},
					},
				}
				var updateResInst models.SessionModel
				helpers.SessionCollection.FindOneAndUpdate(context.Background(), plExistFilter, plExistUpdate, GlobalUpdateOption).Decode(&updateResInst)

				sendResp := types.TCheckRoomResponse{
					ScoreSheet:      sessionParamInst.ScoreSheet,
					TileRack:        sessionParamInst.TileRack,
					TileBoard:       sessionParamInst.TileBoard,
					PlayerTurn:      sessionParamInst.PlayerTurn,
					TileBag:         sessionParamInst.TileBag,
					NumberOfPlayers: sessionParamInst.NumberOfPlayers,
					Timed:           sessionParamInst.Timed,
					Time:            sessionParamInst.Time,
					Spectators:      sessionParamInst.Spectators,
					Initiator:       updateResInst.Initiator,
					Players:         updateResInst.Players,
					Parameters:      updateResInst.Parameters,
					RoomId:          updateResInst.RoomId,
					OwnPlayer: types.TPlayer{
						PlayerName:  newPlayerInst.PlayerName,
						PlayerEmail: newPlayerInst.PlayerEmail,
						PlayerId:    newPlayerInst.ObjectId,
					},
				}

				json.NewEncoder(w).Encode(sendResp)
				mu.Unlock()
			}
		} else {
			if !helpers.CheckIsEmpty(errPl) {
				fmt.Println("Case 3")

				/* Storing new game parameters in DB */
				sessionParamUpdateResp, _ := helpers.SeseionParameterCollection.InsertOne(context.Background(), reqBody.Parameters)
				var newSessionParamsInst models.SessionParameterModel
				helpers.SeseionParameterCollection.FindOne(context.Background(), bson.M{"_id": sessionParamUpdateResp.InsertedID}).Decode(&newSessionParamsInst)

				chatUpdateResp, _ := helpers.ChatCollection.InsertOne(context.Background(), models.ChatModel{
					Session: types.TCreateChat{
						RoomId:   reqBody.RoomId,
						Messages: []types.TMessage{},
					},
				})
				var newChatInst models.ChatModel
				helpers.ChatCollection.FindOne(context.Background(), bson.M{"_id": chatUpdateResp.InsertedID}).Decode(&newChatInst)

				/* Creating a new session */
				sessionUpdateResp, _ := helpers.SessionCollection.InsertOne(context.Background(), models.SessionModel{
					Initiator: types.TPlayer{
						PlayerName:  playerInst.PlayerName,
						PlayerEmail: playerInst.PlayerEmail,
						PlayerId:    playerInst.ObjectId,
					},
					Players: []types.TPlayer{
						types.TPlayer{
							PlayerName:  playerInst.PlayerName,
							PlayerEmail: playerInst.PlayerEmail,
							PlayerId:    playerInst.ObjectId,
						},
					},
					Parameters: newSessionParamsInst.ObjectID,
					Chat:       newChatInst.ObjectID,
					RoomId:     reqBody.RoomId,
				})

				var newSessionUpdateInst models.SessionModel
				sessionInstFilter := bson.M{"_id": sessionUpdateResp.InsertedID}
				helpers.SessionCollection.FindOne(context.Background(), sessionInstFilter).Decode(&newSessionUpdateInst)

				/* Returning the response back to the front-end */
				sendResp := types.TCheckRoomResponse{
					ScoreSheet:      newSessionParamsInst.ScoreSheet,
					TileRack:        newSessionParamsInst.TileRack,
					TileBoard:       newSessionParamsInst.TileBoard,
					PlayerTurn:      newSessionParamsInst.PlayerTurn,
					TileBag:         newSessionParamsInst.TileBag,
					NumberOfPlayers: newSessionParamsInst.NumberOfPlayers,
					Timed:           newSessionParamsInst.Timed,
					Time:            newSessionParamsInst.Time,
					Spectators:      newSessionParamsInst.Spectators,
					Initiator:       newSessionUpdateInst.Initiator,
					Players:         newSessionUpdateInst.Players,
					Parameters:      newSessionUpdateInst.Parameters,
					RoomId:          newSessionUpdateInst.RoomId,
					OwnPlayer: types.TPlayer{
						PlayerName:  playerInst.PlayerName,
						PlayerEmail: playerInst.PlayerEmail,
						PlayerId:    playerInst.ObjectId,
					},
				}

				json.NewEncoder(w).Encode(sendResp)
				mu.Unlock()

			} else {
				fmt.Println("Case 4")

				// var newPlayerInst models.PlayerModel;
				playerInsertResult, err := helpers.PlayerCollection.InsertOne(context.Background(), models.PlayerModel{
					PlayerName: reqBody.Name, PlayerEmail: reqBody.Email, PlayerPhotoUrl: reqBody.PlayerPhotoUrl,
					CreatedAt: time.Now(), UpdatedAt: time.Now(),
				})
				helpers.HandlePanicError(err)

				var newPlayerInst models.PlayerModel
				findPlayerFilter := bson.M{"_id": playerInsertResult.InsertedID}
				helpers.PlayerCollection.FindOne(context.Background(), findPlayerFilter).Decode(&newPlayerInst)

				/* Storing new game parameters in DB */
				sessionParamUpdateResp, _ := helpers.SeseionParameterCollection.InsertOne(context.Background(), reqBody.Parameters)
				var newSessionParamsInst models.SessionParameterModel
				helpers.SeseionParameterCollection.FindOne(context.Background(), bson.M{"_id": sessionParamUpdateResp.InsertedID}).Decode(&newSessionParamsInst)

				chatUpdateResp, _ := helpers.ChatCollection.InsertOne(context.Background(), models.ChatModel{
					Session: types.TCreateChat{
						RoomId:   reqBody.RoomId,
						Messages: []types.TMessage{},
					},
				})
				var newChatInst models.ChatModel
				helpers.ChatCollection.FindOne(context.Background(), bson.M{"_id": chatUpdateResp.InsertedID}).Decode(&newChatInst)

				/* Creating a new session */
				sessionUpdateResp, _ := helpers.SessionCollection.InsertOne(context.Background(), models.SessionModel{
					Initiator: types.TPlayer{
						PlayerName:  newPlayerInst.PlayerName,
						PlayerEmail: newPlayerInst.PlayerEmail,
						PlayerId:    newPlayerInst.ObjectId,
					},
					Players: []types.TPlayer{
						types.TPlayer{
							PlayerName:  newPlayerInst.PlayerName,
							PlayerEmail: newPlayerInst.PlayerEmail,
							PlayerId:    newPlayerInst.ObjectId,
						},
					},
					Parameters: newSessionParamsInst.ObjectID,
					Chat:       newChatInst.ObjectID,
					RoomId:     reqBody.RoomId,
				})

				var newSessionUpdateInst models.SessionModel
				sessionInstFilter := bson.M{"_id": sessionUpdateResp.InsertedID}
				helpers.SessionCollection.FindOne(context.Background(), sessionInstFilter).Decode(&newSessionUpdateInst)

				/* Returning the response back to the front-end */
				sendResp := types.TCheckRoomResponse{
					ScoreSheet:      newSessionParamsInst.ScoreSheet,
					TileRack:        newSessionParamsInst.TileRack,
					TileBoard:       newSessionParamsInst.TileBoard,
					PlayerTurn:      newSessionParamsInst.PlayerTurn,
					TileBag:         newSessionParamsInst.TileBag,
					NumberOfPlayers: newSessionParamsInst.NumberOfPlayers,
					Timed:           newSessionParamsInst.Timed,
					Time:            newSessionParamsInst.Time,
					Spectators:      newSessionParamsInst.Spectators,
					Initiator:       newSessionUpdateInst.Initiator,
					Players:         newSessionUpdateInst.Players,
					Parameters:      newSessionUpdateInst.Parameters,
					RoomId:          newSessionUpdateInst.RoomId,
					OwnPlayer: types.TPlayer{
						PlayerName:  newPlayerInst.PlayerName,
						PlayerEmail: newPlayerInst.PlayerEmail,
						PlayerId:    newPlayerInst.ObjectId,
					},
				}

				json.NewEncoder(w).Encode(sendResp)
				mu.Unlock()
			}
		}
	}
}

/* ================================================== */
/* Handler To Obtain All Running Sessions of A Player */
/* ================================================== */
func GetRunningSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	filter := bson.M{
		"players": bson.M{
			"$elemMatch": bson.M{
				"playerEmail": params["email"],
			},
		},
	}

	cursor, err := helpers.SessionCollection.Find(context.Background(), filter)
	helpers.HandlePanicError(err)
	defer cursor.Close(context.Background())

	var sessions []primitive.M
	for cursor.Next(context.Background()) {
		var sessionInst bson.M
		err := cursor.Decode(&sessionInst)
		helpers.HandlePanicError(err)
		sessions = append(sessions, sessionInst)
	}

	json.NewEncoder(w).Encode(sessions)
}

/* ================================================ */
/* Handler To Ascertain If A Session Already Exists */
/* ================================================ */
func CheckSessionId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var sessionInst models.SessionModel
	filter := bson.M{"roomId": params["roomId"]}
	err := helpers.SessionCollection.FindOne(context.Background(), filter).Decode(&sessionInst)

	// fmt.Println(sessionInst);
	if helpers.CheckIsEmpty(err) {
		w.Write([]byte(`{"response": false}`))
	} else {
		w.Write([]byte(`{"response": true}`))
	}
}

/* ======================================== */
/* Handler To Submit Play For A Single Turn */
/* ======================================== */
func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type TBody struct {
		RoomId           string                       `json:"roomId"`
		SessionParmeters models.SessionParameterModel `json:"sessionParameters"`
	}

	var rBodyInst TBody
	json.NewDecoder(r.Body).Decode(&rBodyInst)
	r.Body.Close()

	/* Fetching the session with the provided roomId */
	var sessionInst models.SessionModel
	sessInstFilter := bson.M{"roomId": rBodyInst.RoomId}
	sessErr := helpers.SessionCollection.FindOne(context.Background(), sessInstFilter).Decode(&sessionInst)
	helpers.HandlePanicError(sessErr)

	var restInst models.SessionParameterModel

	updateFilter := bson.M{"_id": sessionInst.Parameters}
	update := bson.M{"$set": models.SessionParameterModel{
		ScoreSheet:             rBodyInst.SessionParmeters.ScoreSheet,
		TileRack:               rBodyInst.SessionParmeters.TileRack,
		TileBoard:              rBodyInst.SessionParmeters.TileBoard,
		PlayerTurn:             rBodyInst.SessionParmeters.PlayerTurn,
		TileBag:                rBodyInst.SessionParmeters.TileBag,
		NumberOfPlayers:        rBodyInst.SessionParmeters.NumberOfPlayers,
		PointerDirectionOffset: rBodyInst.SessionParmeters.PointerDirectionOffset,
		LastDroppedTiles:       rBodyInst.SessionParmeters.LastDroppedTiles,
	}}

	// upsert := true;
	// after := options.After
	// options := options.FindOneAndUpdateOptions{
	// 	ReturnDocument: &after,
	// 	Upsert: &upsert,
	// }

	err := helpers.SeseionParameterCollection.FindOneAndUpdate(context.Background(), updateFilter, update, GlobalUpdateOption).Decode(&restInst)
	helpers.HandlePanicError(err)

	json.NewEncoder(w).Encode(restInst)
}

/* ==================================== */
/* Handler For A Player To Exit Session */
/* ==================================== */
func HandleQuit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	/* Checking if a session already exists in the database */
	var sessionInst models.SessionModel
	helpers.SessionCollection.FindOne(context.Background(), bson.M{"roomId": params["roomId"]}).Decode(&sessionInst)

	/* Checking if a player already exists in the database */
	var playerInst models.PlayerModel
	helpers.PlayerCollection.FindOne(context.Background(), bson.M{"playerEmail": params["playerEmail"]}).Decode(&playerInst)

	/* Obtaining The Game Parameters For The Session */
	var sessionParamInst models.SessionParameterModel
	helpers.SeseionParameterCollection.FindOne(context.Background(), bson.M{"_id": sessionInst.Parameters}).Decode(&sessionParamInst)

	if sessionParamInst.NumberOfPlayers <= 1 && sessionParamInst.ScoreSheet != nil {
		delRes, _ := helpers.SeseionParameterCollection.DeleteOne(context.Background(), bson.M{"_id": sessionInst.Parameters})
		fmt.Println(delRes)
		helpers.SessionCollection.DeleteOne(context.Background(), bson.M{"roomId": params["roomId"]})
		chatDelFilter := bson.M{"session": bson.M{"roomId": params["roomId"]}}
		chDelRes, _ := helpers.ChatCollection.DeleteOne(context.Background(), chatDelFilter)
		fmt.Println(chDelRes)
		w.Write([]byte(`
			{"status": true, "action": "GAME_ENDED", "shouldTimerBeRestarted": false, "message": "Game session sucessfully closed"}
		`))
	} else {
		newScoreSheet := append(sessionParamInst.ScoreSheet[:helpers.ConvertToNum(params["playerIndex"])], sessionParamInst.ScoreSheet[(helpers.ConvertToNum(params["playerIndex"])+1):]...)
		newTileRack := append(sessionParamInst.TileRack[:helpers.ConvertToNum(params["playerIndex"])], sessionParamInst.TileRack[(helpers.ConvertToNum(params["playerIndex"])+1):]...)

		shouldTimerBeRestarted := (helpers.ConvertToNum(params["playerIndex"]) == sessionParamInst.PlayerTurn)
		var newPlayerTurn = sessionParamInst.PlayerTurn
		if sessionParamInst.PlayerTurn == (sessionParamInst.NumberOfPlayers - 1) {
			newPlayerTurn = 0
		}

		var updatedSession models.SessionModel
		sessionUpdate := bson.M{"$pull": bson.M{"players": bson.M{"playerEmail": params["playerEmail"]}}}
		helpers.SessionCollection.FindOneAndUpdate(context.Background(), bson.M{"roomId": params["roomId"]}, sessionUpdate, GlobalUpdateOption).Decode(&updatedSession)

		sessionParamUpdate := bson.M{
			"$set": bson.M{
				"scoresheet":      newScoreSheet,
				"tileRack":        newTileRack,
				"numberOfPlayers": (sessionParamInst.NumberOfPlayers - 1),
				"playerTurn":      newPlayerTurn,
			},
		}
		var updatedSessionParams models.SessionParameterModel
		helpers.SeseionParameterCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": sessionInst.Parameters}, sessionParamUpdate, GlobalUpdateOption).Decode(&updatedSessionParams)

		sendResp := types.TQuitResponse{
			Status:                 true,
			Action:                 "PLAYER_EXITED",
			Message:                "Player sucessfully exited game session",
			ShouldTimerBeRestarted: shouldTimerBeRestarted,
			Parameters: types.TCheckRoomResponse{
				ScoreSheet:      updatedSessionParams.ScoreSheet,
				TileRack:        updatedSessionParams.TileRack,
				TileBoard:       updatedSessionParams.TileBoard,
				PlayerTurn:      updatedSessionParams.PlayerTurn,
				TileBag:         updatedSessionParams.TileBag,
				NumberOfPlayers: updatedSessionParams.NumberOfPlayers,
				Timed:           updatedSessionParams.Timed,
				Time:            updatedSessionParams.Time,
				Spectators:      updatedSessionParams.Spectators,
				Initiator:       updatedSession.Initiator,
				Players:         updatedSession.Players,
				OwnPlayer: types.TPlayer{
					PlayerName:  playerInst.PlayerName,
					PlayerEmail: playerInst.PlayerEmail,
					PlayerId:    playerInst.ObjectId,
				},
			},
		}

		json.NewEncoder(w).Encode(sendResp)
	}
}
