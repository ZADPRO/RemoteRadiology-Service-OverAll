package model

type ForgetPasswordReq struct {
	Email string `json:"email" binding:"required"`
}

type VerifyPasswordReq struct {
	Email string `json:"email" binding:"required"`
	OTP   int    `json:"otp" binding:"required"`
}

type ChangePasswordReq struct {
	Email    string `json:"email" binding:"required"`
	OTP      int    `json:"otp" binding:"required"`
	Password string `json:"password" binding:"required"`
}
