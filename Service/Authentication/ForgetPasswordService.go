package service

import (
	becrypt "AuthenticationService/internal/Helper/Becrypt"
	helper "AuthenticationService/internal/Helper/GenerateOTP"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	mailservice "AuthenticationService/internal/Helper/MailService"
	model "AuthenticationService/internal/Model/Authentication"
	query "AuthenticationService/query/Authentication"
	"strconv"

	"gorm.io/gorm"
)

func ForgetPasswordService(db *gorm.DB, reqVal model.ForgetPasswordReq) model.LoginResponse {

	log := logger.InitLogger()

	var AdminLoginModel []model.AdminLoginModel

	// Execute the raw SQL query with the username (phone number)
	err := db.Raw(query.LoginAdminSQL, reqVal.Email).Scan(&AdminLoginModel).Error
	if err != nil {
		log.Error("LoginService DB Error: " + err.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if len(AdminLoginModel) == 1 {

		user := AdminLoginModel[0]
		log.Warn("LoginService Invalid Credentials(u) for Username: " + reqVal.Email)
		deleteOTP := db.Exec(query.DeleteOTPSQL, strconv.Itoa(user.UserId), 2).Error
		if deleteOTP != nil {
			log.Error("LoginService DB Error: " + deleteOTP.Error())
			return model.LoginResponse{
				Status:  false,
				Message: "Something went wrong, Try Again",
			}
		}

		otp := helper.GenerateOTP()

		otperr := db.Exec(query.CreateOTPSQL, strconv.Itoa(user.UserId), otp, 2).Error
		if otperr != nil {
			log.Error("LoginService DB Error: " + otperr.Error())
			return model.LoginResponse{
				Status:  false,
				Message: "Something went wrong, Try Again",
			}
		}

		htmlContent := mailservice.ForgetPasswordOTPContent(otp)

		subject := "Your Forget Password Passcode"

		emailStatus := mailservice.MailService(user.CODOEmail, htmlContent, subject)

		if !emailStatus {
			log.Error("Sending Mail Meets Error")
			return model.LoginResponse{
				Status:  false,
				Message: "Something went wrong, Try Again",
			}
		}

		return model.LoginResponse{
			Status:  true,
			Message: "Verification Code Send Successfully",
			// RoleType: user.RId,
			// Token:    accesstoken.CreateToken(user.UserId, user.RId, "0"),
		}

	} else {
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid Email",
		}
	}
}

func VerifyForgetPasswordOTPService(db *gorm.DB, reqVal model.VerifyPasswordReq) model.LoginResponse {
	log := logger.InitLogger()

	var AdminLoginModel []model.AdminLoginModel

	// Execute the raw SQL query with the username (phone number)
	err := db.Raw(query.LoginAdminSQL, reqVal.Email).Scan(&AdminLoginModel).Error
	if err != nil {
		log.Error("LoginService DB Error: " + err.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	user := AdminLoginModel[0]

	if len(AdminLoginModel) == 1 {

		var VerifyOTP []model.VerifyOTP

		otpverification := db.Raw(query.VerifyOTPSQL, strconv.Itoa(user.UserId), reqVal.OTP, 2).Scan(&VerifyOTP).Error
		if otpverification != nil {
			log.Error("LoginService DB Error: " + otpverification.Error())
			return model.LoginResponse{
				Status:  false,
				Message: "Something went wrong, Try Again",
			}
		}

		if !VerifyOTP[0].Result {
			return model.LoginResponse{
				Status:  false,
				Message: "Invalid Passcode",
			}
		}

		return model.LoginResponse{
			Status:  true,
			Message: "Passcode Verified Successfully",
		}

	} else {
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid Email",
		}
	}

}

func ChangePasswordService(db *gorm.DB, reqVal model.ChangePasswordReq) model.LoginResponse {
	log := logger.InitLogger()

	var AdminLoginModel []model.AdminLoginModel

	// Execute the raw SQL query with the username (phone number)
	err := db.Raw(query.LoginAdminSQL, reqVal.Email).Scan(&AdminLoginModel).Error
	if err != nil {
		log.Error("LoginService DB Error: " + err.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	user := AdminLoginModel[0]

	if len(AdminLoginModel) == 1 {

		var VerifyOTP []model.VerifyOTP

		otpverification := db.Raw(query.VerifyOTPSQL, strconv.Itoa(user.UserId), reqVal.OTP, 2).Scan(&VerifyOTP).Error
		if otpverification != nil {
			log.Error("LoginService DB Error: " + otpverification.Error())
			return model.LoginResponse{
				Status:  false,
				Message: "Something went wrong, Try Again",
			}
		}

		if !VerifyOTP[0].Result {
			return model.LoginResponse{
				Status:  false,
				Message: "Invalid Passcode",
			}
		}

		deleteOTP := db.Exec(query.DeleteOTPSQL, strconv.Itoa(user.UserId), 2).Error
		if deleteOTP != nil {
			log.Error("LoginService DB Error: " + deleteOTP.Error())
			return model.LoginResponse{
				Status:  false,
				Message: "Something went wrong, Try Again",
			}
		}

		hashPassword, hashPassworderr := becrypt.HashPassword(reqVal.Password)

		if hashPassworderr != nil {
			log.Printf("ERROR: Failed to hash password: %v\n", hashPassworderr)
			return model.LoginResponse{
				Status:  false,
				Message: "Invalid Passcode",
			}
		}

		chnagePassworderr := db.Exec(query.UpdatePasswordSQL, hashdb.Encrypt(reqVal.Password), hashPassword, false, user.UserId).Error
		if chnagePassworderr != nil {
			log.Error("LoginService DB Error: " + chnagePassworderr.Error())
			return model.LoginResponse{
				Status:  false,
				Message: "Invalid Passcode",
			}
		}

		history := model.RefTransHistory{
			TransTypeId: 9,
			THData:      "Password Changed Successfully",
			UserId:      user.UserId,
			THActionBy:  user.UserId,
		}

		errhistory := db.Create(&history).Error
		if errhistory != nil {
			log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
			return model.LoginResponse{
				Status:  false,
				Message: "Internal server error",
			}
		}

		return model.LoginResponse{
			Status:  true,
			Message: "Password Changed Successfully",
		}

	} else {
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid Email",
		}
	}
}
