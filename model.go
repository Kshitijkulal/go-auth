package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name     string             `json:"name" bson:"name"`
    Email    string             `json:"email" bson:"email"`
    Password string             `json:"password" bson:"password"`
}
