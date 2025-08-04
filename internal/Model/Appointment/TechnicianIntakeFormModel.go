package model

import "time"

type UpdateAnswersReqModel struct {
	ITFId             int    `json:"refITFId" binding:"required" mapstructure:"refITFId"`
	QuestionId        int    `json:"questionId" binding:"required" mapstructure:"questionId"`
	Answer            string `json:"answer" binding:"required" mapstructure:"answer"`
	TechinicianStatus bool   `json:"verifyTechnician" mapstructure:"verifyTechnician"`
}

type Change struct {
	QuestionID int `json:"questionId"`
	// Add other fields here if needed
}

type DicomFile struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
	Side        string `json:"side" binding:"required" mapstructure:"side"`
}

type AddTechnicianIntakeFormReq struct {
	PatientId         int                     `json:"patientId" binding:"required" mapstructure:"patientId"`
	CategoryId        int                     `json:"categoryId" mapstructure:"categoryId"`
	AppointmentId     int                     `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	Priority          string                  `json:"priority" binding:"required" mapstructure:"priority"`
	UpdatedAnswers    []UpdateAnswersReqModel `json:"updatedAnswers" binding:"required" mapstructure:"updatedAnswers"`
	TechnicianAnswers []AnswersReqModel       `json:"technicianAnswers" binding:"required" mapstructure:"technicianAnswers"`
	// DicomFiles        []DicomFile             `json:"dicom_files" binding:"required" mapstructure:"dicom_files"`
}

type SaveDicomReq struct {
	PatientId     int         `json:"patientId" binding:"required" mapstructure:"patientId"`
	AppointmentId int         `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	DicomFiles    []DicomFile `json:"dicom_files" binding:"required" mapstructure:"dicom_files"`
}

type ViewTechnicianIntakeFormReq struct {
	PatientId     int `json:"patientId" binding:"required" mapstructure:"patientId"`
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type GetCategoryIdModel struct {
	CategoryId int `json:"refCategoryId" gorm:"column:refCategoryId"`
}

type DeleteDicomReq struct {
	DFId []int `json:"refDFId" binding:"required" mapstructure:"refDFId"`
}

type DicomFileModel struct {
	DFId          int       `gorm:"primaryKey;autoIncrement;column:refDFId"`
	UserId        int       `gorm:"column:refUserId"`
	AppointmentId int       `gorm:"column:refAppointmentId"`
	FileName      string    `gorm:"column:refDFFilename"`
	CreatedAt     time.Time `gorm:"column:refDFCreatedAt"`
	CreatedBy     int       `gorm:"column:refDFCreatedBy"`
	Side          string    `gorm:"column:refDFSide"`
}

func (DicomFileModel) TableName() string {
	return "dicom.refDicomFiles"
}

type DownloadDicomReq struct {
	FileId int `json:"fileId" binding:"required" mapstructure:"fileId"`
}

type OneDownloadDicomReq struct {
	UserId        int    `json:"UserId" binding:"required" mapstructure:"UserId"`
	AppointmentId int    `json:"AppointmentId" binding:"required" mapstructure:"AppointmentId"`
	Side          string `json:"Side" binding:"required" mapstructure:"Side"`
}

type TechIntakeModel struct {
	TITFId     int    `json:"refTITFId" gorm:"column:refTITFId"`
	TITFQId    int    `json:"questionId" gorm:"column:refTITFQId"`
	TITFAnswer string `json:"answer" gorm:"column:refTITFAnswer"`
}

type OverAllDicomModel struct {
	AppointmentId []int `json:"refAppointmentId" mapstructure:"refAppointmentId"`
}
