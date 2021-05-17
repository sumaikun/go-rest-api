package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Models "github.com/sumaikun/go-rest-api/models"

	Helpers "github.com/sumaikun/go-rest-api/helpers"

	C "github.com/sumaikun/go-rest-api/config"
)

//-----------------------------  Auth functions --------------------------------------------------

func authentication(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var userType int

	response := &Models.TokenResponse{Token: "", User: nil, UserType: 0}

	var creds Models.Credentials

	copyBody := r.Body

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(copyBody).Decode(&creds)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := Models.Users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || !Helpers.CheckPasswordHash(creds.Password, expectedPassword) {

		user, err := dao.FindOneByKEY("users", "email", creds.Username)

		fmt.Println("user", user)
		//fmt.Println(user)

		if user == nil {

			fmt.Println("user not found trying doctor")

			user, err = dao.FindOneByKEY("doctors", "email", creds.Username)

			//fmt.Println("user", user)

			if user == nil {

				fmt.Println("user not found trying patient")

				user, err = dao.FindOneByKEY("patients", "email", creds.Username)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if user != nil {
					userType = 3
				}

			} else {
				userType = 2
			}

			if err != nil {

				fmt.Println("err", err)

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

		} else {
			userType = 1
		}

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		match := Helpers.CheckPasswordHash(creds.Password, user.(bson.M)["password"].(string))

		if !match {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user.(bson.M)["state"] != nil && user.(bson.M)["state"].(string) != "ACTIVE" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		response.User = user.(bson.M)

	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(8 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Header().Set("Content-type", "application/json")

	//Generate json response for get the token
	response.Token = tokenString

	response.UserType = userType

	json.NewEncoder(w).Encode(response)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}

func createInititalUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	users, err := dao.FindAll("users")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if users == nil {

		var user Models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user.ID = bson.NewObjectId()
		user.Date = time.Now().String()
		user.UpdateDate = time.Now().String()

		if len(user.Password) != 0 {
			user.Password, _ = Helpers.HashPassword(user.Password)
		}

		if err := dao.Insert("users", user, []string{"email"}); err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusCreated, user)

	} else {
		Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "can not create initial users again"})
	}

}

func resetPassword(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, reset := resetPasswordValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	claims := jwt.MapClaims{}
	_, err2 := jwt.ParseWithClaims(reset.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwtkey), nil
	})

	if err2 != nil {
		Helpers.RespondWithJSON(w, http.StatusForbidden, map[string]string{"result": "Error decoding jwt"})
		//log.Fatal("Error decoding jwt")
	}

	//fmt.Println(claims)

	//claims["username"].(string)

	if claims["type"].(string) == "forgot-password" {

		user, _ := dao.FindOneByKEY("users", "email", claims["username"].(string))

		doctor, _ := dao.FindOneByKEY("doctors", "email", claims["username"].(string))

		patient, _ := dao.FindOneByKEY("contacts", "email", claims["username"].(string))

		fmt.Println("doctor", doctor)

		if user != nil {
			parsedUser := user.(bson.M)
			parsedUser["password"], _ = Helpers.HashPassword(reset.Password)
			parsedUser["state"] = "ACTIVE"
			if err := dao.Update("users", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if doctor != nil {
			parsedUser := doctor.(bson.M)
			parsedUser["password"], _ = Helpers.HashPassword(reset.Password)
			parsedUser["state"] = "ACTIVE"
			if err := dao.Update("doctors", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if patient != nil {
			parsedUser := patient.(bson.M)
			parsedUser["password"], _ = Helpers.HashPassword(reset.Password)
			parsedUser["state"] = "ACTIVE"
			if err := dao.Update("contacts", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "password reseted"})

	} else {
		user, err3 := dao.FindOneByKEY(claims["type"].(string)+"s", "email", claims["username"].(string))

		if err3 != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err3.Error())
			return
		}

		//fmt.Println(user)

		parsedUser := user.(bson.M)

		//fmt.Println(parsedUser["state"])

		if parsedUser["state"].(string) == "CHANGE_PASSWORD" {
			parsedUser["password"], _ = Helpers.HashPassword(reset.Password)
			parsedUser["state"] = "ACTIVE"
			//fmt.Println(parsedUser)
			if err := dao.Update(claims["type"].(string)+"s", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

		} else {
			Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "can't change password of this account"})
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "password reseted"})
	}

}

