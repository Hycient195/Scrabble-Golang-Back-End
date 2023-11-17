package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"scrabble_backend/helpers"
	"scrabble_backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rBodyInst models.ReviewModel
	json.NewDecoder(r.Body).Decode(&rBodyInst)
	r.Body.Close()

	reviewData := models.ReviewModel{
		Name:      rBodyInst.Name,
		Email:     rBodyInst.Email,
		Title:     rBodyInst.Title,
		Text:      rBodyInst.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err := helpers.ReviewCollection.InsertOne(context.Background(), reviewData)
	helpers.HandleDisplayError(err)

	fmt.Println(result)
	fmt.Println("Review Just added")
	if result.InsertedID != nil {
		w.Write([]byte(`{ "status": true, "message": "Review sucessfully submitted" }`))
		return
	}
	w.Write([]byte(`{ "status": false, "message": "Error adding review" }`))
}

func GetReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reviews []primitive.M
	cursor, err := helpers.ReviewCollection.Find(context.Background(), bson.M{})
	helpers.HandlePanicError(err)
	for cursor.Next(context.Background()) {
		var review bson.M
		err := cursor.Decode(&review)
		helpers.HandleDisplayError(err)
		reviews = append(reviews, review)
	}
	json.NewEncoder(w).Encode(reviews)
}
