package model

type ReqVal struct {
	EncryptedData []string `json:"encryptedData"`
}

type AnswersReqModel struct {
	// CategoryId int    `json:"categoryId" binding:"required" mapstructure:"categoryId"`
	QuestionId int    `json:"questionId" binding:"required" mapstructure:"questionId"`
	Answer     string `json:"answer" binding:"required" mapstructure:"answer"`
}

type AddIntakeFormReq struct {
	CategoryId      int               `json:"categoryId" mapstructure:"categoryId"`
	AppointmentId   int               `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	Answers         []AnswersReqModel `json:"answers" binding:"required" mapstructure:"answers"`
	Consent         string            `json:"consent" binding:"required" mapstructure:"consent"`
	OverrideRequest bool              `json:"overriderequest" mapstructure:"overriderequest"`
}

func (RefTransHistory) TableName() string {
	return "aduit.refTransHistory"
}

type ViewIntakeReq struct {
	UserId        int `json:"userId" binding:"required" mapstructure:"userId"`
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type FileData struct {
	Base64Data  string `json:"base64Data"`  // base64-encoded file content
	ContentType string `json:"contentType"` // e.g., "image/jpeg"
}

type GetViewIntakeData struct {
	IntakeId         int       `json:"refITFId" gorm:"column:refITFId"`
	QuestionId       int       `json:"questionId" gorm:"column:refITFQId"`
	Answer           string    `json:"answer" gorm:"column:refITFAnswer"`
	File             *FileData `json:"file" gorm:"-"`
	VerifyTechnician bool      `json:"verifyTechnician" gorm:"column:refITFVerifiedTechnician"`
}

type OverrideRequestModel struct {
	OVId           int    `json:"refOVId" gorm:"primaryKey;autoIncrement;column:refOVId"`
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	AppointmentId  int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	ApprovedStatus string `json:"refApprovedStatus" gorm:"column:refApprovedStatus"`
}

func (OverrideRequestModel) TableName() string {
	return "notes.refOverRide"
}

type VerifyIntakeFormReq struct {
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type UpdateReqModel struct {
	ITFId              int    `json:"ITFId" binding:"required" mapstructure:"ITFId"`
	QuestionId         int    `json:"questionId" binding:"required" mapstructure:"questionId"`
	Answer             string `json:"answer" binding:"required" mapstructure:"answer"`
	VerifiedTechnician int    `json:"verifyTechnician" mapstructure:"verifyTechnician"`
}

type UpdateIntakeFormReq struct {
	PatientId     int              `json:"patientId" binding:"required" mapstructure:"patientId"`
	CategoryId    int              `json:"categoryId" mapstructure:"categoryId"`
	AppointmentId int              `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	UserId        int              `json:"refUserId" gorm:"column:refUserId"`
	Answers       []UpdateReqModel `json:"answers" binding:"required" mapstructure:"answers"`
}

type AduitModel struct {
	TransTypeId int    `json:"transTypeId" gorm:"column:transTypeId"`
	THData      string `json:"refTHData" gorm:"column:refTHData"`
	CreatedAt   string `json:"refTHTime" gorm:"column:refTHTime"`
	UserId      int    `json:"refUserId" gorm:"column:refUserId"`
	THActionBy  int    `json:"refTHActionBy" gorm:"column:refTHActionBy"`
}

type TechnicianModel struct {
	FirstName string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	CustId    string `json:"refUserCustId" gorm:"column:refUserCustId"`
}

type PatientReq struct {
	Id            int `json:"id" binding:"id"  mapstructure:"id"`
	AppointmentId int `json:"appintmentId" binding:"appintmentId"  mapstructure:"appintmentId"`
}

type PatientResponse struct {
	AppointmentId int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	RTCText       string `json:"refRTCText" gorm:"column:refRTCText"`
}

type GetViewReportReq struct {
	AppointmentId []int `json:"appintmentId" binding:"appintmentId"  mapstructure:"appintmentId"`
}