func forgotPassword(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, reset := forgotPasswordValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: reset.Email,
		Type:     "forgot-password",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendForgotPasswordEmail(tokenString, reset.Email)
	go Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func confirmAccount(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, reset := confirmAccountValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	claims := jwt.MapClaims{}
	_, err2 := jwt.ParseWithClaims(reset.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwtkey), nil
	})

	if err2 != nil {
		Helpers.RespondWithJSON(w, http.StatusForbidden, map[string]string{"result": "Error decoding jwt"})
		//log.Fatal("Error decoding jwt")
	}

	fmt.Println(claims)

	//claims["username"].(string)

	user, err3 := dao.FindOneByKEY(claims["type"].(string)+"s", "email", claims["username"].(string))

	if err3 != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err3.Error())
		return
	}

	//fmt.Println(user)

	parsedUser := user.(bson.M)

	fmt.Println("parsedUser", parsedUser)

	parsedUser["state"] = "ACTIVE"

	//fmt.Println(parsedUser)
	if err := dao.PartialUpdate(claims["type"].(string)+"s", parsedUser["_id"].(bson.ObjectId).Hex(), bson.M{"state": "ACTIVE"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "Account confirmed"})

}

func registerDoctor(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, user := userRegisterValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	user.ID = bson.NewObjectId()
	user.Date = time.Now().String()
	user.UpdateDate = time.Now().String()
	user.State = "INACTIVE"
	user.Password, _ = Helpers.HashPassword(user.Password)

	if err := dao.Insert("doctors", user, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: user.Email,
		Type:     "doctor",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendConfirmationEmail(tokenString, user.Email)

	go Helpers.RespondWithJSON(w, http.StatusCreated, user)

}

func registerContact(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, user := userRegisterValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	user.ID = bson.NewObjectId()
	user.Date = time.Now().String()
	user.UpdateDate = time.Now().String()
	user.State = "INACTIVE"
	user.Password, _ = Helpers.HashPassword(user.Password)

	if err := dao.Insert("contacts", user, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: user.Email,
		Type:     "patient",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendConfirmationEmail(tokenString, user.Email)

	go Helpers.RespondWithJSON(w, http.StatusCreated, user)
}

//-----------------------------  Users functions --------------------------------------------------

func allUsersEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	users, err := dao.FindAll("users")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, users)
}

func createUsersEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	usera := context.Get(r, "user")

	userParsed := usera.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, user := userValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	user.ID = bson.NewObjectId()
	user.Date = time.Now().String()
	user.UpdateDate = time.Now().String()
	user.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	user.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if len(user.Password) != 0 {
		user.Password, _ = Helpers.HashPassword(user.Password)
	}

	if err := dao.Insert("users", user, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, user)

}

func findUserEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	user, err := dao.FindByID("users", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, user)

}

func removeUserEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("users", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateUserEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	usera := context.Get(r, "user")

	userParsed := usera.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, user := userValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevUser, err2 := dao.FindByID("users", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	parsedData := prevUser.(bson.M)

	user.ID = parsedData["_id"].(bson.ObjectId)

	user.Date = parsedData["date"].(string)

	user.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		user.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		user.CreatedBy = parsedData["createdBy"].(string)
	}

	user.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if len(user.Password) == 0 {
		user.Password = parsedData["password"].(string)
	} else {
		user.Password, _ = Helpers.HashPassword(user.Password)
	}

	if err := dao.Update("users", user.ID, user); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------------- Products Functions ----------------------------------

func allProductsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	products, err := dao.FindAll("products")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, products)
}

func createProductEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, product := productValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	product.ID = bson.NewObjectId()
	product.Date = time.Now().String()
	product.UpdateDate = time.Now().String()
	product.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	product.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("products", product, []string{"name"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, product)

}

func findProductEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	product, err := dao.FindByID("products", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, product)

}

func removeProductEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("products", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateProductEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, product := productValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("products", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	parsedData := prevData.(bson.M)

	product.ID = parsedData["_id"].(bson.ObjectId)

	product.Date = parsedData["date"].(string)

	product.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		product.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		product.CreatedBy = parsedData["createdBy"].(string)
	}

	product.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("products", product.ID, product); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------------- Contacts functions ----------------------------------

func allContactsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	contacts, err := dao.FindAll("contacts")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, contacts)
}

func createContactEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, contact := contactValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	contact.ID = bson.NewObjectId()
	contact.Date = time.Now().String()
	contact.UpdateDate = time.Now().String()
	contact.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	contact.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("contacts", contact, []string{"name", "identification", "email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, contact)

}

func findContactEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	contact, err := dao.FindByID("contacts", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Contact ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, contact)

}

func removeContactEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("contacts", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Contact ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateContactEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, contact := contactValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("contacts", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Contact ID")
		return
	}

	parsedData := prevData.(bson.M)

	contact.ID = parsedData["_id"].(bson.ObjectId)

	contact.Date = parsedData["date"].(string)

	contact.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		contact.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		contact.CreatedBy = parsedData["createdBy"].(string)
	}

	contact.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("contacts", contact.ID, contact); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------------- Pet functions ----------------------------------

func allPetsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	pets, err := dao.FindAll("pets")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, pets)
}

func createPetEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, pet := petValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	pet.ID = bson.NewObjectId()
	pet.Date = time.Now().String()
	pet.UpdateDate = time.Now().String()
	pet.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	pet.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("pets", pet, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, pet)

}

func findPetEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("pets", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Pet ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePetEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("pets", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Pet ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePetEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, pet := petValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("pets", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Pet ID")
		return
	}

	parsedData := prevData.(bson.M)

	pet.ID = parsedData["_id"].(bson.ObjectId)

	pet.Date = parsedData["date"].(string)

	pet.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		pet.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		pet.CreatedBy = parsedData["createdBy"].(string)
	}

	pet.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("pets", pet.ID, pet); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------------- Parameters Functions --------------------------------

func createParameterEndPoint(w http.ResponseWriter, r *http.Request) {

	entity := strings.Replace(r.URL.Path, "/", "", -1)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, parameter, uniqueKeys := validatorSelector(r, entity)

	//fmt.Println(parameter)

	if len(err["validationError"].(url.Values)) > 0 {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := dao.Insert(entity, parameter, uniqueKeys); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, parameter)

}

func allParametersEndPoint(w http.ResponseWriter, r *http.Request) {
	entity := strings.Replace(r.URL.Path, "/", "", -1)
	w.Header().Set("Content-type", "application/json")

	parameters, err := dao.FindAll(entity)
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, parameters)

}

func findParameterEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	entity := strings.Replace(r.URL.Path, "/"+params["id"], "", -1)
	entity = strings.Replace(entity, "/", "", -1)

	parameter, err := dao.FindByID(entity, params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Parameter ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, parameter)
}

func deleteParameterEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	entity := strings.Replace(r.URL.Path, "/"+params["id"], "", -1)
	entity = strings.Replace(entity, "/", "", -1)
	err := dao.DeleteByID(entity, params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Parameter ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)
}

func updateParameterEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)
	entity := strings.Replace(r.URL.Path, "/"+params["id"], "", -1)
	entity = strings.Replace(entity, "/", "", -1)
	w.Header().Set("Content-type", "application/json")

	prevData, err2 := dao.FindByID(entity, params["id"])

	if err2 != nil {
		fmt.Println(err2)
		fmt.Println(params["id"])
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Parameter ID")
		return
	}

	parsedData := prevData.(bson.M)

	err, data, dataID := validatorSelectorUpdate(r, entity, parsedData)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := dao.Update(entity, dataID, data); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, "invalid")
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------------- file Upload -----------------------------------------

