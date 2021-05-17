package models

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// Credentials is the request body of credential input
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Users are test users for generate jwt token
var Users = map[string]string{
	"sumaikun": "$2a$14$6NWsioRmg3dogylbm0j3e.0RVDAN2dybn2HzecrFCNex9PPxsEJLi",
	"user2":    "$2a$14$6NWsioRmg3dogylbm0j3e.0RVDAN2dybn2HzecrFCNex9PPxsEJLi",
}

// Claims represents the struct of jwt token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// TokenResponse represents json response after succesfully auth
type TokenResponse struct {
	Token    string `json:"token"`
	User     bson.M `json:"user"`
	UserType int    `json:"userType"`
}

// JwtKey is the sample jwt secret
//var JwtKey = []byte()

// TypeClaims represents the struct of jwt token
type TypeClaims struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.StandardClaims
}

// ResetPassword is the request body of credential input
type ResetPassword struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

//ConfirmAccount body for confirm user account
type ConfirmAccount struct {
	Token string `json:"token"`
}

//ForgotPassword body for reset account
type ForgotPassword struct {
	Email string `json:"email"`
}

//UserRegister body for register user
type UserRegister struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	LastName   string        `bson:"lastName" json:"lastName"`
	Phone      string        `bson:"phone" json:"phone"`
	City       string        `bson:"city" json:"city"`
	Email      string        `bson:"email" json:"email"`
	Password   string        `bson:"password" json:"password"`
	Confirmed  bool          `bson:"confirmed" json:"confirmed"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
	State      string        `bson:"state" json:"state"`
}
