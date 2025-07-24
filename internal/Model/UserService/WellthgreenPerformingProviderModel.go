package model

type CreateWGPPDomainModel struct {
	RDId             int    `json:"refWGPPId" gorm:"primaryKey;autoIncrement;column:refWGPPId"`
	UserId           int    `json:"refUserId" gorm:"column:refUserId"`
	MBBSRegNo        string `json:"refWGPPMBBSRegNo" gorm:"column:refWGPPMBBSRegNo"`
	MDRegNo          string `json:"refWGPPMDRegNo" gorm:"column:refWGPPMDRegNo"`
	Specialization   string `json:"refWGPPSpecialization" gorm:"column:refWGPPSpecialization"`
	Pan              string `json:"refWGPPPan" gorm:"column:refWGPPPan"`
	Aadhar           string `json:"refWGPPAadhar" gorm:"column:refWGPPAadhar"`
	DrivingLicense   string `json:"refWGPPDrivingLicense" gorm:"column:refWGPPDrivingLicense"`
	DigitalSignature string `json:"refWGPPDigitalSignature" gorm:"column:refWGPPDigitalSignature"`
}

func (CreateWGPPDomainModel) TableName() string {
	return "userdomain.refWellthgreenPerformingProvider"
}