func fileUpload(w http.ResponseWriter, r *http.Request) {

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	var extension = filepath.Ext(handler.Filename)

	fmt.Printf("Extension: %+v\n", extension)

	tempFile, err := ioutil.TempFile("files", "upload-*"+extension)

	if err != nil {
		fmt.Println(err)
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
	}

	var tempPath = tempFile.Name()

	fmt.Println("temp file before trim" + tempPath)

	var tempName = strings.Replace(tempPath, "files/", "", -1)

	fmt.Println("tempName " + tempName)

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"filename": tempName})

}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var fileName = params["file"]

	var err = os.Remove("./files/" + fileName)
	if err != nil {
		//log.Fatal(err) // perhaps handle this nicer
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "fileDelete"})
	return

}

func serveImage(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var fileName = params["image"]

	if !strings.Contains(fileName, "png") && !strings.Contains(fileName, "jpg") && !strings.Contains(fileName, "jpeg") && !strings.Contains(fileName, "gif") {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "invalid file extension"})
		return
	}

	img, err := os.Open("./files/" + params["image"])
	if err != nil {
		//log.Fatal(err) // perhaps handle this nicer
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, img)

}

func downloadFile(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var fileName = params["file"]

	/*fmt.Println("fileName " + fileName)

	download, err := os.Open("./files/upload-815043770.pdf")

	if err != nil {

		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	defer download.Close()

	contentType, err := getFileContentType(download)

	if err != nil {
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Println("detected contentType", contentType)

	w.Header().Set("Content-Type", "application/pdf")

	w.Header().Set("Content-Disposition: attachment", "filename=test.pdf")

	_, err = io.Copy(w, download)*/

	http.ServeFile(w, r, "./files/"+fileName)
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// Enums --------------------------------------------------------------------

func userRoles(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [3]string{"admin", "assistant"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func contactStratus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [6]string{"estrato 1", "estrato 2", "estrato 3", "estrato 4", "estrato 5", "estrato 6"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func contactDocumentType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [4]string{"CC", "CE", "Pasaporte", "TI"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func parametersType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := []string{"Tipo de especialización", "Ciudades de atención"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func administrationWayType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [7]string{"Oral", "Intravenosa", "Intramuscular", "Subcutanea", "tópica", "rectal", "inhalatoria"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func presentationType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [7]string{"Jarabes", "Gotas", "Capsulas", "Polvo", "Granulado", "Emulsión", "Bebible"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

//-------------------------------------- Patient Review functions ----------------------------------

func allPatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	patientReviews, err := dao.FindAllWithUsers("patientReviews")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientReviews)
}

func findPatientReviewByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	patientReviews, err := dao.FindManyByKey("patientReviews", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientReviews)

}

func createPatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {

	//fmt.Print("here go the creation of patient review")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()

	w.Header().Set("Content-type", "application/json")

	err, patientReview := patientReviewValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	patientReview.ID = bson.NewObjectId()
	patientReview.Date = time.Now().String()
	patientReview.UpdateDate = time.Now().String()
	patientReview.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	patientReview.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("patientReviews", patientReview, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, patientReview)

}

func findPatientReviewEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("patientReview", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient Review ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePatientReviewEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("patientReview", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient Review ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, patientReview := patientReviewValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("patientReviews", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient Review ID")
		return
	}

	parsedData := prevData.(bson.M)

	patientReview.ID = parsedData["_id"].(bson.ObjectId)

	patientReview.Date = parsedData["date"].(string)

	patientReview.CreatedBy = parsedData["createdBy"].(string)

	patientReview.UpdateDate = time.Now().String()

	patientReview.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("patientReviews", patientReview.ID, patientReview); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//--------------------------------physiological Constants functions ----------------------------------

func allPhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	physiologicalConstant, err := dao.FindAllWithUsers("physiologicalConstants")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, physiologicalConstant)
}

func findPhysiologicalConstantsByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	physiologicalConstant, err := dao.FindManyByKey("physiologicalConstants", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, physiologicalConstant)

}

func createPhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, physiologicalConstant := physiologicalConstantsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	physiologicalConstant.ID = bson.NewObjectId()
	physiologicalConstant.Date = time.Now().String()
	physiologicalConstant.UpdateDate = time.Now().String()
	physiologicalConstant.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	physiologicalConstant.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("physiologicalConstants", physiologicalConstant, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, physiologicalConstant)

}

func findPhysiologicalConstantsEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("physiologicalConstants", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Constant ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePhysiologicalConstantsEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("physiologicalConstants", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Constant ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, physiologicalConstant := physiologicalConstantsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("physiologicalConstants", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Constant ID")
		return
	}

	parsedData := prevData.(bson.M)

	physiologicalConstant.ID = parsedData["_id"].(bson.ObjectId)

	physiologicalConstant.Date = parsedData["date"].(string)

	physiologicalConstant.UpdateDate = time.Now().String()

	physiologicalConstant.CreatedBy = parsedData["createdBy"].(string)

	physiologicalConstant.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("physiologicalConstants", physiologicalConstant.ID, physiologicalConstant); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//--------------------------------Diagnostic Plans functions ----------------------------------

func allDiagnosticPlansEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	diagnosticPlan, err := dao.FindAll("diagnosticPlans")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, diagnosticPlan)
}

func findDiagnosticPlansByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	diagnosticPlans, err := dao.FindManyByKey("diagnosticPlans", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, diagnosticPlans)

}

func createDiagnosticPlansEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, diagnosticPlan := diagnosticPlansValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	diagnosticPlan.ID = bson.NewObjectId()
	diagnosticPlan.Date = time.Now().String()
	diagnosticPlan.UpdateDate = time.Now().String()
	diagnosticPlan.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	diagnosticPlan.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("diagnosticPlans", diagnosticPlan, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, diagnosticPlan)

}

func findDiagnosticPlansEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("diagnosticPlans", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Diagnostic Plan ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removeDiagnosticPlansEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("diagnosticPlans", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Diagnostic Plan ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateDiagnosticPlansEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, diagnosticPlan := diagnosticPlansValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("diagnosticPlans", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Diagnostic Plan ID")
		return
	}

	parsedData := prevData.(bson.M)

	diagnosticPlan.ID = parsedData["_id"].(bson.ObjectId)

	diagnosticPlan.Date = parsedData["date"].(string)

	diagnosticPlan.UpdateDate = time.Now().String()

	diagnosticPlan.CreatedBy = parsedData["createdBy"].(string)

	diagnosticPlan.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("diagnosticPlans", diagnosticPlan.ID, diagnosticPlan); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//--------------------------------TherapeuticPlans functions ----------------------------------

func allTherapeuticPlansEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	therapeuticPlan, err := dao.FindAll("therapeuticPlans")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, therapeuticPlan)
}

func findTherapeuticPlansByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	therapeuticPlans, err := dao.FindManyByKey("therapeuticPlans", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, therapeuticPlans)

}

func createTherapeuticPlansEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, therapeuticPlan := therapeuticPlansValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	therapeuticPlan.ID = bson.NewObjectId()
	therapeuticPlan.Date = time.Now().String()
	therapeuticPlan.UpdateDate = time.Now().String()
	therapeuticPlan.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	therapeuticPlan.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("therapeuticPlans", therapeuticPlan, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, therapeuticPlan)

}

func findTherapeuticPlansEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("therapeuticPlans", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Therapuetic Plan ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removeTherapeuticPlansEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("therapeuticPlans", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Therapuetic Plan ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateTherapeuticPlansEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, therapeuticPlan := therapeuticPlansValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("therapeuticPlans", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Therapuetic Plan ID")
		return
	}

	parsedData := prevData.(bson.M)

	therapeuticPlan.ID = parsedData["_id"].(bson.ObjectId)

	therapeuticPlan.Date = parsedData["date"].(string)

	therapeuticPlan.UpdateDate = time.Now().String()

	therapeuticPlan.CreatedBy = parsedData["createdBy"].(string)

	therapeuticPlan.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("therapeuticPlans", therapeuticPlan.ID, therapeuticPlan); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//--------------------------------Appointments functions ----------------------------------

func allAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	diagnosticPlan, err := dao.FindAllWithPatients("appointments")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, diagnosticPlan)
}

func findAppointmentsByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	appointments, err := dao.FindManyByKey("appointments", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, appointments)

}

func createAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, appointment := appointmentsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	appointment.ID = bson.NewObjectId()
	appointment.Date = time.Now().String()
	appointment.UpdateDate = time.Now().String()
	appointment.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	appointment.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("appointments", appointment, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, appointment)

}

func findAppointmentsEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("appointments", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Appointment ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removeAppointmentsEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("appointments", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Appointment ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, appointment := appointmentsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("appointments", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Appointment ID")
		return
	}

	parsedData := prevData.(bson.M)

	appointment.ID = parsedData["_id"].(bson.ObjectId)

	appointment.Date = parsedData["date"].(string)

	appointment.UpdateDate = time.Now().String()

	appointment.CreatedBy = parsedData["createdBy"].(string)

	appointment.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("appointments", appointment.ID, appointment); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//--------------------------------Appointments functions ----------------------------------

func allDetectedDiseasesEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	detectedDiseases, err := dao.FindAll("detectedDiseases")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, detectedDiseases)
}

func findDetectedDiseasesByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	detectedDiseases, err := dao.FindManyByKey("detectedDiseases", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, detectedDiseases)

}

func createDetectedDiseaseEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, detectedDisease := detectedDiseasesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	detectedDisease.ID = bson.NewObjectId()
	detectedDisease.Date = time.Now().String()
	detectedDisease.UpdateDate = time.Now().String()
	detectedDisease.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	detectedDisease.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("detectedDiseases", detectedDisease, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, detectedDisease)

}

func findDetectedDiseaseEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("detectedDiseases", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Detected Disease ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removeDetectedDiseaseEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("detectedDiseases", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Detected Disease ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateDetectedDiseaseEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, detectedDisease := detectedDiseasesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("detectedDisease", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Detected Disease ID")
		return
	}

	parsedData := prevData.(bson.M)

	detectedDisease.ID = parsedData["_id"].(bson.ObjectId)

	detectedDisease.Date = parsedData["date"].(string)

	detectedDisease.UpdateDate = time.Now().String()

	detectedDisease.CreatedBy = parsedData["createdBy"].(string)

	detectedDisease.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("appointments", detectedDisease.ID, detectedDisease); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------- PatientFiles functions ----------------------------------

func allPatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	patientFiles, err := dao.FindAllWithPatients("patientFiles")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientFiles)
}

func findPatientFilesByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	//fmt.Println("patient log" + params["patient"])

	patientFiles, err := dao.FindManyByKey("patientFiles", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientFiles)

}

func createPatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, patientsFiles := patientsFilesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	patientsFiles.ID = bson.NewObjectId()
	patientsFiles.Date = time.Now().String()
	patientsFiles.UpdateDate = time.Now().String()
	patientsFiles.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	patientsFiles.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("patientFiles", patientsFiles, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, patientsFiles)

}

func findPatientFilesEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("patientFiles", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid PatientsFile ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePatientFilesEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("patientFiles", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid PatientsFile ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	fmt.Println("update log")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, patientsFiles := patientsFilesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("patientFiles", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patients File ID")
		return
	}

	parsedData := prevData.(bson.M)

	patientsFiles.ID = parsedData["_id"].(bson.ObjectId)

	patientsFiles.Date = parsedData["date"].(string)

	patientsFiles.UpdateDate = time.Now().String()

	patientsFiles.CreatedBy = parsedData["createdBy"].(string)

	patientsFiles.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("patientFiles", patientsFiles.ID, patientsFiles); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------- AgendaAnnotations functions ----------------------------------

func allAgendaAnnotationsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	agendaAnnotations, err := dao.FindAllWithPatients("agendaAnnotations")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, agendaAnnotations)
}

func findAgendaAnnotationsByPatientEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	//fmt.Println("patient log" + params["patient"])

	agendaAnnotations, err := dao.FindManyByKey("agendaAnnotations", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, agendaAnnotations)

}

func createAgendaAnnotationEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, agendaAnnotation := agendaAnnotationValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	agendaAnnotation.ID = bson.NewObjectId()
	agendaAnnotation.Date = time.Now().String()
	agendaAnnotation.UpdateDate = time.Now().String()
	agendaAnnotation.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	agendaAnnotation.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("agendaAnnotations", agendaAnnotation, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, agendaAnnotation)

}

func findAgendaAnnotationEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agendaAnnotation, err := dao.FindByID("agendaAnnotations", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid AgendaAnnotation ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, agendaAnnotation)

}

func removeAgendaAnnotationEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("agendaAnnotations", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid AgendaAnnotation ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateAgendaAnnotationEndPoint(w http.ResponseWriter, r *http.Request) {

	//fmt.Printf("agenda update end point")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, agendaAnnotation := agendaAnnotationValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("agendaAnnotations", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid AgendaAnnotation ID")
		return
	}

	parsedData := prevData.(bson.M)

	agendaAnnotation.ID = parsedData["_id"].(bson.ObjectId)

	agendaAnnotation.Date = parsedData["date"].(string)

	agendaAnnotation.UpdateDate = time.Now().String()

	agendaAnnotation.CreatedBy = parsedData["createdBy"].(string)

	agendaAnnotation.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("agendaAnnotations", agendaAnnotation.ID, agendaAnnotation); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-----------------------------  Doctors functions --------------------------------------------------

func allDoctorsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	doctors, err := dao.FindAllWithCities("doctors")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, doctors)
}

func createDoctorsEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, doctor := doctorValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	doctor.ID = bson.NewObjectId()
	doctor.Date = time.Now().String()
	doctor.UpdateDate = time.Now().String()
	doctor.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	doctor.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	doctor.State = "CHANGE_PASSWORD"

	if err := dao.Insert("doctors", doctor, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: doctor.Email,
		Type:     "doctor",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendResetPasswordEmail(tokenString, doctor.Email)

	go Helpers.RespondWithJSON(w, http.StatusCreated, doctor)

}

func findDoctorEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	doctor, err := dao.FindByID("doctors", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, doctor)

}

func inactivateDoctorEndPoint(w http.ResponseWriter, r *http.Request) {

	/*params := mux.Vars(r)
	err := dao.DeleteByID("doctors", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)*/

	params := mux.Vars(r)

	doctor, err2 := dao.FindByID("doctors", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}

	var err error

	if doctor.(bson.M)["state"] != nil && doctor.(bson.M)["state"].(string) == "INACTIVE" {
		err = dao.PartialUpdate("doctors", params["id"], bson.M{"state": "ACTIVE"})
	} else {
		err = dao.PartialUpdate("doctors", params["id"], bson.M{"state": "INACTIVE"})
	}

	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)
}

func updateDoctorEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, doctor := doctorValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevDoctor, err2 := dao.FindByID("doctors", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}

	parsedData := prevDoctor.(bson.M)

	doctor.ID = parsedData["_id"].(bson.ObjectId)

	doctor.State = parsedData["state"].(string)

	doctor.Date = parsedData["date"].(string)

	doctor.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		doctor.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		doctor.CreatedBy = parsedData["createdBy"].(string)
	}

	doctor.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if len(doctor.Password) == 0 {
		doctor.Password = parsedData["password"].(string)
	} else {
		doctor.Password, _ = Helpers.HashPassword(doctor.Password)
	}

	if err := dao.Update("doctors", doctor.ID, doctor); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
