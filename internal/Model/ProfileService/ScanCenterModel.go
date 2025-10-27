package model

type GetAllScaCenter struct {
	SCId               uint      `json:"refSCId" gorm:"column:refSCId"`
	SCCustId           string    `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCProfile          string    `json:"refSCProfile" gorm:"column:refSCProfile"`
	SCName             string    `json:"refSCName" gorm:"column:refSCName"`
	SCAddress          string    `json:"refSCAddress" gorm:"column:refSCAddress"`
	SCPhoneNo1         string    `json:"refSCPhoneNo1" gorm:"column:refSCPhoneNo1"`
	SCEmail            string    `json:"refSCEmail" gorm:"column:refSCEmail"`
	SCWebsite          string    `json:"refSCWebsite" gorm:"column:refSCWebsite"`
	SCAppointments     bool      `json:"refSCAppointments" gorm:"column:refSCAppointments"`
	SCDisclamer        string    `json:"refSCDisclamer" gorm:"column:refSCDisclamer"`
	SCBrouchure        string    `json:"refSCBrouchure" gorm:"column:refSCBrouchure"`
	SCGuidelines       string    `json:"refSCGuidelines" gorm:"column:refSCGuidelines"`
	ProfileImgFile     *FileData `json:"profileImgFile" gorm:"-"`
	SCStatus           bool      `json:"refSCStatus" gorm:"column:refSCStatus"`
	SCConsultantStatus bool      `json:"refSCConsultantStatus" gorm:"column:refSCConsultantStatus"`
	SCConsultantLink   string    `json:"refSCConsultantLink" gorm:"column:refSCConsultantLink"`
}

type Mapping struct {
	SCId     int    `json:"refSCId" gorm:"column:refSCId"`
	SCCustId string `json:"refSCCustId" gorm:"column:refSCCustId"`
}
