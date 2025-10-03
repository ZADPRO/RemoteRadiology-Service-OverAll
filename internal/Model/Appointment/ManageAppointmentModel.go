package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type AddAppointmentReq struct {
	SCId            string `json:"refSCId" binding:"required" mapstructure:"refSCId"`
	AppointmentDate string `json:"refAppointmentDate" binding:"required" mapstructure:"refAppointmentDate"`
	// AppointmentStartTime string `json:"refAppointmentStartTime" binding:"required" mapstructure:"refAppointmentStartTime"`
	// AppointmentEndTime   string `json:"refAppointmentEndTime" binding:"required" mapstructure:"refAppointmentEndTime"`
	// AppointmentUrgency   bool   `json:"refAppointmentUrgency" mapstructure:"refAppointmentUrgency"`
}

type CreateAppointmentModel struct {
	AppointmentId   int    `json:"refAppointmentId" gorm:"primaryKey;autoIncrement;column:refAppointmentId"`
	UserId          int    `json:"refUserId" gorm:"column:refUserId"`
	SCId            int    `json:"refSCId" gorm:"column:refSCId"`
	AppointmentDate string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	// AppointmentStartTime string `json:"refAppointmentStartTime" gorm:"column:refAppointmentStartTime"`
	// AppointmentEndTime   string `json:"refAppointmentEndTime" gorm:"column:refAppointmentEndTime"`
	AppointmentUrgency  bool   `json:"refAppointmentUrgency" gorm:"column:refAppointmentUrgency"`
	AppointmentStatus   bool   `json:"refAppointmentStatus" gorm:"column:refAppointmentStatus"`
	AppointmentComplete string `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
}

func (CreateAppointmentModel) TableName() string {
	return "appointment.refAppointments"
}

type TotalCountModel struct {
	TotalCount int `json:"TotalCount" gorm:"column:TotalCount"`
}

type RefTransHistory struct {
	TransTypeId int    `json:"transTypeId" gorm:"column:transTypeId"`
	THData      string `json:"refTHData" gorm:"column:refTHData"`
	UserId      int    `json:"refUserId" gorm:"column:refUserId"`
	THActionBy  int    `json:"refTHActionBy" gorm:"column:refTHActionBy"`
}

type ScanCenterModel struct {
	SCId       int    `json:"refSCId" gorm:"column:refSCId"`
	SCCustId   string `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCAddress  string `json:"refSCAddress" gorm:"column:refSCAddress"`
	ProfileImg string `json:"refSCProfile" gorm:"column:refSCProfile"`
}

type GetAllUserDetailsModel struct {
	User_Id          string `json:"User_Id" gorm:"column:User_Id"`
	Assigned_Id      string `json:"Assigned_Id" gorm:"column:Assigned_Id"`
	Patient_Id       string `json:"Patient_Id" gorm:"column:Patient_Id"`
	Appointment_date string `json:"appointment_date" gorm:"column:appointment_date"`
}

