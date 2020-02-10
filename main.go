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

	log.Fatal(http.ListenAndServe(":"+port, router))
}
