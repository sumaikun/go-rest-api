package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	middleware "github.com/sumaikun/go-rest-api/middlewares"
	Models "github.com/sumaikun/go-rest-api/models"

	Config "github.com/sumaikun/go-rest-api/config"

	Dao "github.com/sumaikun/go-rest-api/dao"
)

var (
	port   string
	jwtKey []byte
)

var dao = Dao.MongoConnector{}

//Dynamic types

var typeRegistry = make(map[string]reflect.Type)

func registerType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeRegistry[t.PkgPath()+"."+t.Name()] = t
}

func makeInstance(name string) interface{} {
	return reflect.New(typeRegistry[name]).Elem().Interface()
}

// CORSRouterDecorator applies CORS headers to a mux.Router
type CORSRouterDecorator struct {
	R *mux.Router
}

// ServeHTTP wraps the HTTP server enabling CORS headers.
// For more info about CORS, visit https://www.w3.org/TR/cors/
func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//fmt.Println("I am on serve HTTP")

	rw.Header().Set("Access-Control-Allow-Origin", "*")

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		//fmt.Println("I am in options")
		rw.WriteHeader(http.StatusOK)
		return
	}

	c.R.ServeHTTP(rw, req)
}

//-------------------

func init() {

	registerType((*Models.Breeds)(nil))
	registerType((*Models.Species)(nil))

	var config = Config.Config{}
	config.Read()
	//fmt.Println(config.Jwtkey)
	jwtKey = []byte(config.Jwtkey)
	port = config.Port

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	//initEvents()
	fmt.Println("start server in port " + port)
	router := mux.NewRouter().StrictSlash(true)

	/* Authentication */
	router.HandleFunc("/auth", authentication).Methods("POST")
	router.Handle("/exampleHandler", middleware.AuthMiddleware(http.HandlerFunc(exampleHandler))).Methods("GET")

	/* Users Routes */
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(createUsersEndPoint))).Methods("POST")
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(allUsersEndPoint))).Methods("GET")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(findUserEndpoint))).Methods("GET")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeUserEndpoint))).Methods("DELETE")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(updateUserEndPoint))).Methods("PUT")

	/* Products Routes */
	router.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(createProductEndPoint))).Methods("POST")
	router.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(allProductsEndPoint))).Methods("GET")
	router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(findProductEndpoint))).Methods("GET")
	router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeProductEndpoint))).Methods("DELETE")
	router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(updateProductEndPoint))).Methods("PUT")

	/* Contacts Routes */
	router.Handle("/contacts", middleware.AuthMiddleware(http.HandlerFunc(createContactEndPoint))).Methods("POST")
	router.Handle("/contacts", middleware.AuthMiddleware(http.HandlerFunc(allContactsEndPoint))).Methods("GET")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(findContactEndpoint))).Methods("GET")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeContactEndpoint))).Methods("DELETE")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(updateContactEndPoint))).Methods("PUT")

	/* Pets Routes */
	router.Handle("/pets", middleware.AuthMiddleware(http.HandlerFunc(createPetEndPoint))).Methods("POST")
	router.Handle("/pets", middleware.AuthMiddleware(http.HandlerFunc(allPetsEndPoint))).Methods("GET")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPetEndpoint))).Methods("GET")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePetEndpoint))).Methods("DELETE")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(updatePetEndPoint))).Methods("PUT")

	/* Breeds Routes */
	router.Handle("/breeds", middleware.AuthMiddleware(http.HandlerFunc(createParameterEndPoint))).Methods("POST")
	router.Handle("/breeds", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/breeds/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	router.Handle("/breeds/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/breeds/{id}", middleware.AuthMiddleware(http.HandlerFunc(updateParameterEndPoint))).Methods("PUT")

	/* Species Routes */
	router.Handle("/species", middleware.AuthMiddleware(http.HandlerFunc(createParameterEndPoint))).Methods("POST")
	router.Handle("/species", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/species/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	router.Handle("/species/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/species/{id}", middleware.AuthMiddleware(http.HandlerFunc(updateParameterEndPoint))).Methods("PUT")

	/* fileUpload */

	router.Handle("/fileUpload", middleware.AuthMiddleware(http.HandlerFunc(fileUpload))).Methods("POST")
	router.HandleFunc("/serveImage/{image}", serveImage).Methods("GET")

	/* enums */
	router.Handle("/userRoles", middleware.AuthMiddleware(http.HandlerFunc(userRoles))).Methods("GET")
	router.Handle("/contactStratus", middleware.AuthMiddleware(http.HandlerFunc(contactStratus))).Methods("GET")
	router.Handle("/contactDocumentType", middleware.AuthMiddleware(http.HandlerFunc(contactDocumentType))).Methods("GET")
	router.Handle("/parametersType", middleware.AuthMiddleware(http.HandlerFunc(parametersType))).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, &CORSRouterDecorator{router}))
}
