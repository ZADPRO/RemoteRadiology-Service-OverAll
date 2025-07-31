package model

type ListPatientBrochureReq struct {
	ScancenterId int `json:"scancenterId" mapstructure:"scancenterId"`
}

type ListPatientConsentReq struct {
	ScanCenterId []int `json:"scancenterId" binding:"required"`
}

type ListPatientConsent struct {
	Consent string `json:"refAppointmentConsent" gorm:"column:refAppointmentConsent"`
}

type ListWGPatientBrochureModel struct {
	FDId         int    `json:"refFDId" gorm:"column:refFDId"`
	FTId         int    `json:"refFTId" gorm:"column:refFTId"`
	SCId         int    `json:"refSCId" gorm:"column:refSCId"`
	FDData       string `json:"refFDData" gorm:"column:refFDData"`
	FDAccessData bool   `json:"refFDAccessData" gorm:"column:refFDAccessData"`
}

type ListPatientBrochureRes struct {
	Status                 bool   `json:"status"`
	WGPatientBrochure      string `json:"refWGPatientData"`
	SCBrochureAccessStatus bool   `json:"SCAccessStatus"`
	SCPatientBrochure      string `json:"refSCPatientData"`
}

type UpdatePatientBroucherReq struct {
	Brochure     string `json:"data" binding:"required" mapstructure:"data"`
	ScancenterId int    `json:"scancenterId" mapstructure:"scancenterId"`
	AccessStatus bool   `json:"accessStatus" mapstructure:"accessStatus"`
}

type UpdatePatientBroucherResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
