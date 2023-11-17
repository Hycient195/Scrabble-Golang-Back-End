package helpers

import (
	"context"
	"fmt"
	"os"

	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DB_Name = "scrabble-modern"

var PlayerCollection *mongo.Collection
var SessionCollection *mongo.Collection
var SeseionParameterCollection *mongo.Collection
var ChatCollection *mongo.Collection
var ReviewCollection *mongo.Collection

func init() {
	/* Loading Database URI From Environment Variables */

	var DB_URL_STRING string = os.Getenv("MONGODB_URI")

	// err := godotenv.Load()
	// HandlePanicError(err)

	// environment := os.Getenv("ENNVIRONMENT")

	// if environment == "development" {
	// 	DB_URL_STRING = os.Getenv("MONGODB_URL")
	// }

	// fmt.Println("DB URI is", DB_URL_STRING)

	clientOption := options.Client().ApplyURI(DB_URL_STRING)

	client, err := mongo.Connect(context.TODO(), clientOption)
	HandleFatalError(err)

	fmt.Println("Sucessfully connected to database")

	PlayerCollection = client.Database(DB_Name).Collection("players")
	ChatCollection = client.Database(DB_Name).Collection("chats")
	SessionCollection = client.Database(DB_Name).Collection("sessions")
	SeseionParameterCollection = client.Database(DB_Name).Collection("session-parameters")
	ReviewCollection = client.Database(DB_Name).Collection("reviews")

	fmt.Println("Collection instances created")
}
