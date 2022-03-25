package models

import "gopkg.in/mgo.v2/bson"

//User representation on mongo
type User struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Password   string        `bson:"password" json:"password"`
	Email      string        `bson:"email" json:"email"`
	Address    string        `bson:"address" json:"address"`
	Role       string        `bson:"role" json:"role"`
	Phone      string        `bson:"phone" json:"phone"`
	State      string        `bson:"state" json:"state"`
	Picture    string        `bson:"picture" json:"picture"`
	CreatedBy  string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy  string        `bson:"updatedBy" json:"updatedBy"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//Doctor representation on mongo
type Doctor struct {
	ID             bson.ObjectId `bson:"_id" json:"id"`
	Name           string        `bson:"name" json:"name"`
	LastName       string        `bson:"lastName" json:"lastName"`
	City           string        `bson:"city" json:"city"`
	SpecialistType []string      `bson:"specialistType" json:"specialistType"`
	BirthDate      string        `bson:"birthDate" json:"birthDate"`
	Password       string        `bson:"password" json:"password"`
	Email          string        `bson:"email" json:"email"`
	Address        string        `bson:"address" json:"address"`
	Phone          string        `bson:"phone" json:"phone"`
	Phone2         string        `bson:"phone2" json:"phone2"`
	TypeID         string        `bson:"typeId" json:"typeId"`
	Identification string        `bson:"identification" json:"identification"`
	Picture        string        `bson:"picture" json:"picture"`
	State          string        `bson:"state" json:"state"`
	Confirmed      bool          `bson:"confirmed" json:"confirmed"`
	CreatedBy      string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy      string        `bson:"updatedBy" json:"updatedBy"`
	Date           string        `bson:"date" json:"date"`
	UpdateDate     string        `bson:"update_date" json:"update_date"`
	MedicalCenter  string        `bson:"medicalCenter" json:"medicalCenter"`
	Qualification  string        `bson:"qualification" json:"qualification"`
	AboutDoctor    string        `bson:"aboutDoctor" json:"aboutDoctor"`
}

//SpecialistTypes representation on mongo
type SpecialistTypes struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Meta       string        `bson:"meta" json:"meta"`
	CreatedBy  string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy  string        `bson:"updatedBy" json:"updatedBy"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//CitiesTypes representation on mongo
type CitiesTypes struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Meta       string        `bson:"meta" json:"meta"`
	CreatedBy  string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy  string        `bson:"updatedBy" json:"updatedBy"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//Product representation on mongo
type Product struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Value       string        `bson:"value" json:"value"`
	Description string        `bson:"description" json:"description"`
	Picture     string        `bson:"picture" json:"picture"`
	Doctors     []string      `bson:"doctors" json:"doctors"`
	Date        string        `bson:"date" json:"date"`
	UpdateDate  string        `bson:"update_date" json:"update_date"`
	CreatedBy   string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy   string        `bson:"updatedBy" json:"updatedBy"`
}

//Contact representation on mongo
type Contact struct {
	ID             bson.ObjectId `bson:"_id" json:"id"`
	Name           string        `bson:"name" json:"name"`
	Address        string        `bson:"address" json:"address"`
	TypeID         string        `bson:"typeId" json:"typeId"`
	Identification string        `bson:"identification" json:"identification"`
	Stratus        string        `bson:"stratus" json:"stratus"`
	City           string        `bson:"city" json:"city"`
	Password       string        `bson:"password" json:"password"`
	Phone          string        `bson:"phone" json:"phone"`
	Ocupation      string        `bson:"ocupation" json:"ocupation"`
	Email          string        `bson:"email" json:"email"`
	Doctors        []string      `bson:"doctors" json:"doctors"`
	Picture        string        `bson:"picture" json:"picture"`
	Confirmed      bool          `bson:"confirmed" json:"confirmed"`
	CreatedBy      string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy      string        `bson:"updatedBy" json:"updatedBy"`
	Date           string        `bson:"date" json:"date"`
	UpdateDate     string        `bson:"update_date" json:"update_date"`
}

//Pet representation on mongo
type Pet struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Species     string        `bson:"species" json:"species"`
	Breed       string        `bson:"breed" json:"breed"`
	Color       string        `bson:"color" json:"color"`
	Sex         string        `bson:"sex" json:"sex"`
	BirthDate   string        `bson:"birthDate" json:"birthDate"`
	Age         string        `bson:"age" json:"age"`
	Origin      string        `bson:"origin" json:"origin"`
	Description string        `bson:"description" json:"description"`
	Picture     string        `bson:"picture" json:"picture"`
	Doctors     []string      `bson:"doctors" json:"doctors"`
	Contacts    []string      `bson:"contacts" json:"contacts"`
	CreatedBy   string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy   string        `bson:"updatedBy" json:"updatedBy"`
	Date        string        `bson:"date" json:"date"`
	UpdateDate  string        `bson:"update_date" json:"update_date"`
}

//Diseases representation on mongo
type Diseases struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Meta       string        `bson:"meta" json:"meta"`
	CreatedBy  string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy  string        `bson:"updatedBy" json:"updatedBy"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//PatientReview representation on mongo
