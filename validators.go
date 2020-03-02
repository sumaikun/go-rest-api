package main

import (
	"net/http"

	Models "github.com/sumaikun/go-rest-api/models"
	"github.com/thedevsaddam/govalidator"
)

func userValidator(r *http.Request) (map[string]interface{}, Models.User) {

	var user Models.User

	rules := govalidator.MapData{
		"name":    []string{"required"},
		"email":   []string{"required", "email"},
		"phone":   []string{"min:7", "max:10"},
		"address": []string{"required"},
		//"picture": []string{"url"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &user,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, user
}

/*if err := map[string]interface{}{"validationError": e}; len(e) > 0 {
	//fmt.Println(len(e))
	Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
	return
}*/

func productValidator(r *http.Request) (map[string]interface{}, Models.Product) {

	var product Models.Product

	rules := govalidator.MapData{
		"name":        []string{"required"},
		"value":       []string{"required", "numeric"},
		"description": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &product,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, product
}

func contactValidator(r *http.Request) (map[string]interface{}, Models.Contact) {

	var contact Models.Contact

	rules := govalidator.MapData{
		"name":           []string{"required"},
		"email":          []string{"required", "email"},
		"phone":          []string{"min:7", "max:10"},
		"address":        []string{"required"},
		"identification": []string{"required"},
		"stratus":        []string{"required"},
		"city":           []string{"required"},
		"ocupation":      []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &contact,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, contact
}

func petValidator(r *http.Request) (map[string]interface{}, Models.Pet) {

	var pet Models.Pet

	rules := govalidator.MapData{
		"name":        []string{"required"},
		"species":     []string{"required"},
		"breed":       []string{"required"},
		"color":       []string{"required"},
		"sex":         []string{"required"},
		"birthDate":   []string{"required"},
		"age":         []string{"required", "numeric"},
		"origin":      []string{"required"},
		"description": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &pet,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, pet
}
