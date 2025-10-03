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
	Status                 bool   `gorm:"column:status"`
	RefAppointmentAccessId int    `gorm:"column:refAppointmentAccessId"`
	CustID                 string `gorm:"column:userCustId"`
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

type GetOldReport struct {
	RefORCategoryId int    `json:"refORCategoryId" gorm:"column:refORCategoryId"`
	Files           string `json:"files" gorm:"column:files"`
}

type GetReportTextContent struct {
	IntakeId                                     int    `json:"refRTCId" gorm:"column:refRTCId"`
	TextContent                                  string `json:"refRTCText" gorm:"column:refRTCText"`
	RTSyncStatus                                 *bool  `json:"refRTSyncStatus" gorm:"column:refRTSyncStatus"`
	RefRTPatientHistorySyncStatus                *bool  `json:"refRTPatientHistorySyncStatus" gorm:"column:refRTPatientHistorySyncStatus"`
	RefRTBreastImplantSyncStatus                 *bool  `json:"refRTBreastImplantSyncStatus" gorm:"column:refRTBreastImplantSyncStatus"`
	RefRTSymmetrySyncStatus                      *bool  `json:"refRTSymmetrySyncStatus" gorm:"column:refRTSymmetrySyncStatus"`
	RefRTBreastDensityandImageRightSyncStatus    *bool  `json:"refRTBreastDensityandImageRightSyncStatus" gorm:"column:refRTBreastDensityandImageRightSyncStatus"`
	RefRTNippleAreolaSkinRightSyncStatus         *bool  `json:"refRTNippleAreolaSkinRightSyncStatus" gorm:"column:refRTNippleAreolaSkinRightSyncStatus"`
	RefRTLesionsRightSyncStatus                  *bool  `json:"refRTLesionsRightSyncStatus" gorm:"column:refRTLesionsRightSyncStatus"`
	RefRTComparisonPriorSyncStatus               *bool  `json:"refRTComparisonPriorSyncStatus" gorm:"column:refRTComparisonPriorSyncStatus"`
	RefRTGrandularAndDuctalTissueRightSyncStatus *bool  `json:"refRTGrandularAndDuctalTissueRightSyncStatus" gorm:"column:refRTGrandularAndDuctalTissueRightSyncStatus"`
	RefRTLymphNodesRightSyncStatus               *bool  `json:"refRTLymphNodesRightSyncStatus" gorm:"column:refRTLymphNodesRightSyncStatus"`
	RefRTBreastDensityandImageLeftSyncStatus     *bool  `json:"refRTBreastDensityandImageLeftSyncStatus" gorm:"column:refRTBreastDensityandImageLeftSyncStatus"`
	RefRTNippleAreolaSkinLeftSyncStatus          *bool  `json:"refRTNippleAreolaSkinLeftSyncStatus" gorm:"column:refRTNippleAreolaSkinLeftSyncStatus"`
	RefRTLesionsLeftSyncStatus                   *bool  `json:"refRTLesionsLeftSyncStatus" gorm:"column:refRTLesionsLeftSyncStatus"`
	RefRTComparisonPriorLeftSyncStatus           *bool  `json:"refRTComparisonPriorLeftSyncStatus" gorm:"column:refRTComparisonPriorLeftSyncStatus"`
	RefRTGrandularAndDuctalTissueLeftSyncStatus  *bool  `json:"refRTGrandularAndDuctalTissueLeftSyncStatus" gorm:"column:refRTGrandularAndDuctalTissueLeftSyncStatus"`
	RefRTLymphNodesLeftSyncStatus                *bool  `json:"refRTLymphNodesLeftSyncStatus" gorm:"column:refRTLymphNodesLeftSyncStatus"`
	RefRTBreastImplantReportText                 string `json:"refRTBreastImplantReportText" gorm:"column:refRTBreastImplantReportText"`
	RefRTSymmetryReportText                      string `json:"refRTSymmetryReportText" gorm:"column:refRTSymmetryReportText"`
	RefRTBreastDensityandImageRightReportText    string `json:"refRTBreastDensityandImageRightReportText" gorm:"column:refRTBreastDensityandImageRightReportText"`
	RefRTNippleAreolaSkinRightReportText         string `json:"refRTNippleAreolaSkinRightReportText" gorm:"column:refRTNippleAreolaSkinRightReportText"`
	RefRTLesionsRightReportText                  string `json:"refRTLesionsRightReportText" gorm:"column:refRTLesionsRightReportText"`
	RefRTComparisonPriorReportText               string `json:"refRTComparisonPriorReportText" gorm:"column:refRTComparisonPriorReportText"`
	RefRTGrandularAndDuctalTissueRightReportText string `json:"refRTGrandularAndDuctalTissueRightReportText" gorm:"column:refRTGrandularAndDuctalTissueRightReportText"`
	RefRTLymphNodesRightReportText               string `json:"refRTLymphNodesRightReportText" gorm:"column:refRTLymphNodesRightReportText"`
	RefRTBreastDensityandImageLeftReportText     string `json:"refRTBreastDensityandImageLeftReportText" gorm:"column:refRTBreastDensityandImageLeftReportText"`
	RefRTNippleAreolaSkinLeftReportText          string `json:"refRTNippleAreolaSkinLeftReportText" gorm:"column:refRTNippleAreolaSkinLeftReportText"`
	RefRTLesionsLeftReportText                   string `json:"refRTLesionsLeftReportText" gorm:"column:refRTLesionsLeftReportText"`
	RefRTComparisonPriorLeftReportText           string `json:"refRTComparisonPriorLeftReportText" gorm:"column:refRTComparisonPriorLeftReportText"`
	RefRTGrandularAndDuctalTissueLeftReportText  string `json:"refRTGrandularAndDuctalTissueLeftReportText" gorm:"column:refRTGrandularAndDuctalTissueLeftReportText"`
	RefRTLymphNodesLeftReportText                string `json:"refRTLymphNodesLeftReportText" gorm:"column:refRTLymphNodesLeftReportText"`
}

