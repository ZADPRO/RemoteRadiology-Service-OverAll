package model

type GetDailyListModel struct {
	FromDate string `json:"fromDate" mapstructure:"fromDate"`
	ToDate   string `json:"toDate" mapstructure:"toDate"`
}

type DailyListReponse struct {
	AppointmentDate                   string `json:"AppointmentDate" gorm:"column:AppointmentDate"`
	RefUserCustId                     string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RefUserFirstName                  string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	RefCategoryId                     string `json:"refCategoryId" gorm:"column:refCategoryId"`
	ScanSide                          string `json:"scanSide" gorm:"column:scanSide"`
	HandlerName                       string `json:"handlerName" gorm:"column:handlerName"`
	RefAppointmentImpression          string `json:"refAppointmentImpression" gorm:"column:refAppointmentImpression"`
	RefAppointmentRecommendation      string `json:"refAppointmentRecommendation" gorm:"column:refAppointmentRecommendation"`
	RefAppointmentImpressionRight     string `json:"refAppointmentImpressionRight" gorm:"column:refAppointmentImpressionRight"`
	RefAppointmentRecommendationRight string `json:"refAppointmentRecommendationRight" gorm:"column:refAppointmentRecommendationRight"`
}
