package models

import "gopkg.in/mgo.v2/bson"

//User representation on mongo
type User struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Password   string        `bson:"password" json:"password"`
	Email      string        `bson:"email" json:"email"`
	Address    string        `bson:"address" json:"address"`
	Phone      string        `bson:"phone" json:"phone"`
	Picture    string        `bson:"picture" json:"picture"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}
