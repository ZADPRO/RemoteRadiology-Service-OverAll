package model

type GetAllPatientReq struct {
	SCId int `json:"SCId" mapstructure:"SCId"`
}

type GetAllPatient struct {
	RefUserId   uint   `json:"refUserId" gorm:"column:refUserId"`
	CustId      string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId      int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName   string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName    string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserStatus  string `json:"refUserStatus" gorm:"column:refUserStatus"`
	PhoneNumber string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email       string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
}

type GetPatientReq struct {
	SCId   int `json:"refSCId" mapstructure:"refSCId"`
	UserId int `json:"refUserId" mapstructure:"refUserId"`
}

type PatientAppointmentModel struct {
	AppointmentId   int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	AppointmentDate string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	RefSCId         int    `json:"refSCId" gorm:"column:refSCId"`
	RefSCCustId     string `json:"refSCCustId" gorm:"column:refSCCustId"`
}

type PatientOneModel struct {
	RefUserId              uint                      `json:"refUserId" gorm:"column:refUserId"`
	CustId                 string                    `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId                 int                       `json:"refRTId" gorm:"column:refRTId"`
	FirstName              string                    `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName               string                    `json:"refUserLastName" gorm:"column:refUserLastName"`
	ProfileImg             string                    `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB                    string                    `json:"refUserDOB" gorm:"column:refUserDOB"`
	Gender                 string                    `json:"refUserGender" gorm:"column:refUserGender"`
	UserStatus             bool                      `json:"refUserStatus" gorm:"column:refUserStatus"`
	UserAgreement          bool                      `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
	PhoneNumberCountryCode string                    `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	PhoneNumber            string                    `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email                  string                    `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	SCCustId               string                    `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCId                   int                       `json:"refSCId" gorm:"column:refSCId"`
	Appointments           []PatientAppointmentModel `json:"appointments" gorm:"-"`
	ProfileImgFile         *FileData                 `json:"profileImgFile" gorm:"-"`
}
