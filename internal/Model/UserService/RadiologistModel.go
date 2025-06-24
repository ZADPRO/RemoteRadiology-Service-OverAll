package model

type LicenseFiles struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
}

type CVFiles struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
}

type EducationCertificateFiles struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
}

type MalPracticeInsureance struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
}

type RadiologistRegisterReq struct {
	FirstName              string                   `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName               string                   `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg             string                   `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                    string                   `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode     string                   `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo                string                   `json:"phone" binding:"required" mapstructure:"phone"`
	Email                  string                   `json:"email" binding:"required" mapstructure:"email"`
	MBBSRegisterNumber     string                   `json:"mbbs_register_number" binding:"required" mapstructure:"mbbs_register_number"`
	MDRegisterNumber       string                   `json:"md_register_number" binding:"required" mapstructure:"md_register_number"`
	Specialization         string                   `json:"specialization" binding:"required" mapstructure:"specialization"`
	Pan                    string                   `json:"pan" binding:"required" mapstructure:"pan"`
	Aadhar                 string                   `json:"aadhar" binding:"required" mapstructure:"aadhar"`
	DriversLicense         string                   `json:"drivers_license" binding:"required" mapstructure:"drivers_license"`
	MedicalLicenseSecurity []MedicalLicenseSecurity `json:"medical_license_security" binding:"required" mapstructure:"medical_license_security"`
	MalPracticeInsureance  []MalPracticeInsureance  `json:"malpracticeinsureance_files" binding:"required" mapstructure:"malpracticeinsureance_files"`
	CVFiles                []CVFiles                `json:"cv_files" binding:"required" mapstructure:"cv_files"`
	LicenseFiles           []LicenseFiles           `json:"license_files" binding:"required" mapstructure:"license_files"`
	DigitalSignature       string                   `json:"digital_signature" binding:"required" mapstructure:"digital_signature"`
}

type CreateRadiologyModel struct {
	UserId         int    `json:"refUserId" gorm:"primaryKey;autoIncrement;column:refUserId"`
	UserCustId     string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId         int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName      string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName       string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserProfileImg string `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB            string `json:"refUserDOB" gorm:"column:refUserDOB"`
	Status         bool   `json:"refUserStatus" gorm:"column:refUserStatus"`
	UserAgreement  bool   `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
}

func (CreateRadiologyModel) TableName() string {
	return "Users"
}

type CreateRadiologistCommunicationModel struct {
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

func (CreateRadiologistCommunicationModel) TableName() string {
	return "userdomain.refCommunicationDomain"
}

type CreateRadiologistAuthModel struct {
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	Password       string `json:"refADPassword" gorm:"column:refADPassword"`
	HashPassword   string `json:"refADHashPass" gorm:"column:refADHashPass"`
	PasswordStatus bool   `json:"refAHPassChangeStatus" gorm:"column:refAHPassChangeStatus"`
}

func (CreateRadiologistAuthModel) TableName() string {
	return "userdomain.refAuthDomain"
}

type CreateRadiologistDomainModel struct {
	RDId             int    `json:"refRADId" gorm:"primaryKey;autoIncrement;column:refRADId"`
	UserId           int    `json:"refUserId" gorm:"column:refUserId"`
	MBBSRegNo        string `json:"refRAMBBSRegNo" gorm:"column:refRAMBBSRegNo"`
	MDRegNo          string `json:"refRAMDRegNo" gorm:"column:refRAMDRegNo"`
	Specialization   string `json:"refRASpecialization" gorm:"column:refRASpecialization"`
	Pan              string `json:"refRAPan" gorm:"column:refRAPan"`
	Aadhar           string `json:"refRAAadhar" gorm:"column:refRAAadhar"`
	DrivingLicense   string `json:"refRADrivingLicense" gorm:"column:refRADrivingLicense"`
	DigitalSignature string `json:"refRADigitalSignature" gorm:"column:refRADigitalSignature"`
}

func (CreateRadiologistDomainModel) TableName() string {
	return "userdomain.refRadiologistDomain"
}

type UpdateCVFiles struct {
	Id          int    `json:"id" binding:"required" mapstructure:"id"`
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
	Status      string `json:"status" binding:"required" mapstructure:"status"`
}

type UpdateEducationCertificateFiles struct {
	Id          int    `json:"id" binding:"required" mapstructure:"id"`
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
	Status      string `json:"status" binding:"required" mapstructure:"status"`
}

type UpdateLicenseFiles struct {
	Id          int    `json:"id" binding:"required" mapstructure:"id"`
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
	Status      string `json:"status" binding:"required" mapstructure:"status"`
}

type UpdateMalpracticeInsureanceFiles struct {
	Id          int    `json:"id" binding:"required" mapstructure:"id"`
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
	Status      string `json:"status" binding:"required" mapstructure:"status"`
}

type UpdateRadiologistReq struct {
	ID                          int                                 `json:"id" binding:"required" mapstructure:"id"`
	FirstName                   string                              `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName                    string                              `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg                  string                              `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                         string                              `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode          string                              `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo                     string                              `json:"phone" binding:"required" mapstructure:"phone"`
	Email                       string                              `json:"email" binding:"required" mapstructure:"email"`
	MBBSRegisterNumber          string                              `json:"mbbs_register_number" binding:"required" mapstructure:"mbbs_register_number"`
	MDRegisterNumber            string                              `json:"md_register_number" binding:"required" mapstructure:"md_register_number"`
	Specialization              string                              `json:"specialization" binding:"required" mapstructure:"specialization"`
	Pan                         string                              `json:"pan" binding:"required" mapstructure:"pan"`
	Aadhar                      string                              `json:"aadhar" binding:"required" mapstructure:"aadhar"`
	DriversLicense              string                              `json:"drivers_license" binding:"required" mapstructure:"drivers_license"`
	DigitalSignature            string                              `json:"digital_signature" binding:"required" mapstructure:"digital_signature"`
	MedicalLicenseSecurity      []UpdateMedicalLicenseSecurityModel `json:"medical_license_security" binding:"required" mapstructure:"medical_license_security"`
	MalpracticeInsuranceDetails []UpdateMalpracticeInsureanceFiles  `json:"malpracticeinsureance_files" binding:"required" mapstructure:"malpracticeinsureance_files"`
	CVFiles                     []UpdateCVFiles                     `json:"cv_files" binding:"required" mapstructure:"cv_files"`
	LicenseFiles                []UpdateLicenseFiles                `json:"license_files" binding:"required" mapstructure:"license_files"`
	Status                      bool                                `json:"status" mapstructure:"status"`
}