type PatientReview struct {
	ID                     bson.ObjectId `bson:"_id" json:"id"`
	Patient                string        `bson:"patient" json:"patient"`
	PvcVaccine             bool          `bson:"pvcVaccine" json:"pvcVaccine"`
	PvcVaccineDate         string        `bson:"pvcVaccineDate" json:"pvcVaccineDate"`
	TripleVaccine          bool          `bson:"tripleVaccine" json:"tripleVaccine"`
	TripleVaccineDate      string        `bson:"tripleVaccineDate" json:"tripleVaccineDate"`
	RabiesVaccine          bool          `bson:"rabiesVaccine" json:"rabiesVaccine"`
	RabiesVaccineDate      string        `bson:"rabiesVaccineDate" json:"rabiesVaccineDate"`
	DesparasitationProduct string        `bson:"desparasitationProduct" json:"desparasitationProduct"`
	LastDesparasitation    string        `bson:"lastDesparasitation" json:"lastDesparasitation"`
	FeedingType            string        `bson:"feedingType" json:"feedingType"`
	ReproductiveState      string        `bson:"reproductiveState" json:"reproductiveState"`
	PreviousIllnesses      string        `bson:"previousIllnesses" json:"previousIllnesses"`
	Surgeris               string        `bson:"surgeris" json:"surgeris"`
	FamilyBackground       string        `bson:"familyBackground" json:"familyBackground"`
	Habitat                string        `bson:"habitat" json:"habitat"`
	CreatedBy              string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy              string        `bson:"updatedBy" json:"updatedBy"`
	Date                   string        `bson:"date" json:"date"`
	UpdateDate             string        `bson:"update_date" json:"update_date"`
}

