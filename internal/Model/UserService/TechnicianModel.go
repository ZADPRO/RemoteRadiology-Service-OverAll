package model

type TechnicianRegisterReq struct {
	FirstName          string         `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName           string         `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg         string         `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                string         `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode string         `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo            string         `json:"phone" binding:"required" mapstructure:"phone"`
	Email              string         `json:"email" binding:"required" mapstructure:"email"`
	SSNo               string         `json:"social_security_no" binding:"required" mapstructure:"social_security_no"`
	DriversLicenseNo   string         `json:"drivers_license" mapstructure:"drivers_license"`
	TrainedEaseQT      bool           `json:"trained_ease_qt" binding:"required" mapstructure:"trained_ease_qt"`
	DigitalSignature   string         `json:"digital_signature" mapstructure:"digital_signature"`
	LicenseFiles       []LicenseFiles `json:"license_files" mapstructure:"license_files"`
	ScanCenterId       int            `json:"scan_center_id" binding:"required" mapstructure:"scan_center_id"`
}

type VerifyData struct {
	UserId       int    `json:"refUserId" gorm:"column:refUserId"`
	UserCustId   string `json:"refUserCustId" gorm:"column:refUserCustId"`
	PhoneNumber1 string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	PhoneNumber2 string `json:"refCODOPhoneNo2" gorm:"column:refCODOPhoneNo2"`
	Email        string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
}

type ScanCenterVerifyData struct {
	SCId       int    `json:"refSCId" gorm:"column:refSCId"`
	SCCustId   string `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCPhoneNo1 string `json:"refSCPhoneNo1" gorm:"column:refSCPhoneNo1"`
	SCEmail    string `json:"refSCEmail" gorm:"column:refSCEmail"`
}

// type Expreience struct {
// 	FromDate       string `json:"from_date" binding:"required"`
// 	ToDate         string `json:"to_date" binding:"required"`
// 	HospitalName   string `json:"hospital_name" binding:"required"`
// 	Designation    string `json:"designation" binding:"required"`
// 	Specialization string `json:"specialization" binding:"required"`
// 	Address        string `json:"address" binding:"required"`
// }

type TotalCount struct {
	TotalCount int    `json:"TotalCount" gorm:"column:TotalCount"`
	SCCustId   string `json:"refSCCustId" gorm:"column:refSCCustId"`
}

type TotalCountModel struct {
	TotalCount int `json:"TotalCount" gorm:"column:TotalCount"`
}

type TechnicianMapsCreate struct {
	ScanCenterId int  `json:"refSCId" binding:"required"  mapstructure:"refSCId"`
	RoleId       int  `json:"refRTId" binding:"required"  mapstructure:"refRTId"`
	Status       bool `json:"status"  mapstructure:"status"`
}

type CreateTechnicianReq struct {
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
	SocialSecurityCenter        string         `json:"social_security_center" binding:"required" mapstructure:"social_security_center"`
	NPI                         string         `json:"npi" binding:"required" mapstructure:"npi"`
	LicenseFiles                []LicenseFiles `json:"license_files" binding:"required" mapstructure:"license_files"`
}

type CreateTechnicianModel struct {
	UserId          int    `json:"refUserId" gorm:"primaryKey;autoIncrement;column:refUserId"`
	UserCustId      string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId          int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName       string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName        string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserProfileImg  string `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB             string `json:"refUserDOB" gorm:"column:refUserDOB"`
	Gender          string `json:"refUserGender" gorm:"column:refUserGender"`
	Status          bool   `json:"refUserStatus" gorm:"column:refUserStatus"`
	AgreementStatus bool   `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
}

func (CreateTechnicianModel) TableName() string {
	return "Users"
}

