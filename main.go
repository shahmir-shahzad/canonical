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
	fmt.Println("*** Starting Server ***")

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

	log.Println("---> Testing DB: Successful")

	router := mux.NewRouter()

	// User Routes
	router.HandleFunc("/register", controllers.RegisterUserHandler(client)).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUserHandler(client)).Methods("POST")

	// Event Routes
	router.HandleFunc("/event", controllers.AuthenticateAndHandle(controllers.CreateEventHandler(client))).Methods("POST")
	router.HandleFunc("/events/{id}", controllers.AuthenticateAndHandle(controllers.FetchEvent(client))).Methods("GET")
	router.HandleFunc("/events", controllers.AuthenticateAndHandle(controllers.FetchAllEvents(client))).Methods("GET")

	// Middleware attempt
	router.HandleFunc("/", controllers.AuthenticateAndHandle(home))

	log.Println("---> Feel free to initiate the API calls now :)")

	portAddress := fmt.Sprintf(":%s", config.Server.Port)
	log.Fatal(http.ListenAndServe(portAddress, router))
}
