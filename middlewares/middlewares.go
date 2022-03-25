package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"gopkg.in/mgo.v2/bson"

	C "click-al-vet/config"

	Dao "click-al-vet/dao"

	"github.com/gorilla/context"
)

var dao = Dao.MongoConnector{}

// AuthMiddleware verify
func AuthMiddleware(next http.Handler) http.Handler {

	var config = C.Config{}
	config.Read()

	var JwtKey = []byte(config.Jwtkey)

	if len(JwtKey) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader,
			jwtmiddleware.FromParameter("token")),
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)

}

// UserMiddleware get user from request
func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var config = C.Config{}
		config.Read()

		ua := r.Header.Get("Authorization")

		ua = strings.Replace(ua, "Bearer ", "", 1)

		tokenString := ua
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Jwtkey), nil
		})
		// ... error handling

		if err != nil {
			log.Fatal("Error decoding jwt")
		}

		log.Println("claims username", claims["username"])

		var userType int

		user, err := dao.FindOneByKEY("users", "email", claims["username"].(string))

		if user == nil {

			fmt.Println("user not found trying doctor")

			user, err = dao.FindOneByKEY("doctors", "email", claims["username"].(string))

			//fmt.Println("user", user)

			if user == nil {

				user, err = dao.FindOneByKEY("contacts", "email", claims["username"].(string))
				if err != nil {
					log.Fatal("Can not get user from token")
					return
				} else {
					userType = 3
				}
			}
			if err != nil {

				log.Fatal("Can not get user from token")
				return
			} else {
				userType = 2
			}

		} else {
			userType = 1
		}

		if err != nil {
			log.Fatal("Can not get user from token")
			return
		}

		context.Set(r, "user", user)

		context.Set(r, "userType", userType)

		//log.Println(user)

		next.ServeHTTP(w, r)

		//log.Println("Executing middlewareOne again")
	})
}

// OnlyAdminMiddleware can execute request if is admin
func OnlyAdminMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := context.Get(r, "user")

		userParsed := user.(bson.M)

		if userParsed["role"] == "admin" {
			next.ServeHTTP(w, r)
		} else {
			return
		}

	})

}
