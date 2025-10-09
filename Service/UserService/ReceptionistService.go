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

func PostReceptionistService(db *gorm.DB, reqVal model.ReceptionistRegisterReq, idValue int) (bool, string) {
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

	fmt.Println("Verification COmplete!")

	var TotalCount []model.TotalCount

	fmt.Println("Step 1")
	err := db.Raw(query.GetUsersScanCountSQL, 3, reqVal.ScanCenterId).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return false, "Something went wrong, Try Again"
	}

	fmt.Println(TotalCount[0])
	fmt.Println(strconv.Itoa(TotalCount[0].TotalCount + 1))

	CustId := "C" + TotalCount[0].SCCustId + "M" + strconv.Itoa(TotalCount[0].TotalCount+1)
	fmt.Println("Step 2")

	UserData := model.CreateRadiologyModel{
		UserCustId:     CustId,
		RoleId:         3,
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

	ReceptionistDomain := model.CreateReceptionstDomainModel{
		UserId:         int(UserData.UserId),
		SSId:           hashdb.Encrypt(reqVal.SSNo),
		DrivingLicense: hashdb.Encrypt(reqVal.DriversLicenseNo),
	}

	ReceptionistDomainerr := tx.Create(&ReceptionistDomain).Error
	if ReceptionistDomainerr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain: %v\n", ReceptionistDomainerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	MapScanCenter := model.MapScanCenterModel{
		UserId:    int(UserData.UserId),
		SCId:      int(reqVal.ScanCenterId),
		RTId:      3,
		SCMStatus: true,
	}

	MapScanCentererr := tx.Create(&MapScanCenter).Error
	if MapScanCentererr != nil {
		log.Printf("ERROR: Failed to create Receptionist Domain: %v\n", MapScanCentererr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	history := model.RefTransHistory{
		TransTypeId: 7,
		THData:      "Account Created Successfully",
		UserId:      UserData.UserId,
		THActionBy:  idValue,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return false, "Something went wrong, Try Again"
	}

	htmlContent := mailservice.RegistrationMailContent(reqVal.FirstName+" "+reqVal.LastName, CustId, reqVal.Email, reqVal.DOB, "Scan Center Admin")

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

func PatchReceptionistService(db *gorm.DB, reqVal model.UpdateReceptionistReq, idValue int) (bool, string) {
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

	PreviousData := model.GetAllReceptionistData{}

	errPrev := tx.Raw(query.GetAllReceptionistDataSQL, reqVal.ID).Scan(&PreviousData).Error
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
		"Driving License":        hashdb.Decrypt(PreviousData.DriversLicenseNo),
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
		"Driving License":        hashdb.Decrypt(PreviousData.DriversLicenseNo),
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

		transData := 8

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

	ReceptionistDomainerr := tx.Exec(
		query.UpdateReceptionistDomainSQL,
		hashdb.Encrypt(reqVal.SocialSecurityNo),
		hashdb.Encrypt(reqVal.DriversLicenseNo),
		reqVal.ID,
	).Error
	if ReceptionistDomainerr != nil {
		log.Printf("ERROR: Failed to update Doctor Domain: %v\n", ReceptionistDomainerr)
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
