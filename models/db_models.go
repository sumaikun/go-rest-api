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
	Picture    string        `bson:"picture" json:"picture"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//Product representation on mongo
type Product struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	Name              string        `bson:"name" json:"name"`
	Value             string        `bson:"value" json:"value"`
	Description       string        `bson:"description" json:"description"`
	Picture           string        `bson:"picture" json:"picture"`
	Date              string        `bson:"date" json:"date"`
	UpdateDate        string        `bson:"update_date" json:"update_date"`
	AdministrationWay string        `bson:"administrationWay" json:"administrationWay"`
	Presentation      string        `bson:"presentation" json:"presentation"`
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
	Phone          string        `bson:"phone" json:"phone"`
	Ocupation      string        `bson:"ocupation" json:"ocupation"`
	Email          string        `bson:"email" json:"email"`
	Picture        string        `bson:"picture" json:"picture"`
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
	Date        string        `bson:"date" json:"date"`
	UpdateDate  string        `bson:"update_date" json:"update_date"`
}

//Breeds representation on mongo
type Breeds struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Species    string        `bson:"species" json:"species"`
	Meta       string        `bson:"meta" json:"meta"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//Species representation on mongo
type Species struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Meta       string        `bson:"meta" json:"meta"`
	Date       string        `bson:"date" json:"date"`
	UpdateDate string        `bson:"update_date" json:"update_date"`
}

//PatientReview representation on mongo
type PatientReview struct {
	ID                     bson.ObjectId `bson:"_id" json:"id"`
	Patient                string        `bson:"patient" json:"patient"`
	PvcVaccine             bool          `bson:"pvcVaccine" json:"pvcVaccine"`
	PvcVaccineDate         string        `bson:"pvcVaccineDate" json:"pvcVaccineDate"`
	TripleVaccine          bool          `bson:"tipleVaccine" json:"tripleVaccine"`
	TripleVaccineDate      string        `bson:"tipleVaccineDate" json:"tripleVaccineDate"`
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
	TLIC                          string        `bson:"tlic" json:"tlic"`
	HeartRate                     string        `bson:"heartRate" json:"heartRate"`
	RespiratoryRate               string        `bson:"respiratoryRate" json:"respiratoryRate"`
	HeartBeat                     string        `bson:"HeartBeat" json:"HeartBeat"`
	Temperature                   string        `bson:"temperature" json:"temperature"`
	Weight                        string        `bson:"weight" json:"weight"`
	Attitude                      string        `bson:"attitude" json:"attitude"`
	BodyCondition                 string        `bson:"bodyCondition" json:"bodyCondition"`
	HidrationStatus               string        `bson:"hidrationStatus" json:"hidrationStatus"`
	ConjuntivalMucosa             string        `bson:"conjuntivalMucosa" json:"conjuntivalMucosa"`
	OralMucosa                    string        `bson:"oralMucosa" json:"oralMucosa"`
	VulvalMucosa                  string        `bson:"vulvalMucosa" json:"vulvallMucosa"`
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

//DiagnosticPlans representation on mongo
type DiagnosticPlans struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	TypeOfExam   string        `bson:"typeOfExam" json:"typeOfExam"`
	Descrjption  string        `bson:"description" json:"description"`
	RealizedData string        `bson:"realizedData" json:"realizedData"`
	Laboratory   string        `bson:"laboratory" json:"laboratory"`
	Results      string        `bson:"results" json:"results"`
	CreatedBy    string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy    string        `bson:"updatedBy" json:"updatedBy"`
	Date         string        `bson:"date" json:"date"`
	UpdateDate   string        `bson:"update_date" json:"update_date"`
}

//TerapeuticPlans representation on mongo
type TerapeuticPlans struct {
	ID                          bson.ObjectId `bson:"_id" json:"id"`
	TypeOfPlsn                  string        `bson:"typeOfPlan" json:"typeOfPlan"`
	ActiveSubstanceToAdminister string        `bson:"activeSubstanceToAdminister " json:"activeSubstanceToAdminister "`
	Posology                    string        `bson:"posology" json:"posology"`
	TotalDose                   string        `bson:"totalDose" json:"totalDose"`
	FrecuencyAndDuration        string        `bson:"frecuencyAndDuration" json:"frecuencyAndDuration"`
	CreatedBy                   string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy                   string        `bson:"updatedBy" json:"updatedBy"`
	Date                        string        `bson:"date" json:"date"`
	UpdateDate                  string        `bson:"update_date" json:"update_date"`
}

//Appointments representation on mongo
type Appointments struct {
	ID                     bson.ObjectId `bson:"_id" json:"id"`
	ReasonForConsultation  string        `bson:"reasonForConsultation" json:"reasonForConsultation"`
	ResultsForConsultation string        `bson:"resultsForConsultation" json:"resultsForConsultation"`
	RelatedProcesses       string        `bson:"relatedProcesses" json:"relatedProcesses"`
	Product                string        `bson:"product" json:"product"`
	CreatedBy              string        `bson:"createdBy" json:"createdBy"`
	UpdatedBy              string        `bson:"updatedBy" json:"updatedBy"`
	Date                   string        `bson:"date" json:"date"`
	UpdateDate             string        `bson:"update_date" json:"update_date"`
}