type GetReportHistory struct {
	RHId              int    `json:"refRHId" gorm:"column:refRHId"`
	HandledUserId     int    `json:"refRHHandledUserId" gorm:"column:refRHHandledUserId"`
	HandleStartTime   string `json:"refRHHandleStartTime" gorm:"column:refRHHandleStartTime"`
	HandleEndTime     string `json:"refRHHandleEndTime" gorm:"column:refRHHandleEndTime"`
	HandleUserName    string `json:"HandleUserName" gorm:"column:HandleUserName"`
	HandleStatus      string `json:"refRHHandleStatus" gorm:"column:refRHHandleStatus"`
	HandleContentText string `json:"refRHHandleContentText" gorm:"column:refRHHandleContentText"`
	HandlerRTId       int    `json:"HandlerRTId" gorm:"column:HandlerRTId"`
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
	ReportQuestion                         []int `json:"reportQuestion" mapstructure:"reportQuestion"`
	ReportTextContent                      bool  `json:"reportTextContent" mapstructure:"reportTextContent"`
	SyncStatus                             bool  `json:"syncStatus" mapstructure:"syncStatus"`
	Impression                             bool  `json:"impression" mapstructure:"impression"`
	Recommendation                         bool  `json:"recommendation" mapstructure:"recommendation"`
	ImpressionAddtional                    bool  `json:"impressionaddtional" mapstructure:"impressionaddtional"`
	RecommendationAddtional                bool  `json:"recommendationaddtional" mapstructure:"recommendationaddtional"`
	CommonImpressionRecommendation         bool  `json:"commonImpressionRecommendation" mapstructure:"commonImpressionRecommendation"`
	ImpressionRight                        bool  `json:"impressionRight" mapstructure:"impressionRight"`
	RecommendationRight                    bool  `json:"recommendationRight" mapstructure:"recommendationRight"`
	ImpressionAddtionalRight               bool  `json:"impressionaddtionalRight" mapstructure:"impressionaddtionalRight"`
	RecommendationAddtionalRight           bool  `json:"recommendationaddtionalRight" mapstructure:"recommendationaddtionalRight"`
	CommonImpressionRecommendationRight    bool  `json:"commonImpressionRecommendationRight" mapstructure:"commonImpressionRecommendationRight"`
	ArtificatsLeft                         bool  `json:"artificatsLeft" mapstructure:"artificatsLeft"`
	ArtificatsRight                        bool  `json:"artificatsRight" mapstructure:"artificatsRight"`
	PatientHistory                         bool  `json:"patienthistory" mapstructure:"patienthistory"`
	BreastImplantImageText                 bool  `json:"breastimplantImageText" mapstructure:"breastimplantImageText"`
	SymmetryImageText                      bool  `json:"symmetryImageText" mapstructure:"symmetryImageText"`
	BreastDensityImageText                 bool  `json:"breastdensityImageText" mapstructure:"breastdensityImageText"`
	NippleAreolaImageText                  bool  `json:"nippleareolaImageText" mapstructure:"nippleareolaImageText"`
	GlandularImageText                     bool  `json:"glandularImageText" mapstructure:"glandularImageText"`
	LymphNodesImageText                    bool  `json:"lymphnodesImageText" mapstructure:"lymphnodesImageText"`
	BreastDensityImageTextLeft             bool  `json:"breastdensityImageTextLeft" mapstructure:"breastdensityImageTextLeft"`
	NippleAreolaImageTextLeft              bool  `json:"nippleareolaImageTextLeft" mapstructure:"nippleareolaImageTextLeft"`
	GlandularImageTextLeft                 bool  `json:"glandularImageTextLeft" mapstructure:"glandularImageTextLeft"`
	LymphNodesImageTextLeft                bool  `json:"lymphnodesImageTextLeft" mapstructure:"lymphnodesImageTextLeft"`
	BreastImplantSyncStatus                bool  `json:"breastImplantSyncStatus" mapstructure:"breastImplantSyncStatus"`
	SymmetrySyncStatus                     bool  `json:"symmetrySyncStatus" mapstructure:"symmetrySyncStatus"`
	BreastDensitySyncStatus                bool  `json:"breastDensityandImageRightSyncStatus" mapstructure:"breastDensityandImageRightSyncStatus"`
	NippleAreolaSyncStatus                 bool  `json:"nippleAreolaSkinRightSyncStatus" mapstructure:"nippleAreolaSkinRightSyncStatus"`
	GlandularSyncStatus                    bool  `json:"grandularAndDuctalTissueRightSyncStatus" mapstructure:"grandularAndDuctalTissueRightSyncStatus"`
	LymphNodesSyncStatus                   bool  `json:"LymphNodesRightSyncStatus" mapstructure:"LymphNodesRightSyncStatus"`
	LesionsSyncStatus                      bool  `json:"LesionsRightSyncStatus" mapstructure:"LesionsRightSyncStatus"`
	ComparisonPriorSyncStatus              bool  `json:"ComparisonPriorSyncStatus" mapstructure:"ComparisonPriorSyncStatus"`
	BreastDensityLeftSyncStatus            bool  `json:"breastDensityandImageLeftSyncStatus" mapstructure:"breastDensityandImageLeftSyncStatus"`
	NippleAreolaLeftSyncStatus             bool  `json:"nippleAreolaSkinLeftSyncStatus" mapstructure:"nippleAreolaSkinLeftSyncStatus"`
	GlandularLeftSyncStatus                bool  `json:"grandularAndDuctalTissueLeftSyncStatus" mapstructure:"grandularAndDuctalTissueLeftSyncStatus"`
	LymphNodesLeftSyncStatus               bool  `json:"LymphNodesLeftSyncStatus" mapstructure:"LymphNodesLeftSyncStatus"`
	LesionsLeftSyncStatus                  bool  `json:"LesionsLeftSyncStatus" mapstructure:"LesionsLeftSyncStatus"`
	ComparisonPriorLeftSyncStatus          bool  `json:"ComparisonPriorLeftSyncStatus" mapstructure:"ComparisonPriorLeftSyncStatus"`
	BreastImplantReportText                bool  `json:"breastImplantReportText" mapstructure:"breastImplantReportText"`
	SymmetryReportText                     bool  `json:"symmetryReportText" mapstructure:"symmetryReportText"`
	BreastDensityReportText                bool  `json:"breastDensityandImageRightReportText" mapstructure:"breastDensityandImageRightReportText"`
	NippleAreolaReportText                 bool  `json:"nippleAreolaSkinRightReportText" mapstructure:"nippleAreolaSkinRightReportText"`
	LesionsReportText                      bool  `json:"LesionsRightReportText" mapstructure:"LesionsRightReportText"`
	ComparisonPriorReportText              bool  `json:"ComparisonPriorReportText" mapstructure:"ComparisonPriorReportText"`
	GrandularAndDuctalTissueReportText     bool  `json:"grandularAndDuctalTissueRightReportText" mapstructure:"grandularAndDuctalTissueRightReportText"`
	LymphNodesReportText                   bool  `json:"LymphNodesRightReportText" mapstructure:"LymphNodesRightReportText"`
	BreastDensityLeftReportText            bool  `json:"breastDensityandImageLeftReportText" mapstructure:"breastDensityandImageLeftReportText"`
	NippleAreolaLeftReportText             bool  `json:"nippleAreolaSkinLeftReportText" mapstructure:"nippleAreolaSkinLeftReportText"`
	LesionsLeftReportText                  bool  `json:"LesionsLeftReportText" mapstructure:"LesionsLeftReportText"`
	ComparisonPriorLeftReportText          bool  `json:"ComparisonPriorLeftReportText" mapstructure:"ComparisonPriorLeftReportText"`
	GrandularAndDuctalTissueLeftReportText bool  `json:"grandularAndDuctalTissueLeftReportText" mapstructure:"grandularAndDuctalTissueLeftReportText"`
	LymphNodesLeftReportText               bool  `json:"LymphNodesLeftReportText" mapstructure:"LymphNodesLeftReportText"`
}

