package model

type AdminScanCenterModel struct {
	Month             string `json:"month" gorm:"column:month"`
	MonthName         string `json:"month_name" gorm:"column:month_name"`
	TotalAppointments int    `json:"total_appointments" gorm:"column:total_appointments"`
}

type AdminOverallScanIndicatesAnalayticsModel struct {
	TotalAppointments int `json:"total_appointments" gorm:"column:total_appointments"`
	SForm             int `json:"SForm" gorm:"column:SForm"`
	DaForm            int `json:"DaForm" gorm:"column:DaForm"`
	DbForm            int `json:"DbForm" gorm:"column:DbForm"`
	DcForm            int `json:"DcForm" gorm:"column:DcForm"`
}

type UserAccessTimingModel struct {
	TotalMinutes string `json:"total_minutes" gorm:"column:total_minutes"`
	TotalHours   string `json:"total_hours" gorm:"column:total_hours"`
}

type GetAllScaCenter struct {
	SCId           uint   `json:"refSCId" gorm:"column:refSCId"`
	SCCustId       string `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCProfile      string `json:"refSCProfile" gorm:"column:refSCProfile"`
	SCName         string `json:"refSCName" gorm:"column:refSCName"`
	SCAddress      string `json:"refSCAddress" gorm:"column:refSCAddress"`
	SCPhoneNo1     string `json:"refSCPhoneNo1" gorm:"column:refSCPhoneNo1"`
	SCEmail        string `json:"refSCEmail" gorm:"column:refSCEmail"`
	SCWebsite      string `json:"refSCWebsite" gorm:"column:refSCWebsite"`
	SCAppointments bool   `json:"refSCAppointments" gorm:"column:refSCAppointments"`
}

type UserListIdsData struct {
	UserId     int    `json:"refUserId" gorm:"column:refUserId"`
	UserCustId string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId     int    `json:"refRTId" gorm:"column:refRTId"`
}

type Artificates struct {
	Leftartifacts  int `json:"leftartifacts" gorm:"column:leftartifacts"`
	Rightartifacts int `json:"rightartifacts" gorm:"column:rightartifacts"`
	Bothartifacts  int `json:"bothartifacts" gorm:"column:bothartifacts"`
}

type OverAllScancenterAnalaytics struct {
	SCId                     int    `json:"refSCId" gorm:"column:refSCId"`
	SCCustId                 string `json:"refSCCustId" gorm:"column:refSCCustId"`
	TotalCase                int    `json:"totalcase" gorm:"column:totalcase"`
	TotalSForm               int    `json:"totalsform" gorm:"column:totalsform"`
	TotalDaForm              int    `json:"totaldaform" gorm:"column:totaldaform"`
	TotalDbForm              int    `json:"totaldbform" gorm:"column:totaldbform"`
	TotalDcForm              int    `json:"totaldcform" gorm:"column:totaldcform"`
	TechArtificatsLeft       int    `json:"techartificatsleft" gorm:"column:techartificatsleft"`
	TechArtificateRight      int    `json:"techartificatsright" gorm:"column:techartificatsright"`
	ReportArtificatesLeft    int    `json:"reportartificatsleft" gorm:"column:reportartificatsleft"`
	ReportArtificatesRight   int    `json:"reportartificatsright" gorm:"column:reportartificatsright"`
	Leftannualscreening      int    `json:"leftannualscreening" gorm:"column:leftannualscreening"`
	Leftusgsfu               int    `json:"leftusgsfu" gorm:"column:leftusgsfu"`
	LeftBiopsy               int    `json:"leftBiopsy" gorm:"column:leftBiopsy"`
	LeftBreastradiologist    int    `json:"leftBreastradiologist" gorm:"column:leftBreastradiologist"`
	LeftClinicalCorrelation  int    `json:"leftClinicalCorrelation" gorm:"column:leftClinicalCorrelation"`
	LeftOncoConsult          int    `json:"leftOncoConsult" gorm:"column:leftOncoConsult"`
	LeftRedo                 int    `json:"leftRedo" gorm:"column:leftRedo"`
	Rightannualscreening     int    `json:"rightannualscreening" gorm:"column:rightannualscreening"`
	Rightusgsfu              int    `json:"rightusgsfu" gorm:"column:rightusgsfu"`
	RightBiopsy              int    `json:"rightBiopsy" gorm:"column:rightBiopsy"`
	RightBreastradiologist   int    `json:"rightBreastradiologist" gorm:"column:rightBreastradiologist"`
	RightClinicalCorrelation int    `json:"rightClinicalCorrelation" gorm:"column:rightClinicalCorrelation"`
	RightOncoConsult         int    `json:"rightOncoConsult" gorm:"column:rightOncoConsult"`
	RightRedo                int    `json:"rightRedo" gorm:"column:rightRedo"`
}

type AdminOverallAnalyticsResponse struct {
	AdminScanCenterModel                     []AdminScanCenterModel
	AdminOverallScanIndicatesAnalayticsModel []AdminOverallScanIndicatesAnalayticsModel
	GetAllScaCenter                          []GetAllScaCenter
	UserListIdsData                          []UserListIdsData
	RightRecommendation                      []RecommendationUserSQL
	LeftRecommendation                       []RecommendationUserSQL
	TechArtificats                           []Artificates
	ReportArtificats                         []Artificates
	OverAllScancenterAnalaytics              []OverAllScancenterAnalaytics
}

type ListScanAppointmentCountModel struct {
	SCId              int    `json:"refSCId" gorm:"column:refSCId"`
	SCName            string `json:"refSCName" gorm:"column:refSCName"`
	TotalAppointments int    `json:"total_appointments" gorm:"column:total_appointments"`
}

type AdminOverallOneAnalyticsReq struct {
	SCId int `json:"SCId" mapstructure:"SCId"`
	// Monthyear string `json:"monthnyear" mapstructure:"monthnyear"`
	StartDate string `json:"startDate" mapstructure:"startDate"`
	EndDate   string `json:"endDate" mapstructure:"endDate"`
}

type OneUserReq struct {
	UserId int `json:"userId" binding:"required" mapstructure:"userId"`
	RoleId int `json:"roleId" binding:"required" mapstructure:"roleId"`
	// Monthyear string `json:"monthnyear" binding:"required" mapstructure:"monthnyear"`
	StartDate string `json:"startDate" mapstructure:"startDate"`
	EndDate   string `json:"endDate" mapstructure:"endDate"`
}

type TotalCorrectEditModel struct {
	TotalCorrect int `json:"totalCorrect" gorm:"column:totalCorrect"`
	TotalEdit    int `json:"totalEdit" gorm:"column:totalEdit"`
}

type RecommendationUserSQL struct {
	GroupName  string `json:"group_name" gorm:"column:group_name"`
	TotalCount string `json:"total_count" gorm:"column:total_count"`
}

type ImpressionModel struct {
	Impression string `json:"impression" gorm:"column:impression"`
	Count      int    `json:"count" gorm:"column:count"`
}

type DurationBucketModel struct {
	Le1Day   int `json:"le_1_day" gorm:"column:le_1_day"`
	Le3Days  int `json:"le_3_days" gorm:"column:le_3_days"`
	Le7Days  int `json:"le_7_days" gorm:"column:le_7_days"`
	Le10Days int `json:"le_10_days" gorm:"column:le_10_days"`
	Gt10Days int `json:"gt_10_days" gorm:"column:gt_10_days"`
}

type OneUserReponse struct {
	AdminScanCenterModel                     []AdminScanCenterModel
	AdminOverallScanIndicatesAnalayticsModel []AdminOverallScanIndicatesAnalayticsModel
	UserAccessTimingModel                    []UserAccessTimingModel
	ListScanAppointmentCountModel            []ListScanAppointmentCountModel
	TotalCorrectEdit                         []TotalCorrectEditModel
	LeftRecommendation                       []RecommendationUserSQL
	RightRecommendation                      []RecommendationUserSQL
	DurationBucketModel                      []DurationBucketModel
	TechArtificats                           []Artificates
	ReportArtificats                         []Artificates
	OverAllAnalaytics                        []UsersOverAllAnalyticsModel
}

type UsersOverAllAnalyticsModel struct {
	UserId                   int     `json:"refUserId" gorm:"column:refUserId"`
	UserCustId               string  `json:"refUserCustId" gorm:"column:refUserCustId"`
	TotalCase                int     `json:"totalcase" gorm:"column:totalcase"`
	TotalSForm               int     `json:"totalsform" gorm:"column:totalsform"`
	TotalDaForm              int     `json:"totaldaform" gorm:"column:totaldaform"`
	TotalDbForm              int     `json:"totaldbform" gorm:"column:totaldbform"`
	TotalDcForm              int     `json:"totaldcform" gorm:"column:totaldcform"`
	TechArtificatsLeft       int     `json:"techartificatsleft" gorm:"column:techartificatsleft"`
	TechArtificateRight      int     `json:"techartificatsright" gorm:"column:techartificatsright"`
	ReportArtificatesLeft    int     `json:"reportartificatsleft" gorm:"column:reportartificatsleft"`
	ReportArtificatesRight   int     `json:"reportartificatsright" gorm:"column:reportartificatsright"`
	TotalTiming              float64 `json:"totaltiming" gorm:"column:totaltiming"`
	TotalReportCorrect       int     `json:"totalreportcorrect" gorm:"column:totalreportcorrect"`
	TotalReportEdit          int     `json:"totalreportedit" gorm:"column:totalreportedit"`
	Leftannualscreening      int     `json:"leftannualscreening" gorm:"column:leftannualscreening"`
	Leftusgsfu               int     `json:"leftusgsfu" gorm:"column:leftusgsfu"`
	LeftBiopsy               int     `json:"leftBiopsy" gorm:"column:leftBiopsy"`
	LeftBreastradiologist    int     `json:"leftBreastradiologist" gorm:"column:leftBreastradiologist"`
	LeftClinicalCorrelation  int     `json:"leftClinicalCorrelation" gorm:"column:leftClinicalCorrelation"`
	LeftOncoConsult          int     `json:"leftOncoConsult" gorm:"column:leftOncoConsult"`
	LeftRedo                 int     `json:"leftRedo" gorm:"column:leftRedo"`
	Rightannualscreening     int     `json:"rightannualscreening" gorm:"column:rightannualscreening"`
	Rightusgsfu              int     `json:"rightusgsfu" gorm:"column:rightusgsfu"`
	RightBiopsy              int     `json:"rightBiopsy" gorm:"column:rightBiopsy"`
	RightBreastradiologist   int     `json:"rightBreastradiologist" gorm:"column:rightBreastradiologist"`
	RightClinicalCorrelation int     `json:"rightClinicalCorrelation" gorm:"column:rightClinicalCorrelation"`
	RightOncoConsult         int     `json:"rightOncoConsult" gorm:"column:rightOncoConsult"`
	RightRedo                int     `json:"rightRedo" gorm:"column:rightRedo"`
}
