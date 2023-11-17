package router

import (
	"net/http"
	"scrabble_backend/controllers"
	"scrabble_backend/helpers"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

func Router() (*mux.Router, *socketio.Server) {
	var r *mux.Router = mux.NewRouter()

	// r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Access-Control-Allow-Origin", "")
	// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT,DELETE")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	// })

	r.HandleFunc("/check-room", controllers.CheckRoom).Methods("POST")
	r.HandleFunc("/get-running-sessions/{email}", controllers.GetRunningSessions).Methods("GET")
	r.HandleFunc("/check-session-id/{roomId}", controllers.CheckSessionId).Methods("GET")
	r.HandleFunc("/get-players", controllers.GetPlayers).Methods("GET")
	r.HandleFunc("/verify-player-count", controllers.CheckPlayerCount).Methods("POST")
	r.HandleFunc("/submit", controllers.HandleSubmit).Methods("POST")
	r.HandleFunc("/send-chat", controllers.SendChat).Methods("PATCH")
	r.HandleFunc("/get-chats/{roomId}", controllers.GetChats).Methods("GET")
	r.HandleFunc("/quit/{roomId}/{playerEmail}/{playerIndex}", controllers.HandleQuit).Methods("DELETE")
	r.HandleFunc("/add-review", controllers.AddReview).Methods("POST")
	// r.HandleFunc("/get-reviews", controllers.GetReviews).Methods("GET")

	socketServer := helpers.WebSocket()
	r.Handle("/socket.io/", socketServer)

	/* These routes are directed to the front-end to handle */
	r.HandleFunc("/setup", index)
	r.HandleFunc("/game", index)
	r.HandleFunc("/leaderboard", index)
	r.HandleFunc("/join-session/{roomId}", index)
	r.HandleFunc("/reviews", index)
	r.HandleFunc("/how-to-play", index)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/build/")))

	return r, socketServer
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/build/index.html")
}
