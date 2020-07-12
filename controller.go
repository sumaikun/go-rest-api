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

	user.CreatedBy = parsedData["createdBy"].(string)

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

	product.CreatedBy = parsedData["createdBy"].(string)

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

	contact.CreatedBy = parsedData["createdBy"].(string)

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

	pet.CreatedBy = parsedData["createdBy"].(string)

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

func parametersType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := []string{"Especie", "Raza", "Tipo de examen", "Tipo de plan", "Enfermedades"}

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

	diagnosticPlan, err := dao.FindAll("appointments")
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

func allDetectedDeseasesEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	detectedDeseases, err := dao.FindAll("detectedDeseases")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, detectedDeseases)
}

func createDetectedDeseaseEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, detectedDesease := detectedDeseasesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	detectedDesease.ID = bson.NewObjectId()
	detectedDesease.Date = time.Now().String()
	detectedDesease.UpdateDate = time.Now().String()
	detectedDesease.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	detectedDesease.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("detectedDeseases", detectedDesease, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, detectedDesease)

}

func findDetectedDeseaseEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("detectedDeseases", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Detected Desease ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removeDetectedDeseaseEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("detectedDeseases", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Detected Desease ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateDetectedDeseaseEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, detectedDesease := detectedDeseasesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("detectedDesease", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Detected Desease ID")
		return
	}

	parsedData := prevData.(bson.M)

	detectedDesease.ID = parsedData["_id"].(bson.ObjectId)

	detectedDesease.Date = parsedData["date"].(string)

	detectedDesease.UpdateDate = time.Now().String()

	detectedDesease.CreatedBy = parsedData["createdBy"].(string)

	detectedDesease.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("appointments", detectedDesease.ID, detectedDesease); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------- PatientFiles functions ----------------------------------

func allPatientsFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	patientFiles, err := dao.FindAll("patientFiles")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientFiles)
}

func createPatientsFilesEndPoint(w http.ResponseWriter, r *http.Request) {

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

	if err := dao.Insert("detectedDeseases", patientsFiles, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, patientsFiles)

}

func findPatientsFilesEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("patientsFiles", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid PatientsFile ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePatientsFilesEndpoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("patientsFiles", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid PatientsFile ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePatientsFilesEndPoint(w http.ResponseWriter, r *http.Request) {

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

	prevData, err2 := dao.FindByID("patientsFiles", params["id"])
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

	if err := dao.Update("patientsFiles", patientsFiles.ID, patientsFiles); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//-------------------------------- AgendaAnnotations functions ----------------------------------

func allAgendaAnnotationsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	agendaAnnotations, err := dao.FindAll("agendaAnnotations")
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

	if err := dao.Update("patientsFiles", agendaAnnotation.ID, agendaAnnotation); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
