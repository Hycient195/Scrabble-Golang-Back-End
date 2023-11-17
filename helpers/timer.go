package helpers;

import (
	"scrabble_backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"scrabble_backend/models"
	"github.com/googollee/go-socket.io"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	"sync"
	"time"
)

// Timer represents a timer for a game room.
type Timer struct {
	roomID             string
	initialTimeSeconds int
	timeSeconds        int
	timerInterval      *time.Ticker
	stopTimer          chan bool
	lock               sync.Mutex
	callback					 func(arg int)
	socketServer			 *socketio.Server
}

// NewTimer creates a new Timer instance.
func NewTimer(roomID string, initialTimeSeconds int, socketServer *socketio.Server, callback func(arg int)) *Timer {
	return &Timer{
		roomID:             roomID,
		initialTimeSeconds: initialTimeSeconds,
		timeSeconds:        initialTimeSeconds,
		stopTimer:          make(chan bool),
		callback:						callback,
		socketServer:				socketServer,
	}
}

// Start starts the timer.
func (t *Timer) Start(currentTime int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if currentTime > 0 {
		t.timeSeconds = currentTime
	}

	if t.timeSeconds > 0 && t.timerInterval == nil {
		t.timerInterval = time.NewTicker(time.Second)
		go func() {
			for {
				select {
				case <-t.timerInterval.C:
					t.timeSeconds--					
					t.callback(t.timeSeconds)
					fmt.Println(t.timeSeconds)
					if t.timeSeconds <= 0 {
						t.HandleTimeout()
					}
				case <-t.stopTimer:
					return
				}
			}
		}()
	}
}

// GetCurrentTime returns the current time in seconds.
func (t *Timer) GetCurrentTime() int {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.timeSeconds
}

// Stop stops the timer.
func (t *Timer) Stop() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.timerInterval != nil {
		t.timerInterval.Stop()
		t.timerInterval = nil
	}
}

// Pause pauses the timer (same as stopping it).
func (t *Timer) Pause() {
	t.Stop()
}

// Resume resumes the timer.
func (t *Timer) Resume() {
	t.Start(t.timeSeconds)
}

// Reset resets the timer to its initial time.
func (t *Timer) Reset() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Stop()
	t.timeSeconds = t.initialTimeSeconds
}

// Restart restarts the timer.
func (t *Timer) Restart() {
	t.Reset()
	t.Start(t.initialTimeSeconds);
}

// HandleTimeout handles timer timeout.
func (t *Timer) HandleTimeout() {
	fmt.Println("Timer has timed out") // Replace with desired logic

	/* Checking if a session already exists in the database */
	var sessionInst models.SessionModel;
	SessionCollection.FindOne(context.Background(), bson.M{"roomId": t.roomID}).Decode(&sessionInst)
	
	/* Obtaining The Game Parameters For The Session */
	var sessionParamInst models.SessionParameterModel;
	SeseionParameterCollection.FindOne(context.Background(), bson.M{"_id": sessionInst.Parameters}).Decode(&sessionParamInst);

	updatedScoreSheet := sessionParamInst.ScoreSheet;

	updatedScoreSheet[sessionParamInst.PlayerTurn][len(updatedScoreSheet[sessionParamInst.PlayerTurn]) - 1] = []types.ScoreEntry{types.ScoreEntry{ Word: "~timeout", Score: 0 }};
	updatedScoreSheet[sessionParamInst.PlayerTurn] = append(updatedScoreSheet[sessionParamInst.PlayerTurn], []types.ScoreEntry{});

	filter := bson.M{"_id": sessionInst.Parameters}
	var playerTurn int32 = 0;
	if sessionParamInst.PlayerTurn == (sessionParamInst.NumberOfPlayers - 1) {
			playerTurn = 0
	} else {
		playerTurn = (sessionParamInst.PlayerTurn + 1)
	}
	update := bson.M{
		"$set": bson.M{
			"playerTurn": playerTurn,
			"scoresheet": updatedScoreSheet,
			"tileRack":   sessionParamInst.TileRack,
			"tileBoard":  sessionParamInst.TileBoard,
		},
	}

	var updatedGameSession models.SessionParameterModel
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	SeseionParameterCollection.FindOneAndUpdate(context.TODO(), filter, update, option).Decode(&updatedGameSession)

	t.socketServer.BroadcastToRoom("/", t.roomID, "remotely-pass-turn", updatedGameSession)
	t.Start(t.initialTimeSeconds);
}