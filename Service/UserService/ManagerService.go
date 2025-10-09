package service

import (
	becrypt "AuthenticationService/internal/Helper/Becrypt"
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	mailservice "AuthenticationService/internal/Helper/MailService"
	model "AuthenticationService/internal/Model/UserService"
	query "AuthenticationService/query/UserService"
	"encoding/json"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func PostManagerService(db *gorm.DB, reqVal model.ManagerRegisterReq, idValue int) (bool, string) {
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
		query.VerifyDataSQL,
		reqVal.PhoneNo,
		reqVal.Email,
	).Scan(&verifyData).Error

	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return false, "Something went wrong, Try Again"
	}

	if len(verifyData) > 0 {

		if verifyData[0].Email == reqVal.Email {
			return false, "Email Already Exists"
		} else {
			return false, "Mobile Number Already Exists"
		}

	}

	var TotalCount []model.TotalCountModel

	err := db.Raw(query.GetUsersCountSQL, 9).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return false, "Something went wrong, Try Again"
	}

	CustId := "WGM" + strconv.Itoa(TotalCount[0].TotalCount+1)

	UserData := model.CreateRadiologyModel{
		UserCustId:     CustId,
		RoleId:         9,
		FirstName:      hashdb.Encrypt(reqVal.FirstName),
		LastName:       hashdb.Encrypt(reqVal.LastName),
		UserProfileImg: hashdb.Encrypt(reqVal.ProfileImg),
		DOB:            hashdb.Encrypt(reqVal.DOB),
		Status:         true,
		UserAgreement:  false,
	}

	UserDataerr := tx.Create(&UserData).Error
	if UserDataerr != nil {
		log.Printf("ERROR: Failed to create Receptionist User Data: %v\n", UserDataerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	CommunicationData := model.CreateRadiologistCommunicationModel{
		UserId:             int(UserData.UserId),
		PhoneNoCountryCode: reqVal.PhoneNoCountryCode,
		PhoneNo:            reqVal.PhoneNo,
		Email:              reqVal.Email,
	}

	CommunicationDataerr := tx.Create(&CommunicationData).Error
	if CommunicationDataerr != nil {
		log.Printf("ERROR: Failed to create Receptionist Communication Data: %v\n", CommunicationDataerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	hashPassword, hashPassworderr := becrypt.HashPassword(reqVal.DOB)

	if hashPassworderr != nil {
		log.Printf("ERROR: Failed to hash password: %v\n", hashPassworderr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	AuthData := model.CreateRadiologistAuthModel{
		UserId: int(UserData.UserId),
		// Password:       hashdb.Encrypt(reqVal.DOB),
		Password:       hashdb.Encrypt("test@123"),
		HashPassword:   hashPassword,
		PasswordStatus: true,
	}

	AuthDataerr := tx.Create(&AuthData).Error
	if AuthDataerr != nil {
		log.Printf("ERROR: Failed to create Receptionist Auth Data: %v\n", AuthDataerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	ManagerDomain := model.CreateManagerDomainModel{
		UserId:         int(UserData.UserId),
		Pan:            hashdb.Encrypt(reqVal.Pan),
		Aadhar:         hashdb.Encrypt(reqVal.Aadhar),
		DrivingLicense: hashdb.Encrypt(reqVal.DriversLicense),
	}

	ManagerDomainerr := tx.Create(&ManagerDomain).Error
	if ManagerDomainerr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain: %v\n", ManagerDomainerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	for _, file := range reqVal.EducationCertificate {
		CVFiles := model.CreateEducationCertificateModel{
			UserId:        int(UserData.UserId),
			ECFileName:    file.FilesName,
			ECOldFileName: file.OldFileName,
			ECStatus:      true,
		}

		errCV := tx.Create(&CVFiles).Error
		if errCV != nil {
			log.Printf("ERROR: Failed to create Radiologist CV File: %v\n", errCV)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	history := model.RefTransHistory{
		TransTypeId: 15,
		THData:      "Account Created Successfully",
		UserId:      UserData.UserId,
		THActionBy:  idValue,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return false, "Something went wrong, Try Again"
	}

	htmlContent := mailservice.RegistrationMailContent(reqVal.FirstName+" "+reqVal.LastName, CustId, reqVal.Email, reqVal.DOB, "Wellthgreen Manager")

	subject := "Welcome to Wellthgreen Report Portal – Your User ID & Login Details Inside"

	emailStatus := mailservice.MailService(reqVal.Email, htmlContent, subject)

	if !emailStatus {
		log.Error("Sending Mail Meets Error")
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Account Created"
}

func PatchManagerService(db *gorm.DB, reqVal model.UpdateManagerReq, idValue int) (bool, string) {
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
		query.VerifyDataSQL,
		reqVal.PhoneNo,
		reqVal.Email,
	).Scan(&verifyData).Error

	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return false, "Something went wrong, Try Again"
	}

	if len(verifyData) > 0 {
		if verifyData[0].Email == reqVal.Email && verifyData[0].UserId != reqVal.ID {
			return false, "Email Already Exists"
		} else if verifyData[0].UserId != reqVal.ID {
			return false, "Mobile Number Already Exists"
		}
	}

	PreviousData := model.GetAllManagerData{}

	errPrev := tx.Raw(query.GetAllManagerDataSQL, reqVal.ID).Scan(&PreviousData).Error
	if errPrev != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", errPrev)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	oldData := map[string]interface{}{
		"FirstName":            hashdb.Decrypt(PreviousData.FirstName),
		"LastName":             hashdb.Decrypt(PreviousData.LastName),
		"UserProfileImg":       hashdb.Decrypt(PreviousData.UserProfileImg),
		"DOB":                  hashdb.Decrypt(PreviousData.DOB),
		"Active Status":        PreviousData.Status,
		"PhoneNo Country Code": PreviousData.PhoneNoCountryCode,
		"PhoneNo":              PreviousData.PhoneNo,
		"Email":                PreviousData.Email,
		"Pan":                  hashdb.Decrypt(PreviousData.MDPan),
		"Aadhar":               hashdb.Decrypt(PreviousData.MDAadhar),
		"Driving License":      hashdb.Decrypt(PreviousData.MDDrivingLicense),
	}

	updatedData := map[string]interface{}{
		"FirstName":            reqVal.FirstName,
		"LastName":             reqVal.LastName,
		"UserProfileImg":       reqVal.ProfileImg,
		"DOB":                  reqVal.DOB,
		"Active Status":        reqVal.Status,
		"PhoneNo Country Code": reqVal.PhoneNoCountryCode,
		"PhoneNo":              reqVal.PhoneNo,
		"Email":                reqVal.Email,
		"Pan":                  reqVal.Pan,
		"Aadhar":               reqVal.Aadhar,
		"Driving License":      reqVal.DriversLicense,
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

		transData := 16

		errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
		if errTrans != nil {
			log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	usererr := tx.Exec(
		query.UpdateUserSQL,
		hashdb.Encrypt(reqVal.FirstName),
		hashdb.Encrypt(reqVal.LastName),
		hashdb.Encrypt(reqVal.ProfileImg),
		hashdb.Encrypt(reqVal.DOB),
		reqVal.Status,
		reqVal.ID,
	).Error
	if usererr != nil {
		log.Printf("ERROR: Failed to update User: %v\n", usererr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	communicationerr := tx.Exec(
		query.UpdateCommunicationSQL,
		reqVal.PhoneNoCountryCode,
		reqVal.PhoneNo,
		reqVal.Email,
		reqVal.ID,
	).Error
	if communicationerr != nil {
		log.Printf("ERROR: Failed to update Communication: %v\n", communicationerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	ManagerDomainerr := tx.Exec(
		query.UpdateManagerDomainSQL,
		hashdb.Encrypt(reqVal.Pan),
		hashdb.Encrypt(reqVal.Aadhar),
		hashdb.Encrypt(reqVal.DriversLicense),
		reqVal.ID,
	).Error
	if ManagerDomainerr != nil {
		log.Printf("ERROR: Failed to update Manager Domain: %v\n", ManagerDomainerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	fmt.Println("*********************", reqVal.EducationCertificate)

	for _, file := range reqVal.EducationCertificate {

		fmt.Println("**************************************", file.Status)

		switch file.Status {
		case "new":

			oldData := map[string]interface{}{
				"FileUpdated": "",
				"UniqueFile":  "",
			}

			updatedData := map[string]interface{}{
				"FileUpdated": file.OldFileName,
				"UniqueFile":  file.FilesName,
			}

			ChangesData := helper.GetChanges(updatedData, oldData)

			if len(ChangesData) > 0 {

				CVFiles := model.CreateEducationCertificateModel{
					UserId:        int(reqVal.ID),
					ECFileName:    file.FilesName,
					ECOldFileName: file.OldFileName,
					ECStatus:      true,
				}

				errCV := tx.Create(&CVFiles).Error

				if errCV != nil {
					log.Printf("ERROR: Failed to create Radiologist CV File: %v\n", errCV)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

				var ChangesDataJSON []byte
				var errChange error
				ChangesDataJSON, errChange = json.Marshal(ChangesData)
				if errChange != nil {
					// Corrected log message
					log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", errChange)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

				transData := 16

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

		case "update":

			PrevData := model.GetEducationCertificateFilesModel{}
			errPrev := tx.Raw(query.GetEducationCertificateFilesSQL, file.Id).Scan(&PrevData).Error
			if errPrev != nil {
				log.Printf("ERROR: Failed to Get Exprience: %v\n", PrevData)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			oldData := map[string]interface{}{
				"FileUpdated": PrevData.ECOldFileName,
				"UniqueFile":  PrevData.ECFileName,
			}

			updatedData := map[string]interface{}{
				"FileUpdated": file.OldFileName,
				"UniqueFile":  file.FilesName,
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

				transData := 16

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			CVFileserr := tx.Exec(
				query.UpdateEductionCertificateFilesSQL,
				file.FilesName,
				file.OldFileName,
				true,
				file.Id,
			).Error

			if CVFileserr != nil {
				log.Printf("ERROR: Failed to update CV File: %v\n", CVFileserr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

		case "delete":
			PrevData := model.GetEducationCertificateFilesModel{}
			errPrev := tx.Raw(query.GetEducationCertificateFilesSQL, file.Id).Scan(&PrevData).Error
			if errPrev != nil {
				log.Printf("ERROR: Failed to Get Exprience: %v\n", PrevData)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			oldData := map[string]interface{}{
				"FileUpdated": PrevData.ECOldFileName,
				"UniqueFile":  PrevData.ECFileName,
			}

			updatedData := map[string]interface{}{
				"FileUpdated": "",
				"UniqueFile":  "",
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

				transData := 16

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			CVFileserr := tx.Exec(
				query.UpdateEductionCertificateFilesSQL,
				file.FilesName,
				file.OldFileName,
				false,
				file.Id,
			).Error

			if CVFileserr != nil {
				log.Printf("ERROR: Failed to delete CV File: %v\n", CVFileserr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Account Updated"
}
