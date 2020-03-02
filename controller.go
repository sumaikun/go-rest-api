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
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Models "github.com/sumaikun/go-rest-api/models"

	Helpers "github.com/sumaikun/go-rest-api/helpers"
)

//-----------------------------  Auth functions --------------------------------------------------

func authentication(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	response := &Models.TokenResponse{Token: "", User: nil}

	var creds Models.Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

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

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {

			match := Helpers.CheckPasswordHash(creds.Password, user.(bson.M)["password"].(string))

			if !match {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			response.User = user.(bson.M)

		}

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

	json.NewEncoder(w).Encode(response)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
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

	if err := dao.Update("pets", pet.ID, pet); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
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

	var tempName = strings.Trim(tempFile.Name(), "files/")

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

// Enums --------------------------------------------------------------------

func userRoles(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [3]string{"developer", "doctor", "assistant"}

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
