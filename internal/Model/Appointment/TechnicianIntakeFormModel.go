package model

import "time"

type UpdateAnswersReqModel struct {
	ITFId             int    `json:"ITFId" binding:"required" mapstructure:"ITFId"`
	QuestionId        int    `json:"questionId" binding:"required" mapstructure:"questionId"`
	Answer            string `json:"answer" binding:"required" mapstructure:"answer"`
	TechinicianStatus bool   `json:"techinicianStatus" binding:"required" mapstructure:"techinicianStatus"`
}

type Change struct {
	QuestionID int `json:"questionId"`
	// Add other fields here if needed
}

type DicomFile struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
}

type AddTechnicianIntakeFormReq struct {
	PatientId         int                     `json:"patientId" binding:"required" mapstructure:"patientId"`
	CategoryId        int                     `json:"categoryId" mapstructure:"categoryId"`
	AppointmentId     int                     `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	UpdatedAnswers    []UpdateAnswersReqModel `json:"updatedAnswers" binding:"required" mapstructure:"updatedAnswers"`
	TechnicianAnswers []AnswersReqModel       `json:"technicianAnswers" binding:"required" mapstructure:"technicianAnswers"`
	DicomFiles        []DicomFile             `json:"dicom_files" binding:"required" mapstructure:"dicom_files"`
}

type GetCategoryIdModel struct {
	CategoryId int `json:"refCategoryId" gorm:"column:refCategoryId"`
}

type DicomFileModel struct {
	DFId          int       `gorm:"primaryKey;autoIncrement;column:refDFId"`
	UserId        int       `gorm:"column:refUserId"`
	AppointmentId int       `gorm:"column:refAppointmentId"`
	FileName      string    `gorm:"column:refDFFilename"`
	CreatedAt     time.Time `gorm:"column:refDFCreatedAt"`
}

func (DicomFileModel) TableName() string {
	return "dicom.refDicomFiles"
}
