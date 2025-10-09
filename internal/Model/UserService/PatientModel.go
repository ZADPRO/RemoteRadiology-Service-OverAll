package model

type RegisterPatientReq struct {
	FirstName          string `json:"firstname" binding:"required"  mapstructure:"firstname"`
	LastName           string `json:"lastname"  mapstructure:"lastname"`
	Email              string `json:"email" binding:"required"  mapstructure:"email"`
	PhoneNoCountryCode string `json:"phoneCountryCode" binding:"required"  mapstructure:"phoneCountryCode"`
	PhoneNo            string `json:"phone" binding:"required"  mapstructure:"phone"`
	Password           string `json:"password" binding:"required"  mapstructure:"password"`
	OTP                int    `json:"otp" binding:"required"  mapstructure:"otp"`
}

type RegisterPatientRes struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type RegisterPatientUserModel struct {
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

type GetOtpPatient struct {
	FirstName          string `json:"firstname"`
	LastName           string `json:"lastname"`
	Email              string `json:"email"`
	PhoneNoCountryCode string `json:"phoneCountryCode"`
	PhoneNo            string `json:"phone"`
}

type VerifyOtpPatient struct {
	Email string `json:"email"`
	OTP   int    `json:"otp"`
}

type VerifyOTP struct {
	Result bool `json:"result" gorm:"column:result"`
}

