package service

import (
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/UserService"
	query "AuthenticationService/query/UserService"
	"encoding/json"

	"gorm.io/gorm"
)

func PostScanCenterService(db *gorm.DB, reqVal model.ScanCenterRegisterReq, idValue int) (bool, string) {
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

	var verifyData []model.ScanCenterVerifyData

	verifyDataerr := db.Raw(query.ScanCenterVerifyDataSQL, reqVal.Telephone, reqVal.Email, reqVal.CustId).Scan(&verifyData).Error
	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return false, "Something went wrong, Try Again"
	}

	if len(verifyData) > 0 {

		if verifyData[0].SCEmail == reqVal.Email {
			return false, "Email Already Exists"
		} else if verifyData[0].SCCustId == reqVal.CustId {
			return false, "Customer Id Already Exists"
		} else {
			return false, "Mobile Number Already Exists"
		}

	}

	// var TotalCount []model.TotalCount

	// err := db.Raw(query.GetScanCenterCountSQL).Scan(&TotalCount).Error
	// if err != nil {
	// 	log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
	// 	return false, "Something went wrong, Try Again"
	// }

	// CustId := "SC" + strconv.Itoa(TotalCount[0].TotalCount+100001)

	ScanCenter := model.ScanCenterModel{
		CustId:       reqVal.CustId,
		Logo:         hashdb.Encrypt(reqVal.Logo),
		Name:         hashdb.Encrypt(reqVal.Name),
		Address:      hashdb.Encrypt(reqVal.Address),
		PhoneNo1:     reqVal.Telephone,
		Email:        reqVal.Email,
		Website:      hashdb.Encrypt(reqVal.Website),
		Appointments: reqVal.Appointments,
	}

	ScanCenterErr := db.Create(&ScanCenter).Error
	if ScanCenterErr != nil {
		log.Error("Scan Center INSERT ERROR at Trnasaction: " + ScanCenterErr.Error())
		return false, "Something went wrong, Try Again"
	}

	history := model.RefTransHistory{
		TransTypeId: 2,
		THData:      "Scan Center Created Successfully",
		UserId:      ScanCenter.Id,
		THActionBy:  idValue,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Account Created"
}

func PatchScanCenterService(db *gorm.DB, reqVal model.UpdateScanCentertReq, idValue int) (bool, string) {
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

	var verifyData []model.ScanCenterVerifyData

	verifyDataerr := db.Raw(query.UpdateScanCenterVerifyDataSQL, reqVal.Telephone, reqVal.Email).Scan(&verifyData).Error
	if verifyDataerr != nil {
		log.Printf("ERROR: Failed to fetch Verify data: %v", verifyDataerr)
		return false, "Something went wrong, Try Again"
	}

	if len(verifyData) > 0 {

		if verifyData[0].SCEmail == reqVal.Email && verifyData[0].SCId != reqVal.ID {
			return false, "Email Already Exists"
		} else if verifyData[0].SCId != reqVal.ID {
			return false, "Mobile Number Already Exists"
		}

	}

	PreviousData := model.ScanCenterModel{}

	errPrev := tx.Raw(query.GetScancenterOneDataSQL, reqVal.ID).Scan(&PreviousData).Error
	if errPrev != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", errPrev)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	oldData := map[string]interface{}{
		"Name":               hashdb.Decrypt(PreviousData.Name),
		"Address":            hashdb.Decrypt(PreviousData.Address),
		"Email":              PreviousData.Email,
		"Telephone":          PreviousData.PhoneNo1,
		"Logo":               hashdb.Decrypt(PreviousData.Logo),
		"Website":            hashdb.Decrypt(PreviousData.Website),
		"Appointment Status": PreviousData.Appointments,
		"Disclaimer":         PreviousData.Disclaimer,
		"Brouchure":          PreviousData.Brouchure,
		"Guidelines":         PreviousData.Guidelines,
	}

	updatedData := map[string]interface{}{
		"Name":               reqVal.Name,
		"Address":            reqVal.Address,
		"Email":              reqVal.Email,
		"Telephone":          reqVal.Telephone,
		"Logo":               reqVal.Logo,
		"Website":            reqVal.Website,
		"Appointment Status": reqVal.Appointments,
		"Disclaimer":         reqVal.Disclaimer,
		"Brouchure":          reqVal.Brouchure,
		"Guidelines":         PreviousData.Guidelines,
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

		transData := 3

		errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.ID), int(idValue), string(ChangesDataJSON)).Error
		if errTrans != nil {
			log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	ScanCenterErr := tx.Exec(
		query.UpdateScancenterSQL,
		hashdb.Encrypt(reqVal.Logo),
		hashdb.Encrypt(reqVal.Name),
		hashdb.Encrypt(reqVal.Address),
		reqVal.Telephone,
		reqVal.Email,
		hashdb.Encrypt(reqVal.Website),
		reqVal.Appointments,
		reqVal.Disclaimer,
		reqVal.Brouchure,
		reqVal.Guidelines,
		reqVal.ID,
	).Error
	if ScanCenterErr != nil {
		log.Printf("ERROR: Failed to update Scan Center: %v\n", ScanCenterErr)
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
