package main

import (
	"canonical/controllers"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

//type Event struct {
//	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
//	Type    string             `json:"type,omitempty" bson:"type,omitempty"`
//	Details string             `json:"details,omitempty" bson:"details,omitempty"`
//}

var (
	client *mongo.Client
	//err    error
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Welcome to this demo app!"))
}

func main() {
	fmt.Println("Starting Server ***")

	// Opening config file
	configFile, err := os.Open("configuration/config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file")
	}

	defer configFile.Close()

	var config *Config

	// Creating a Yaml decoder
	decoder := yaml.NewDecoder(configFile)

	// Decoding the yaml file into a configuration struct
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding config file")
	}

	log.Printf("The config is : %+v", config)

	dbUser := config.Database.User
	dbPassword := config.Database.Password
	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.zbkg4.mongodb.net/?retryWrites=true&w=majority", dbUser, dbPassword)

	// Creating connection string for mongodb connection
	opts := options.Client().ApplyURI(connectionString)

	// Create a client and connect to the server
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("Error occurred in connecting to db: %v", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error occurred in disconnecting the datatabase: %v", err)
		}
	}()

	if err := client.Database("canonical").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatalf("Error occurred in pinging the datatabase: %v", err)
	}

	log.Println("Pinged successfully")

	router := mux.NewRouter()

	// User Routes
	router.HandleFunc("/register", controllers.RegisterUserHandler(client)).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUserHandler(client)).Methods("POST")

	// Event Routes
	router.HandleFunc("/event", controllers.AuthenticateAndHandle(controllers.CreateEventHandler(client))).Methods("POST")
	router.HandleFunc("/events/{id}", controllers.AuthenticateAndHandle(controllers.FetchEvent(client))).Methods("GET")
	router.HandleFunc("/events", controllers.AuthenticateAndHandle(controllers.FetchAllEvents(client))).Methods("GET")
	//router.HandleFunc("/events", controllers.FetchAllEvents(client)).Methods("GET")

	// Middleware attempt
	router.HandleFunc("/", controllers.AuthenticateAndHandle(home))

	log.Fatal(http.ListenAndServe(":8080", router))
}

//func CreateEvent(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//	var event Event
//	// Reading the request body
//	//_ = json.NewDecoder(r.Body).Decode(&prod)
//	byteData, err := io.ReadAll(r.Body)
//	if err != nil {
//		log.Printf("Error reading from the body: %v", err)
//	}
//	unmarshalErr := json.Unmarshal(byteData, &event)
//	if unmarshalErr != nil {
//		log.Printf("Error unmarshalling the data: %v", unmarshalErr)
//	}
//
//	log.Printf("Event Created is: %+v", event)
//
//	collection := client.Database("canonical").Collection("events")
//	result, insErr := collection.InsertOne(context.TODO(), event)
//	if insErr != nil {
//		log.Printf("Error inserting the document: %v", insErr)
//	}
//	log.Printf("Data Inserted successfully: %+v", result)
//	json.NewEncoder(w).Encode(result)
//}

//func FetchEvent(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//	vars := mux.Vars(r)
//	eventId, err := primitive.ObjectIDFromHex(vars["id"])
//	if err != nil {
//		log.Printf("Error in creating the object ID from string. Err: %v", err)
//	}
//	//filter := bson.D{"_id": bson.D{"$eq": eventId}}
//	log.Printf("Fetch the event with ID : %v", eventId)
//
//	var event Event
//
//	collection := client.Database("canonical").Collection("events")
//	err = collection.FindOne(context.TODO(), Event{ID: eventId}).Decode(&event)
//	if err != nil {
//		log.Printf("Error in fetching and decoding the document with matching ID. Err: %v", err)
//	}
//	err = json.NewEncoder(w).Encode(event)
//	if err != nil {
//		log.Printf("Error in returning the fetched object. Err: %v", err)
//	}
//}

//func FetchAllEvents(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//
//	var allEvents []interface{}
//
//	collection := client.Database("canonical").Collection("events")
//
//	cursor, err := collection.Find(context.TODO(), bson.M{})
//	if err != nil {
//		log.Printf("Error in finding the objects in db. Err: %v", err)
//	}
//	defer cursor.Close(context.TODO())
//	for cursor.Next(context.TODO()) {
//		var event Event
//		cursor.Decode(&event)
//		allEvents = append(allEvents, event)
//	}
//	json.NewEncoder(w).Encode(allEvents)
//
//}
