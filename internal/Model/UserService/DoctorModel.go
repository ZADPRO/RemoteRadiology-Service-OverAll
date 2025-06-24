package model

type MedicalLicenseSecurity struct {
	State                    string `json:"State" binding:"required" mapstructure:"State"`
	MedicalLicenseSecurityNo string `json:"MedicalLicenseSecurityNo" binding:"required" mapstructure:"MedicalLicenseSecurityNo"`
}

type UpdateMedicalLicenseSecurityModel struct {
	MLSId    int    `json:"refMLSId" binding:"required" mapstructure:"refMLSId"`
	MLSState string `json:"refMLSState" binding:"required" mapstructure:"refMLSState"`
	MLSNo    string `json:"refMLSNo" binding:"required" mapstructure:"refMLSNo"`
	MLStatus string `json:"refMLStatus" binding:"required" mapstructure:"refMLStatus"`
}

type MalpracticeInsureance struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
}

type UpdateMalpracticeInsureance struct {
	Id          int    `json:"id" binding:"required" mapstructure:"id"`
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
	Status      string `json:"status" binding:"required" mapstructure:"status"`
}

type DoctorRegisterReq struct {
	FirstName                   string                   `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName                    string                   `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg                  string                   `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                         string                   `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode          string                   `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo                     string                   `json:"phone" binding:"required" mapstructure:"phone"`
	Email                       string                   `json:"email" binding:"required" mapstructure:"email"`
	SocialSecurityNo            string                   `json:"social_security_no" binding:"required" mapstructure:"social_security_no"`
	DriversLicenseNo            string                   `json:"drivers_license" binding:"required" mapstructure:"drivers_license"`
	Specialization              string                   `json:"Specialization" binding:"required" mapstructure:"Specialization"`
	NPI                         string                   `json:"npi" binding:"required" mapstructure:"npi"`
	MedicalLicenseSecurity      []MedicalLicenseSecurity `json:"medical_license_security" binding:"required" mapstructure:"medical_license_security"`
	LicenseFiles                []LicenseFiles           `json:"license_files" binding:"required" mapstructure:"license_files"`
	MalpracticeInsuranceDetails []MalpracticeInsureance  `json:"malpracticeinsureance_files" binding:"required" mapstructure:"malpracticeinsureance_files"`
	DigitalSignature            string                   `json:"digital_signature" binding:"required" mapstructure:"digital_signature"`
	ScanCenterId                int                      `json:"refSCId" binding:"required"  mapstructure:"refSCId"`
}

type CreateDoctorModel struct {
	UserId         int    `json:"refUserId" gorm:"primaryKey;autoIncrement;column:refUserId"`
	UserCustId     string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId         int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName      string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName       string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserProfileImg string `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB            string `json:"refUserDOB" gorm:"column:refUserDOB"`
	Gender         string `json:"refUserGender" gorm:"column:refUserGender"`
	Status         bool   `json:"refUserStatus" gorm:"column:refUserStatus"`
	UserAgreement  bool   `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
}

func (CreateDoctorModel) TableName() string {
	return "Users"
}

type CreateDoctorCommunicationModel struct {
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

func (CreateDoctorCommunicationModel) TableName() string {
	return "userdomain.refCommunicationDomain"
}

type CreateDoctorAuthModel struct {
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	Password       string `json:"refADPassword" gorm:"column:refADPassword"`
	HashPassword   string `json:"refADHashPass" gorm:"column:refADHashPass"`
	PasswordStatus bool   `json:"refAHPassChangeStatus" gorm:"column:refAHPassChangeStatus"`
}

func (CreateDoctorAuthModel) TableName() string {
	return "userdomain.refAuthDomain"
}

type CreateDoctorDomainModel struct {
	DDId             int    `json:"refDDId" gorm:"primaryKey;autoIncrement;column:refDDId"`
	UserId           int    `json:"refUserId" gorm:"column:refUserId"`
	SocialSecurityNo string `json:"refDDSocialSecurityNo" gorm:"column:refDDSocialSecurityNo"`
	NPI              string `json:"refDDNPI" gorm:"column:refDDNPI"`
	DrivingLicense   string `json:"refDDDrivingLicense" gorm:"column:refDDDrivingLicense"`
	DigitalSignature string `json:"refDDDigitalSignature" gorm:"column:refDDDigitalSignature"`
	Specialization   string `json:"refDDSpecialization" gorm:"column:refDDSpecialization"`
}

func (CreateDoctorDomainModel) TableName() string {
	return "userdomain.refDoctorDomain"
}

type MedicalLicenseSecurityModel struct {
	MLSId    int    `json:"refMLSId" gorm:"primaryKey;autoIncrement;column:refMLSId"`
	UserId   int    `json:"refUserId" gorm:"column:refUserId"`
	MLSState string `json:"refMLSState" gorm:"column:refMLSState"`
	MLSNo    string `json:"refMLSNo" gorm:"column:refMLSNo"`
	MLStatus bool   `json:"refMLStatus" gorm:"column:refMLStatus"`
}

func (MedicalLicenseSecurityModel) TableName() string {
	return "userdomain.refMedicalLicenseSecurity"
}

type UpdateDoctorReq struct {
	ID                          int                                 `json:"id" binding:"required" mapstructure:"id"`
	FirstName                   string                              `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName                    string                              `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg                  string                              `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                         string                              `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode          string                              `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo                     string                              `json:"phone" binding:"required" mapstructure:"phone"`
	Email                       string                              `json:"email" binding:"required" mapstructure:"email"`
	SocialSecurityNo            string                              `json:"social_security_no" binding:"required" mapstructure:"social_security_no"`
	DriversLicenseNo            string                              `json:"drivers_license" binding:"required" mapstructure:"drivers_license"`
	Specialization              string                              `json:"Specialization" binding:"required" mapstructure:"Specialization"`
	NPI                         string                              `json:"npi" binding:"required" mapstructure:"npi"`
	MedicalLicenseSecurity      []UpdateMedicalLicenseSecurityModel `json:"medical_license_security" binding:"required" mapstructure:"medical_license_security"`
	LicenseFiles                []UpdateLicenseFiles                `json:"license_files" binding:"required" mapstructure:"license_files"`
	MalpracticeInsuranceDetails []UpdateMalpracticeInsureanceFiles  `json:"malpracticeinsureance_files" binding:"required" mapstructure:"malpracticeinsureance_files"`
	DigitalSignature            string                              `json:"digital_signature" binding:"required" mapstructure:"digital_signature"`
	Status                      bool                                `json:"status" binding:"required" mapstructure:"status"`
}

type GetAllDoctorData struct {
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
	SocialSecurityNo   string `json:"refDDSocialSecurityNo" gorm:"column:refDDSocialSecurityNo"`
	DrivingLicense     string `json:"refDDDrivingLicense" gorm:"column:refDDDrivingLicense"`
	Specialization     string `json:"refDDSpecialization" gorm:"column:refDDSpecialization"`
	DigitalSignature   string `json:"refDDDigitalSignature" gorm:"column:refDDDigitalSignature"`
	NPI                string `json:"refDDNPI" gorm:"column:refDDNPI"`
}

type GetMedicalLicenseSecurityModel struct {
	MLSId    int    `json:"refMLSId" gorm:"primaryKey;autoIncrement;column:refMLSId"`
	UserId   int    `json:"refUserId" gorm:"column:refUserId"`
	MLSState string `json:"refMLSState" gorm:"column:refMLSState"`
	MLSNo    string `json:"refMLSNo" gorm:"column:refMLSNo"`
	MLStatus bool   `json:"refMLStatus" gorm:"column:refMLStatus"`
}