//PhysiologicalConstants representation on mongo
type PhysiologicalConstants struct {
	ID                            bson.ObjectId `bson:"_id" json:"id"`
	Patient                       string        `bson:"patient" json:"patient"`
	TLIC                          string        `bson:"tlic" json:"tlic"`
	HeartRate                     string        `bson:"heartRate" json:"heartRate"`
	RespiratoryRate               string        `bson:"respiratoryRate" json:"respiratoryRate"`
	HeartBeat                     string        `bson:"heartBeat" json:"heartBeat"`
	Temperature                   string        `bson:"temperature" json:"temperature"`
	Weight                        string        `bson:"weight" json:"weight"`
	Attitude                      string        `bson:"attitude" json:"attitude"`
	BodyCondition                 string        `bson:"bodyCondition" json:"bodyCondition"`
	HidrationStatus               string        `bson:"hidrationStatus" json:"hidrationStatus"`
	ConjuntivalMucosa             string        `bson:"conjuntivalMucosa" json:"conjuntivalMucosa"`
	OralMucosa                    string        `bson:"oralMucosa" json:"oralMucosa"`
	VulvalMucosa                  string        `bson:"vulvalMucosa" json:"vulvalMucosa"`
	RectalMucosa                  string        `bson:"rectalMucosa" json:"rectalMucosa"`
	PhysicalsEye                  string        `bson:"physicalsEye" json:"physicalsEye"`
	PhysicalsEars                 string        `bson:"physicalsEars" json:"physicalsEars"`
	PhysicalsLinfaticmodules      string        `bson:"physicalsLinfaticmodules" json:"physicalsLinfaticmodules"`
	PhysicalsSkinandanexes        string        `bson:"physicalsSkinandanexes" json:"physicalsSkinandanexes"`
	PhysicalsLocomotion           string        `bson:"physicalsLocomotion" json:"physicalsLocomotion"`
	PhysicalsMusclesqueletal      string        `bson:"physicalsMusclesqueletal" json:"physicalsMusclesqueletal"`
	PhysicalsNervoussystem        string        `bson:"physicalsNervoussystem" json:"physicalsNervoussystem"`
	PhysicalsCardiovascularsystem string        `bson:"physicalsCardiovascularsystem" json:"physicalsCardiovascularsystem"`
	PhysicalsRespiratorysystem    string        `bson:"physicalsRespiratorysystem" json:"physicalsRespiratorysystem"`
	PhysicalsDigestivesystem      string        `bson:"physicalsDigestivesystem" json:"physicalsDigestivesystem"`
	PhysicalsGenitourinarysystem  string        `bson:"physicalsGenitourinarysystem" json:"physicalsGenitourinarysystem"`
	CreatedBy                     string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy                     string        `bson:"updatedBy" json:"updatedBy"`
	Date                          string        `bson:"date" json:"date"`
	UpdateDate                    string        `bson:"update_date" json:"update_date"`
}

//Appointments representation on mongo
type Appointments struct {
	ID                           bson.ObjectId `bson:"_id" json:"id"`
	Patient                      string        `bson:"patient" json:"patient"`
	Doctor                       string        `bson:"doctor" json:"doctor"`
	ReasonForConsultation        string        `bson:"reasonForConsultation" json:"reasonForConsultation"`
	ResultsForConsultation       string        `bson:"resultsForConsultation" json:"resultsForConsultation"`
	MedicalReasonForConsultation string        `bson:"medicalReasonForConsultation" json:"medicalReasonForConsultation"`
	AgendaAnnotation             string        `bson:"agendaAnnotation" json:"agendaAnnotation"`
	AppointmentDate              string        `bson:"appointmentDate" json:"appointmentDate"`
	ExamType                     string        `bson:"examType" json:"examType"`
	ExamDescription              string        `bson:"examDescription" json:"examDescription"`
	Laboratory                   string        `bson:"laboratory" json:"laboratory"`
	LaboratoryAddress            string        `bson:"laboratoryAddress" json:"laboratoryAddress"`
	State                        string        `bson:"state" json:"state"`
	Qualification                string        `bson:"qualification" json:"qualification"`
	PatientComments              string        `bson:"patientComments" json:"patientComments"`
	Products                     []string      `bson:"products" json:"products"`
	CreatedBy                    string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy                    string        `bson:"updatedBy" json:"updatedBy"`
	Date                         string        `bson:"date" json:"date"`
	UpdateDate                   string        `bson:"update_date" json:"update_date"`
}

//DetectedDiseases  representation on mongo
type DetectedDiseases struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Patient      string        `bson:"patient" json:"patient"`
	Doctor       string        `bson:"doctor" json:"doctor"`
	Disease      string        `bson:"disease" json:"disease"`
	Criteria     string        `bson:"criteria" json:"criteria"`
	Observations string        `bson:"observations" json:"observations"`
	CreatedBy    string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy    string        `bson:"updatedBy" json:"updatedBy"`
	Date         string        `bson:"date" json:"date"`
	UpdateDate   string        `bson:"update_date" json:"update_date"`
}

