package model

import "time"

type GetCheckAccessModel struct {
	Status bool `json:"status" gorm:"column:status"`
}

type CheckAccessReq struct {
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type AssignGetReportReq struct {
	AppointmentId  int  `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId      int  `json:"patientId" binding:"required" mapstructure:"patientId"`
	ReadOnlyStatus bool `json:"readOnly" binding:"required" mapstructure:"readOnly"`
}

type AccessStatus struct {
	Status                       bool   `gorm:"column:status"`
	RefAppointmentAccessId       int    `gorm:"column:refAppointmentAccessId"`
	CustID                       string `gorm:"column:userCustId"`
}

type GetTechnicianIntakeData struct {
	IntakeId   int    `json:"refTITFId" gorm:"column:refTITFId"`
	QuestionId int    `json:"questionId" gorm:"column:refTITFQId"`
	Answer     string `json:"answer" gorm:"column:refTITFAnswer"`
}

type GetReportIntakeData struct {
	IntakeId   int    `json:"refRITFId" gorm:"column:refRITFId"`
	QuestionId int    `json:"questionId" gorm:"column:refRITFQId"`
	Answer     string `json:"answer" gorm:"column:refRITFAnswer"`
}

type GetReportTextContent struct {
	IntakeId     int    `json:"refRTCId" gorm:"column:refRTCId"`
	TextContent  string `json:"refRTCText" gorm:"column:refRTCText"`
	RTSyncStatus bool   `json:"refRTSyncStatus" gorm:"column:refRTSyncStatus"`
}

type GetReportHistory struct {
	RHId              int    `json:"refRHId" gorm:"column:refRHId"`
	HandledUserId     int    `json:"refRHHandledUserId" gorm:"column:refRHHandledUserId"`
	HandleStartTime   string `json:"refRHHandleStartTime" gorm:"column:refRHHandleStartTime"`
	HandleEndTime     string `json:"refRHHandleEndTime" gorm:"column:refRHHandleEndTime"`
	HandleUserName    string `json:"HandleUserName" gorm:"column:HandleUserName"`
	HandleStatus      string `json:"refRHHandleStatus" gorm:"column:refRHHandleStatus"`
	HandleContentText string `json:"refRHHandleContentText" gorm:"column:refRHHandleContentText"`
}

type GetReportComments struct {
	RCId     int    `json:"refRCId" gorm:"column:refRCId"`
	RCFor    int    `json:"refRCFor" gorm:"column:refRCFor"`
	RCBy     int    `json:"refRCBy" gorm:"column:refRCBy"`
	Status   string `json:"refRCStatus" gorm:"column:refRCStatus"`
	Comments string `json:"refRCComments" gorm:"column:refRCComments"`
	UserFor  string `json:"UserForName" gorm:"column:UserForName"`
	UserBy   string `json:"UserByName" gorm:"column:UserByName"`
}

type AppointmentModel struct {
	AppointmentId                 int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	UserId                        int    `json:"refUserId" gorm:"column:refUserId"`
	SCId                          int    `json:"refSCId" gorm:"column:refSCId"`
	CategoryId                    int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	AppointmentDate               string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	AppointmentComplete           string `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
	AppointmentPriority           string `json:"refAppointmentPriority" gorm:"column:refAppointmentPriority"`
	AppointmentAccessId           int    `json:"refAppointmentAccessId" gorm:"column:refAppointmentAccessId"`
	AppointmentAccessStatus       bool   `json:"refAppointmentAccessStatus" gorm:"column:refAppointmentAccessStatus"`
	AppointmentScribeAccessId     int    `json:"refAppointmentScribeAccessId" gorm:"column:refAppointmentScribeAccessId"`
	AppointmentScribeAccessStatus bool   `json:"refAppointmentScribeAccessStatus" gorm:"column:refAppointmentScribeAccessStatus"`
}

type AnswerReportIntakeReq struct {
	AppointmentId int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	QuestionId    int    `json:"questionId" binding:"required" mapstructure:"questionId"`
	Answer        string `json:"answer" binding:"required" mapstructure:"answer"`
	PatientId     int    `json:"patientId" binding:"required" mapstructure:"patientId"`
}

type AnswerTextContentReq struct {
	AppointmentId int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	TextContent   string `json:"textContent" binding:"required" mapstructure:"textContent"`
	PatientId     int    `json:"patientId" binding:"required" mapstructure:"patientId"`
	SyncStatus    bool   `json:"syncStatus" binding:"required" mapstructure:"syncStatus"`
}

type GetReportTextContentModel struct {
	RTCId   int    `json:"refRTCId" gorm:"column:refRTCId"`
	RTCText string `json:"refRTCText" gorm:"column:refRTCText"`
}

type AddCommentReq struct {
	AppointmentId int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId     int    `json:"patientId" binding:"required" mapstructure:"patientId"`
	AssignId      int    `json:"assignto" binding:"required" mapstructure:"assignto"`
	Status        string `json:"status" binding:"required" mapstructure:"status"`
	Comments      string `json:"comments" binding:"required" mapstructure:"comments"`
}

type CompleteReportReq struct {
	AppointmentId int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId     int    `json:"patientId" binding:"required" mapstructure:"patientId"`
	MovedStatus   string `json:"movedStatus" binding:"required" mapstructure:"movedStatus"`
	CurrentStatus string `json:"currentStatus" binding:"required" mapstructure:"currentStatus"`
}

type AnswerReqModel struct {
	QuestionId int    `json:"questionId" binding:"required" mapstructure:"questionId"`
	Answer     string `json:"answer" binding:"required" mapstructure:"answer"`
}

type SubmitReportReq struct {
	ReportIntakeForm                    []AnswerReqModel `json:"reportIntakeForm" binding:"required" mapstructure:"reportIntakeForm"`
	ReportTextContent                   string           `json:"reportTextContent" binding:"required" mapstructure:"reportTextContent"`
	AppointmentId                       int              `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId                           int              `json:"patientId" binding:"required" mapstructure:"patientId"`
	MovedStatus                         string           `json:"movedStatus" binding:"required" mapstructure:"movedStatus"`
	CurrentStatus                       string           `json:"currentStatus" binding:"required" mapstructure:"currentStatus"`
	SyncStatus                          bool             `json:"syncStatus" mapstructure:"syncStatus"`
	EditStatus                          bool             `json:"editStatus" mapstructure:"editStatus"`
	Impression                          string           `json:"impression" mapstructure:"impression"`
	Recommendation                      string           `json:"recommendation" mapstructure:"recommendation"`
	PatientMailStatus                   bool             `json:"patientMailStatus" mapstructure:"patientMailStatus"`
	ManagerMailStatus                   bool             `json:"managerMailStatus" mapstructure:"managerMailStatus"`
	ImpressionAddtional                 string           `json:"impressionaddtional" mapstructure:"impressionaddtional"`
	RecommendationAddtional             string           `json:"recommendationaddtional" mapstructure:"recommendationaddtional"`
	CommonImpressionRecommendation      string           `json:"commonImpressionRecommendation" mapstructure:"commonImpressionRecommendation"`
	ImpressionRight                     string           `json:"impressionRight" mapstructure:"impressionRight"`
	RecommendationRight                 string           `json:"recommendationRight" mapstructure:"recommendationRight"`
	ImpressionAddtionalRight            string           `json:"impressionaddtionalRight" mapstructure:"impressionaddtionalRight"`
	RecommendationAddtionalRight        string           `json:"recommendationaddtionalRight" mapstructure:"recommendationaddtionalRight"`
	CommonImpressionRecommendationRight string           `json:"commonImpressionRecommendationRight" mapstructure:"commonImpressionRecommendationRight"`
	ScanCenterProfileImg                *FileData        `json:"scancenterProfileImg" gorm:"-"`
	LeaveStatus                         bool             `json:"leaveStatus" mapstructure:"leaveStatus"`
	ArtificatsLeft                      bool             `json:"artificatsLeft" mapstructure:"artificatsLeft"`
	ArtificatsRight                     bool             `json:"artificatsRight" mapstructure:"artificatsRight"`
	PatientHistory                      string           `json:"patienthistory" mapstructure:"patienthistory"`
	BreastImplantsImagetext             string           `json:"breastimplantImageText" mapstructure:"breastimplantImageText"`
	SymmetryImageText                   string           `json:"symmetryImageText" mapstructure:"symmetryImageText"`
	BreastdensityImageText              string           `json:"breastdensityImageText" mapstructure:"breastdensityImageText"`
	NippleAreolaImageText               string           `json:"nippleareolaImageText" mapstructure:"nippleareolaImageText"`
	GlandularImageText                  string           `json:"glandularImageText" mapstructure:"glandularImageText"`
	LymphnodesImageText                 string           `json:"lymphnodesImageText" mapstructure:"lymphnodesImageText"`
	BreastdensityImageTextLeft          string           `json:"breastdensityImageTextLeft" mapstructure:"breastdensityImageTextLeft"`
	NippleAreolaImageTextLeft           string           `json:"nippleareolaImageTextLeft" mapstructure:"nippleareolaImageTextLeft"`
	GlandularImageTextLeft              string           `json:"glandularImageTextLeft" mapstructure:"glandularImageTextLeft"`
	LymphnodesImageTextLeft             string           `json:"lymphnodesImageTextLeft" mapstructure:"lymphnodesImageTextLeft"`
}

type ChangedOneState struct {
	ReportQuestion                      []int `json:"reportQuestion"` // or []string if your IDs are strings
	ReportTextContent                   bool  `json:"reportTextContent"`
	SyncStatus                          bool  `json:"syncStatus"`
	Impression                          bool  `json:"impression"`
	Recommendation                      bool  `json:"recommendation"`
	ImpressionAddtional                 bool  `json:"impressionaddtional"`
	RecommendationAddtional             bool  `json:"recommendationaddtional"`
	CommonImpressionRecommendation      bool  `json:"commonImpressionRecommendation"`
	ImpressionRight                     bool  `json:"impressionRight"`
	RecommendationRight                 bool  `json:"recommendationRight"`
	ImpressionAddtionalRight            bool  `json:"impressionaddtionalRight"`
	RecommendationAddtionalRight        bool  `json:"recommendationaddtionalRight"`
	CommonImpressionRecommendationRight bool  `json:"commonImpressionRecommendationRight"`
	ArtificatsLeft                      bool  `json:"artificatsLeft"`
	ArtificatsRight                     bool  `json:"artificatsRight"`
	PatientHistory                      bool  `json:"patienthistory"`
	BreastImplantImageText              bool  `json:"breastimplantImageText"`
	SymmetryImageText                   bool  `json:"symmetryImageText"`
	BreastDensityImageText              bool  `json:"breastdensityImageText"`
	NippleAreolaImageText               bool  `json:"nippleareolaImageText"`
	GlandularImageText                  bool  `json:"glandularImageText"`
	LymphNodesImageText                 bool  `json:"lymphnodesImageText"`
	BreastDensityImageTextLeft          bool  `json:"breastdensityImageTextLeft"`
	NippleAreolaImageTextLeft           bool  `json:"nippleareolaImageTextLeft"`
	GlandularImageTextLeft              bool  `json:"glandularImageTextLeft"`
	LymphNodesImageTextLeft             bool  `json:"lymphnodesImageTextLeft"`
}

type AutoSubmitReportReq struct {
	ChangedOneState                     ChangedOneState  `json:"changedOne" binding:"required" mapstructure:"changedOne"`
	ReportIntakeForm                    []AnswerReqModel `json:"reportIntakeForm" binding:"required" mapstructure:"reportIntakeForm"`
	ReportTextContent                   string           `json:"reportTextContent" binding:"required" mapstructure:"reportTextContent"`
	AppointmentId                       int              `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId                           int              `json:"patientId" binding:"required" mapstructure:"patientId"`
	SyncStatus                          bool             `json:"syncStatus" mapstructure:"syncStatus"`
	Impression                          string           `json:"impression" mapstructure:"impression"`
	Recommendation                      string           `json:"recommendation" mapstructure:"recommendation"`
	ImpressionAddtional                 string           `json:"impressionaddtional" mapstructure:"impressionaddtional"`
	RecommendationAddtional             string           `json:"recommendationaddtional" mapstructure:"recommendationaddtional"`
	CommonImpressionRecommendation      string           `json:"commonImpressionRecommendation" mapstructure:"commonImpressionRecommendation"`
	ImpressionRight                     string           `json:"impressionRight" mapstructure:"impressionRight"`
	RecommendationRight                 string           `json:"recommendationRight" mapstructure:"recommendationRight"`
	ImpressionAddtionalRight            string           `json:"impressionaddtionalRight" mapstructure:"impressionaddtionalRight"`
	RecommendationAddtionalRight        string           `json:"recommendationaddtionalRight" mapstructure:"recommendationaddtionalRight"`
	CommonImpressionRecommendationRight string           `json:"commonImpressionRecommendationRight" mapstructure:"commonImpressionRecommendationRight"`
	ArtificatsLeft                      bool             `json:"artificatsLeft" mapstructure:"artificatsLeft"`
	ArtificatsRight                     bool             `json:"artificatsRight" mapstructure:"artificatsRight"`
	PatientHistory                      string           `json:"patienthistory" mapstructure:"patienthistory"`
	BreastImplantsImagetext             string           `json:"breastimplantImageText" mapstructure:"breastimplantImageText"`
	SymmetryImageText                   string           `json:"symmetryImageText" mapstructure:"symmetryImageText"`
	BreastdensityImageText              string           `json:"breastdensityImageText" mapstructure:"breastdensityImageText"`
	NippleAreolaImageText               string           `json:"nippleareolaImageText" mapstructure:"nippleareolaImageText"`
	GlandularImageText                  string           `json:"glandularImageText" mapstructure:"glandularImageText"`
	LymphnodesImageText                 string           `json:"lymphnodesImageText" mapstructure:"lymphnodesImageText"`
	BreastdensityImageTextLeft          string           `json:"breastdensityImageTextLeft" mapstructure:"breastdensityImageTextLeft"`
	NippleAreolaImageTextLeft           string           `json:"nippleareolaImageTextLeft" mapstructure:"nippleareolaImageTextLeft"`
	GlandularImageTextLeft              string           `json:"glandularImageTextLeft" mapstructure:"glandularImageTextLeft"`
	LymphnodesImageTextLeft             string           `json:"lymphnodesImageTextLeft" mapstructure:"lymphnodesImageTextLeft"`
}

type UpdateRemarkReq struct {
	AppointmentId int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId     int    `json:"patientId" binding:"required" mapstructure:"patientId"`
	Remark        string `json:"remark" binding:"required" mapstructure:"remark"`
}

type ListRemarkReq struct {
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type ListRemarkModel struct {
	RId            int    `json:"refRId" gorm:"column:refRId"`
	AppointmentId  int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	CustId         string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RemarksMessage string `json:"refRemarksMessage" gorm:"column:refRemarksMessage"`
	RCreatedAt     string `json:"refRCreatedAt" gorm:"column:refRCreatedAt"`
}

type UploadReportFormateReq struct {
	Name            string `json:"name" binding:"required" mapstructure:"name"`
	FormateTemplate string `json:"formateTemplate" binding:"required" mapstructure:"formateTemplate"`
}

type GetUserDetails struct {
	CustId         string `json:"refUserCustId" gorm:"column:refUserCustId"`
	FirstName      string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	Email          string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	Specialization string `json:"specialization" gorm:"column:specialization"`
	Department     string `json:"department" gorm:"column:department"`
}

type GetReportFormateReq struct {
	Id int `json:"id" binding:"required" mapstructure:"id"`
}

type ReportFormateModel struct {
	RFId        int       `json:"refRFId" gorm:"column:refRFId"`
	RFName      string    `json:"refRFName" gorm:"column:refRFName"`
	RFCreatedAt time.Time `json:"refRFCreatedAt" gorm:"column:refRFCreatedAt"`
	RFCreatedBy int       `json:"refRFCreatedBy" gorm:"column:refRFCreatedBy"`
}

type ReportTextFormateModel struct {
	RFId        int       `json:"refRFId" gorm:"column:refRFId"`
	RFName      string    `json:"refRFName" gorm:"column:refRFName"`
	RFText      string    `json:"refRFText" gorm:"column:refRFText"`
	RFCreatedAt time.Time `json:"refRFCreatedAt" gorm:"column:refRFCreatedAt"`
	RFCreatedBy int       `json:"refRFCreatedBy" gorm:"column:refRFCreatedBy"`
}

type GetOneUserAppointmentModel struct {
	AppointmentId                                  int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	UserId                                         int    `json:"refUserId" gorm:"column:refUserId"`
	SCId                                           int    `json:"refSCId" gorm:"column:refSCId"`
	SCCustId                                       string `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCName                                         string `json:"refSCName" gorm:"column:refSCName"`
	CategoryId                                     int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	AppointmentDate                                string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	AppointmentComplete                            string `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
	AppointmentPriority                            string `json:"refAppointmentPriority" gorm:"column:refAppointmentPriority"`
	AppointmentAccessId                            int    `json:"refAppointmentAccessId" gorm:"column:refAppointmentAccessId"`
	AppointmentAccessStatus                        bool   `json:"refAppointmentAccessStatus" gorm:"column:refAppointmentAccessStatus"`
	AppointmentAssignedUserId                      int    `json:"refAppointmentAssignedUserId" gorm:"column:refAppointmentAssignedUserId"`
	AppointmentRemarks                             string `json:"refAppointmentRemarks" gorm:"column:refAppointmentRemarks"`
	AppointmentImpression                          string `json:"refAppointmentImpression" gorm:"column:refAppointmentImpression"`
	AppointmentRecommendation                      string `json:"refAppointmentRecommendation" gorm:"column:refAppointmentRecommendation"`
	AppointmentImpressionAdditional                string `json:"refAppointmentImpressionAdditional" gorm:"column:refAppointmentImpressionAdditional"`
	AppointmentRecommendationAdditional            string `json:"refAppointmentRecommendationAdditional" gorm:"column:refAppointmentRecommendationAdditional"`
	AppointmentCommonImpressionRecommendation      string `json:"refAppointmentCommonImpressionRecommendation" gorm:"column:refAppointmentCommonImpressionRecommendation"`
	AppointmentImpressionRight                     string `json:"refAppointmentImpressionRight" gorm:"column:refAppointmentImpressionRight"`
	AppointmentRecommendationRight                 string `json:"refAppointmentRecommendationRight" gorm:"column:refAppointmentRecommendationRight"`
	AppointmentImpressionAdditionalRight           string `json:"refAppointmentImpressionAdditionalRight" gorm:"column:refAppointmentImpressionAdditionalRight"`
	AppointmentRecommendationAdditionalRight       string `json:"refAppointmentRecommendationAdditionalRight" gorm:"column:refAppointmentRecommendationAdditionalRight"`
	AppointmentCommonImpressionRecommendationRight string `json:"refAppointmentCommonImpressionRecommendationRight" gorm:"column:refAppointmentCommonImpressionRecommendationRight"`
	AppointmentPatientHistory                      string `json:"refAppointmentPatietHistory" gorm:"column:refAppointmentPatietHistory"`
	AppointmentBreastImplantImageText              string `json:"refAppointmentBreastImplantImageText" gorm:"column:refAppointmentBreastImplantImageText"`
	AppointmentSymmetryImageText                   string `json:"refAppointmentSymmetryImageText" gorm:"column:refAppointmentSymmetryImageText"`
	AppointmentBreastdensityImageText              string `json:"refAppointmentBreastdensityImageText" gorm:"column:refAppointmentBreastdensityImageText"`
	AppointmentNippleAreolaImageText               string `json:"refAppointmentNippleAreolaImageText" gorm:"column:refAppointmentNippleAreolaImageText"`
	AppointmentGlandularImageText                  string `json:"refAppointmentGlandularImageText" gorm:"column:refAppointmentGlandularImageText"`
	AppointmentLymphnodeImageText                  string `json:"refAppointmentLymphnodeImageText" gorm:"column:refAppointmentLymphnodeImageText"`
	AppointmentBreastdensityImageTextLeft          string `json:"refAppointmentBreastdensityImageTextLeft" gorm:"column:refAppointmentBreastdensityImageTextLeft"`
	AppointmentNippleAreolaImageTextLeft           string `json:"refAppointmentNippleAreolaImageTextLeft" gorm:"column:refAppointmentNippleAreolaImageTextLeft"`
	AppointmentGlandularImageTextLeft              string `json:"refAppointmentGlandularImageTextLeft" gorm:"column:refAppointmentGlandularImageTextLeft"`
	AppointmentLymphnodeImageTextLeft              string `json:"refAppointmentLymphnodeImageTextLeft" gorm:"column:refAppointmentLymphnodeImageTextLeft"`
}

type PatientCustId struct {
	CustId string `json:"refUserCustId" gorm:"column:refUserCustId"`
}

type ListUserModel struct {
	RHId int `json:"refRHId" gorm:"column:refRHId"`
}

type Patientdata struct {
	UserFirstName   string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	CustId          string `json:"refUserCustId" gorm:"column:refUserCustId"`
	AppointmentDate string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	SCCustId        string `json:"refSCCustId" gorm:"column:refSCCustId"`
	Email           string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
}

type ManagerData struct {
	SCCustId string `json:"refSCCustId" gorm:"column:refSCCustId"`
	Email    string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
}

type SendMailReportReq struct {
	AppointmentId     int  `json:"appintmentId" binding:"appintmentId"  mapstructure:"appintmentId"`
	PatientId         int  `json:"patientId" binding:"required" mapstructure:"patientId"`
	PatientMailStatus bool `json:"patientMailStatus" mapstructure:"patientMailStatus"`
	ManagerMailStatus bool `json:"managerMailStatus" mapstructure:"managerMailStatus"`
}

type DoctorReportAccessStatus struct {
	DDEaseQTReportAccess *bool `json:"refDDEaseQTReportAccess" gorm:"column:refDDEaseQTReportAccess"`
}

type CoDoctorReportAccessStatus struct {
	CDEaseQTReportAccess *bool `json:"refCDEaseQTReportAccess" gorm:"column:refCDEaseQTReportAccess"`
}

type AddAddendumReq struct {
	AddAddendumText string `json:"addAddendumText" binding:"required" mapstructure:"addAddendumText"`
	AppointmentId   int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type AddAddendumModel struct {
	ADID          int    `json:"refADID" gorm:"column:refADID"`
	AppointmentId string `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	CustId        string `json:"refUserCustId" gorm:"column:refUserCustId"`
	ADText        string `json:"refADText" gorm:"column:refADText"`
	ADCreatedAt   string `json:"refADCreatedAt" gorm:"column:refADCreatedAt"`
}

type DownloadReportReq struct {
	Id int `json:"id" binding:"required" mapstructure:"id"`
}
