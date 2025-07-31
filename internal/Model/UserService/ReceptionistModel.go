package model

type ReceptionistRegisterReq struct {
	FirstName          string `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName           string `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg         string `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                string `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode string `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo            string `json:"phone" binding:"required" mapstructure:"phone"`
	Email              string `json:"email" binding:"required" mapstructure:"email"`
	SSNo               string `json:"social_security_no" binding:"required" mapstructure:"social_security_no"`
	DriversLicenseNo   string `json:"drivers_license" mapstructure:"drivers_license"`
	ScanCenterId       int    `json:"refSCId" binding:"required"  mapstructure:"refSCId"`
}

type UpdateReceptionistReq struct {
	ID                 int    `json:"id" binding:"required" mapstructure:"id"`
	FirstName          string `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName           string `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg         string `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                string `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode string `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo            string `json:"phone" binding:"required" mapstructure:"phone"`
	Email              string `json:"email" binding:"required" mapstructure:"email"`
	SocialSecurityNo   string `json:"social_security_no" binding:"required" mapstructure:"social_security_no"`
	DriversLicenseNo   string `json:"drivers_license_no" mapstructure:"drivers_license_no"`
	Status             bool   `json:"status" mapstructure:"status"`
}

type CreateReceptionstDomainModel struct {
	RDId           int    `json:"refRDId" gorm:"primaryKey;autoIncrement;column:refRDId"`
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	SSId           string `json:"refRDSSId" gorm:"column:refRDSSId"`
	DrivingLicense string `json:"refRDDrivingLicense" gorm:"column:refRDDrivingLicense"`
}

func (CreateReceptionstDomainModel) TableName() string {
	return "userdomain.refReceptionstDomain"
}

type MapScanCenterModel struct {
	SCMId     int  `json:"refSCMId" gorm:"primaryKey;autoIncrement;column:refSCMId"`
	UserId    int  `json:"refUserId" gorm:"column:refUserId"`
	SCId      int  `json:"refSCId" gorm:"column:refSCId"`
	RTId      int  `json:"refRTId" gorm:"column:refRTId"`
	SCMStatus bool `json:"refSCMStatus" gorm:"column:refSCMStatus"`
}

func (MapScanCenterModel) TableName() string {
	return "map.refScanCenterMap"
}

type GetAllReceptionistData struct {
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
	SocialSecurityNo   string `json:"refRDSSId" gorm:"column:refRDSSId"`
	DriversLicenseNo   string `json:"refRDDrivingLicense" gorm:"column:refRDDrivingLicense"`
}