type GetAllRadiologistData struct {
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
	MBBSRegisterNumber string `json:"refRAMBBSRegNo" gorm:"column:refRAMBBSRegNo"`
	MDRegisterNumber   string `json:"refRAMDRegNo" gorm:"column:refRAMDRegNo"`
	Specialization     string `json:"refRASpecialization" gorm:"column:refRASpecialization"`
	Pan                string `json:"refRAPan" gorm:"column:refRAPan"`
	Aadhar             string `json:"refRAAadhar" gorm:"column:refRAAadhar"`
	DrivingLicense     string `json:"refRADrivingLicense" gorm:"column:refRADrivingLicense"`
	DigitalSignature   string `json:"refRADigitalSignature" gorm:"column:refRADigitalSignature"`
}

type CreateRadiologistMalpracticeModel struct {
	LId          int    `json:"refMPId" gorm:"primaryKey;autoIncrement;column:refMPId"`
	UserId       int    `json:"refUserId" gorm:"column:refUserId"`
	LFileName    string `json:"refMPFileName" gorm:"column:refMPFileName"`
	LOldFileName string `json:"refMPOldFileName" gorm:"column:refMPOldFileName"`
	LStatus      bool   `json:"refMPStatus" gorm:"column:refMPStatus"`
}

func (CreateRadiologistMalpracticeModel) TableName() string {
	return "userdomain.refMalpractice"
}

type CreateRadiologistLicenseModel struct {
	LId          int    `json:"refLId" gorm:"primaryKey;autoIncrement;column:refLId"`
	UserId       int    `json:"refUserId" gorm:"column:refUserId"`
	LFileName    string `json:"refLFileName" gorm:"column:refLFileName"`
	LOldFileName string `json:"refLOldFileName" gorm:"column:refLOldFileName"`
	LStatus      bool   `json:"refLStatus" gorm:"column:refLStatus"`
}

func (CreateRadiologistLicenseModel) TableName() string {
	return "userdomain.refLicences"
}

type CreateRadiologistCVFilesModel struct {
	CVId          int    `json:"refCVID" gorm:"primaryKey;autoIncrement;column:refCVID"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	CVFileName    string `json:"refCVFileName" gorm:"column:refCVFileName"`
	CVOldFileName string `json:"refCVOldFileName" gorm:"column:refCVOldFileName"`
	CVStatus      bool   `json:"refCVStatus" gorm:"column:refCVStatus"`
}

func (CreateRadiologistCVFilesModel) TableName() string {
	return "userdomain.refCV"
}

type CreateEducationCertificateModel struct {
	ECId          int    `json:"refECId" gorm:"primaryKey;autoIncrement;column:refECId"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	ECFileName    string `json:"refECFileName" gorm:"column:refECFileName"`
	ECOldFileName string `json:"refECOldFileName" gorm:"column:refECOldFileName"`
	ECStatus      bool   `json:"refECStatus" gorm:"column:refECStatus"`
}

func (CreateEducationCertificateModel) TableName() string {
	return "userdomain.refEducationCertificate"
}

type GetCVFilesModel struct {
	CVId          int    `json:"refCVID" gorm:"primaryKey;autoIncrement;column:refCVID"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	CVFileName    string `json:"refCVFileName" gorm:"column:refCVFileName"`
	CVOldFileName string `json:"refCVOldFileName" gorm:"column:refCVOldFileName"`
	CVStatus      bool   `json:"refCVStatus" gorm:"column:refCVStatus"`
}

type GetEducationCertificateFilesModel struct {
	ECId          int    `json:"refECId" gorm:"primaryKey;autoIncrement;column:refECId"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	ECFileName    string `json:"refECFileName" gorm:"column:refECFileName"`
	ECOldFileName string `json:"refECOldFileName" gorm:"column:refECOldFileName"`
	ECStatus      bool   `json:"refECStatus" gorm:"column:refECStatus"`
}

type GetLicenseFilesModel struct {
	LId          int    `json:"refLId" gorm:"primaryKey;autoIncrement;column:refLId"`
	UserId       int    `json:"refUserId" gorm:"column:refUserId"`
	LFileName    string `json:"refLFileName" gorm:"column:refLFileName"`
	LOldFileName string `json:"refLOldFileName" gorm:"column:refLOldFileName"`
	LStatus      bool   `json:"refLStatus" gorm:"column:refLStatus"`
}

type GetMalpracticeFilesModel struct {
	LId          int    `json:"refMPId" gorm:"primaryKey;autoIncrement;column:refMPId"`
	UserId       int    `json:"refUserId" gorm:"column:refUserId"`
	LFileName    string `json:"refMPFileName" gorm:"column:refMPFileName"`
	LOldFileName string `json:"refMPOldFileName" gorm:"column:refMPOldFileName"`
	LStatus      bool   `json:"refMPStatus" gorm:"column:refMPStatus"`
}
