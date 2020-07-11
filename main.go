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

	Helpers "github.com/sumaikun/go-rest-api/helpers"

	"github.com/thedevsaddam/govalidator"
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

}

func main() {
	//initEvents()
	fmt.Println("start server in port " + port)
	router := mux.NewRouter().StrictSlash(true)

	/* Authentication */
	router.HandleFunc("/auth", authentication).Methods("POST")
	router.Handle("/exampleHandler", middleware.AuthMiddleware(http.HandlerFunc(exampleHandler))).Methods("GET")

	/* Users Routes */
	router.Handle("/users", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createUsersEndPoint)))).Methods("POST")
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(allUsersEndPoint))).Methods("GET")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(findUserEndpoint))).Methods("GET")
	router.Handle("/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeUserEndpoint))).Methods("DELETE")
	router.Handle("/users/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateUserEndPoint)))).Methods("PUT")

	/* Products Routes */
	router.Handle("/products", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createProductEndPoint)))).Methods("POST")
	router.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(allProductsEndPoint))).Methods("GET")
	router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(findProductEndpoint))).Methods("GET")
	router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeProductEndpoint))).Methods("DELETE")
	router.Handle("/products/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateProductEndPoint)))).Methods("PUT")

	/* Contacts Routes */
	router.Handle("/contacts", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createContactEndPoint)))).Methods("POST")
	router.Handle("/contacts", middleware.AuthMiddleware(http.HandlerFunc(allContactsEndPoint))).Methods("GET")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(findContactEndpoint))).Methods("GET")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeContactEndpoint))).Methods("DELETE")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateContactEndPoint)))).Methods("PUT")

	/* Pets Routes */
	router.Handle("/pets", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPetEndPoint)))).Methods("POST")
	router.Handle("/pets", middleware.AuthMiddleware(http.HandlerFunc(allPetsEndPoint))).Methods("GET")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPetEndpoint))).Methods("GET")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePetEndpoint))).Methods("DELETE")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePetEndPoint)))).Methods("PUT")

	/* Breeds Routes */
	router.Handle("/breeds", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/breeds", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/breeds/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	router.Handle("/breeds/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/breeds/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	/* Species Routes */
	router.Handle("/species", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/species", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/species/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	router.Handle("/species/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/species/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

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
	router.Handle("/planTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/planTypes/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	/* Diseases Routes */
	router.Handle("/diseases", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/diseases", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/diseases/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	router.Handle("/diseases/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
	router.Handle("/diseases/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateParameterEndPoint)))).Methods("PUT")

	/* fileUpload */

	router.Handle("/fileUpload", middleware.AuthMiddleware(http.HandlerFunc(fileUpload))).Methods("POST")
	router.HandleFunc("/serveImage/{image}", serveImage).Methods("GET")

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
	router.Handle("/patientReviews/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPatientReviewEndpoint))).Methods("GET")
	router.Handle("/patientReviews/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePatientReviewEndpoint))).Methods("DELETE")
	router.Handle("/patientReviews/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePatientReviewEndPoint)))).Methods("PUT")

	/* physiological Constants */

	router.Handle("/physiologicalConstants", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPhysiologicalConstantsEndPoint)))).Methods("POST")
	router.Handle("/physiologicalConstants", middleware.AuthMiddleware(http.HandlerFunc(allPhysiologicalConstantsEndPoint))).Methods("GET")
	router.Handle("/physiologicalConstants/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findPhysiologicalConstantsByPatientEndpoint))).Methods("GET")
	router.Handle("/physiologicalConstants/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPhysiologicalConstantsEndpoint))).Methods("GET")
	router.Handle("/physiologicalConstants/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePhysiologicalConstantsEndpoint))).Methods("DELETE")
	router.Handle("/physiologicalConstants/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePhysiologicalConstantsEndPoint)))).Methods("PUT")

	/*  diagnostic Plan */

	router.Handle("/diagnosticPlans", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createDiagnosticPlansEndPoint)))).Methods("POST")
	router.Handle("/diagnosticPlans", middleware.AuthMiddleware(http.HandlerFunc(allDiagnosticPlansEndPoint))).Methods("GET")
	router.Handle("/diagnosticPlans/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findDiagnosticPlansByPatientEndpoint))).Methods("GET")
	router.Handle("/diagnosticPlans/{id}", middleware.AuthMiddleware(http.HandlerFunc(findDiagnosticPlansEndpoint))).Methods("GET")
	router.Handle("/diagnosticPlans/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeDiagnosticPlansEndpoint))).Methods("DELETE")
	router.Handle("/diagnosticPlans/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateDiagnosticPlansEndPoint)))).Methods("PUT")

	/* therapeutic Plan */

	router.Handle("/therapeuticPlans", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createTherapeuticPlansEndPoint)))).Methods("POST")
	router.Handle("/therapeuticPlans", middleware.AuthMiddleware(http.HandlerFunc(allTherapeuticPlansEndPoint))).Methods("GET")
	router.Handle("/therapeuticPlans/{patient}", middleware.AuthMiddleware(http.HandlerFunc(findTherapeuticPlansByPatientEndpoint))).Methods("GET")
	router.Handle("/therapeuticPlans/{id}", middleware.AuthMiddleware(http.HandlerFunc(findTherapeuticPlansEndpoint))).Methods("GET")
	router.Handle("/therapeuticPlans/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeTherapeuticPlansEndpoint))).Methods("DELETE")
	router.Handle("/therapeuticPlans/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateTherapeuticPlansEndPoint)))).Methods("PUT")

	/* appointment */

	router.Handle("/appointments", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createAppointmentsEndPoint)))).Methods("POST")
	router.Handle("/appointments", middleware.AuthMiddleware(http.HandlerFunc(allAppointmentsEndPoint))).Methods("GET")
	router.Handle("/appointments/{id}", middleware.AuthMiddleware(http.HandlerFunc(findAppointmentsEndpoint))).Methods("GET")
	router.Handle("/appointments/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeAppointmentsEndpoint))).Methods("DELETE")
	router.Handle("/appointments/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateAppointmentsEndPoint)))).Methods("PUT")

	/* agendaAnnotations */

	router.Handle("/agendaAnnotations", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createAgendaAnnotationEndPoint)))).Methods("POST")
	router.Handle("/agendaAnnotations", middleware.AuthMiddleware(http.HandlerFunc(allAgendaAnnotationsEndPoint))).Methods("GET")
	router.Handle("/agendaAnnotations/{id}", middleware.AuthMiddleware(http.HandlerFunc(findAgendaAnnotationEndpoint))).Methods("GET")
	router.Handle("/agendaAnnotations/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeAgendaAnnotationEndpoint))).Methods("DELETE")
	router.Handle("/agendaAnnotations/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateAgendaAnnotationEndPoint)))).Methods("PUT")

	log.Fatal(http.ListenAndServe(":"+port, &CORSRouterDecorator{router}))
}
