package service

import (
	becrypt "AuthenticationService/internal/Helper/Becrypt"
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	mailservice "AuthenticationService/internal/Helper/MailService"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	model "AuthenticationService/internal/Model/UserService"
	query "AuthenticationService/query/UserService"
	"encoding/json"

	"gorm.io/gorm"
)

func PostCheckPatientService(db *gorm.DB, reqVal model.PatientCheckReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var verifyData []model.VerifyData

	verifyDataerr := db.Raw(
		query.CheckpatientExits,
		reqVal.EmailId,
		reqVal.PhoneNo,
		reqVal.PatientId,
	).Scan(&verifyData).Error

	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return false, "Something went wrong, Try Again"
	}

	if len(verifyData) > 0 {
		if verifyData[0].Email == reqVal.EmailId {
			return false, "Email Already Exists"
		} else if verifyData[0].PhoneNumber1 == reqVal.PhoneNo {
			return false, "Mobile Number Already Exists"
		} else {
			return false, "Patient ID Already Exists"
		}
	}

	return true, "Succcessfully Checked"
}

func PostPatientService(db *gorm.DB, reqVal model.RegisterNewPatientReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	PatientData := model.CreateTechnicianModel{
		UserCustId:      reqVal.PatientId,
		RoleId:          4,
		FirstName:       hashdb.Encrypt(reqVal.Firstname),
		LastName:        hashdb.Encrypt(""),
		Gender:          reqVal.Gender,
		DOB:             hashdb.Encrypt(reqVal.DOB),
		UserProfileImg:  hashdb.Encrypt(reqVal.Profile_img),
		Status:          true,
		AgreementStatus: true,
	}

	PatientDataErr := tx.Create(&PatientData).Error
	if PatientDataErr != nil {
		log.Printf("ERROR: Failed to create Patient User Data: %v\n", PatientDataErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	CommunicationData := model.CreateTechnicianCommunicationModel{
		UserId:             int(PatientData.UserId),
		PhoneNoCountryCode: reqVal.PhoneCountryCode,
		PhoneNo:            reqVal.PhoneNo,
		Email:              reqVal.EmailId,
	}

	CommunicationDataerr := tx.Create(&CommunicationData).Error
	if CommunicationDataerr != nil {
		log.Printf("ERROR: Failed to create Technician Communication Data: %v\n", CommunicationDataerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	hashPassword, hashPassworderr := becrypt.HashPassword(reqVal.DOB)
	if hashPassworderr != nil {
		log.Printf("ERROR: Failed to create Technician Domain Data: %v\n", hashPassworderr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	AuthData := model.CreateTechnicianAuthModel{
		UserId:         int(PatientData.UserId),
		Password:       hashdb.Encrypt(reqVal.DOB),
		HashPassword:   hashPassword,
		PasswordStatus: true,
	}

	AuthDataerr := tx.Create(&AuthData).Error
	if AuthDataerr != nil {
		log.Printf("ERROR: Failed to create Technician Auth Data: %v\n", AuthDataerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
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
		return false, "Something went wrong, Try Again"
	}

	MapScanCenter := model.MapScanCenterPatientModel{
		UserId: int(PatientData.UserId),
		SCId:   int(reqVal.SCId),
	}

	MapScanCentererr := tx.Create(&MapScanCenter).Error
	if MapScanCentererr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain: %v\n", MapScanCentererr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	historyErr := model.RefTransHistory{
		TransTypeId: 4,
		THData:      "Account Created Successfully",
		UserId:      PatientData.UserId,
		THActionBy:  idValue,
	}

	errhistoryERr := tx.Create(&historyErr).Error
	if errhistoryERr != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistoryERr.Error())
		return false, "Something went wrong, Try Again"
	}

	Appointment := model.CreateAppointmentModel{
		UserId:          PatientData.UserId,
		SCId:            reqVal.SCId,
		AppointmentDate: reqVal.DateofAppointment,
		// AppointmentStartTime: reqVal.AppointmentStartTime,
		// AppointmentEndTime:   reqVal.AppointmentEndTime,
		// AppointmentUrgency: reqVal.AppointmentUrgency,
		AppointmentStatus:   true,
		AppointmentComplete: "fillform",
	}

	Appointmenterr := db.Create(&Appointment).Error
	if Appointmenterr != nil {
		log.Printf("ERROR: Failed to create Appointment: %v\n", Appointmenterr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	InsertReportQuestionsErr := tx.Exec(query.InsertReportIntakeAllSQL, int(PatientData.UserId), int(Appointment.AppointmentId), timeZone.GetPacificTime(), int(idValue)).Error
	if InsertReportQuestionsErr != nil {
		log.Printf("ERROR: Failed to Report Question: %v\n", InsertReportQuestionsErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	InsertReportTextContentErr := tx.Exec(
		query.InsertNewReportTextContentSQL,
		int(PatientData.UserId),
		int(Appointment.AppointmentId),
		timeZone.GetPacificTime(),
		int(idValue),
	).Error

	if InsertReportTextContentErr != nil {
		log.Printf("ERROR: Failed to Report Text Content: %v\n", InsertReportTextContentErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if reqVal.Mailoption == "sendbywellthgreen" {

		htmlContent := mailservice.RegistrationMailContent(reqVal.Firstname, reqVal.PatientId, reqVal.EmailId, reqVal.DOB, "Patient")

		subject := "Welcome – Your Appointment at " + reqVal.SCCustId + " Scan Centre"

		emailStatus := mailservice.MailService(reqVal.EmailId, htmlContent, subject)

		if !emailStatus {
			log.Error("Sending Mail Meets Error")
			return false, "Something went wrong, Try Again"
		}

	}
	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Account Created"
}

func PatchPatientService(db *gorm.DB, reqVal model.UpdatePatientReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var verifyData []model.VerifyData

	verifyDataerr := db.Raw(
		query.VerifyUserDataSQL,
		reqVal.PhoneNumber,
		reqVal.Email,
		reqVal.RefUserCustId,
	).Scan(&verifyData).Error

	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return false, "Something went wrong, Try Again"
	}

	if len(verifyData) > 0 {
		if verifyData[0].Email == reqVal.Email && verifyData[0].UserId != int(reqVal.RefUserId) {
			return false, "Email Already Exists"
		} else if verifyData[0].PhoneNumber1 == reqVal.PhoneNumber && verifyData[0].UserId != int(reqVal.RefUserId) {
			return false, "Mobile Number Already Exists"
		} else if verifyData[0].UserCustId == reqVal.RefUserCustId && verifyData[0].UserId != int(reqVal.RefUserId) {
			return false, "Patient ID Already Exists"
		}
	}

	PreviousData := model.UpdatePatientReq{}

	errPrev := tx.Raw(query.GetAllPatientDataQuery, reqVal.RefUserId).Scan(&PreviousData).Error
	if errPrev != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", errPrev)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	oldData := map[string]interface{}{
		"Cust Id":              hashdb.Decrypt(PreviousData.RefUserCustId),
		"FirstName":            hashdb.Decrypt(PreviousData.FirstName),
		"UserProfileImg":       hashdb.Decrypt(PreviousData.ProfileImg),
		"DOB":                  hashdb.Decrypt(PreviousData.DOB),
		"Gender":               hashdb.Decrypt(PreviousData.Gender),
		"Active Status":        PreviousData.ActiveStatus,
		"PhoneNo Country Code": PreviousData.PhoneNumberCode,
		"PhoneNo":              PreviousData.PhoneNumber,
		"Email":                PreviousData.Email,
	}

	updatedData := map[string]interface{}{
		"Cust Id":              PreviousData.RefUserCustId,
		"FirstName":            reqVal.FirstName,
		"UserProfileImg":       reqVal.ProfileImg,
		"DOB":                  reqVal.DOB,
		"Gender":               reqVal.Gender,
		"Active Status":        reqVal.ActiveStatus,
		"PhoneNo Country Code": reqVal.PhoneNumberCode,
		"PhoneNo":              reqVal.PhoneNumber,
		"Email":                PreviousData.Email,
	}

	ChangesData := helper.GetChanges(updatedData, oldData)

	if len(ChangesData) > 0 {
		var ChangesDataJSON []byte
		var errChange error
		ChangesDataJSON, errChange = json.Marshal(ChangesData)
		if errChange != nil {
			// Corrected log message
			log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", errChange)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		transData := 5

		errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.RefUserId), int(idValue), string(ChangesDataJSON)).Error
		if errTrans != nil {
			log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	usererr := tx.Exec(
		query.UpdatePatientQuery,
		reqVal.RefUserCustId,
		hashdb.Encrypt(reqVal.FirstName),
		hashdb.Encrypt(reqVal.ProfileImg),
		hashdb.Encrypt(reqVal.DOB),
		reqVal.Gender,
		reqVal.ActiveStatus,
		reqVal.RefUserId,
	).Error
	if usererr != nil {
		log.Printf("ERROR: Failed to update User: %v\n", usererr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	communicationerr := tx.Exec(
		query.UpdateCommunicationSQL,
		reqVal.PhoneNumberCode,
		reqVal.PhoneNumber,
		reqVal.Email,
		reqVal.RefUserId,
	).Error
	if communicationerr != nil {
		log.Printf("ERROR: Failed to update Communication: %v\n", communicationerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Account Updated"
}

func PostCreatePatientService(db *gorm.DB, reqVal model.CreateAppointmentPatientReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	Appointment := model.CreateAppointmentModel{
		UserId:          reqVal.RefUserId,
		SCId:            reqVal.SCId,
		AppointmentDate: reqVal.DateofAppointment,
		// AppointmentStartTime: reqVal.AppointmentStartTime,
		// AppointmentEndTime:   reqVal.AppointmentEndTime,
		// AppointmentUrgency: reqVal.AppointmentUrgency,
		AppointmentStatus:   true,
		AppointmentComplete: "fillform",
	}

	Appointmenterr := db.Create(&Appointment).Error
	if Appointmenterr != nil {
		log.Printf("ERROR: Failed to create Appointment: %v\n", Appointmenterr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if reqVal.Mailoption == "sendbywellthgreen" {

		htmlContent := mailservice.RegistrationMailContent(reqVal.Firstname, reqVal.SCCustId, reqVal.EmailId, "Use your Password", "Patient")

		subject := "Welcome – Your Appointment at " + reqVal.SCCustId + " Scan Centre"

		emailStatus := mailservice.MailService(reqVal.EmailId, htmlContent, subject)

		if !emailStatus {
			log.Error("Sending Mail Meets Error")
			return false, "Something went wrong, Try Again"
		}

	}

	InsertReportQuestionsErr := tx.Exec(query.InsertReportIntakeAllSQL, int(reqVal.RefUserId), int(Appointment.AppointmentId), timeZone.GetPacificTime(), int(idValue)).Error
	if InsertReportQuestionsErr != nil {
		log.Printf("ERROR: Failed to Report Question: %v\n", InsertReportQuestionsErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	InsertReportTextContentErr := tx.Exec(
		query.InsertNewReportTextContentSQL,
		int(reqVal.RefUserId),
		int(Appointment.AppointmentId),
		timeZone.GetPacificTime(),
		int(idValue),
	).Error
	if InsertReportTextContentErr != nil {
		log.Printf("ERROR: Failed to Report Text Content: %v\n", InsertReportTextContentErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Account Created"
}

func PostSendMailPatientService(db *gorm.DB, reqVal model.CreateAppointmentPatientReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	if reqVal.Mailoption == "sendbywellthgreen" {

		htmlContent := mailservice.RegistrationMailContent(reqVal.Firstname, reqVal.SCCustId, reqVal.EmailId, "Use your password", "Patient")

		subject := "Welcome – Scan Appointment at " + reqVal.SCCustId + " Scan Centre"

		emailStatus := mailservice.MailService(reqVal.EmailId, htmlContent, subject)

		if !emailStatus {
			log.Error("Sending Mail Meets Error")
			return false, "Something went wrong, Try Again"
		}

	}
	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Account Created"
}

func PostCancelResheduleAppointmentService(db *gorm.DB, reqVal model.CancelResheduleAppointmentReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	switch reqVal.AccessMethod {
	case "delete":

		DeleteAppointmentErr := tx.Exec(query.DeleteAppointmentSQL, false, reqVal.AppointmentId).Error
		if DeleteAppointmentErr != nil {
			log.Error(DeleteAppointmentErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

	case "reschedule":

		RescheduleAppointmentErr := tx.Exec(query.RescheduleAppointmentSQL, reqVal.AppointmentDate, reqVal.AppointmentId).Error
		if RescheduleAppointmentErr != nil {
			log.Error(RescheduleAppointmentErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	switch reqVal.AccessMethod {
	case "delete":
		return true, "Succcessfully Appointment Canceled"
	case "reschedule":
		return true, "Succcessfully Appointment Rescheduled"
	default:
		return false, "Something went wrong, Try Again"
	}
}