type AutosyncSyncStatus struct {
	BreastImplantSyncStatus       bool `json:"breastImplantSyncStatus" mapstructure:"breastImplantSyncStatus"`
	SymmetrySyncStatus            bool `json:"symmetrySyncStatus" mapstructure:"symmetrySyncStatus"`
	BreastDensitySyncStatus       bool `json:"breastDensitySyncStatus" mapstructure:"breastDensitySyncStatus"`
	NippleAreolaSyncStatus        bool `json:"nippleAreolaSyncStatus" mapstructure:"nippleAreolaSyncStatus"`
	GlandularSyncStatus           bool `json:"glandularSyncStatus" mapstructure:"glandularSyncStatus"`
	LymphNodesSyncStatus          bool `json:"lymphnodesSyncStatus" mapstructure:"lymphnodesSyncStatus"`
	LesionsSyncStatus             bool `json:"lesionsSyncStatus" mapstructure:"lesionsSyncStatus"`
	ComparisonPriorSyncStatus     bool `json:"comparisonPriorSyncStatus" mapstructure:"comparisonPriorSyncStatus"`
	BreastDensityLeftSyncStatus   bool `json:"breastDensityLeftSyncStatus" mapstructure:"breastDensityLeftSyncStatus"`
	NippleAreolaLeftSyncStatus    bool `json:"nippleAreolaLeftSyncStatus" mapstructure:"nippleAreolaLeftSyncStatus"`
	GlandularLeftSyncStatus       bool `json:"glandularLeftSyncStatus" mapstructure:"glandularLeftSyncStatus"`
	LymphNodesLeftSyncStatus      bool `json:"lymphnodesLeftSyncStatus" mapstructure:"lymphnodesLeftSyncStatus"`
	LesionsLeftSyncStatus         bool `json:"lesionsLeftSyncStatus" mapstructure:"lesionsLeftSyncStatus"`
	ComparisonPriorLeftSyncStatus bool `json:"comparisonPriorLeftSyncStatus" mapstructure:"comparisonPriorLeftSyncStatus"`
}

