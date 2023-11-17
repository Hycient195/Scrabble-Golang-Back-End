package helpers

import (
	"strings"
	"strconv"
	"fmt"
	"scrabble_backend/models"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func HandleFatalError(e error) {
	if e != nil {
		log.Fatal(e);
	}
}

func HandlePanicError(e error) {
	if e != nil {
		panic(e);
	}
}

func HandleDisplayError(e error) {
	if e != nil {
		fmt.Println(e);
	}
}

func CheckIsEmpty(e error) bool {
	if e != nil && e == mongo.ErrNoDocuments {
		return true;
	} else {
		return false;
	}
}

func ContainsValue(slice []interface{}, valueToCheck interface{}) bool {
	for _, element := range slice {
		if element == valueToCheck {
				return true;
		}
	}
	return false;
}

func PlayerExistsInSession(slice models.SessionModel, valueToCheck interface{}) bool {
	for _, player := range slice.Players {
		if player.PlayerId == valueToCheck {
				return true;
		}
	}
	return false;
}

func ConvertToNum(input string) int32 {
	convertedInput, convErr := strconv.ParseFloat(strings.TrimSpace(input), 64);
	if (convErr != nil) {
		fmt.Println(convErr);
	}
	return int32(convertedInput);
}