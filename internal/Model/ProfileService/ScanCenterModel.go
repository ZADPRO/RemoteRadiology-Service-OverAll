package model

type GetAllScaCenter struct {
	SCId           uint      `json:"refSCId" gorm:"column:refSCId"`
	SCCustId       string    `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCProfile      string    `json:"refSCProfile" gorm:"column:refSCProfile"`
	SCName         string    `json:"refSCName" gorm:"column:refSCName"`
	SCAddress      string    `json:"refSCAddress" gorm:"column:refSCAddress"`
	SCPhoneNo1     string    `json:"refSCPhoneNo1" gorm:"column:refSCPhoneNo1"`
	SCEmail        string    `json:"refSCEmail" gorm:"column:refSCEmail"`
	SCWebsite      string    `json:"refSCWebsite" gorm:"column:refSCWebsite"`
	SCAppointments bool      `json:"refSCAppointments" gorm:"column:refSCAppointments"`
	ProfileImgFile *FileData `json:"profileImgFile" gorm:"-"`
}

type Mapping struct {
	SCId int `json:"refSCId" gorm:"column:refSCId"`
}
