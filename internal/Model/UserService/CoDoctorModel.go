package model

type GetAllCODoctorData struct {
	UserId             int    `json:"refUserId" gorm:"primaryKey;autoIncrement;column:refUserId"`
	UserCustId         string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId             int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName          string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName           string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserProfileImg     string `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB                string `json:"refUserDOB" gorm:"column:refUserDOB"`
	Status             bool   `json:"refUserStatus" gorm:"column:refUserStatus"`
	PhoneNoCountryCode string `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	PhoneNo            string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email              string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	SocialSecurityNo   string `json:"refCDSocialSecurityNo" gorm:"column:refCDSocialSecurityNo"`
	DrivingLicense     string `json:"refCDDrivingLicense" gorm:"column:refCDDrivingLicense"`
	Specialization     string `json:"refCDSpecialization" gorm:"column:refCDSpecialization"`
	DigitalSignature   string `json:"refCDDigitalSignature" gorm:"column:refCDDigitalSignature"`
	NPI                string `json:"refCDNPI" gorm:"column:refCDNPI"`
}

type CreateCoDoctorDomainModel struct {
	DDId             int    `json:"refCDDId" gorm:"primaryKey;autoIncrement;column:refCDDId"`
	UserId           int    `json:"refUserId" gorm:"column:refUserId"`
	SocialSecurityNo string `json:"refCDSocialSecurityNo" gorm:"column:refCDSocialSecurityNo"`
	NPI              string `json:"refCDNPI" gorm:"column:refCDNPI"`
	DrivingLicense   string `json:"refCDDrivingLicense" gorm:"column:refCDDrivingLicense"`
	DigitalSignature string `json:"refCDDigitalSignature" gorm:"column:refCDDigitalSignature"`
	Specialization   string `json:"refCDSpecialization" gorm:"column:refCDSpecialization"`
}

func (CreateCoDoctorDomainModel) TableName() string {
	return "userdomain.refCoDoctorDomain"
}

type CreateMalpractice struct {
	LId          int    `json:"refMPId" gorm:"primaryKey;autoIncrement;column:refMPId"`
	UserId       int    `json:"refUserId" gorm:"column:refUserId"`
	LFileName    string `json:"refMPFileName" gorm:"column:refMPFileName"`
	LOldFileName string `json:"refMPOldFileName" gorm:"column:refMPOldFileName"`
	LStatus      bool   `json:"refMPStatus" gorm:"column:refMPStatus"`
}

func (CreateMalpractice) TableName() string {
	return "userdomain.refMalpractice"
}
