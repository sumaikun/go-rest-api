package models

import "gopkg.in/mgo.v2/bson"

//User representation on mongo
type User struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Password   string        `bson:"password" json:"password"`
	Email      string        `bson:"email" json:"email"`
	Address    string        `bson:"address" json:"address"`
	Role       string        `bson:"role" json:"role"`
	Phone      string        `bson:"phone" json:"phone"`
	Picture    string        `bson:"picture" json:"picture"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//Product representation on mongo
type Product struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	Name              string        `bson:"name" json:"name"`
	Value             string        `bson:"value" json:"value"`
	Description       string        `bson:"description" json:"description"`
	Picture           string        `bson:"picture" json:"picture"`
	Date              string        `bson:"date" json:"date"`
	UpdateDate        string        `bson:"update_date" json:"update_date"`
	AdministrationWay string        `bson:"administrationWay" json:"administrationWay"`
	Presentation      string        `bson:"presentation" json:"presentation"`
}

//Contact representation on mongo
type Contact struct {
	ID             bson.ObjectId `bson:"_id" json:"id"`
	Name           string        `bson:"name" json:"name"`
	Address        string        `bson:"address" json:"address"`
	TypeID         string        `bson:"typeId" json:"typeId"`
	Identification string        `bson:"identification" json:"identification"`
	Stratus        string        `bson:"stratus" json:"stratus"`
	City           string        `bson:"city" json:"city"`
	Phone          string        `bson:"phone" json:"phone"`
	Ocupation      string        `bson:"ocupation" json:"ocupation"`
	Email          string        `bson:"email" json:"email"`
	Picture        string        `bson:"picture" json:"picture"`
	Date           string        `bson:"date" json:"date"`
	UpdateDate     string        `bson:"update_date" json:"update_date"`
}

//Pet representation on mongo
type Pet struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Species     string        `bson:"species" json:"species"`
	Breed       string        `bson:"breed" json:"breed"`
	Color       string        `bson:"color" json:"color"`
	Sex         string        `bson:"sex" json:"sex"`
	BirthDate   string        `bson:"birthDate" json:"birthDate"`
	Age         string        `bson:"age" json:"age"`
	Origin      string        `bson:"origin" json:"origin"`
	Description string        `bson:"description" json:"description"`
	Picture     string        `bson:"picture" json:"picture"`
	Date        string        `bson:"date" json:"date"`
	UpdateDate  string        `bson:"update_date" json:"update_date"`
}

//Breeds representation on mongo
type Breeds struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Species    string        `bson:"species" json:"species"`
	Meta       string        `bson:"meta" json:"meta"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//Species representation on mongo
type Species struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Meta       string        `bson:"meta" json:"meta"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}
