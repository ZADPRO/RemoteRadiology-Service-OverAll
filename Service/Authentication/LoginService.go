package service

import (
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	becrypt "AuthenticationService/internal/Helper/Becrypt"
	helper "AuthenticationService/internal/Helper/GenerateOTP"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	mailservice "AuthenticationService/internal/Helper/MailService"
	model "AuthenticationService/internal/Model/Authentication"
	query "AuthenticationService/query/Authentication"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func LoginServices(db *gorm.DB, reqVal model.LoginReq) model.LoginResponse {
	log := logger.InitLogger()
	var AdminLoginModel []model.AdminLoginModel

	// Execute the raw SQL query with the username (phone number)
	err := db.Raw(query.LoginAdminSQL, reqVal.Username).Scan(&AdminLoginModel).Error
	if err != nil {
		log.Error("LoginService DB Error: " + err.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	// Check if any user found
	if len(AdminLoginModel) == 0 {
		log.Warn("LoginService Invalid Credentials(u) for Username: " + reqVal.Username)
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid username or password",
		}
	}

	// Password verification
	user := AdminLoginModel[0]
	match := becrypt.ComparePasswords(user.ADHashPass, reqVal.Password)

	if !match {
		log.Warn("LoginService Invalid Credentials(p) for Username: " + reqVal.Username)
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid username or password",
		}
	}

	log.Info("LoginService Logined Successfully for Username: " + reqVal.Username)

	otp := helper.GenerateOTP()

	deleteOTP := db.Exec(query.DeleteOTPSQL, strconv.Itoa(user.UserId), 1).Error
	if deleteOTP != nil {
		log.Error("LoginService DB Error: " + deleteOTP.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	otperr := db.Exec(query.CreateOTPSQL, strconv.Itoa(user.UserId), otp, 1).Error
	if otperr != nil {
		log.Error("LoginService DB Error: " + otperr.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	htmlContent := mailservice.LoginOTPContent(otp)

	subject := "Your Login Passcode"

	emailStatus := mailservice.MailService(user.CODOEmail, htmlContent, subject)

	if !emailStatus {
		log.Error("Sending Mail Meets Error")
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	return model.LoginResponse{
		Status:   true,
		Message:  "Verification Code Send Successfully",
		RoleType: user.RId,
		Email:    user.CODOEmail,
		// Token:    accesstoken.CreateToken(user.UserId, user.RId, "0"),
	}
}

func VerifyOTPService(db *gorm.DB, reqVal model.VerifyReq) model.LoginResponse {
	log := logger.InitLogger()
	var AdminLoginModel []model.AdminLoginModel

	// Execute the raw SQL query with the username (phone number)
	err := db.Raw(query.LoginAdminSQL, reqVal.Username).Scan(&AdminLoginModel).Error
	if err != nil {
		log.Error("LoginService DB Error: " + err.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Internal server error",
		}
	}

	// Check if any user found
	if len(AdminLoginModel) == 0 {
		log.Warn("LoginService Invalid Credentials(u) for Username: " + reqVal.Username)
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid username or password",
		}
	}

	// Password verification
	user := AdminLoginModel[0]
	match := becrypt.ComparePasswords(user.ADHashPass, reqVal.Password)

	if !match {
		log.Warn("LoginService Invalid Credentials(p) for Username: " + reqVal.Username)
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid username or password",
		}
	}

	var VerifyOTP []model.VerifyOTP

	otpverification := db.Raw(query.VerifyOTPSQL, strconv.Itoa(user.UserId), reqVal.OTP, 1).Scan(&VerifyOTP).Error
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

	deleteOTP := db.Exec(query.DeleteOTPSQL, strconv.Itoa(user.UserId), 1).Error
	if deleteOTP != nil {
		log.Error("LoginService DB Error: " + deleteOTP.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	history := model.RefTransHistory{
		TransTypeId: 1,
		THData:      "Logged In Successfully",
		UserId:      user.UserId,
		THActionBy:  user.UserId,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	return model.LoginResponse{
		Status:         true,
		Message:        "Logedin Successfully",
		RoleType:       user.RId,
		PasswordStatus: user.AHPassChangeStatus,
		Token:          accesstoken.CreateToken(user.UserId, user.RId),
	}

}

func UserChangePasswordService(db *gorm.DB, reqVal model.UserChnagePasswordReq, userId any) model.LoginResponse {
	log := logger.InitLogger()

	fmt.Print(reqVal.Password)

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.LoginResponse{
			Status:  false,
			Message: "Something Went Wrong, Try Again",
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	hashPassword, hashPassworderr := becrypt.HashPassword(reqVal.Password)

	if hashPassworderr != nil {
		log.Printf("ERROR: Failed to hash password: %v\n", hashPassworderr)
		return model.LoginResponse{
			Status:  false,
			Message: "Something Went Wrong, Try Again",
		}
	}

	changePassword := db.Exec(query.UpdateUserDataSQL, true, userId).Error
	if changePassword != nil {
		log.Error("UserChangePasswordService DB Error: " + changePassword.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something Went Wrong, Try Again",
		}
	}

	chnagePassworderr := db.Exec(query.UpdatePasswordSQL, hashdb.Encrypt(reqVal.Password), hashPassword, false, userId).Error
	if chnagePassworderr != nil {
		log.Error("UserChangePasswordService DB Error: " + chnagePassworderr.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something Went Wrong, Try Again",
		}
	}

	floatId, ok := userId.(float64)
	if !ok {
		log.Error("UserChangePasswordService DB Error: userId is not a float64")
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid user ID format",
		}
	}
	uid := int(floatId)

	history := model.RefTransHistory{
		TransTypeId: 9,
		THData:      "Password Changed Successfully and Accept the Agreement",
		UserId:      uid,
		THActionBy:  uid,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	return model.LoginResponse{
		Status:  true,
		Message: "Password Changed Successfully",
	}

}
