package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type    string             `json:"type,omitempty" bson:"type,omitempty"`
	Details string             `json:"details,omitempty" bson:"details,omitempty"`
}

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:",omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
