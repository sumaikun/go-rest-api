package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	Helpers "click-al-vet/helpers"
	Models "click-al-vet/models"

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

func resetPasswordValidator(r *http.Request) (map[string]interface{}, Models.ResetPassword) {

	var reset Models.ResetPassword

	rules := govalidator.MapData{
		"password": []string{"required"},
		"token":    []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &reset,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, reset
}

func forgotPasswordValidator(r *http.Request) (map[string]interface{}, Models.ForgotPassword) {

	var forgot Models.ForgotPassword

	rules := govalidator.MapData{
		"email": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &forgot,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, forgot
}

func confirmAccountValidator(r *http.Request) (map[string]interface{}, Models.ConfirmAccount) {

	var confirm Models.ConfirmAccount

	rules := govalidator.MapData{
		"token": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &confirm,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, confirm
}

func productValidator(r *http.Request) (map[string]interface{}, Models.Product) {

	var product Models.Product

	rules := govalidator.MapData{
		"name":        []string{"required"},
		"value":       []string{"numeric"},
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

func specialistTypesValidator(r *http.Request) (map[string]interface{}, Models.SpecialistTypes) {

	var parameters Models.SpecialistTypes

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

func cityTypesValidator(r *http.Request) (map[string]interface{}, Models.CitiesTypes) {

	var parameters Models.CitiesTypes

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

	case "specialistTypes":
		err, data := specialistTypesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = bson.NewObjectId()
			data.Date = time.Now().String()
			data.UpdateDate = time.Now().String()
		}
		fmt.Println(data)
		return err, data, []string{"name"}

	case "cityTypes":
		err, data := cityTypesValidator(r)
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

	case "specialistTypes":
		err, data := specialistTypesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = prevData["_id"].(bson.ObjectId)
			data.Date = prevData["date"].(string)
			data.UpdateDate = time.Now().String()
		}

		return err, data, data.ID

	case "cityTypes":
		err, data := cityTypesValidator(r)
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
		"patient":           []string{"required"},
		"rabiesVaccine":     []string{"bool"},
		"rabiesVaccineDate": []string{"date"},
		"feedingType":       []string{"required", "feedingTypeEnum"},
		"reproductiveState": []string{"required", "reproductiveStateEnum"},
		"previousIllnesses": []string{"required"},
		"surgeris":          []string{"required"},
		"familyBackground":  []string{"required"},
		"habitat":           []string{"required"},
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
		"patient":         []string{"required"},
		"tlic":            []string{"required"},
		"heartRate":       []string{"required"},
		"respiratoryRate": []string{"required"},
		"heartBeat":       []string{"required"},
		"temperature":     []string{"required"},
		"weight":          []string{"required"},
		"attitude":        []string{"required", "attitudeEnum"},
		"bodyCondition":   []string{"required", "bodyConditionEnum"},
		"hidrationStatus": []string{"required", "hidrationStatusEnum"},
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

//////////////////////////////////////////////////////////////////////

func appointmentsValidator(r *http.Request) (map[string]interface{}, Models.Appointments) {

	var appointments Models.Appointments

	rules := govalidator.MapData{
		"patient":                      []string{"required"},
		"reasonForConsultation":        []string{"required"},
		"resultsForConsultation":       []string{"required"},
		"appointmentDate":              []string{"required"},
		"state":                        []string{"required"},
		"medicalReasonForConsultation": []string{"required"},
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

//////////////////////////////////////////////////////////////////////

func appointmentsScheduleValidator(r *http.Request) (map[string]interface{}, Models.Appointments) {

	var appointments Models.Appointments

	rules := govalidator.MapData{
		"patient":          []string{"required"},
		"appointmentDate":  []string{"required"},
		"state":            []string{"required"},
		"agendaAnnotation": []string{"required"},
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

//////////////////////////////////////////////////////////////////////

func medicinesValidator(r *http.Request) (map[string]interface{}, Models.Medicines) {

	var medicines Models.Medicines

	rules := govalidator.MapData{
		"patient":           []string{"required"},
		"administrationWay": []string{"required"},
		"duration":          []string{"required"},
		"posology":          []string{"required"},
		"presentation":      []string{"required"},
		"product":           []string{"required"},
		"appointment":       []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &medicines,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, medicines
}

//////////////////////////////////////////////////////////////////////

func detectedDiseasesValidator(r *http.Request) (map[string]interface{}, Models.DetectedDiseases) {

	var detectedDisease Models.DetectedDiseases

	rules := govalidator.MapData{
		"patient":  []string{"required"},
		"disease":  []string{"required"},
		"criteria": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &detectedDisease,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, detectedDisease
}

//////////////////////////////////////////////////////////////////////

func patientsFilesValidator(r *http.Request) (map[string]interface{}, Models.PatientFiles) {

	var patientFile Models.PatientFiles

	rules := govalidator.MapData{
		"patient":     []string{"required"},
		"filePath":    []string{"required"},
		"description": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &patientFile,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, patientFile
}

//////////////////////////////////////////////////////////////////////

func agendaAnnotationValidator(r *http.Request) (map[string]interface{}, Models.AgendaAnnotation) {

	var agendaAnnotation Models.AgendaAnnotation

	rules := govalidator.MapData{
		"annotationDate":   []string{"required"},
		"annotationToDate": []string{"required"},
		"description":      []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &agendaAnnotation,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, agendaAnnotation
}

func userRegisterValidator(r *http.Request) (map[string]interface{}, Models.UserRegister) {

	var userRegister Models.UserRegister

	rules := govalidator.MapData{
		"name":      []string{"required"},
		"lastName":  []string{"required"},
		"email":     []string{"required", "email"},
		"phone":     []string{"required", "min:7", "max:10"},
		"password":  []string{"required"},
		"confirmed": []string{"bool"},
		"city":      []string{"required", "cityParam"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &userRegister,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, userRegister
}

func doctorValidator(r *http.Request) (map[string]interface{}, Models.Doctor) {

	var user Models.Doctor

	rules := govalidator.MapData{
		"name":           []string{"required"},
		"lastName":       []string{"required"},
		"email":          []string{"required", "email"},
		"phone":          []string{"required", "min:7", "max:10"},
		"address":        []string{"required"},
		"birthDate":      []string{"required"},
		"city":           []string{"required", "cityParam"},
		"specialistType": []string{"required", "specialistTypeParam"},
		"typeId":         []string{"required", "documentTypeEnum"},
		"identification": []string{"required"},
		"state":          []string{"stateEnum"},
		"medicalCenter":  []string{"medicalCenterParam"},
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

//////////////////////////////////////////////////////////////////////

func doctorSettingsValidator(r *http.Request) (map[string]interface{}, Models.DoctorSettings) {

	var doctorSettings Models.DoctorSettings

	rules := govalidator.MapData{
		"hoursRange":   []string{"required", "hoursRangeType"},
		"daysRange":    []string{"required", "daysRangeType"},
		"isScheduling": []string{"bool"},
		"doctor":       []string{"required", "doctorParam"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &doctorSettings,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, doctorSettings
}

//////////////////////////////////////////////////////////////////////
func init_validators() {
	govalidator.AddCustomRule("presentationEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Jarabes", "Gotas", "Capsulas", "Polvo", "Granulado", "Emulsión", "Bebible"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for presentation Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("administrationWayEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Oral", "Intravenosa", "Intramuscular", "Subcutanea", "tópica", "rectal", "inhalatoria"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for administrationWay Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("feedingTypeEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Balanceada", "Casera", "Mixta"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for feedingType Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("reproductiveStateEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Castrado", "Gestacion", "Entero", "Lactancia"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for reproductiveState Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("habitatEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Casa", "Lote", "Finca", "Taller", "Apartamento"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for reproductiveState Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("attitudeEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Astenico", "Apopletico", "Linfatico"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for reproductiveState Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("bodyConditionEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Caquetico", "Delgado", "Normal", "Obeso", "Sobrepeso"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for reproductiveState Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("hidrationStatusEnum", func(field string, rule string, message string, value interface{}) error {

		x := []string{"Normal", "0-5%", "6-7%", "8-9%", "+10%"}

		val := Helpers.Contains(x, value.(string))

		if val != true {
			return fmt.Errorf("The %s field must be a valid value for reproductiveState Enum", field)
		}
		return nil
	})

	govalidator.AddCustomRule("cityParam", func(field string, rule string, message string, value interface{}) error {

		fmt.Println("city value " + value.(string))
		if len(value.(string)) > 0 {
			_, err := dao.FindByID("cityTypes", value.(string))
			if err != nil {
				return fmt.Errorf("The %s field must be a valid value must have a valid city ID", field)
			}
		}
		return nil
	})

	govalidator.AddCustomRule("specialistTypeParam", func(field string, rule string, message string, value interface{}) error {

		for _, element := range value.([]string) {
			fmt.Println("specialist Type value " + element)

			if len(element) > 0 {
				_, err := dao.FindByID("specialistTypes", element)
				if err != nil {
					return fmt.Errorf("The %s field must be a valid value must have a valid specialist type ID", field)
				}
			}
		}

		return nil
	})

	govalidator.AddCustomRule("stateEnum", func(field string, rule string, message string, value interface{}) error {
		if len(value.(string)) > 0 {
			x := []string{"ACTIVE", "INACTIVE", "PENDING", "CHANGE_PASSWORD"}

			val := Helpers.Contains(x, value.(string))

			if val != true {
				return fmt.Errorf("The %s field must be a valid value for state Enum", field)
			}
		}
		return nil
	})

	govalidator.AddCustomRule("documentTypeEnum", func(field string, rule string, message string, value interface{}) error {

		if len(value.(string)) > 0 {
			x := []string{"CC", "CE", "PS"}

			val := Helpers.Contains(x, value.(string))

			if val != true {
				return fmt.Errorf("The %s field must be a valid value for documentType Enum", field)
			}
		}
		return nil
	})

	govalidator.AddCustomRule("medicalCenterParam", func(field string, rule string, message string, value interface{}) error {

		if len(value.(string)) > 0 {
			_, err := dao.FindByID("medicalCenters", value.(string))
			if err != nil {
				return fmt.Errorf("The %s field must be a valid value must have a valid medical center ID", field)
			}
		}

		return nil
	})

	govalidator.AddCustomRule("doctorParam", func(field string, rule string, message string, value interface{}) error {

		if len(value.(string)) > 0 {
			_, err := dao.FindByID("doctors", value.(string))
			if err != nil {
				return fmt.Errorf("The %s field must be a valid value must have a valid medical center ID", field)
			}
		}

		return nil
	})

	govalidator.AddCustomRule("appointmentTypeEnum", func(field string, rule string, message string, value interface{}) error {

		if len(value.(string)) > 0 {
			x := []string{"DONE", "PENDING", "CONFIRMED", "CANCELLED", "PENDING DOCTOR", "DUE"}

			val := Helpers.Contains(x, value.(string))

			if val != true {
				return fmt.Errorf("The %s field must be a valid value for appointmentType Enum", field)
			}
		}
		return nil
	})

	govalidator.AddCustomRule("sexTypeEnum", func(field string, rule string, message string, value interface{}) error {

		if len(value.(string)) > 0 {
			x := []string{"M", "F"}

			val := Helpers.Contains(x, value.(string))

			if val != true {
				return fmt.Errorf("The %s field must be a valid value for sexTypeEnum Enum", field)
			}
		}
		return nil
	})

	govalidator.AddCustomRule("hoursRangeType", func(field string, rule string, message string, value interface{}) error {

		parsedRangeType, ok := value.([]int)

		if ok == true {
			if len(parsedRangeType) != 2 {
				return fmt.Errorf("The %s field must be a valid string array of two positions", field)
			} else {
				/*hour1, err := strconv.Atoi(parsedRangeType[0])
				if err != nil {
					return fmt.Errorf("internal error converting data", field)
				}
				hour2, err := strconv.Atoi(parsedRangeType[1])
				if err != nil {
					return fmt.Errorf("internal error converting data", field)
				}*/

				hour1 := parsedRangeType[0]

				hour2 := parsedRangeType[1]

				if hour2 <= hour1 {
					return fmt.Errorf("The %s field have invalid inputs final time must no be greater than initial times", field)
				}

			}
		} else {
			return fmt.Errorf("The %s field must be a valid string array", field)
		}

		return nil

	})

	govalidator.AddCustomRule("daysRangeType", func(field string, rule string, message string, value interface{}) error {

		parsedRangeType, ok := value.([]string)

		if ok == true {

			x := []string{"Mon", "Tues", "Wed", "Thurs", "Frid", "Sat", "Sun"}

			for _, day := range parsedRangeType {

				val := Helpers.Contains(x, day)

				if val != true {
					return fmt.Errorf("The %s field must have a valid day value", field)
				}

			}

		} else {
			return fmt.Errorf("The %s field must be a valid string array", field)
		}

		return nil
	})
}
