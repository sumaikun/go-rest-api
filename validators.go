package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	Models "github.com/sumaikun/go-rest-api/models"
	"github.com/thedevsaddam/govalidator"
	"gopkg.in/mgo.v2/bson"
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
		"name":              []string{"required"},
		"value":             []string{"numeric"},
		"description":       []string{"required"},
		"presentation":      []string{"required", "presentationEnum"},
		"administrationWay": []string{"required", "administrationWayEnum"},
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
		"name":    []string{"required"},
		"species": []string{"required"},
		"breed":   []string{"required"},
		"color":   []string{"required"},
		"sex":     []string{"required"},
		"age":     []string{"required"},
		"origin":  []string{"required"},
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

func breedsValidator(r *http.Request) (map[string]interface{}, Models.Breeds) {

	var parameters Models.Breeds

	rules := govalidator.MapData{
		"name":    []string{"required"},
		"species": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &parameters,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, parameters
}

func speciesValidator(r *http.Request) (map[string]interface{}, Models.Species) {

	var parameters Models.Species

	rules := govalidator.MapData{
		"name": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &parameters,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, parameters
}

func validatorSelector(r *http.Request, entity string) (map[string]interface{}, interface{}, []string) {

	var err map[string]interface{} = nil

	switch entity {
	case "breeds":
		err, data := breedsValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = bson.NewObjectId()
			data.Date = time.Now().String()
			data.UpdateDate = time.Now().String()
		}
		return err, data, []string{"name"}

	case "species":
		err, data := speciesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = bson.NewObjectId()
			data.Date = time.Now().String()
			data.UpdateDate = time.Now().String()
		}
		fmt.Println(data)
		return err, data, []string{"name"}
	}

	return err, nil, nil

}

func validatorSelectorUpdate(r *http.Request, entity string, prevData bson.M) (map[string]interface{}, interface{}, interface{}) {

	var err map[string]interface{} = nil

	switch entity {
	case "breeds":
		err, data := breedsValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = prevData["_id"].(bson.ObjectId)
			data.Date = prevData["date"].(string)
			data.UpdateDate = time.Now().String()
		}
		return err, data, data.ID

	case "species":
		err, data := speciesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = prevData["_id"].(bson.ObjectId)
			data.Date = prevData["date"].(string)
			data.UpdateDate = time.Now().String()
		}

		return err, data, data.ID
	}

	return err, nil, nil

}

///////////////////////////////////////////////////////////////////////

func patientReviewValidator(r *http.Request) (map[string]interface{}, Models.PatientReview) {

	var patientReview Models.PatientReview

	rules := govalidator.MapData{
		"patient":                []string{"required", "string"},
		"pvcVaccine":             []string{"bool"},
		"pvcVaccineDate":         []string{"date"},
		"tripleVaccine":          []string{"bool"},
		"tripleVaccineDate":      []string{"date"},
		"rabiesVaccine":          []string{"bool"},
		"rabiesVaccineDate":      []string{"date"},
		"desparasitationProduct": []string{"string"},
		"lastDesparasitation":    []string{"string"},
		"feedingType":            []string{"required", "feedingTypeEnum"},
		"reproductiveState":      []string{"required", "reproductiveStateEnum"},
		"previousIllnesses":      []string{"required", "string"},
		"surgeris":               []string{"required", "string"},
		"familyBackground":       []string{"required", "string"},
		"habitat":                []string{"required", "string"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &patientReview,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, patientReview
}

///////////////////////////////////////////////////////////////////////

func physiologicalConstantsValidator(r *http.Request) (map[string]interface{}, Models.PhysiologicalConstants) {

	var physiologicalConstants Models.PhysiologicalConstants

	rules := govalidator.MapData{
		"tlic":                          []string{"string"},
		"heartRate":                     []string{"string"},
		"respiratoryRate":               []string{"string"},
		"HeartBeat":                     []string{"string"},
		"temperature":                   []string{"string"},
		"weight":                        []string{"string"},
		"attitude":                      []string{"attitudeEnum"},
		"bodyCondition":                 []string{"bodyConditionEnum"},
		"hidrationStatus":               []string{"hidrationStatusEnum"},
		"conjuntivalMucosa":             []string{"string"},
		"oralMucosa":                    []string{"string"},
		"vulvallMucosa":                 []string{"string"},
		"rectalMucosa":                  []string{"string"},
		"physicalsEye":                  []string{"string"},
		"physicalsEars":                 []string{"string"},
		"physicalsLinfaticmodules":      []string{"string"},
		"physicalsSkinandanexes":        []string{"string"},
		"physicalsLocomotion":           []string{"string"},
		"physicalsMusclesqueletal":      []string{"string"},
		"physicalsNervoussystem":        []string{"string"},
		"physicalsCardiovascularsystem": []string{"string"},
		"physicalsRespiratorysystem":    []string{"string"},
		"physicalsDigestivesystem":      []string{"string"},
		"physicalsGenitourinarysystem":  []string{"string"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &physiologicalConstants,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, physiologicalConstants
}

///////////////////////////////////////////////////////////////////////

func diagnosticPlansValidator(r *http.Request) (map[string]interface{}, Models.DiagnosticPlans) {

	var diagnosticPlans Models.DiagnosticPlans

	rules := govalidator.MapData{
		"typeOfExam":        []string{"string"},
		"description":       []string{"string"},
		"examDate":          []string{"string"},
		"laboratory":        []string{"string"},
		"laboratoryAddress": []string{"string"},
		"results":           []string{"string"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &diagnosticPlans,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, diagnosticPlans
}

//////////////////////////////////////////////////////////////////////

func terapeuticPlansValidator(r *http.Request) (map[string]interface{}, Models.TerapeuticPlans) {

	var terapeuticPlans Models.TerapeuticPlans

	rules := govalidator.MapData{
		"typeOfPlan":                  []string{"string"},
		"activeSubstanceToAdminister": []string{"string"},
		"posology":                    []string{"string"},
		"totalDose":                   []string{"string"},
		"frecuencyAndDuration":        []string{"string"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &terapeuticPlans,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, terapeuticPlans
}

//////////////////////////////////////////////////////////////////////

func appointmentsValidator(r *http.Request) (map[string]interface{}, Models.Appointments) {

	var appointments Models.Appointments

	rules := govalidator.MapData{
		"reasonForConsultation":  []string{"string"},
		"resultsForConsultation": []string{"string"},
		"relatedProcesses":       []string{"string"},
		"product":                []string{"string"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &appointments,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, appointments
}
