package model

type ScribeRegisterReq struct {
	FirstName            string                      `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName             string                      `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg           string                      `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                  string                      `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode   string                      `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo              string                      `json:"phone" binding:"required" mapstructure:"phone"`
	Email                string                      `json:"email" binding:"required" mapstructure:"email"`
	DriversLicense       string                      `json:"drivers_license" binding:"required" mapstructure:"drivers_license"`
	Pan                  string                      `json:"pan" binding:"required" mapstructure:"pan"`
	Aadhar               string                      `json:"aadhar" binding:"required" mapstructure:"aadhar"`
	EducationCertificate []EducationCertificateFiles `json:"education_certificate" binding:"required" mapstructure:"education_certificate"`
}

type CreateScribeDomainModel struct {
	SDId           int    `json:"refSDId" gorm:"primaryKey;autoIncrement;column:refSDId"`
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	Pan            string `json:"refSDPan" gorm:"column:refSDPan"`
	Aadhar         string `json:"refSDAadhar" gorm:"column:refSDAadhar"`
	DrivingLicense string `json:"refSDDrivingLicense" gorm:"column:refSDDrivingLicense"`
}

func (CreateScribeDomainModel) TableName() string {
	return "userdomain.refScribeDomain"
}

type ScribeUpdateReq struct {
	FirstName                   string          `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName                    string          `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg                  string          `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                         string          `json:"dob" binding:"required" mapstructure:"dob"`
	Gender                      string          `json:"gender" binding:"required" mapstructure:"gender"`
	PhoneNoCountryCode          string          `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo                     string          `json:"phone" binding:"required" mapstructure:"phone"`
	Email                       string          `json:"email" binding:"required" mapstructure:"email"`
	EmergencyPhoneNoCountryCode string          `json:"emergency_phoneCountryCode" binding:"required" mapstructure:"emergency_phoneCountryCode"`
	EmergencyPhone              string          `json:"emergency_phone" binding:"required" mapstructure:"emergency_phone"`
	TimeZone                    string          `json:"time_zone" binding:"required" mapstructure:"time_zone"`
	Shifts                      string          `json:"shifts" binding:"required" mapstructure:"shifts"`
	CVFiles                     []UpdateCVFiles `json:"cv_files" binding:"required" mapstructure:"cv_files"`
	TrainedEaseQT               bool            `json:"trained_ease_qt" binding:"required" mapstructure:"trained_ease_qt"`
}

type UpdateScribeReq struct {
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

type GetAllScribeData struct {
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
	Pan                string `json:"refSDPan" gorm:"column:refSDPan"`
	Aadhar             string `json:"refSDAadhar" gorm:"column:refSDAadhar"`
	DrivingLicense     string `json:"refSDDrivingLicense" gorm:"column:refSDDrivingLicense"`
}
