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

func PostCoDoctorService(db *gorm.DB, reqVal model.DoctorRegisterReq, idValue int) (bool, string) {
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

	err := db.Raw(query.GetUsersScanCountSQL, 8, reqVal.ScanCenterId).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return false, "Something went wrong, Try Again"
	}

	CustId := TotalCount[0].SCCustId + "R" + strconv.Itoa(TotalCount[0].TotalCount+1)

	UserData := model.CreateDoctorModel{
		UserCustId:     CustId,
		RoleId:         8,
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

	CommunicationData := model.CreateDoctorCommunicationModel{
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

	AuthData := model.CreateDoctorAuthModel{
		UserId:         int(UserData.UserId),
		Password:       hashdb.Encrypt(reqVal.DOB),
		HashPassword:   hashPassword,
		PasswordStatus: true,
	}

	AuthDataerr := tx.Create(&AuthData).Error
	if AuthDataerr != nil {
		log.Printf("ERROR: Failed to create Receptionist Auth Data: %v\n", AuthDataerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	DomainData := model.CreateCoDoctorDomainModel{
		UserId:           int(UserData.UserId),
		SocialSecurityNo: hashdb.Encrypt(reqVal.SocialSecurityNo),
		NPI:              hashdb.Encrypt(reqVal.NPI),
		DrivingLicense:   hashdb.Encrypt(reqVal.DriversLicenseNo),
		DigitalSignature: hashdb.Encrypt(reqVal.DigitalSignature),
		Specialization:   hashdb.Encrypt(reqVal.Specialization),
	}

	DomainDataerr := tx.Create(&DomainData).Error
	if DomainDataerr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain Data: %v\n", DomainDataerr)
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

	for _, file := range reqVal.MedicalLicenseSecurity {
		MedicalLicenseSecurity := model.MedicalLicenseSecurityModel{
			UserId:   int(UserData.UserId),
			MLSState: hashdb.Encrypt(file.State),
			MLSNo:    hashdb.Encrypt(file.MedicalLicenseSecurityNo),
			MLStatus: true,
		}

		errMedicalLicense := tx.Create(&MedicalLicenseSecurity).Error
		if errMedicalLicense != nil {
			log.Printf("ERROR: Failed to create Radiologist License File: %v\n", errMedicalLicense)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	for _, file := range reqVal.MalpracticeInsuranceDetails {
		LicenseFiles := model.CreateMalpractice{
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
		RTId:      8,
		SCMStatus: true,
	}

	MapScanCentererr := tx.Create(&MapScanCenter).Error
	if MapScanCentererr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain: %v\n", MapScanCentererr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	history := model.RefTransHistory{
		TransTypeId: 19,
		THData:      "Account Created Successfully",
		UserId:      UserData.UserId,
		THActionBy:  idValue,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return false, "Something went wrong, Try Again"
	}

	htmlContent := mailservice.RegistrationMailContent(reqVal.FirstName+" "+reqVal.LastName, CustId, reqVal.Email, reqVal.DOB, "Scan Center Reviewer")

	subject := "Welcome to Wellthgreen HealthCare Pvt Ltd – Your User ID & Login Details Inside"

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

func PatchCoDoctorService(db *gorm.DB, reqVal model.UpdateDoctorReq, idValue int) (bool, string) {
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

	PreviousData := model.GetAllCODoctorData{}

	errPrev := tx.Raw(query.GetAllCoDoctorDataSQL, reqVal.ID).Scan(&PreviousData).Error
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
		"Social Security Number": hashdb.Decrypt(PreviousData.SocialSecurityNo),
		"NPI":                    hashdb.Decrypt(PreviousData.NPI),
		"Driving License":        hashdb.Decrypt(PreviousData.DrivingLicense),
		"Specialization":         hashdb.Decrypt(PreviousData.Specialization),
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
		"Email":                  reqVal.Email,
		"Social Security Number": reqVal.SocialSecurityNo,
		"NPI":                    reqVal.NPI,
		"Driving License":        reqVal.DriversLicenseNo,
		"Specialization":         reqVal.Specialization,
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

		transData := 20

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

	DoctorDomainerr := tx.Exec(
		query.UpdateCoDoctorDomainSQL,
		hashdb.Encrypt(reqVal.SocialSecurityNo),
		hashdb.Encrypt(reqVal.NPI),
		hashdb.Encrypt(reqVal.DriversLicenseNo),
		hashdb.Encrypt(reqVal.DigitalSignature),
		hashdb.Encrypt(reqVal.Specialization),
		reqVal.ID,
	).Error
	if DoctorDomainerr != nil {
		log.Printf("ERROR: Failed to update Doctor Domain: %v\n", DoctorDomainerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	for _, medicalsecurity := range reqVal.MedicalLicenseSecurity {

		switch medicalsecurity.MLStatus {
		case "new":
			oldData := map[string]interface{}{
				"State":                  "",
				"Medical License Number": "",
			}

			updatedData := map[string]interface{}{
				"State":                  medicalsecurity.MLSState,
				"Medical License Number": medicalsecurity.MLSNo,
			}

			ChangesData := helper.GetChanges(updatedData, oldData)

			if len(ChangesData) > 0 {

				MedicalLicenseSecurity := model.MedicalLicenseSecurityModel{
					UserId:   int(reqVal.ID),
					MLSState: hashdb.Encrypt(medicalsecurity.MLSState),
					MLSNo:    hashdb.Encrypt(medicalsecurity.MLSNo),
					MLStatus: true,
				}

				errMedicalLicense := tx.Create(&MedicalLicenseSecurity).Error
				if errMedicalLicense != nil {
					log.Printf("ERROR: Failed to create Radiologist License File: %v\n", errMedicalLicense)
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

				transData := 20

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}
		case "update":
			PrevData := model.GetMedicalLicenseSecurityModel{}
			errPrev := tx.Raw(query.GetMedicalLicenseSecuritySQL, medicalsecurity.MLSId).Scan(&PrevData).Error
			if errPrev != nil {
				log.Printf("ERROR: Failed to Get Exprience: %v\n", PrevData)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			oldData := map[string]interface{}{
				"State":                  PrevData.MLSState,
				"Medical License Number": PrevData.MLSNo,
			}

			updatedData := map[string]interface{}{
				"State":                  medicalsecurity.MLSState,
				"Medical License Number": medicalsecurity.MLSNo,
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

				transData := 20

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			MedicalLicenseSecurityerr := tx.Exec(
				query.UpdateMedicalLicenseSecuritySQL,
				hashdb.Encrypt(medicalsecurity.MLSState),
				hashdb.Encrypt(medicalsecurity.MLSNo),
				true,
				medicalsecurity.MLSId,
			).Error

			if MedicalLicenseSecurityerr != nil {
				log.Printf("ERROR: Failed to update CV File: %v\n", MedicalLicenseSecurityerr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}
		case "delete":
			PrevData := model.GetMedicalLicenseSecurityModel{}
			errPrev := tx.Raw(query.GetMedicalLicenseSecuritySQL, medicalsecurity.MLSId).Scan(&PrevData).Error
			if errPrev != nil {
				log.Printf("ERROR: Failed to Get Exprience: %v\n", PrevData)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			oldData := map[string]interface{}{
				"State":                  PrevData.MLSState,
				"Medical License Number": PrevData.MLSNo,
			}

			updatedData := map[string]interface{}{
				"State":                  "",
				"Medical License Number": "",
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

				transData := 20

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			MedicalLicenseSecurityerr := tx.Exec(
				query.UpdateMedicalLicenseSecuritySQL,
				hashdb.Encrypt(medicalsecurity.MLSState),
				hashdb.Encrypt(medicalsecurity.MLSNo),
				false,
				medicalsecurity.MLSId,
			).Error

			if MedicalLicenseSecurityerr != nil {
				log.Printf("ERROR: Failed to update CV File: %v\n", MedicalLicenseSecurityerr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}
		}

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

				transData := 20

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

				transData := 20

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

				transData := 20

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

	for _, file := range reqVal.MalpracticeInsuranceDetails {

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

				LicenseFiles := model.CreateRadiologistMalpracticeModel{
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

				transData := 20

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

		case "update":

			PrevData := model.GetMalpracticeFilesModel{}
			errPrev := tx.Raw(query.GetMalpracticeFilesSQL, file.Id).Scan(&PrevData).Error
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

				transData := 20

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			LicenseFileerr := tx.Exec(
				query.UpdateMalpracticeFilesSQL,
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
			PrevData := model.GetMalpracticeFilesModel{}
			errPrev := tx.Raw(query.GetMalpracticeFilesSQL, file.Id).Scan(&PrevData).Error
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

				transData := 20

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}

			LicenseFileerr := tx.Exec(
				query.UpdateMalpracticeFilesSQL,
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
