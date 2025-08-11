package model

type ManagerRegisterReq struct {
	FirstName            string                      `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName             string                      `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg           string                      `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                  string                      `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode   string                      `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo              string                      `json:"phone" binding:"required" mapstructure:"phone"`
	Email                string                      `json:"email" binding:"required" mapstructure:"email"`
	DriversLicense       string                      `json:"drivers_license" mapstructure:"drivers_license"`
	Pan                  string                      `json:"pan" binding:"required" mapstructure:"pan"`
	Aadhar               string                      `json:"aadhar" binding:"required" mapstructure:"aadhar"`
	EducationCertificate []EducationCertificateFiles `json:"education_certificate" binding:"required" mapstructure:"education_certificate"`
}

type CreateManagerModel struct {
	UserId         int    `json:"refUserId" gorm:"primaryKey;autoIncrement;column:refUserId"`
	UserCustId     string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId         int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName      string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName       string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserProfileImg string `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB            string `json:"refUserDOB" gorm:"column:refUserDOB"`
	Gender         string `json:"refUserGender" gorm:"column:refUserGender"`
	Status         bool   `json:"refUserStatus" gorm:"column:refUserStatus"`
}

func (CreateManagerModel) TableName() string {
	return "Users"
}

type CreateManagerCommunicationModel struct {
	CODOId                      int    `json:"refCODOId" gorm:"primaryKey;autoIncrement;column:refCODOId"`
	UserId                      int    `json:"refUserId" gorm:"column:refUserId"`
	PhoneNoCountryCode          string `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	PhoneNo                     string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	EmergencyPhoneNoCountryCode string `json:"refCODOPhoneNo2CountryCode" gorm:"column:refCODOPhoneNo2CountryCode"`
	EmergencyPhone              string `json:"refCODOPhoneNo2" gorm:"column:refCODOPhoneNo2"`
	Email                       string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	DoorNo                      string `json:"refCODODoorNo" gorm:"column:refCODODoorNo"`
	Street                      string `json:"refCODOStreet" gorm:"column:refCODOStreet"`
	District                    string `json:"refCODODistrict" gorm:"column:refCODODistrict"`
	State                       string `json:"refCODOState" gorm:"column:refCODOState"`
	Country                     string `json:"refCODOCountry" gorm:"column:refCODOCountry"`
	Pincode                     string `json:"refCODOPincode" gorm:"column:refCODOPincode"`
}

func (CreateManagerCommunicationModel) TableName() string {
	return "userdomain.refCommunicationDomain"
}

type CreateManagerAuthModel struct {
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	Password       string `json:"refADPassword" gorm:"column:refADPassword"`
	HashPassword   string `json:"refADHashPass" gorm:"column:refADHashPass"`
	PasswordStatus bool   `json:"refAHPassChangeStatus" gorm:"column:refAHPassChangeStatus"`
}

func (CreateManagerAuthModel) TableName() string {
	return "userdomain.refAuthDomain"
}

type CreateManagerDomainModel struct {
	MDId           int    `json:"refMDId" gorm:"primaryKey;autoIncrement;column:refMDId"`
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	Pan            string `json:"refMDPan" gorm:"column:refMDPan"`
	Aadhar         string `json:"refMDAadhar" gorm:"column:refMDAadhar"`
	DrivingLicense string `json:"refMDDrivingLicense" gorm:"column:refMDDrivingLicense"`
}

func (CreateManagerDomainModel) TableName() string {
	return "userdomain.refManagerDomain"
}

type UpdateManagerReq struct {
	ID                   int                               `json:"id" binding:"required" mapstructure:"id"`
	FirstName            string                            `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName             string                            `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg           string                            `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                  string                            `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode   string                            `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo              string                            `json:"phone" binding:"required" mapstructure:"phone"`
	Email                string                            `json:"email" binding:"required" mapstructure:"email"`
	Pan                  string                            `json:"pan" binding:"required" mapstructure:"pan"`
	Aadhar               string                            `json:"aadhar" binding:"required" mapstructure:"aadhar"`
	DriversLicense       string                            `json:"drivers_license" binding:"required" mapstructure:"drivers_license"`
	EducationCertificate []UpdateEducationCertificateFiles `json:"education_certificate" binding:"required" mapstructure:"education_certificate"`
	Status               bool                              `json:"status" binding:"required" mapstructure:"status"`
}

type GetAllManagerData struct {
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
	MDPan              string `json:"refMDPan" gorm:"column:refMDPan"`
	MDAadhar           string `json:"refMDAadhar" gorm:"column:refMDAadhar"`
	MDDrivingLicense   string `json:"refMDDrivingLicense" gorm:"column:refMDDrivingLicense"`
}
