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
	OverrideRequest bool              `json:"overriderequest" mapstructure:"overriderequest"`
}

func (RefTransHistory) TableName() string {
	return "aduit.refTransHistory"
}

type ViewIntakeReq struct {
	UserId        int `json:"userId" binding:"required" mapstructure:"userId"`
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type GetViewIntakeData struct {
	IntakeId         int    `json:"refITFId" gorm:"column:refITFId"`
	CategoryId       int    `json:"categoryId" gorm:"column:refCategoryId"`
	QuestionId       int    `json:"questionId" gorm:"column:refITFQId"`
	Answer           string `json:"answer" gorm:"column:refITFAnswer"`
	VerifyTechnician int    `json:"verifyTechnician" gorm:"column:refITFVerifiedTechnician"`
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
	UserId  int              `json:"refUserId" gorm:"column:refUserId"`
	Answers []UpdateReqModel `json:"answers" binding:"required" mapstructure:"answers"`
}
