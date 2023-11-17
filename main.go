package main;

import (
	"scrabble_backend/router"
	"net/http"
	"log"
)

func main() {
	router, server := router.Router();
	go server.Serve();
	defer server.Close();
	log.Fatal(http.ListenAndServe(":3001", router));
}