package service

import (
	becrypt "AuthenticationService/internal/Helper/Becrypt"
	helper "AuthenticationService/internal/Helper/GenerateOTP"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	mailservice "AuthenticationService/internal/Helper/MailService"
	model "AuthenticationService/internal/Model/UserService"
	query "AuthenticationService/query/UserService"
	"strconv"

	"gorm.io/gorm"
)

func PostGetOtpPatientService(db *gorm.DB, reqVal model.GetOtpPatient) model.RegisterPatientRes {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var verifyData []model.VerifyData

	verifyDataerr := db.Raw(query.RegisterUserVerifyData, reqVal.PhoneNo, reqVal.PhoneNo, reqVal.Email).Scan(&verifyData).Error
	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if len(verifyData) > 0 {

		if verifyData[0].Email == reqVal.Email {
			return model.RegisterPatientRes{
				Status:  false,
				Message: "Email Already Exists",
			}
		} else {
			return model.RegisterPatientRes{
				Status:  false,
				Message: "Mobile Number Already Exists",
			}
		}

	}

	otp := helper.GenerateOTP()

	deleteOTP := db.Exec(query.DeleteOTPSQL, reqVal.Email, 3).Error
	if deleteOTP != nil {
		log.Error("LoginService DB Error: " + deleteOTP.Error())
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	otperr := db.Exec(query.CreateOTPSQL, reqVal.Email, otp, 3).Error
	if otperr != nil {
		log.Error("LoginService DB Error: " + otperr.Error())
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	userName := reqVal.FirstName + " " + reqVal.LastName

	subject := "Verify OTP"

	htmlContent := mailservice.GetOTPMailContent(userName, otp)

	emailStatus := mailservice.MailService(reqVal.Email, htmlContent, subject)

	if !emailStatus {
		log.Error("Sending Mail Meets Error")
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	return model.RegisterPatientRes{
		Status:  true,
		Message: "Succcessfully OTP sended",
	}

}

func PostCheckOTPPatientService(db *gorm.DB, reqVal model.VerifyOtpPatient) model.RegisterPatientRes {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var VerifyOTP []model.VerifyOTP

	otpverification := db.Raw(query.VerifyOTPSQL, reqVal.Email, reqVal.OTP, 3).Scan(&VerifyOTP).Error
	if otpverification != nil {
		log.Error("LoginService DB Error: " + otpverification.Error())
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if !VerifyOTP[0].Result {
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Invalid OTP",
		}
	}

	return model.RegisterPatientRes{
		Status:  true,
		Message: "OTP Successfully Verified",
	}

}

func PostRegisterPatientService(db *gorm.DB, reqVal model.RegisterPatientReq) model.RegisterPatientRes {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var VerifyOTP []model.VerifyOTP

	otpverification := db.Raw(query.VerifyOTPSQL, reqVal.Email, reqVal.OTP, 3).Scan(&VerifyOTP).Error
	if otpverification != nil {
		log.Error("LoginService DB Error: " + otpverification.Error())
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if !VerifyOTP[0].Result {
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Invalid OTP",
		}
	}

	var verifyData []model.VerifyData

	verifyDataerr := db.Raw(query.RegisterUserVerifyData, reqVal.PhoneNo, reqVal.PhoneNo, reqVal.Email).Scan(&verifyData).Error
	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if len(verifyData) > 0 {

		if verifyData[0].Email == reqVal.Email {
			return model.RegisterPatientRes{
				Status:  false,
				Message: "Email Already Exists",
			}
		} else {
			return model.RegisterPatientRes{
				Status:  false,
				Message: "Mobile Number Already Exists",
			}
		}

	}

	var TotalCount []model.TotalCount

	err := db.Raw(query.GetUsersCountSQL, 4).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	CustId := "P" + strconv.Itoa(TotalCount[0].TotalCount+100000)

	PatientData := model.CreateTechnicianModel{
		UserCustId:      CustId,
		RoleId:          4,
		FirstName:       hashdb.Encrypt(reqVal.FirstName),
		LastName:        hashdb.Encrypt(reqVal.LastName),
		Status:          true,
		AgreementStatus: true,
	}

	PatientDataErr := tx.Create(&PatientData).Error
	if PatientDataErr != nil {
		log.Printf("ERROR: Failed to create Patient User Data: %v\n", PatientDataErr)
		tx.Rollback()
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	CommunicationData := model.CreateTechnicianCommunicationModel{
		UserId:             int(PatientData.UserId),
		PhoneNoCountryCode: reqVal.PhoneNoCountryCode,
		PhoneNo:            reqVal.PhoneNo,
		Email:              reqVal.Email,
	}

	CommunicationDataerr := tx.Create(&CommunicationData).Error
	if CommunicationDataerr != nil {
		log.Printf("ERROR: Failed to create Technician Communication Data: %v\n", CommunicationDataerr)
		tx.Rollback()
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	hashPassword, hashPassworderr := becrypt.HashPassword(reqVal.Password)
	if hashPassworderr != nil {
		log.Printf("ERROR: Failed to create Technician Domain Data: %v\n", hashPassworderr)
		tx.Rollback()
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	AuthData := model.CreateTechnicianAuthModel{
		UserId:         int(PatientData.UserId),
		Password:       hashdb.Encrypt(reqVal.Password),
		HashPassword:   hashPassword,
		PasswordStatus: false,
	}

	AuthDataerr := tx.Create(&AuthData).Error
	if AuthDataerr != nil {
		log.Printf("ERROR: Failed to create Technician Auth Data: %v\n", AuthDataerr)
		tx.Rollback()
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	history := model.RefTransHistory{
		TransTypeId: 10,
		THData:      "Account Created Successfully",
		UserId:      PatientData.UserId,
		THActionBy:  PatientData.UserId,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	userName := reqVal.FirstName + " " + reqVal.LastName

	subject := "Welcome to Wellthgreen HealthCare Pvt Ltd"

	htmlContent := mailservice.RegisterationMailContent(userName)

	emailStatus := mailservice.MailService(reqVal.Email, htmlContent, subject)

	if !emailStatus {
		log.Error("Sending Mail Meets Error")
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return model.RegisterPatientRes{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	return model.RegisterPatientRes{
		Status:  true,
		Message: "Succcessfully Account Created",
	}
}
