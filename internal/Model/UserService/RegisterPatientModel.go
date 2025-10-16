package model

type PatientCheckReq struct {
	PatientId string `json:"patientId" mapstructure:"patientId"`
	EmailId   string `json:"emailId" mapstructure:"emailId"`
	PhoneNo   string `json:"phoneNo" mapstructure:"phoneNo"`
}

type RegisterNewPatientReq struct {
	Firstname         string `json:"firstname" mapstructure:"firstname"`
	Profile_img       string `json:"profile_img" mapstructure:"profile_img"`
	DOB               string `json:"dob" mapstructure:"dob"`
	PhoneCountryCode  string `json:"phoneCountryCode" mapstructure:"phoneCountryCode"`
	PhoneNo           string `json:"phoneNo" mapstructure:"phoneNo"`
	EmailId           string `json:"emailId" mapstructure:"emailId"`
	Gender            string `json:"gender" mapstructure:"gender"`
	DateofAppointment string `json:"dateofAppointment" mapstructure:"dateofAppointment"`
	PatientId         string `json:"patientId" mapstructure:"patientId"`
	Mailoption        string `json:"mailoption" mapstructure:"mailoption"`
	SCId              int    `json:"refSCId" mapstructure:"refSCId"`
	SCCustId          string `json:"refSCustId" mapstructure:"refSCustId"`
}

type UpdatePatientReq struct {
	RefUserId       uint   `json:"refUserId" gorm:"column:refUserId" mapstructure:"refUserId"`
	RefUserCustId   string `json:"refUserCustId" gorm:"column:refUserCustId" mapstructure:"refUserCustId"`
	ActiveStatus    bool   `json:"refUserStatus" gorm:"column:refUserStatus" mapstructure:"refUserStatus"`
	FirstName       string `json:"refUserFirstName" gorm:"column:refUserFirstName" mapstructure:"refUserFirstName"`
	ProfileImg      string `json:"refUserProfileImg" gorm:"column:refUserProfileImg" mapstructure:"refUserProfileImg"`
	DOB             string `json:"refUserDOB" gorm:"column:refUserDOB" mapstructure:"refUserDOB"`
	Gender          string `json:"refUserGender" gorm:"column:refUserGender" mapstructure:"refUserGender"`
	PhoneNumberCode string `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode" mapstructure:"refCODOPhoneNo1CountryCode"`
	PhoneNumber     string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1" mapstructure:"refCODOPhoneNo1"`
	Email           string `json:"refCODOEmail" gorm:"column:refCODOEmail" mapstructure:"refCODOEmail"`
}

type CreateAppointmentPatientReq struct {
	RefUserId         int    `json:"refUserId" mapstructure:"refUserId"`
	Firstname         string `json:"firstname" mapstructure:"firstname"`
	Profile_img       string `json:"profile_img" mapstructure:"profile_img"`
	DOB               string `json:"dob" mapstructure:"dob"`
	PhoneCountryCode  string `json:"phoneCountryCode" mapstructure:"phoneCountryCode"`
	PhoneNo           string `json:"phoneNo" mapstructure:"phoneNo"`
	EmailId           string `json:"emailId" mapstructure:"emailId"`
	Gender            string `json:"gender" mapstructure:"gender"`
	DateofAppointment string `json:"dateofAppointment" mapstructure:"dateofAppointment"`
	PatientId         int    `json:"patientId" mapstructure:"patientId"`
	Mailoption        string `json:"mailoption" mapstructure:"mailoption"`
	SCId              int    `json:"refSCId" mapstructure:"refSCId"`
	SCCustId          string `json:"refSCustId" mapstructure:"refSCustId"`
}

type CancelResheduleAppointmentReq struct {
	AppointmentId   int    `json:"appointmentId" binding:"required" mapstructure:"appointmentId"`
	AppointmentDate string `json:"appointmentDate" mapstructure:"appointmentDate"`
	AccessMethod    string `json:"accessMethod" binding:"required" mapstructure:"accessMethod"`
}
