package model

import "time"

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Userdata struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

type LoginResponse struct {
	Status         bool   `json:"status"`
	Message        string `json:"message"`
	RoleType       int    `json:"roleType"`
	Token          string `json:"token"`
	Email          string `json:"email"`
	PasswordStatus bool   `json:"passwordStatus"`
}

type AdminLoginModel struct {
	UserId             int    `json:"refUserId" gorm:"column:refUserId"`
	CustUserId         string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RId                int    `json:"refRTId" gorm:"column:refRTId"`
	UserFirstName      string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	UserLastName       string `json:"refUserLastName" gorm:"column:refUserLastName"`
	RTName             string `json:"refRTName" gorm:"column:refRTName"`
	ADHashPass         string `json:"refADHashPass" gorm:"column:refADHashPass"`
	ADPassword         string `json:"refADPassword" gorm:"column:refADPassword"`
	AHPassChangeStatus bool   `json:"refAHPassChangeStatus" gorm:"column:refAHPassChangeStatus"`
	CODOPhoneNo1       string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	CODOEmail          string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
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

type VerifyReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	OTP      int    `json:"otp" binding:"required"`
}

type VerifyOTP struct {
	Result bool `json:"result" gorm:"column:result"`
}

type UserChnagePasswordReq struct {
	Password string `json:"password" binding:"required" mapstructure:"password"`
}
