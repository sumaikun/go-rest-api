package main

import (
	Config "click-al-vet/config"
	Dao "click-al-vet/dao"
	middleware "click-al-vet/middlewares"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port        string
	jwtKey      []byte
	logo        string
	frontEndUrl string
)

var dao = Dao.MongoConnector{}

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

//------------------- INIT function -------------------------------------------------

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
	init_validators()
}

//------------------- main function -------------------------------------------------

func main() {
	fmt.Println("start server in port " + port)
	router := mux.NewRouter().StrictSlash(true)
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
	//router.Handle("/doctors/{id}", middleware.AuthMiddleware(http.HandlerFunc(inactivateDoctorEndPoint))).Methods("DELETE")
	router.Handle("/doctors/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateDoctorEndPoint)))).Methods("PUT")

	/* Products Routes */
	router.Handle("/products", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createProductEndPoint)))).Methods("POST")
	router.Handle("/products", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(allProductsEndPoint)))).Methods("GET")
	router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(findProductEndpoint))).Methods("GET")
	//router.Handle("/products/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeProductEndpoint))).Methods("DELETE")
	router.Handle("/products/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateProductEndPoint)))).Methods("PUT")

	/* Contacts Routes */
	router.Handle("/contacts", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createContactEndPoint)))).Methods("POST")
	router.Handle("/contacts", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(allContactsEndPoint)))).Methods("GET")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(findContactEndpoint))).Methods("GET")
	//router.Handle("/contacts/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeContactEndpoint))).Methods("DELETE")
	router.Handle("/contacts/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateContactEndPoint)))).Methods("PUT")

	/* Pets Routes */
	router.Handle("/pets", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPetEndPoint)))).Methods("POST")
	router.Handle("/pets", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(allPetsEndPoint)))).Methods("GET")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPetEndpoint))).Methods("GET")
	//router.Handle("/pets/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePetEndpoint))).Methods("DELETE")
	router.Handle("/updatePetContactsEndPoint/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePetContactsEndPoint)))).Methods("PUT")
	router.Handle("/pets/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePetEndPoint)))).Methods("PUT")

	/* Examtypes Routes */
	router.Handle("/examTypes", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createParameterEndPoint)))).Methods("POST")
	router.Handle("/examTypes", middleware.AuthMiddleware(http.HandlerFunc(allParametersEndPoint))).Methods("GET")
	router.Handle("/examTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(findParameterEndPoint))).Methods("GET")
	//router.Handle("/examTypes/{id}", middleware.AuthMiddleware(http.HandlerFunc(deleteParameterEndPoint))).Methods("DELETE")
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
	router.Handle("/patientReviews/{pet}", middleware.AuthMiddleware(http.HandlerFunc(findPatientReviewByPatientEndpoint))).Methods("GET")
	//router.Handle("/patientReview/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPatientReviewEndpoint))).Methods("GET")
	//router.Handle("/patientReviews/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePatientReviewEndpoint))).Methods("DELETE")
	router.Handle("/patientReviews/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePatientReviewEndPoint)))).Methods("PUT")

	/* physiological Constants */

	router.Handle("/physiologicalConstants", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPhysiologicalConstantsEndPoint)))).Methods("POST")
	router.Handle("/physiologicalConstants", middleware.AuthMiddleware(http.HandlerFunc(allPhysiologicalConstantsEndPoint))).Methods("GET")
	router.Handle("/physiologicalConstants/{pet}", middleware.AuthMiddleware(http.HandlerFunc(findPhysiologicalConstantsByPatientEndpoint))).Methods("GET")
	//router.Handle("/physiologicalConstant/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPhysiologicalConstantsEndpoint))).Methods("GET")
	//router.Handle("/physiologicalConstants/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePhysiologicalConstantsEndpoint))).Methods("DELETE")
	router.Handle("/physiologicalConstants/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePhysiologicalConstantsEndPoint)))).Methods("PUT")

	/* appointment */

	router.Handle("/appointments", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createAppointmentsEndPoint)))).Methods("POST")
	router.Handle("/appointments", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(allAppointmentsEndPoint)))).Methods("GET")
	router.Handle("/appointmentsByPatient/{pet}", middleware.AuthMiddleware(http.HandlerFunc(findAppointmentsByPatientEndpoint))).Methods("GET")
	router.Handle("/appointmentsByPatientAndDate/{pet}/{date}", middleware.AuthMiddleware(http.HandlerFunc(appointmentsByPatientAndDateEndPoint))).Methods("GET")
	router.Handle("/appointment/{id}", middleware.AuthMiddleware(http.HandlerFunc(findAppointmentsEndpoint))).Methods("GET")
	//router.Handle("/appointments/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeAppointmentsEndpoint))).Methods("DELETE")
	router.Handle("/appointments/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateAppointmentsEndPoint)))).Methods("PUT")
	router.Handle("/confirmAppointment/{email}/{appointment}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(confirmPatientAppointment)))).Methods("GET")
	router.Handle("/cancelAppointment/{email}/{appointment}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(cancelPatientAppointment)))).Methods("GET")

	/* medicines */

	router.Handle("/medicines", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createMedicinesEndPoint)))).Methods("POST")
	router.Handle("/medicinesByPatient/{pet}", middleware.AuthMiddleware(http.HandlerFunc(findMedicinesByPatientEndPoint))).Methods("GET")
	router.Handle("/medicinesByAppointment/{appointment}", middleware.AuthMiddleware(http.HandlerFunc(findMedicinesByAppointmentEndPoint))).Methods("GET")
	router.Handle("/medicines/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateMedicinesEndPoint)))).Methods("PUT")

	/* detectedDiseases */

	router.Handle("/detectedDiseases", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createDetectedDiseaseEndPoint)))).Methods("POST")
	router.Handle("/detectedDiseases", middleware.AuthMiddleware(http.HandlerFunc(allDetectedDiseasesEndPoint))).Methods("GET")
	router.Handle("/detectedDiseases/{pet}", middleware.AuthMiddleware(http.HandlerFunc(findDetectedDiseasesByPatientEndpoint))).Methods("GET")
	router.Handle("/detectedDisease/{id}", middleware.AuthMiddleware(http.HandlerFunc(findDetectedDiseaseEndpoint))).Methods("GET")
	//router.Handle("/detectedDiseases/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeDetectedDiseaseEndpoint))).Methods("DELETE")
	router.Handle("/detectedDiseases/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateDetectedDiseaseEndPoint)))).Methods("PUT")

	/* patientFiles */

	router.Handle("/patientFiles", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createPatientFilesEndPoint)))).Methods("POST")
	router.Handle("/patientFiles", middleware.AuthMiddleware(http.HandlerFunc(allPatientFilesEndPoint))).Methods("GET")
	router.Handle("/patientFiles/{pet}", middleware.AuthMiddleware(http.HandlerFunc(findPatientFilesByPatientEndpoint))).Methods("GET")
	router.Handle("/patientFile/{id}", middleware.AuthMiddleware(http.HandlerFunc(findPatientFilesEndpoint))).Methods("GET")
	//router.Handle("/patientFiles/{id}", middleware.AuthMiddleware(http.HandlerFunc(removePatientFilesEndpoint))).Methods("DELETE")
	router.Handle("/patientFiles/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updatePatientFilesEndPoint)))).Methods("PUT")

	/* agendaAnnotations */

	router.Handle("/agendaAnnotations", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createAgendaAnnotationEndPoint)))).Methods("POST")
	router.Handle("/agendaAnnotations", middleware.AuthMiddleware(http.HandlerFunc(allAgendaAnnotationsEndPoint))).Methods("GET")
	router.Handle("/agendaAnnotations/{pet}", middleware.AuthMiddleware(http.HandlerFunc(findPatientFilesByPatientEndpoint))).Methods("GET")
	router.Handle("/agendaAnnotation/{id}", middleware.AuthMiddleware(http.HandlerFunc(findAgendaAnnotationEndpoint))).Methods("GET")
	//router.Handle("/agendaAnnotations/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeAgendaAnnotationEndpoint))).Methods("DELETE")
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

	/* Doctors Settings Routes */

	router.Handle("/doctorSettings", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(createDoctorSettingEndPoint)))).Methods("POST")
	router.HandleFunc("/doctorSettings", allDoctorSettingsEndPoint).Methods("GET")
	router.Handle("/doctorSettingsByDoctor/{doctor}", middleware.AuthMiddleware(http.HandlerFunc(findDoctorSettingsByDoctorEndPoint))).Methods("GET")
	router.Handle("/doctorSettings/{id}", middleware.AuthMiddleware(http.HandlerFunc(findDoctorSettingsEndPoint))).Methods("GET")
	//router.Handle("/doctorSettings/{id}", middleware.AuthMiddleware(http.HandlerFunc(removeDoctorSettingsEndPoint))).Methods("DELETE")
	router.Handle("/doctorSettings/{id}", middleware.AuthMiddleware(middleware.UserMiddleware(http.HandlerFunc(updateDoctorSettingsEndPoint)))).Methods("PUT")
	log.Fatal(http.ListenAndServe(":"+port, &CORSRouterDecorator{router}))
}
