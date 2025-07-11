package model

type AdminScanCenterModel struct {
	Month             string `json:"month" gorm:"column:month"`
	MonthName         string `json:"month_name" gorm:"column:month_name"`
	TotalAppointments string `json:"total_appointments" gorm:"column:total_appointments"`
}

type AdminOverallScanIndicatesAnalayticsModel struct {
	TotalAppointments string `json:"total_appointments" gorm:"column:total_appointments"`
	SForm             string `json:"SForm" gorm:"column:SForm"`
	DaForm            string `json:"DaForm" gorm:"column:DaForm"`
	DbForm            string `json:"DbForm" gorm:"column:DbForm"`
	DcForm            string `json:"DcForm" gorm:"column:DcForm"`
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

type AdminOverallAnalyticsResponse struct {
	AdminScanCenterModel                     []AdminScanCenterModel
	AdminOverallScanIndicatesAnalayticsModel []AdminOverallScanIndicatesAnalayticsModel
	GetAllScaCenter                          []GetAllScaCenter
}

type AdminOverallOneAnalyticsReq struct {
	SCId      int    `json:"SCId" binding:"required" mapstructure:"SCId"`
	Monthyear string `json:"monthnyear" binding:"required" mapstructure:"monthnyear"`
}