type AutoReportText struct {
	BreastImplantReportText                string `json:"breastImplantReportText" mapstructure:"breastImplantReportText"`
	SymmetryReportText                     string `json:"symmetryReportText" mapstructure:"symmetryReportText"`
	BreastDensityReportText                string `json:"breastDensityReportText" mapstructure:"breastDensityReportText"`
	NippleAreolaReportText                 string `json:"nippleAreolaReportText" mapstructure:"nippleAreolaReportText"`
	LesionsReportText                      string `json:"lesionsReportText" mapstructure:"lesionsReportText"`
	ComparisonPriorReportText              string `json:"comparisonPriorReportText" mapstructure:"comparisonPriorReportText"`
	GrandularAndDuctalTissueReportText     string `json:"grandularAndDuctalTissueReportText" mapstructure:"grandularAndDuctalTissueReportText"`
	LymphNodesReportText                   string `json:"lymphnodesReportText" mapstructure:"lymphnodesReportText"`
	BreastDensityReportTextLeft            string `json:"breastDensityReportTextLeft" mapstructure:"breastDensityReportTextLeft"`
	NippleAreolaReportTextLeft             string `json:"nippleAreolaReportTextLeft" mapstructure:"nippleAreolaReportTextLeft"`
	LesionsReportTextLeft                  string `json:"lesionsReportTextLeft" mapstructure:"lesionsReportTextLeft"`
	ComparisonPriorReportTextLeft          string `json:"comparisonPriorReportTextLeft" mapstructure:"comparisonPriorReportTextLeft"`
	GrandularAndDuctalTissueReportTextLeft string `json:"grandularAndDuctalTissueReportTextLeft" mapstructure:"grandularAndDuctalTissueReportTextLeft"`
	LymphNodesReportTextLeft               string `json:"lymphnodesReportTextLeft" mapstructure:"lymphnodesReportTextLeft"`
}

