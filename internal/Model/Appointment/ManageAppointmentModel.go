package model

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
	SCId     int    `json:"refSCId" gorm:"column:refSCId"`
	SCCustId string `json:"refSCCustId" gorm:"column:refSCCustId"`
}

type ViewPatientHistoryModel struct {
	AppointmentId       int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	SCId                int    `json:"refSCId" gorm:"column:refSCId"`
	CategoryId          int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	AppointmentDate     string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	Remarks             string `json:"refAppointmentRemarks" gorm:"column:refAppointmentRemarks"`
	AppointmentComplete string `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
	CustSCId            string `json:"refSCCustId" gorm:"column:refSCCustId"`
}

type GetDicomFile struct {
	RefDFId          int    `json:"refDFId" gorm:"column:refDFId"`
	RefUserId        int    `json:"refUserId" gorm:"column:refUserId"`
	RefAppointmentId int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	RefDFFilename    string `json:"refDFFilename" gorm:"column:refDFFilename"`
	RefDFCreatedAt   string `json:"refDFCreatedAt" gorm:"column:refDFCreatedAt"`
	RefDFSide        string `json:"refDFSide" gorm:"column:refDFSide"`
}

type ViewTechnicianPatientQueueModel struct {
	AppointmentId       int                 `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	Remarks             string              `json:"refAppointmentRemarks" gorm:"column:refAppointmentRemarks"`
	AppointmentDate     string              `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	Username            string              `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	UserCustId          string              `json:"refUserCustId" gorm:"column:refUserCustId"`
	CategoryId          int                 `json:"refCategoryId" gorm:"column:refCategoryId"`
	UserId              int                 `json:"refUserId" gorm:"column:refUserId"`
	AssignedUserId      int                 `json:"refAppointmentAssignedUserId" gorm:"column:refAppointmentAssignedUserId"`
	AppointmentComplete string              `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
	ScanCenterCustId    string              `json:"refSCCustId" gorm:"column:refSCCustId"`
	DicomFiles          []GetDicomFile      `json:"dicomFiles" gorm:"-"`
	GetCorrectEditModel GetCorrectEditModel `json:"GetCorrectEditModel" gorm:"-"`
}

type StaffAvailableModel struct {
	UserId     int    `json:"refUserId" gorm:"column:refUserId"`
	UserCustId string `json:"refUserCustId" gorm:"column:refUserCustId"`
	Username   string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
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