//PatientFiles  representation on mongo
type PatientFiles struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Patient     string        `bson:"patient" json:"patient"`
	Doctor      string        `bson:"doctor" json:"doctor"`
	Name        string        `bson:"name" json:"name"`
	FilePath    string        `bson:"filePath" json:"filePath"`
	Description string        `bson:"description" json:"description"`
	CreatedBy   string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy   string        `bson:"updatedBy" json:"updatedBy"`
	Date        string        `bson:"date" json:"date"`
	UpdateDate  string        `bson:"update_date" json:"update_date"`
}

//AgendaAnnotation  representation on mongo
type AgendaAnnotation struct {
	ID               bson.ObjectId `bson:"_id" json:"id"`
	AnnotationDate   string        `bson:"annotationDate" json:"annotationDate"`
	AnnotationToDate string        `bson:"annotationToDate" json:"annotationToDate"`
	Title            string        `bson:"title" json:"title"`
	Description      string        `bson:"description" json:"description"`
	Patient          string        `bson:"patient" json:"patient"`
	Doctor           string        `bson:"doctor" json:"doctor"`
	CreatedBy        string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy        string        `bson:"updatedBy" json:"updatedBy"`
	Date             string        `bson:"date" json:"date"`
	UpdateDate       string        `bson:"update_date" json:"update_date"`
}

//DoctorSettings  representation on mongo
type DoctorSettings struct {
	ID              bson.ObjectId `bson:"_id" json:"id"`
	HoursRange      []int         `bson:"hoursRange" json:"hoursRange"`
	DaysRange       []string      `bson:"daysRange" json:"daysRange"`
	IsScheduling    bool          `bson:"isScheduling" json:"isScheduling"`
	AppointmentTime string        `bson:"appointmentTime" json:"appointmentTime"`
	Doctor          string        `bson:"doctor" json:"doctor"`
	CreatedBy       string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy       string        `bson:"updatedBy" json:"updatedBy"`
	Date            string        `bson:"date" json:"date"`
	UpdateDate      string        `bson:"update_date" json:"update_date"`
}

//DoctorSchedule Response
type DoctorSchedule struct {
	Appointments []Appointments     `json:"appointments"`
	Annotations  []AgendaAnnotation `json:"annotation"`
}

//PatientAppointment representation on mongo
type PatientAppointment struct {
	Name            string `bson:"name" json:"name"`
	LastName        string `bson:"lastName" json:"lastName"`
	Doctor          string `bson:"doctor" json:"doctor"`
	TypeID          string `bson:"typeId" json:"typeId"`
	Identification  string `bson:"identification" json:"identification"`
	City            string `bson:"city" json:"city"`
	Phone           string `bson:"phone" json:"phone"`
	Ocupation       string `bson:"ocupation" json:"ocupation"`
	AppointmentDate string `bson:"appointmentDate" json:"appointmentDate"`
	Email           string `bson:"email" json:"email"`
}

//Medicines representation on mongo
type Medicines struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	Patient           string        `bson:"patient" json:"patient"`
	Appointment       string        `bson:"appointment" json:"appointment"`
	AdministrationWay string        `bson:"administrationWay" json:"administrationWay"`
	Duration          string        `bson:"duration" json:"duration"`
	Posology          string        `bson:"posology" json:"posology"`
	Presentation      string        `bson:"presentation" json:"presentation"`
	Product           string        `bson:"product" json:"product"`
	CreatedBy         string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy         string        `bson:"updatedBy" json:"updatedBy"`
	Date              string        `bson:"date" json:"date"`
	UpdateDate        string        `bson:"update_date" json:"update_date"`
}
