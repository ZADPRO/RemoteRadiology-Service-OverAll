package model

type ScanCenterRegisterReq struct {
	CustId       string `json:"cust_id" binding:"required" mapstructure:"cust_id"`
	Name         string `json:"name" binding:"required" mapstructure:"name"`
	Address      string `json:"address" binding:"required" mapstructure:"address"`
	Email        string `json:"email" binding:"required" mapstructure:"email"`
	Website      string `json:"website" binding:"required" mapstructure:"website"`
	Telephone    string `json:"telephone" binding:"required" mapstructure:"telephone"`
	Logo         string `json:"logo" binding:"required" mapstructure:"logo"`
	Appointments bool   `json:"appointments" mapstructure:"appointments"`
}

type ScanCenterModel struct {
	Id                  int    `json:"refSCId" gorm:"column:refSCId;primaryKey;autoIncrement"`
	CustId              string `json:"refSCCustId" gorm:"column:refSCCustId"`
	Logo                string `json:"refSCProfile" gorm:"column:refSCProfile"`
	Name                string `json:"refSCName" gorm:"column:refSCName"`
	Address             string `json:"refSCAddress" gorm:"column:refSCAddress"`
	PhoneNo1CountryCode string `json:"refSCPhoneNo1CountryCode" gorm:"column:refSCPhoneNo1CountryCode"`
	PhoneNo1            string `json:"refSCPhoneNo1" gorm:"column:refSCPhoneNo1"`
	PhoneNo2CountryCode string `json:"refSCPhoneNo2CountryCode" gorm:"column:refSCPhoneNo2CountryCode"`
	PhoneNo2            string `json:"refSCPhoneNo2" gorm:"column:refSCPhoneNo2"`
	Email               string `json:"refSCEmail" gorm:"column:refSCEmail"`
	Website             string `json:"refSCWebsite" gorm:"column:refSCWebsite"`
	Appointments        bool   `json:"refSCAppointments" gorm:"column:refSCAppointments"`
	SCStatus            bool   `json:"refSCStatus" gorm:"column:refSCStatus"`
}

func (ScanCenterModel) TableName() string {
	return "public.ScanCenter"
}

type UpdateScanCentertReq struct {
	ID           int    `json:"id" binding:"required" mapstructure:"id"`
	Name         string `json:"name" binding:"required" mapstructure:"name"`
	Address      string `json:"address" binding:"required" mapstructure:"address"`
	Email        string `json:"email" binding:"required" mapstructure:"email"`
	Website      string `json:"website" binding:"required" mapstructure:"website"`
	Telephone    string `json:"telephone" binding:"required" mapstructure:"telephone"`
	Logo         string `json:"logo" binding:"required" mapstructure:"logo"`
	Appointments bool   `json:"appointments" mapstructure:"appointments"`
	Status       bool   `json:"status" mapstructure:"status"`
}
