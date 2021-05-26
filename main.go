package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	middleware "github.com/sumaikun/go-rest-api/middlewares"

	Config "github.com/sumaikun/go-rest-api/config"

	Dao "github.com/sumaikun/go-rest-api/dao"

	Helpers "github.com/sumaikun/go-rest-api/helpers"

	"github.com/thedevsaddam/govalidator"

	gomail "gopkg.in/mail.v2"
)

var (
	port        string
	jwtKey      []byte
	logo        string
	frontEndUrl string
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

	var config = Config.Config{}
	config.Read()
	//fmt.Println(config.Jwtkey)
	frontEndUrl = config.FrontEndUrl
	jwtKey = []byte(config.Jwtkey)
	port = config.Port
	logo = config.LogoUrl

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()

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

func main() {
	//initEvents()
	fmt.Println("start server in port " + port)
	router := mux.NewRouter().StrictSlash(true)

	//testEmail()

	/* Authentication */
	router.HandleFunc("/auth", authentication).Methods("POST")
	router.Handle("/exampleHandler", middleware.AuthMiddleware(http.HandlerFunc(exampleHandler))).Methods("GET")
	router.HandleFunc("/createInitialUser", createInititalUser).Methods("POST")
	router.Handle("/resetPassword", middleware.AuthMiddleware(http.HandlerFunc(resetPassword))).Methods("POST")
	router.HandleFunc("/forgotPassword", forgotPassword).Methods("POST")
	router.Handle("/confirmAccount", middleware.AuthMiddleware(http.HandlerFunc(confirmAccount))).Methods("POST")
	router.HandleFunc("/registerDoctor", registerDoctor).Methods("POST")
	router.HandleFunc("/registerContact", registerContact).Methods("POST")

	/* Users Routes */
	router.Handle("/users", middleware.AuthMiddleware(middleware.UserMiddleware(middleware.OnlyAdminMiddleware(http.HandlerFunc(createUsersEndPoint))))).Methods("POST")
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(allUsersEndPoint))).Methods("GET")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(findUserEndpoint))).Methods("GET")
	//router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeUserEndpoint))).Methods("DELETE")
	router.Handle("/users/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(middleware.OnlyAdminMiddleware(http.HandlerFunc(updateUserEndPoint))))).Methods("PUT")

	/* Doctors Routes */
	router.Handle("/doctors", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createDoctorsEndPoint)))).Methods("POST")
	router.Handle("/doctors", middleware.AuthMiddleware(http.HandlerFunc(allDoctorsEndPoint))).Methods("GET")
	router.Handle("/doctors/{id}", middleware.AuthMiddleware(http.HandlerFunc(findDoctorEndPoint))).Methods("GET")
	router.Handle("/doctors/{id}", middleware.AuthMiddleware(http.HandlerFunc(inactivateDoctorEndPoint))).Methods("DELETE")
	router.Handle("/doctors/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateDoctorEndPoint)))).Methods("PUT")

	/* Products Routes */
	router.Handle("/products", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createProductEndPoint)))).Methods("POST")
	router.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(allProductsEndPoint))).Methods("GET")
	router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(findProductEndpoint))).Methods("GET")
	//router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeProductEndpoint))).Methods("DELETE")
	router.Handle("/products/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateProductEndPoint)))).Methods("PUT")

	/* Contacts Routes */
	router.Handle("/contacts", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createContactEndPoint)))).Methods("POST")
	router.Handle("/contacts", middleware.AuthMiddleware(http.HandlerFunc(allContactsEndPoint))).Methods("GET")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(findContactEndpoint))).Methods("GET")
	//router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeContactEndpoint))).Methods("DELETE")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateContactEndPoint)))).Methods("PUT")

	/* Pets Routes */
	router.Handle("/pets", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPetEndPoint)))).Methods("POST")
	router.Handle("/pets", middleware.AuthMiddleware(http.HandlerFunc(allPetsEndPoint))).Methods("GET")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPetEndpoint))).Methods("GET")
	//router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePetEndpoint))).Methods("DELETE")
	router.Handle("/updatePetContactsEndPoint/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePetContactsEndPoint)))).Methods("PUT")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePetEndPoint)))).Methods("PUT")

	/* Examtypes Routes */
	router.Handle("/examTypes", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/examTypes", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/examTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	router.Handle("/examTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/examTypes/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	/* Plantypes Routes */
	router.Handle("/planTypes", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/planTypes", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/planTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	//router.Handle("/planTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/planTypes/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	/* Diseases Routes */
	router.Handle("/diseases", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/diseases", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/diseases/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	//router.Handle("/diseases/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/diseases/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	/* fileUpload */

	router.Handle("/fileUpload", middleware.AuthMiddleware(http.HandlerFunc(fileUpload))).Methods("POST")
	router.HandleFunc("/serveImage/{image}", serveImage).Methods("GET")
	router.Handle("/deleteFile/{file}", middleware.AuthMiddleware(http.HandlerFunc(deleteImage))).Methods("DELETE")
	router.Handle("/downloadFile/{file}", middleware.AuthMiddleware(http.HandlerFunc(downloadFile))).Methods("GET")

	/* enums */
	router.Handle("/userRoles", middleware.AuthMiddleware(http.HandlerFunc(userRoles))).Methods("GET")
	router.Handle("/contactStratus", middleware.AuthMiddleware(http.HandlerFunc(contactStratus))).Methods("GET")
	router.Handle("/contactDocumentType", middleware.AuthMiddleware(http.HandlerFunc(contactDocumentType))).Methods("GET")
	router.Handle("/parametersType", middleware.AuthMiddleware(http.HandlerFunc(parametersType))).Methods("GET")
	router.Handle("/administrationWays", middleware.AuthMiddleware(http.HandlerFunc(administrationWayType))).Methods("GET")
	router.Handle("/presentations", middleware.AuthMiddleware(http.HandlerFunc(presentationType))).Methods("GET")

	/* patientReviews */

	router.Handle("/patientReviews", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPatientReviewEndPoint)))).Methods("POST")
	router.Handle("/patientReviews", middleware.AuthMiddleware(http.HandlerFunc(allPatientReviewEndPoint))).Methods("GET")
	router.Handle("/patientReviews/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findPatientReviewByPatientEndpoint))).Methods("GET")
	router.Handle("/patientReview/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPatientReviewEndpoint))).Methods("GET")
	router.Handle("/patientReviews/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePatientReviewEndpoint))).Methods("DELETE")
	router.Handle("/patientReviews/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePatientReviewEndPoint)))).Methods("PUT")

	/* physiological Constants */

	router.Handle("/physiologicalConstants", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPhysiologicalConstantsEndPoint)))).Methods("POST")
	router.Handle("/physiologicalConstants", middleware.AuthMiddleware(http.HandlerFunc(allPhysiologicalConstantsEndPoint))).Methods("GET")
	router.Handle("/physiologicalConstants/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findPhysiologicalConstantsByPatientEndpoint))).Methods("GET")
	router.Handle("/physiologicalConstant/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPhysiologicalConstantsEndpoint))).Methods("GET")
	router.Handle("/physiologicalConstants/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePhysiologicalConstantsEndpoint))).Methods("DELETE")
	router.Handle("/physiologicalConstants/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePhysiologicalConstantsEndPoint)))).Methods("PUT")

	/*  diagnostic Plan */

	router.Handle("/diagnosticPlans", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createDiagnosticPlansEndPoint)))).Methods("POST")
	router.Handle("/diagnosticPlans", middleware.AuthMiddleware(http.HandlerFunc(allDiagnosticPlansEndPoint))).Methods("GET")
	router.Handle("/diagnosticPlans/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findDiagnosticPlansByPatientEndpoint))).Methods("GET")
	router.Handle("/diagnosticPlan/{id}", middleware.AuthMiddleware(http.HandlerFunc(findDiagnosticPlansEndpoint))).Methods("GET")
	router.Handle("/diagnosticPlans/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeDiagnosticPlansEndpoint))).Methods("DELETE")
	router.Handle("/diagnosticPlans/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateDiagnosticPlansEndPoint)))).Methods("PUT")

	/* therapeutic Plan */

	router.Handle("/therapeuticPlans", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createTherapeuticPlansEndPoint)))).Methods("POST")
	router.Handle("/therapeuticPlans", middleware.AuthMiddleware(http.HandlerFunc(allTherapeuticPlansEndPoint))).Methods("GET")
	router.Handle("/therapeuticPlans/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findTherapeuticPlansByPatientEndpoint))).Methods("GET")
	router.Handle("/therapeuticPlan/{id}", middleware.AuthMiddleware(http.HandlerFunc(findTherapeuticPlansEndpoint))).Methods("GET")
	router.Handle("/therapeuticPlans/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeTherapeuticPlansEndpoint))).Methods("DELETE")
	router.Handle("/therapeuticPlans/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateTherapeuticPlansEndPoint)))).Methods("PUT")

	/* appointment */

	router.Handle("/appointments", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createAppointmentsEndPoint)))).Methods("POST")
	router.Handle("/appointments", middleware.AuthMiddleware(http.HandlerFunc(allAppointmentsEndPoint))).Methods("GET")
	router.Handle("/appointments/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findAppointmentsByPatientEndpoint))).Methods("GET")
	router.Handle("/appointment/{id}", middleware.AuthMiddleware(http.HandlerFunc(findAppointmentsEndpoint))).Methods("GET")
	router.Handle("/appointments/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeAppointmentsEndpoint))).Methods("DELETE")
	router.Handle("/appointments/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateAppointmentsEndPoint)))).Methods("PUT")

	/* detectedDiseases */

	router.Handle("/detectedDiseases", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createDetectedDiseaseEndPoint)))).Methods("POST")
	router.Handle("/detectedDiseases", middleware.AuthMiddleware(http.HandlerFunc(allDetectedDiseasesEndPoint))).Methods("GET")
	router.Handle("/detectedDiseases/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findDetectedDiseasesByPatientEndpoint))).Methods("GET")
	router.Handle("/detectedDisease/{id}", middleware.AuthMiddleware(http.HandlerFunc(findDetectedDiseaseEndpoint))).Methods("GET")
	router.Handle("/detectedDiseases/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeDetectedDiseaseEndpoint))).Methods("DELETE")
	router.Handle("/detectedDiseases/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateDetectedDiseaseEndPoint)))).Methods("PUT")

	/* patientFiles */

	router.Handle("/patientFiles", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPatientFilesEndPoint)))).Methods("POST")
	router.Handle("/patientFiles", middleware.AuthMiddleware(http.HandlerFunc(allPatientFilesEndPoint))).Methods("GET")
	router.Handle("/patientFiles/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findPatientFilesByPatientEndpoint))).Methods("GET")
	router.Handle("/patientFile/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPatientFilesEndpoint))).Methods("GET")
	router.Handle("/patientFiles/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePatientFilesEndpoint))).Methods("DELETE")
	router.Handle("/patientFiles/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePatientFilesEndPoint)))).Methods("PUT")

	/* agendaAnnotations */

	router.Handle("/agendaAnnotations", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createAgendaAnnotationEndPoint)))).Methods("POST")
	router.Handle("/agendaAnnotations", middleware.AuthMiddleware(http.HandlerFunc(allAgendaAnnotationsEndPoint))).Methods("GET")
	router.Handle("/agendaAnnotations/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findPatientFilesByPatientEndpoint))).Methods("GET")
	router.Handle("/agendaAnnotation/{id}", middleware.AuthMiddleware(http.HandlerFunc(findAgendaAnnotationEndpoint))).Methods("GET")
	router.Handle("/agendaAnnotations/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeAgendaAnnotationEndpoint))).Methods("DELETE")
	router.Handle("/agendaAnnotations/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateAgendaAnnotationEndPoint)))).Methods("PUT")

	/* SpecialistTypes Routes */
	router.Handle("/specialistTypes", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/specialistTypes", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/specialistTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	//router.Handle("/specialistTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/specialistTypes/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	/* CityTypes Routes */
	router.Handle("/cityTypes", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.HandleFunc("/cityTypes", allParametersEndPoint).Methods("GET")
	router.Handle("/cityTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	//router.Handle("/cityTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/cityTypes/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	log.Fatal(http.ListenAndServe(":"+port, &CORSRouterDecorator{router}))
}

var htmlContent = "<div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>!Hola ¡Javier gil!, los detalles de tu cita son los siguientes: </span><br/><br/><table style='border-collapse: collapse; width:100%; border: 1px solid black;'  ><tbody><tr><td style='border: 1px solid black' ><b>Diagnostico:</b></td><td style='border: 1px solid black'  >Infeccion debida a coronavirus, sin otra especificacion</td></tr><tr><td style='border: 1px solid black'  ><b>Observaciones:</b></td><td style='border: 1px solid black'  >Se hicieron diferentes pruebas y se determino que el diagnostico es debido a ...</td></tr></tbody></table><br/><span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Los medicamentos recetados son los siguientes:</span><br/><br/><table style='border-collapse: collapse; width:100%; border: 1px solid black;'  ><thead><tr><th style='color: #54ace2;font-size: 16px;font-weight: bold;'>Medicamento</th><th  style='color: #54ace2;font-size: 16px;font-weight: bold;'>Presentación</th><th  style='color: #54ace2;font-size: 16px;font-weight: bold;'>Posología</th><th  style='color: #54ace2;font-size: 16px;font-weight: bold;'>Duración</th></tr></thead><tbody><tr><td>dsd</td><td>dsd</td><td>dsd</td><td>dsd</td></tr></tbody></table></div></div>"

func testEmail() {

	var config = Config.Config{}
	config.Read()

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	//m.SetHeader("To", "solucionesitecnologia@gmail.com")
	m.SetHeader("To", "ventas.javc@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Bienvenido a clickal medic, confirma tu contraseña")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContent)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendResetPasswordEmail(token string, mail string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = " <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Habilita tu usuario con este </span><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/recover-password?tokenizer=" + token + "' >Enlace</a></div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Saludos de clickal medic, confirma tu contraseña")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendConfirmationEmail(token string, mail string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = " <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Habilita tu usuario con este </span><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/confirm-account?tokenizer=" + token + "' >Enlace</a></div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Bienvenido a clickal medic, confirma tu usuario")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendForgotPasswordEmail(token string, mail string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = " <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Habilita tu usuario con este </span><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/recover-password?tokenizer=" + token + "' >Enlace</a></div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Recupera tu usario con el siguiente enlace")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendAppointmentConfirmationEmail(token string, mail string, appointment string, doctorName string, appointmentDate string, appointmentHour string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = "  <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div><span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>El doctor </span><span style='color: #54ace2;font-size: 20px;font-weight: bold;'> " + doctorName + " </span><span style='color: #0f76b0;font-size: 20px;font-weight: bold;'> agendo su cita para:</span><span style='color: #54ace2;font-size: 20px;font-weight: bold;'> " + appointmentDate + " a las " + appointmentHour + " </span><br/><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/confirm-appointment?tokenizer=" + token + "&appointment= " + appointment + "' >Confirmar</a><a style='color: red;font-weight: bold;font-size: 20px;margin-left:5px' href='" + frontEndUrl + "/cencel-appointment?tokenizer=" + token + "&appointment= " + appointment + "  ' >Cancelar</a></div></div> "

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "¡Tienes una cita pendiente!")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendEmailConfirmationToPatient(mail string, phone string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = "<div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Se ha confirmado tu cita</span><br/>No dudes en llamar al " + phone + "</div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "¡Tienes una cita pendiente!")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendEmailCancelationToPatient(mail string, phone string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = "<div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: red;font-size: 20px;font-weight: bold;'>Se ha cancelado tu cita</span><br/>Confirma que paso al " + phone + "</div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "¡Tienes una cita pendiente!")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}
