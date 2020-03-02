package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	middleware "github.com/sumaikun/go-rest-api/middlewares"

	Config "github.com/sumaikun/go-rest-api/config"

	Dao "github.com/sumaikun/go-rest-api/dao"
)

var (
	port   string
	jwtKey []byte
)

var dao = Dao.MongoConnector{}

func init() {

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

	/* fileUpload */

	router.Handle("/fileUpload", middleware.AuthMiddleware(http.HandlerFunc(fileUpload))).Methods("POST")
	router.HandleFunc("/serveImage/{image}", serveImage).Methods("GET")

	/* enums */
	router.Handle("/userRoles", middleware.AuthMiddleware(http.HandlerFunc(userRoles))).Methods("GET")
	router.Handle("/contactStratus", middleware.AuthMiddleware(http.HandlerFunc(contactStratus))).Methods("GET")
	router.Handle("/contactDocumentType", middleware.AuthMiddleware(http.HandlerFunc(contactDocumentType))).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
