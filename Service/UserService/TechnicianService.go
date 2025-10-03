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
	"strconv"

	"gorm.io/gorm"
)

func PostTechnicianService(db *gorm.DB, reqVal model.TechnicianRegisterReq, idValue int) (bool, string) {
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

	var TotalCount []model.TotalCount

	err := db.Raw(query.GetUsersScanCountSQL, 2, reqVal.ScanCenterId).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return false, "Something went wrong, Try Again"
	}

	CustId := "C" + TotalCount[0].SCCustId + "T" + strconv.Itoa(TotalCount[0].TotalCount+1)

	UserData := model.CreateRadiologyModel{
		UserCustId:     CustId,
		RoleId:         2,
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

	TechnicianDomain := model.CreateTechnicianDomainModel{
		UserId:                int(UserData.UserId),
		TDTrainedEaseQTStatus: reqVal.TrainedEaseQT,
		TDSSNo:                hashdb.Encrypt(reqVal.SSNo),
		TDDrivingLicense:      hashdb.Encrypt(reqVal.DriversLicenseNo),
		TDDigitalSignature:    hashdb.Encrypt(reqVal.DigitalSignature),
	}

	TechnicianDomainerr := tx.Create(&TechnicianDomain).Error
	if TechnicianDomainerr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain: %v\n", TechnicianDomainerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	for _, file := range reqVal.LicenseFiles {
		LicenseFiles := model.CreateRadiologistLicenseModel{
			UserId:       int(UserData.UserId),
			LFileName:    file.FilesName,
			LOldFileName: file.OldFileName,
			LStatus:      true,
		}

		errLicense := tx.Create(&LicenseFiles).Error
		if errLicense != nil {
			log.Printf("ERROR: Failed to create Radiologist License File: %v\n", errLicense)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	MapScanCenter := model.MapScanCenterModel{
		UserId:    int(UserData.UserId),
		SCId:      int(reqVal.ScanCenterId),
		RTId:      2,
		SCMStatus: true,
	}

	MapScanCentererr := tx.Create(&MapScanCenter).Error
	if MapScanCentererr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain: %v\n", MapScanCentererr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	history := model.RefTransHistory{
		TransTypeId: 4,
		THData:      "Account Created Successfully",
		UserId:      UserData.UserId,
		THActionBy:  idValue,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return false, "Something went wrong, Try Again"
	}

	htmlContent := mailservice.RegistrationMailContent(reqVal.FirstName+" "+reqVal.LastName, CustId, reqVal.Email, reqVal.DOB, "Scan Center Technician")

	subject := "Welcome to easeQT – Your User ID & Login Details Inside"

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

func PatchTechnicianService(db *gorm.DB, reqVal model.UpdateTechnicianReq, idValue int) (bool, string) {
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

	PreviousData := model.GetAllTechnicianData{}

	errPrev := tx.Raw(query.GetAllTechnicianDataSQL, reqVal.ID).Scan(&PreviousData).Error
	if errPrev != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", errPrev)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	oldData := map[string]interface{}{
		"FirstName":              hashdb.Decrypt(PreviousData.FirstName),
		"LastName":               hashdb.Decrypt(PreviousData.LastName),
		"UserProfileImg":         hashdb.Decrypt(PreviousData.UserProfileImg),
		"DOB":                    hashdb.Decrypt(PreviousData.DOB),
		"Active Status":          PreviousData.Status,
		"PhoneNo Country Code":   PreviousData.PhoneNoCountryCode,
		"PhoneNo":                PreviousData.PhoneNo,
		"Email":                  PreviousData.Email,
		"Tranined By Ease QT":    PreviousData.TrainedEaseQT,
		"Social Security Number": hashdb.Decrypt(PreviousData.SSNo),
		"Driving License":        hashdb.Decrypt(PreviousData.DrivingLicense),
		"Digital Signature":      hashdb.Decrypt(PreviousData.DigitalSignature),
	}

	updatedData := map[string]interface{}{
		"FirstName":              reqVal.FirstName,
		"LastName":               reqVal.LastName,
		"UserProfileImg":         reqVal.ProfileImg,
		"DOB":                    reqVal.DOB,
		"Active Status":          reqVal.Status,
		"PhoneNo Country Code":   reqVal.PhoneNoCountryCode,
		"PhoneNo":                reqVal.PhoneNo,
		"Email":                  PreviousData.Email,
		"Tranined By Ease QT":    reqVal.TrainedEaseQT,
		"Social Security Number": reqVal.SSNo,
		"Driving License":        reqVal.DriversLicenseNo,
		"Digital Signature":      reqVal.DigitalSignature,
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

	TechnicianDomainerr := tx.Exec(
		query.UpdateTechnicianDomainSQL,
		reqVal.TrainedEaseQT,
		reqVal.SSNo,
		reqVal.DriversLicenseNo,
		reqVal.DigitalSignature,
		reqVal.ID,
	).Error
	if TechnicianDomainerr != nil {
		log.Printf("ERROR: Failed to update Doctor Domain: %v\n", TechnicianDomainerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	for _, file := range reqVal.LicenseFiles {

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

				LicenseFiles := model.CreateRadiologistLicenseModel{
					UserId:       int(reqVal.ID),
					LFileName:    file.FilesName,
					LOldFileName: file.OldFileName,
					LStatus:      true,
				}

				errLicense := tx.Create(&LicenseFiles).Error
				if errLicense != nil {
					log.Printf("ERROR: Failed to create Radiologist License File: %v\n", errLicense)
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

				transData := 5

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

		case "update":

			PrevData := model.GetLicenseFilesModel{}
			errPrev := tx.Raw(query.GetLicenseFilesSQL, file.Id).Scan(&PrevData).Error
			if errPrev != nil {
				log.Printf("ERROR: Failed to Get Exprience: %v\n", PrevData)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			oldData := map[string]interface{}{
				"FileUpdated": PrevData.LOldFileName,
				"UniqueFile":  PrevData.LFileName,
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

				transData := 5

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			LicenseFileerr := tx.Exec(
				query.UpdateLicenseFilesSQL,
				file.FilesName,
				file.OldFileName,
				true,
				file.Id,
			).Error

			if LicenseFileerr != nil {
				log.Printf("ERROR: Failed to update CV File: %v\n", LicenseFileerr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

		case "delete":
			PrevData := model.GetLicenseFilesModel{}
			errPrev := tx.Raw(query.GetLicenseFilesSQL, file.Id).Scan(&PrevData).Error
			if errPrev != nil {
				log.Printf("ERROR: Failed to Get Exprience: %v\n", PrevData)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			oldData := map[string]interface{}{
				"FileUpdated": PrevData.LOldFileName,
				"UniqueFile":  PrevData.LFileName,
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

				transData := 5

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			LicenseFileerr := tx.Exec(
				query.UpdateLicenseFilesSQL,
				file.FilesName,
				file.OldFileName,
				false,
				file.Id,
			).Error

			if LicenseFileerr != nil {
				log.Printf("ERROR: Failed to update CV File: %v\n", LicenseFileerr)
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