type AutoSubmitReportReq struct {
	ChangedOneState                     ChangedOneState    `json:"changedOne" binding:"required" mapstructure:"changedOne"`
	ReportIntakeForm                    []AnswerReqModel   `json:"reportIntakeForm" binding:"required" mapstructure:"reportIntakeForm"`
	ReportTextContent                   string             `json:"reportTextContent" binding:"required" mapstructure:"reportTextContent"`
	AppointmentId                       int                `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId                           int                `json:"patientId" binding:"required" mapstructure:"patientId"`
	SyncStatus                          bool               `json:"syncStatus" mapstructure:"syncStatus"`
	Impression                          string             `json:"impression" mapstructure:"impression"`
	Recommendation                      string             `json:"recommendation" mapstructure:"recommendation"`
	ImpressionAddtional                 string             `json:"impressionaddtional" mapstructure:"impressionaddtional"`
	RecommendationAddtional             string             `json:"recommendationaddtional" mapstructure:"recommendationaddtional"`
	CommonImpressionRecommendation      string             `json:"commonImpressionRecommendation" mapstructure:"commonImpressionRecommendation"`
	ImpressionRight                     string             `json:"impressionRight" mapstructure:"impressionRight"`
	RecommendationRight                 string             `json:"recommendationRight" mapstructure:"recommendationRight"`
	ImpressionAddtionalRight            string             `json:"impressionaddtionalRight" mapstructure:"impressionaddtionalRight"`
	RecommendationAddtionalRight        string             `json:"recommendationaddtionalRight" mapstructure:"recommendationaddtionalRight"`
	CommonImpressionRecommendationRight string             `json:"commonImpressionRecommendationRight" mapstructure:"commonImpressionRecommendationRight"`
	ArtificatsLeft                      bool               `json:"artificatsLeft" mapstructure:"artificatsLeft"`
	ArtificatsRight                     bool               `json:"artificatsRight" mapstructure:"artificatsRight"`
	PatientHistory                      string             `json:"patienthistory" mapstructure:"patienthistory"`
	BreastImplantsImagetext             string             `json:"breastimplantImageText" mapstructure:"breastimplantImageText"`
	SymmetryImageText                   string             `json:"symmetryImageText" mapstructure:"symmetryImageText"`
	BreastdensityImageText              string             `json:"breastdensityImageText" mapstructure:"breastdensityImageText"`
	NippleAreolaImageText               string             `json:"nippleareolaImageText" mapstructure:"nippleareolaImageText"`
	GlandularImageText                  string             `json:"glandularImageText" mapstructure:"glandularImageText"`
	LymphnodesImageText                 string             `json:"lymphnodesImageText" mapstructure:"lymphnodesImageText"`
	BreastdensityImageTextLeft          string             `json:"breastdensityImageTextLeft" mapstructure:"breastdensityImageTextLeft"`
	NippleAreolaImageTextLeft           string             `json:"nippleareolaImageTextLeft" mapstructure:"nippleareolaImageTextLeft"`
	GlandularImageTextLeft              string             `json:"glandularImageTextLeft" mapstructure:"glandularImageTextLeft"`
	LymphnodesImageTextLeft             string             `json:"lymphnodesImageTextLeft" mapstructure:"lymphnodesImageTextLeft"`
	ReportSyncStatus                    AutosyncSyncStatus `json:"reportSyncStatus" mapstructure:"reportSyncStatus"`
	AutoReportText                      AutoReportText     `json:"autoReportText" mapstructure:"autoReportText"`
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

type AddedumCountModel struct {
	Count int `json:"count" gorm:"column:count"`
}

type UploadReportFormateReq struct {
	Name            string `json:"name" binding:"required" mapstructure:"name"`
	FormateTemplate string `json:"formateTemplate" binding:"required" mapstructure:"formateTemplate"`
	AccessStatus    string `json:"accessStatus" binding:"required" mapstructure:"accessStatus"`
}

type DeleteReportFormateReq struct {
	Id int `json:"id" binding:"required" mapstructure:"id"`
}

type UpdateReportFormateReq struct {
	Id           int    `json:"id" binding:"required" mapstructure:"id"`
	AccessStatus string `json:"accessStatus" binding:"required" mapstructure:"accessStatus"`
}

type ListOldReportReq struct {
	CategoryId    int `json:"categoryId" binding:"required" mapstructure:"categoryId"`
	PatientId     int `json:"patientId" binding:"required" mapstructure:"patientId"`
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
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
	RFId              int       `json:"refRFId" gorm:"column:refRFId"`
	RFName            string    `json:"refRFName" gorm:"column:refRFName"`
	RFCreatedAt       time.Time `json:"refRFCreatedAt" gorm:"column:refRFCreatedAt"`
	RFCreatedBy       int       `json:"refRFCreatedBy" gorm:"column:refRFCreatedBy"`
	RefUserCustId     string    `json:"refUserCustId" gorm:"column:refUserCustId"`
	RefRFAccessStatus string    `json:"refRFAccessStatus" gorm:"column:refRFAccessStatus"`
}

type ReportTextFormateModel struct {
	RFId        int       `json:"refRFId" gorm:"column:refRFId"`
	RFName      string    `json:"refRFName" gorm:"column:refRFName"`
	RFText      string    `json:"refRFText" gorm:"column:refRFText"`
	RFCreatedAt time.Time `json:"refRFCreatedAt" gorm:"column:refRFCreatedAt"`
	RFCreatedBy int       `json:"refRFCreatedBy" gorm:"column:refRFCreatedBy"`
}

type ReportFormateCreateModel struct {
	RFId          int    `json:"refRFId" gorm:"column:refRFId"`
	RefUserCustId string `json:"refUserCustId" gorm:"column:refUserCustId"`
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
	CustId        string `json:"refUserCustId" gorm:"column:refUserCustId"`
	UserFirstName string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	UserDOB       string `json:"refUserDOB" gorm:"column:refUserDOB"`
	UserGender    string `json:"refUserGender" gorm:"column:refUserGender"`
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
	DDEaseQTReportAccess   *bool `json:"refDDEaseQTReportAccess" gorm:"column:refDDEaseQTReportAccess"`
	DDNAsystemReportAccess *bool `json:"refDDNAsystemReportAccess" gorm:"column:refDDNAsystemReportAccess"`
}

type CoDoctorReportAccessStatus struct {
	CDEaseQTReportAccess   *bool `json:"refCDEaseQTReportAccess" gorm:"column:refCDEaseQTReportAccess"`
	CDNAsystemReportAccess *bool `json:"refCDNAsystemReportAccess" gorm:"column:refCDNAsystemReportAccess"`
}

type AddAddendumReq struct {
	AddAddendumText   string `json:"addAddendumText" binding:"required" mapstructure:"addAddendumText"`
	AppointmentId     int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientMailStatus bool   `json:"patientMailStatus" mapstructure:"patientMailStatus"`
	ManagerMailStatus bool   `json:"managerMailStatus" mapstructure:"managerMailStatus"`
	PatientId         int    `json:"patientId" binding:"required" mapstructure:"patientId"`
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

type ViewReportReq struct {
	FileName string `json:"filename" binding:"required" mapstructure:"filename"`
}

type ViewReportRes struct {
	File *FileData `json:"file" gorm:"-"`
}

type ListOldReportModel struct {
	ORId          int    `json:"refORId" gorm:"column:refORId"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	AppointmentId int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	ORCategoryId  int    `json:"refORCategoryId" gorm:"column:refORCategoryId"`
	ORFilename    string `json:"refORFilename" gorm:"column:refORFilename"`
	ORCreatedAt   string `json:"refORCreatedAt" gorm:"column:refORCreatedAt"`
	ORCreatedBy   string `json:"refORCreatedBy" gorm:"column:refORCreatedBy"`
}

type DeleteOldReportModel struct {
	ORId int `json:"refORId" binding:"required" mapstructure:"refORId"`
}

type AddSignatureReq struct {
	AddSignatureText string `json:"addSignatureText" binding:"required" mapstructure:"addSignatureText"`
	AppointmentId    int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	PatientId        int    `json:"patientId" binding:"required" mapstructure:"patientId"`
}

type ListAllSignatureReq struct {
	AppointmentId int `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
}

type ListAllSignatureModel struct {
	SId           int    `json:"refSId" gorm:"column:refSId"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	AppointmentId int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	SText         string `json:"refSText" gorm:"column:refSText"`
	SCreatedAt    string `json:"refSCreatedAt" gorm:"column:refSCreatedAt"`
}
