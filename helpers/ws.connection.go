package helpers;

import (
	"fmt"
	// "net/http"
	"github.com/googollee/go-socket.io"
	// "github.com/googollee/go-socket.io/engineio"
	// "github.com/googollee/go-socket.io/engineio/transport"
	// "github.com/googollee/go-socket.io/engineio/transport/polling"
	// "github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// var allowOriginFunc = func(r *http.Request) bool {
// 	return true
// }

var allTimers = make(map[string]*Timer)

func WebSocket() *socketio.Server{
	// io := socketio.NewServer(&engineio.Options{
	// 	Transports: []transport.Transport{
	// 			&polling.Transport{
	// 					CheckOrigin: allowOriginFunc,
	// 			},
	// 			&websocket.Transport{
	// 					CheckOrigin: allowOriginFunc,
	// 			},
	// 	},
	// });
	io := socketio.NewServer(nil);

	io.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	});

	io.OnEvent("/", "join-room", func(s socketio.Conn, roomId string) {
		// fmt.Println("Room Joined")
		io.JoinRoom("/", roomId, s);
	})

	io.OnEvent("/", "initiate-user-online", func(s socketio.Conn, roomId string, playerId string) {
		// fmt.Println("User online initiated");
		io.BroadcastToRoom("/", roomId, "respond-user-online", playerId);
	})

	io.OnEvent("/", "initiate-user-offline", func(s socketio.Conn, roomId string, playerId string) {
		// fmt.Println("Offline initiated")
		io.BroadcastToRoom("/", roomId, "respond-user-offline", playerId);
	});

	io.OnEvent("/", "initiate-send-chat", func(s socketio.Conn, roomId string, message interface{}) {
		io.BroadcastToRoom("/", roomId, "respond-send-chat", message);
	})

	io.OnEvent("/", "initiate-user-joined", func(s socketio.Conn, roomId string) {
		// fmt.Println("User joined session")
		io.BroadcastToRoom("/", roomId, "respond-user-joined", roomId)
	});

	io.OnEvent("/", "user-action", func(s socketio.Conn, message interface{}, roomId string) {
		// fmt.Println("action carried out")
		io.BroadcastToRoom("/", roomId, "user-response", message)
	});
	
	/* Pause And Resume Events */
	io.OnEvent("/", "initiate-pause-timer", func(s socketio.Conn, roomId string) {
		runningTimer := allTimers[roomId];
		if runningTimer != nil {
			runningTimer.Pause();
			io.BroadcastToRoom("/", roomId, "respond-pause-timer")
		}
	});

	io.OnEvent("/", "initiate-resume-timer", func(s socketio.Conn, roomId string) {
		runningTimer := allTimers[roomId];
		if runningTimer != nil {
			runningTimer.Resume();
			io.BroadcastToRoom("/", roomId, "respond-resume-timer");
		}
	})

	io.OnEvent("/", "initiate-submit-pause", func(s socketio.Conn, roomId string, callback func()) {
		runningTimer := allTimers[roomId];
		if runningTimer != nil {
			runningTimer.Pause();
		}
		io.BroadcastToRoom("/", roomId, "respond-submit-pause")
	});

	io.OnEvent("/", "initiate-submit-resume", func(s socketio.Conn, roomId string, callback func()) {
		runningTimer := allTimers[roomId];
		if runningTimer != nil {
			runningTimer.Resume();
		}
		io.BroadcastToRoom("/", roomId, "respond-submit-resume")
	})

	io.OnEvent("/", "initiate-quit-pause", func(s socketio.Conn, roomId string, callback func()) {
		runningTimer := allTimers[roomId];
		if runningTimer != nil {
			runningTimer.Pause();
		}
		io.BroadcastToRoom("/", roomId, "respond-quit-pause")
	});
	// No quit result, the quit ws event decides that.

	io.OnEvent("/", "restart-timer-on-submit", func(s socketio.Conn, roomId string) {
		runningTimer := allTimers[roomId];
		if runningTimer != nil {
			runningTimer.Start(runningTimer.initialTimeSeconds);
		}
	})

	io.OnEvent("/", "initiate-count-time", func(s socketio.Conn, roomId string, initialTime int) {
		fmt.Println("Timer started for room:", roomId)
		
		// Create a new timer for the room if it doesn't exist
		if _, ok := allTimers[roomId]; !ok {
				allTimers[roomId] = NewTimer(roomId, initialTime, io, func(time int) {
					io.BroadcastToRoom("/", roomId, "respond-timer", time)
				})
		} else {
			allTimers[roomId].Resume();
		}
		allTimers[roomId].Start(initialTime);
	})

	io.OnEvent("/", "initiate-quit", func(s socketio.Conn, roomId string, parameters interface{}, action string, shouldTimerBeRestarted bool) {
		fmt.Println("Timer should be restarted");
		fmt.Println(shouldTimerBeRestarted);
		runningTimer := allTimers[roomId];
		if runningTimer != nil {
			if action ==  "GAME_ENDED" {
				runningTimer.Stop();
				delete(allTimers, roomId);
			} else if (action == "PLAYER_EXITED") {
				if shouldTimerBeRestarted {
					runningTimer.Start(runningTimer.initialTimeSeconds);
				}
			}
		}
		io.BroadcastToRoom("/", roomId, "respond-quit", parameters, action);
	})

	io.OnError("/", func(s socketio.Conn, e error) {
		// server.Remove(s.ID())
		fmt.Println("meet error:", e)
	})

	io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// Add the Remove session id. Fixed the connection & mem leak
		// io.c(s.ID())
		fmt.Println("Disconnected", reason)
	})

	return io;
}