type ViewPatientHistoryModel struct {
	AppointmentId       int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	SCId                int    `json:"refSCId" gorm:"column:refSCId"`
	CategoryId          int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	AppointmentDate     string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	Remarks             string `json:"refAppointmentRemarks" gorm:"column:refAppointmentRemarks"`
	AppointmentComplete string `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
	CustSCId            string `json:"refSCCustId" gorm:"column:refSCCustId"`
	OverrideStatus      string `json:"refOverrideStatus" gorm:"column:refOverrideStatus"`
	UserId              int    `json:"refUserId" gorm:"column:refUserId"`
	OldReportCount      string `json:"OldReportCount" gorm:"column:OldReportCount"`
}

type GetDicomFile struct {
	RefDFId          int    `json:"refDFId" gorm:"column:refDFId"`
	RefUserId        int    `json:"refUserId" gorm:"column:refUserId"`
	RefAppointmentId int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	RefDFFilename    string `json:"refDFFilename" gorm:"column:refDFFilename"`
	RefDFCreatedAt   string `json:"refDFCreatedAt" gorm:"column:refDFCreatedAt"`
	RefDFSide        string `json:"refDFSide" gorm:"column:refDFSide"`
}

type DicomFileArray []GetDicomFile

// Scan implements the sql.Scanner interface for DicomFileArray.
func (a *DicomFileArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	return json.Unmarshal(bytes, a)
}

// Value implements the driver.Valuer interface for DicomFileArray.
func (a DicomFileArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

type ViewTechnicianPatientQueueModel struct {
	AppointmentId              int                 `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	Remarks                    string              `json:"refAppointmentRemarks" gorm:"column:refAppointmentRemarks"`
	AppointmentDate            string              `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	Username                   string              `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	UserCustId                 string              `json:"refUserCustId" gorm:"column:refUserCustId"`
	CategoryId                 int                 `json:"refCategoryId" gorm:"column:refCategoryId"`
	UserId                     int                 `json:"refUserId" gorm:"column:refUserId"`
	AssignedUserId             int                 `json:"refAppointmentAssignedUserId" gorm:"column:refAppointmentAssignedUserId"`
	AppointmentComplete        string              `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
	ScanCenterCustId           string              `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCAddress                  string              `json:"refSCAddress" gorm:"column:refSCAddress"`
	ScanCenterId               string              `json:"refSCId" gorm:"column:refSCId"`
	DicomFiles                 DicomFileArray      `json:"dicomFiles" gorm:"column:dicomFiles"`
	OverrideStatus             string              `json:"refOverrideStatus" gorm:"column:refOverrideStatus"`
	AppointmentMailSendStatus  string              `json:"refAppointmentMailSendStatus" gorm:"column:refAppointmentMailSendStatus"`
	OldReportCount             string              `json:"OldReportCount" gorm:"column:OldReportCount"`
	GetCorrectEditModel        GetCorrectEditModel `json:"GetCorrectEditModel" gorm:"-"`
	PatientPrivatePublicStatus string              `json:"patientPrivatePublicStatus" gorm:"-"`
	ReportStatus               string              `json:"reportStatus"  gorm:"-"`
}

type ReportUrgentStatusModel struct {
	ReportStatus string `json:"reportStatus" gorm:"column:refTITFAnswer"`
}
type StaffAvailableModel struct {
	UserId     int    `json:"refUserId" gorm:"column:refUserId"`
	RoleId     int    `json:"refRTId" gorm:"column:refRTId"`
	UserCustId string `json:"refUserCustId" gorm:"column:refUserCustId"`
	Username   string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	SCId       int    `json:"refSCId" gorm:"column:refSCId"`
}

type AdditionalFileModel struct {
	FileName    string `json:"refFileName" binding:"required" mapstructure:"refFileName"`
	OldFileName string `json:"refOldFileName" binding:"required" mapstructure:"refOldFileName"`
}

type AddAddtionalFilesReq struct {
	AppointmentId int                   `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	Files         []AdditionalFileModel `json:"files" binding:"required" mapstructure:"files"`
}

type AdditionalFileUploadModel struct {
	ADId          int       `json:"refADId" gorm:"column:refADId"`
	UserId        int       `json:"refUserId" gorm:"column:refUserId"`
	AppointmentId int       `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	FileName      string    `json:"refADFileName"  gorm:"column:refADFileName"`
	OldFileName   string    `json:"refADOldFileName"  gorm:"column:refADOldFileName"`
	FileData      *FileData `json:"fileData" gorm:"column:-"`
	Status        bool      `json:"refADStatus"  gorm:"column:refADStatus"`
	CreatedAt     string    `json:"refADCreatedAt"  gorm:"column:refADCreatedAt"`
}

