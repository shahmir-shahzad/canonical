package controllers

import (
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

		// Uncomment the below line to incorporate the event structure
		//var event models.Event

		// Using a map[string]interface{} (instead of the Event struct) type to store an event related information
		// to provide the flexibility of different/newer fields for new/different
		// event types.
		var event map[string]interface{}

		// Reading the request body
		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading from the body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		unmarshalErr := json.Unmarshal(byteData, &event)
		if unmarshalErr != nil {
			log.Printf("Error unmarshalling the data: %v", unmarshalErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//log.Printf("Event Created is: %+v", event)

		collection := client.Database("canonical").Collection("events")
		result, insErr := collection.InsertOne(context.TODO(), event)
		if insErr != nil {
			log.Printf("Error inserting the document: %v", insErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//log.Printf("Data (New map) Inserted successfully: %+v", result)
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Invalid Object ID provided."))
			return
		}
		//log.Printf("Fetching the event with ID : %v", eventId)

		var event map[string]interface{}

		collection := client.Database("canonical").Collection("events")
		err = collection.FindOne(context.TODO(), bson.M{"_id": bson.M{"$eq": eventId}}).Decode(&event)
		if err != nil {
			log.Printf("Error in fetching and decoding the document with matching ID. Err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(event)
		if err != nil {
			log.Printf("Error in returning the fetched object. Err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	return fn

}

func FetchAllEvents(client *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		// fetching the query parameters to create a filter for events based on type
		queryType := r.URL.Query().Get("type")
		log.Printf("The query extracted from URL is: %v", queryType)

		// default filter, in case of no 'type' filter provided in the request URL
		filter := bson.M{}

		// filter based on event type if a type is provided in the query parameters
		if len(queryType) != 0 {
			filter = bson.M{"type": bson.M{"$eq": queryType}}
		}

		var allEvents []interface{}

		collection := client.Database("canonical").Collection("events")

		cursor, err := collection.Find(context.TODO(), filter)
		if err != nil {
			log.Printf("Error in finding the objects in db. Err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.TODO())
		for cursor.Next(context.TODO()) {
			var event map[string]interface{}
			cursor.Decode(&event)
			allEvents = append(allEvents, event)
		}
		json.NewEncoder(w).Encode(allEvents)
	}

	return fn

}