type CreateTechnicianCommunicationModel struct {
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

func (CreateTechnicianCommunicationModel) TableName() string {
	return "userdomain.refCommunicationDomain"
}

type CreateTechnicianAuthModel struct {
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	Password       string `json:"refADPassword" gorm:"column:refADPassword"`
	HashPassword   string `json:"refADHashPass" gorm:"column:refADHashPass"`
	PasswordStatus bool   `json:"refAHPassChangeStatus" gorm:"column:refAHPassChangeStatus"`
}

func (CreateTechnicianAuthModel) TableName() string {
	return "userdomain.refAuthDomain"
}

type CreateStaffExprienceModel struct {
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	HospitalName   string `json:"refSEHospitalName" gorm:"column:refSEHospitalName"`
	Designation    string `json:"refSEDesignation" gorm:"column:refSEDesignation"`
	Specialization string `json:"refSESpecialization" gorm:"column:refSESpecialization"`
	Address        string `json:"refSEAddress" gorm:"column:refSEAddress"`
	FromDate       string `json:"refSEFrom" gorm:"column:refSEFrom"`
	ToDate         string `json:"refSETo" gorm:"column:refSETo"`
}

func (CreateStaffExprienceModel) TableName() string {
	return "userdomain.refStaffExprience"
}

type UpdateExpreience struct {
	ID             int    `json:"id" binding:"required" mapstructure:"id"`
	FromDate       string `json:"from_date" binding:"required" mapstructure:"from_date"`
	ToDate         string `json:"to_date" binding:"required" mapstructure:"to_date"`
	HospitalName   string `json:"hospital_name" binding:"required" mapstructure:"hospital_name"`
	Designation    string `json:"designation" binding:"required" mapstructure:"designation"`
	Specialization string `json:"specialization" binding:"required" mapstructure:"specialization"`
	Address        string `json:"address" binding:"required" mapstructure:"address"`
	Status         string `json:"status" binding:"required" mapstructure:"status"`
}

type StaffTransactionOne struct {
	UserId         int    `json:"refUserId" gorm:"column:refUserId"`
	HospitalName   string `json:"refSEHospitalName" gorm:"column:refSEHospitalName"`
	Designation    string `json:"refSEDesignation" gorm:"column:refSEDesignation"`
	Specialization string `json:"refSESpecialization" gorm:"column:refSESpecialization"`
	Address        string `json:"refSEAddress" gorm:"column:refSEAddress"`
	FromDate       string `json:"refSEFrom" gorm:"column:refSEFrom"`
	ToDate         string `json:"refSETo" gorm:"column:refSETo"`
}

type TechnicianMaps struct {
	UserId       int  `json:"refUserId" binding:"required"  mapstructure:"refUserId"`
	ScanCenterId int  `json:"refSCId" binding:"required"  mapstructure:"refSCId"`
	RoleId       int  `json:"refRTId" binding:"required"  mapstructure:"refRTId"`
	Status       bool `json:"status"  mapstructure:"status"`
	UserRoleId   int  `json:"userRoleId" binding:"required"  mapstructure:"userRoleId"`
}

type TechnicianMapReq struct {
	TechnicianMaps []TechnicianMaps `json:"technicianMaps" binding:"required" mapstructure:"technicianMaps"`
}

type CreateTechnicianMap struct {
	UserId       int  `json:"refUserId" gorm:"column:refUserId"`
	ScanCenterId int  `json:"refSCId" gorm:"column:refSCId"`
	RoleId       int  `json:"refRTId" gorm:"column:refRTId"`
	Status       bool `json:"refSCMStatus" gorm:"column:refSCMStatus"`
}

func (CreateTechnicianMap) TableName() string {
	return "map.refScanCenterMap"
}

type GetMapData struct {
	MapId        int  `json:"refSCMId" gorm:"primaryKey;autoIncrement;column:refSCMId"`
	UserId       int  `json:"refUserId" gorm:"column:refUserId"`
	ScanCenterId int  `json:"refSCId" gorm:"column:refSCId"`
	RoleId       int  `json:"refRTId" gorm:"column:refRTId"`
	Status       bool `json:"refSCMStatus" gorm:"column:refSCMStatus"`
}

type GetScanCenterName struct {
	ScanCenterName string `json:"refSCName" gorm:"column:refSCName"`
}

type RefTransHistory struct {
	TransTypeId int    `json:"transTypeId" gorm:"column:transTypeId"`
	THData      string `json:"refTHData" gorm:"column:refTHData"`
	UserId      int    `json:"refUserId" gorm:"column:refUserId"`
	THActionBy  int    `json:"refTHActionBy" gorm:"column:refTHActionBy"`
}

func (RefTransHistory) TableName() string {
	return "aduit.refTransHistory"
}

type CreateAppointmentModel struct {
	AppointmentId   int    `json:"refAppointmentId" gorm:"primaryKey;autoIncrement;column:refAppointmentId"`
	UserId          int    `json:"refUserId" gorm:"column:refUserId"`
	SCId            int    `json:"refSCId" gorm:"column:refSCId"`
	AppointmentDate string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	// AppointmentStartTime string `json:"refAppointmentStartTime" gorm:"column:refAppointmentStartTime"`
	// AppointmentEndTime   string `json:"refAppointmentEndTime" gorm:"column:refAppointmentEndTime"`
	AppointmentUrgency  bool   `json:"refAppointmentUrgency" gorm:"column:refAppointmentUrgency"`
	AppointmentStatus   bool   `json:"refAppointmentStatus" gorm:"column:refAppointmentStatus"`
	AppointmentComplete string `json:"refAppointmentComplete" gorm:"column:refAppointmentComplete"`
}

func (CreateAppointmentModel) TableName() string {
	return "appointment.refAppointments"
}

type CreateTechnicianDomainModel struct {
	TDId                  int    `json:"refTDId" gorm:"primaryKey;autoIncrement;column:refTDId"`
	UserId                int    `json:"refUserId" gorm:"column:refUserId"`
	TDTrainedEaseQTStatus bool   `json:"refTDTrainedEaseQTStatus" gorm:"column:refTDTrainedEaseQTStatus"`
	TDSSNo                string `json:"refTDSSNo" gorm:"column:refTDSSNo"`
	TDDrivingLicense      string `json:"refTDDrivingLicense" gorm:"column:refTDDrivingLicense"`
	TDDigitalSignature    string `json:"refTDDigitalSignature" gorm:"column:refTDDigitalSignature"`
}

func (CreateTechnicianDomainModel) TableName() string {
	return "userdomain.refTechnicianDomain"
}

type UpdateTechnicianReq struct {
	ID                 int                  `json:"id" binding:"required" mapstructure:"id"`
	FirstName          string               `json:"firstname" binding:"required" mapstructure:"firstname"`
	LastName           string               `json:"lastname" binding:"required" mapstructure:"lastname"`
	ProfileImg         string               `json:"profile_img" binding:"required" mapstructure:"profile_img"`
	DOB                string               `json:"dob" binding:"required" mapstructure:"dob"`
	PhoneNoCountryCode string               `json:"phoneCountryCode" binding:"required" mapstructure:"phoneCountryCode"`
	PhoneNo            string               `json:"phone" binding:"required" mapstructure:"phone"`
	Email              string               `json:"email" binding:"required" mapstructure:"email"`
	SSNo               string               `json:"social_security_no" binding:"required" mapstructure:"social_security_no"`
	DriversLicenseNo   string               `json:"drivers_license" mapstructure:"drivers_license"`
	TrainedEaseQT      bool                 `json:"trained_ease_qt" binding:"required" mapstructure:"trained_ease_qt"`
	DigitalSignature   string               `json:"digital_signature" mapstructure:"digital_signature"`
	LicenseFiles       []UpdateLicenseFiles `json:"license_files" mapstructure:"license_files"`
	Status             bool                 `json:"status" binding:"required" mapstructure:"status"`
}

type GetAllTechnicianData struct {
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
	TrainedEaseQT      bool   `json:"refTDTrainedEaseQTStatus" gorm:"column:refTDTrainedEaseQTStatus"`
	SSNo               string `json:"refTDSSNo" gorm:"column:refTDSSNo"`
	DrivingLicense     string `json:"refTDDrivingLicense" gorm:"column:refTDDrivingLicense"`
	DigitalSignature   string `json:"refTDDigitalSignature" gorm:"column:refTDDigitalSignature"`
}
