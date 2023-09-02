package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event struct could be used potentially. But in order to fulfill the requirement of the Event document
// object flexible to addition of more fields without changing code;
// I have decided to use a map[string]interface{} instead.
// TODO: figure out a better way that uses structs with dynamic fields somehow. So far, the current implementation fulfills
// the purpose but doesn't feel the most efficient approach.
type Event struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type             string             `json:"type,omitempty" bson:"type,omitempty"`
	Details          string             `json:"details,omitempty" bson:"details,omitempty"`
	CreatedTimestamp string             `json:"createdTimestamp,omitempty" bson:"createdTimestamp,omitempty"`
}

// User struct to store user related info
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:",omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

// Claims object holds info related to an authenticated user and the generated JWT
// For example in this demo we are storing the expiration time of the token.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
