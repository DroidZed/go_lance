package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	FullName string             `json:"fullName" bson:"fullName"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Photo    string             `json:"photo" bson:"photo"`
	Age      int                `json:"age" bson:"age"`
}
