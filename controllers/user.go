package controllers

import (
	"canonical/models"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	jwtKey = []byte("very_secret_key")
)

func RegisterUserHandler(client *mongo.Client) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var user models.User

		// Reading the request body
		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading from the body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		unmarshalErr := json.Unmarshal(byteData, &user)
		if unmarshalErr != nil {
			log.Printf("Error unmarshalling the data: %v", unmarshalErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		encryptedPassword, err := EncryptPassword(user.Password)
		if err != nil {
			log.Printf("Unable to encrypt the password. Err : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Password = encryptedPassword

		//log.Printf("User Created is: %+v", user)

		collection := client.Database("canonical").Collection("users")
		result, insErr := collection.InsertOne(context.TODO(), user)
		if insErr != nil {
			log.Printf("Error inserting the document: %v", insErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//log.Printf("Data Inserted successfully: %+v", result)
		json.NewEncoder(w).Encode(result)

	}

	return fn
}

func LoginUserHandler(client *mongo.Client) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var user models.User

		log.Printf("The Header of the request is %+v", r.Header)

		// Reading the request body
		//_ = json.NewDecoder(r.Body).Decode(&prod)
		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading from the body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		unmarshalErr := json.Unmarshal(byteData, &user)
		if unmarshalErr != nil {
			log.Printf("Error unmarshalling the data: %v", unmarshalErr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var dbUser models.User
		receivedPassword := user.Password
		username := user.Username

		collection := client.Database("canonical").Collection("users")
		err = collection.FindOne(context.TODO(), bson.M{"username": bson.M{"$eq": username}}).Decode(&dbUser)
		if err != nil {
			log.Printf("Unable to find the User. No such user exists Err : %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("Unable to encrypt the password. Err : %v", err)
			return
		}

		passwordVerified := CheckPasswordHash(receivedPassword, dbUser.Password)
		if !passwordVerified {
			log.Printf("Invalid Password : %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(time.Hour * 1)

		claims := models.Claims{
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwtToken, err := token.SignedString(jwtKey)
		if err != nil {
			log.Printf("Unable to sign the jwt token. Err : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		respBody := map[string]interface{}{
			"token": jwtToken,
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login Successful !"))
		json.NewEncoder(w).Encode(respBody)
	}

	return fn
}

func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
