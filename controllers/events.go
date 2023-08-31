package controllers

import (
	"canonical/models"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
)

func CreateEventHandler(client *mongo.Client) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var event models.Event
		// Reading the request body
		//_ = json.NewDecoder(r.Body).Decode(&prod)
		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading from the body: %v", err)
		}
		unmarshalErr := json.Unmarshal(byteData, &event)
		if unmarshalErr != nil {
			log.Printf("Error unmarshalling the data: %v", unmarshalErr)
		}

		log.Printf("Event Created is: %+v", event)

		collection := client.Database("canonical").Collection("events")
		result, insErr := collection.InsertOne(context.TODO(), event)
		if insErr != nil {
			log.Printf("Error inserting the document: %v", insErr)
		}
		log.Printf("Data Inserted successfully: %+v", result)
		json.NewEncoder(w).Encode(result)

	}

	return fn
}

func FetchEvent(client *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		vars := mux.Vars(r)
		eventId, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			log.Printf("Error in creating the object ID from string. Err: %v", err)
		}
		log.Printf("Fetch the event with ID : %v", eventId)

		var event models.Event

		collection := client.Database("canonical").Collection("events")
		err = collection.FindOne(context.TODO(), models.Event{ID: eventId}).Decode(&event)
		if err != nil {
			log.Printf("Error in fetching and decoding the document with matching ID. Err: %v", err)
		}
		err = json.NewEncoder(w).Encode(event)
		if err != nil {
			log.Printf("Error in returning the fetched object. Err: %v", err)
		}
	}

	return fn

}

func FetchAllEvents(client *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		var allEvents []interface{}

		collection := client.Database("canonical").Collection("events")

		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Printf("Error in finding the objects in db. Err: %v", err)
		}
		defer cursor.Close(context.TODO())
		for cursor.Next(context.TODO()) {
			var event models.Event
			cursor.Decode(&event)
			allEvents = append(allEvents, event)
		}
		json.NewEncoder(w).Encode(allEvents)
	}

	return fn

}
