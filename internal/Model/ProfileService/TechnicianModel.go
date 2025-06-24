package model

type GetTechnicianOne struct {
	RefUserId              uint                `json:"refUserId" gorm:"column:refUserId"`
	CustId                 string              `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId                 int                 `json:"refRTId" gorm:"column:refRTId"`
	FirstName              string              `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName               string              `json:"refUserLastName" gorm:"column:refUserLastName"`
	ProfileImg             string              `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB                    string              `json:"refUserDOB" gorm:"column:refUserDOB"`
	UserStatus             bool                `json:"refUserStatus" gorm:"column:refUserStatus"`
	Agreement              bool                `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
	PhoneNumberCountryCode string              `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	PhoneNumber            string              `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email                  string              `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	TrainedEaseQTStatus    bool                `json:"refTDTrainedEaseQTStatus" gorm:"column:refTDTrainedEaseQTStatus"`
	SSNo                   string              `json:"refTDSSNo" gorm:"column:refTDSSNo"`
	DrivingLicense         string              `json:"refTDDrivingLicense" gorm:"column:refTDDrivingLicense"`
	DigitalSignature       string              `json:"refTDDigitalSignature" gorm:"column:refTDDigitalSignature"`
	ProfileImgFile         *FileData           `json:"profileImgFile" gorm:"-"`
	DigitalSignatureFile   *FileData           `json:"digitalSignatureFile" gorm:"-"`
	DrivingLicenseFile     *FileData           `json:"drivingLicenseFile" gorm:"-"`
	LicenseFiles           []LicenseFilesModel `json:"licenseFiles" gorm:"-"`
}

// type GetRoleId struct {
// 	RoleId int `json:"roleId" binding:"required" mapstructure:"roleId"`
// }

// type GetTechnicianReq struct {
// 	UserID int `json:"refUserId" binding:"required" mapstructure:"refUserId"`
// }

// type GetTechnician struct {
// 	RefUserId   uint   `json:"refUserId" gorm:"column:refUserId"`
// 	CustId      string `json:"refUserCustId" gorm:"column:refUserCustId"`
// 	RoleId      int    `json:"refRTId" gorm:"column:refRTId"`
// 	FirstName   string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
// 	LastName    string `json:"refUserLastName" gorm:"column:refUserLastName"`
// 	UserStatus  string `json:"refUserStatus" gorm:"column:refUserStatus"`
// 	PhoneNumber string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
// 	Email       string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
// }

// type GetTechnicianOne struct {
// 	RefUserId            uint              `json:"refUserId" gorm:"column:refUserId"`
// 	CustId               string            `json:"refUserCustId" gorm:"column:refUserCustId"`
// 	RoleId               int               `json:"refRTId" gorm:"column:refRTId"`
// 	FirstName            string            `json:"refUserFirstName" gorm:"column:refUserFirstName"`
// 	LastName             string            `json:"refUserLastName" gorm:"column:refUserLastName"`
// 	ProfileImg           string            `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
// 	DOB                  string            `json:"refUserDOB" gorm:"column:refUserDOB"`
// 	Gender               string            `json:"refUserGender" gorm:"column:refUserGender"`
// 	UserStatus           string            `json:"refUserStatus" gorm:"column:refUserStatus"`
// 	PhoneNumber          string            `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
// 	EmergencyPhoneNumber string            `json:"refCODOPhoneNo2" gorm:"column:refCODOPhoneNo2"`
// 	DoorNo               string            `json:"refCODODoorNo" gorm:"column:refCODODoorNo"`
// 	Street               string            `json:"refCODOStreet" gorm:"column:refCODOStreet"`
// 	District             string            `json:"refCODODistrict" gorm:"column:refCODOCity"`
// 	State                string            `json:"refCODOState" gorm:"column:refCODOState"`
// 	Country              string            `json:"refCODOCountry" gorm:"column:refCODOCountry"`
// 	Pincode              string            `json:"refCODOPincode" gorm:"column:refCODOPincode"`
// 	Email                string            `json:"refCODOEmail" gorm:"column:refCODOEmail"`
// 	Designation          string            `json:"refTDDesignation" gorm:"column:refTDDesignation"`
// 	Specialization       string            `json:"refTDSpecialization" gorm:"column:refTDSpecialization"`
// 	RoleName             string            `json:"refRTName" gorm:"column:refRTName"`
// 	StaffExperience      []StaffExperience `json:"staffExperience" gorm:"-"`
// }

// type StaffExperience struct {
// 	SEID           int    `json:"refSEId" gorm:"column:refSEId"`
// 	HospitalName   string `json:"refSEHospitalName" gorm:"column:refSEHospitalName"`
// 	Designation    string `json:"refSEDesignation" gorm:"column:refSEDesignation"`
// 	Specialization string `json:"refSESpecialization" gorm:"column:refSESpecialization"`
// 	Address        string `json:"refSEAddress" gorm:"column:refSEAddress"`
// 	FromDate       string `json:"refSEFromDate" gorm:"column:refSEFromDate"`
// 	ToDate         string `json:"refSEToDate" gorm:"column:refSEToDate"`
// }

// type GetTechnicianMapList struct {
// 	SCMID          int    `json:"refSCMId" gorm:"column:refSCMId"`
// 	SCID           int    `json:"refSCId" gorm:"column:refSCId"`
// 	RTId           int    `json:"refRTId" gorm:"column:refRTId"`
// 	Status         bool   `json:"refSCMStatus" gorm:"column:refSCMStatus"`
// 	ScanCenterName string `json:"refSCName" gorm:"column:refSCName"`
// 	RoleTpyeName   string `json:"refRTName" gorm:"column:refRTName"`
// }

// type GetTechnicianAudit struct {
// 	THID      int    `json:"refTHId" gorm:"column:refTHId"`
// 	TransId   int    `json:"transTypeId" gorm:"column:transTypeId"`
// 	THData    string `json:"refTHData" gorm:"column:refTHData"`
// 	THTime    string `json:"refTHTime" gorm:"column:refTHTime"`
// 	TransName string `json:"transTypeName" gorm:"column:transTypeName"`
// 	Username  string `json:"Username" gorm:"column:Username"`
// }

type LicenseFiles struct {
	FilesName   string `json:"file_name" binding:"required" mapstructure:"file_name"`
	OldFileName string `json:"old_file_name" binding:"required" mapstructure:"old_file_name"`
}

type RadiologistRegisterReq struct {
	FirstName                   string         `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName                    string         `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg                  string         `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                         string         `json:"dob" binding:"required" mapstructure:"dob"`
	Gender                      string         `json:"gender" binding:"required" mapstructure:"gender"`
	PhoneNoCountryCode          string         `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo                     string         `json:"phone" binding:"required" mapstructure:"phone"`
	Email                       string         `json:"email" binding:"required" mapstructure:"email"`
	EmergencyPhoneNoCountryCode string         `json:"emergency_phoneCountryCode" binding:"required" mapstructure:"emergency_phoneCountryCode"`
	EmergencyPhone              string         `json:"emergency_phone" binding:"required" mapstructure:"emergency_phone"`
	RegisterNumber              string         `json:"register_number" binding:"required" mapstructure:"register_number"`
	LicenseFiles                []LicenseFiles `json:"license_files" binding:"required" mapstructure:"license_files"`
}