type ViewAddtionalFileReq struct {
	UserId        int `json:"refUserId" binding:"required" mapstructure:"refUserId"`
	AppointmentId int `json:"refAppointmentId" binding:"required" mapstructure:"refAppointmentId"`
}

type AssignUserReq struct {
	AssingUserId     int    `json:"assingUserId" binding:"required" mapstructure:"assingUserId"`
	AssingUserCustId string `json:"assingUsercustId" binding:"required" mapstructure:"assingUsercustId"`
	AppointmentId    int    `json:"refAppointmentId" binding:"required" mapstructure:"refAppointmentId"`
	PatientId        int    `json:"patientId" binding:"required" mapstructure:"patientId"`
}

type CorrectEditResponse struct {
	CorrectStatus bool `json:"correctStatus" gorm:"column:correctStatus"`
	EditStatus    bool `json:"editStatus" gorm:"column:editStatus"`
}

type GetCorrectEditModel struct {
	RHHandleCorrect bool `json:"isHandleCorrect" gorm:"column:isHandleCorrect"`
	RHHandleEdit    bool `json:"isHandleEdited" gorm:"column:isHandleEdited"`
}

type Notification struct {
	RefNId        int    `json:"refNId" gorm:"column:refNId"`
	RefUserId     int    `json:"refUserId" gorm:"column:refUserId"`
	RefNMessage   string `json:"refNMessage" gorm:"column:refNMessage"`
	RefNStatus    *bool  `json:"refNStatus" gorm:"column:refNStatus"`
	RefNCreatedAt string `json:"refNCreatedAt" gorm:"column:refNCreatedAt"`
}

type RefAuditTransHistory struct {
	RefTHId       int    `json:"refTHId" gorm:"column:refTHId;primaryKey"`
	TransTypeId   int    `json:"transTypeId" gorm:"column:transTypeId"`
	RefTHData     string `json:"refTHData" gorm:"column:refTHData"`
	RefTHTime     string `json:"refTHTime" gorm:"column:refTHTime"` // can use time.Time instead of string
	RefUserId     int    `json:"refUserId" gorm:"column:refUserId"`
	RefTHActionBy int    `json:"refTHActionBy" gorm:"column:refTHActionBy"`
	RefUserCustId string `json:"refUserCustId" gorm:"column:refUserCustId"`
	// RefUserFirstName string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
}

type ScanCenterConsultantModel struct {
	SCId               uint   `json:"refSCId" gorm:"column:refSCId"`
	SCCustId           string `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCProfile          string `json:"refSCProfile" gorm:"column:refSCProfile"`
	SCName             string `json:"refSCName" gorm:"column:refSCName"`
	SCAddress          string `json:"refSCAddress" gorm:"column:refSCAddress"`
	SCPhoneNo1         string `json:"refSCPhoneNo1" gorm:"column:refSCPhoneNo1"`
	SCEmail            string `json:"refSCEmail" gorm:"column:refSCEmail"`
	SCWebsite          string `json:"refSCWebsite" gorm:"column:refSCWebsite"`
	SCAppointments     bool   `json:"refSCAppointments" gorm:"column:refSCAppointments"`
	SCDisclamer        string `json:"refSCDisclamer" gorm:"column:refSCDisclamer"`
	SCBrouchure        string `json:"refSCBrouchure" gorm:"column:refSCBrouchure"`
	SCGuidelines       string `json:"refSCGuidelines" gorm:"column:refSCGuidelines"`
	SCStatus           bool   `json:"refSCStatus" gorm:"column:refSCStatus"`
	SCConsultantStatus bool   `json:"refSCConsultantStatus" gorm:"column:refSCConsultantStatus"`
}

type MapScanCenterPatientModel struct {
	SCMPId int `json:"refSCMPId" gorm:"primaryKey;autoIncrement;column:refSCMPId"`
	UserId int `json:"refUserId" gorm:"column:refUserId"`
	SCId   int `json:"refSCId" gorm:"column:refSCId"`
}

