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

type AdminOverallAnalyticsResponse struct {
	AdminScanCenterModel                     []AdminScanCenterModel
	AdminOverallScanIndicatesAnalayticsModel []AdminOverallScanIndicatesAnalayticsModel
	GetAllScaCenter                          []GetAllScaCenter
	UserListIdsData                          []UserListIdsData
}

type ListScanAppointmentCountModel struct {
	SCId              int    `json:"refSCId" gorm:"column:refSCId"`
	SCName            string `json:"refSCName" gorm:"column:refSCName"`
	TotalAppointments int    `json:"total_appointments" gorm:"column:total_appointments"`
}

type AdminOverallOneAnalyticsReq struct {
	SCId      int    `json:"SCId" mapstructure:"SCId"`
	Monthyear string `json:"monthnyear" binding:"required" mapstructure:"monthnyear"`
}

type OneUserReq struct {
	UserId    int    `json:"userId" binding:"required" mapstructure:"userId"`
	RoleId    int    `json:"roleId" binding:"required" mapstructure:"roleId"`
	Monthyear string `json:"monthnyear" binding:"required" mapstructure:"monthnyear"`
}

type TotalCorrectEditModel struct {
	TotalCorrect string `json:"totalCorrect" gorm:"column:totalCorrect"`
	TotalEdit    string `json:"totalEdit" gorm:"column:totalEdit"`
}

type ImpressionModel struct {
	Impression string `json:"impression" gorm:"column:impression"`
	Count      int    `json:"count" gorm:"column:count"`
}

type OneUserReponse struct {
	AdminScanCenterModel                     []AdminScanCenterModel
	AdminOverallScanIndicatesAnalayticsModel []AdminOverallScanIndicatesAnalayticsModel
	UserAccessTimingModel                    []UserAccessTimingModel
	ListScanAppointmentCountModel            []ListScanAppointmentCountModel
	TotalCorrectEdit                         []TotalCorrectEditModel
	ImpressionModel                          []ImpressionModel
}
