package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/Analaytics"
	query "AuthenticationService/query/Analaytics"

	"gorm.io/gorm"
)

func GetAmountService(db *gorm.DB) (bool, int, int, []model.ScanCenterModel, []model.UserModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, 0, 0, []model.ScanCenterModel{}, []model.UserModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var AmountModel []model.AmountModel

	// Amount Value
	AmountValueErr := tx.Raw(query.GetAmountSQL).Scan(&AmountModel).Error
	if AmountValueErr != nil {
		log.Fatal(AmountValueErr.Error())
		return false, 0, 0, []model.ScanCenterModel{}, []model.UserModel{}
	}

	//List of Scan Center
	var ScancenterData []model.ScanCenterModel
	ScancenterDataErr := tx.Raw(query.ListAllScanCenter).Scan(&ScancenterData).Error
	if ScancenterDataErr != nil {
		log.Fatal(ScancenterDataErr.Error())
		return false, 0, 0, []model.ScanCenterModel{}, []model.UserModel{}
	}

	//List of User
	var UserData []model.UserModel
	UserDataErr := tx.Raw(query.ListUserSQL).Scan(&UserData).Error
	if UserDataErr != nil {
		log.Fatal(UserDataErr.Error())
		return false, 0, 0, []model.ScanCenterModel{}, []model.UserModel{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, 0, 0, []model.ScanCenterModel{}, []model.UserModel{}
	}

	return true, AmountModel[0].ScancenterAmount, AmountModel[0].UserAmount, ScancenterData, UserData
}

func UpdateAmountService(db *gorm.DB, reqVal model.AmountModel) (bool, string) {
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

	UpdateErr := tx.Exec(query.UpdateAmountSQL, reqVal.ScancenterAmount, reqVal.UserAmount).Error
	if UpdateErr != nil {
		log.Fatal(UpdateErr)
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Updated"
}

func GetInvoiceDataService(db *gorm.DB, reqVal model.GetInvoiceDataReq) model.GetInvoiceDataReponse {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.GetInvoiceDataReponse{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var response model.GetInvoiceDataReponse

	// Amount Value
	AmountValueErr := tx.Raw(query.GetAmountSQL).Scan(&response.AmountModel).Error
	if AmountValueErr != nil {
		log.Fatal(AmountValueErr.Error())
		return model.GetInvoiceDataReponse{}
	}

	switch reqVal.Type {
	case 1:
		// Get Scan Center Data
		ScancenterDataErr := tx.Raw(query.ListOneScanCenter, reqVal.UserId).Scan(&response.ScanCenterModel).Error
		if ScancenterDataErr != nil {
			log.Fatal(ScancenterDataErr.Error())
			return model.GetInvoiceDataReponse{}
		}

		for i, data := range response.ScanCenterModel {
			response.ScanCenterModel[i].SCName = hashdb.Decrypt(data.SCName)
			response.ScanCenterModel[i].Address = hashdb.Decrypt(data.Address)
		}

		// Get Scan Center Total Count
		ScanCenterCountDataErr := tx.Raw(query.GetScanCenterCountPerMonthSQL, reqVal.Monthyear, reqVal.UserId).Scan(&response.GetCountScanCenterMonthModel).Error
		if ScanCenterCountDataErr != nil {
			log.Fatal(ScanCenterCountDataErr.Error())
			return model.GetInvoiceDataReponse{}
		}

	case 2:
		// Get User Details
		UserDataErr := tx.Raw(query.GetOneUserSQL, reqVal.UserId).Scan(&response.GetUserModel).Error
		if UserDataErr != nil {
			log.Fatal(UserDataErr.Error())
			return model.GetInvoiceDataReponse{}
		}

		// Decrypt User Details
		for i, data := range response.GetUserModel {
			response.GetUserModel[i].Name = hashdb.Decrypt(data.Name)
		}

		// Total Count for User
		UserCountDataErr := tx.Raw(query.WellGreenUserIndicatesAnalayticsInvoiceSQL, reqVal.UserId, reqVal.Monthyear).Scan(&response.AdminOverallScanIndicatesAnalayticsModel).Error
		if UserCountDataErr != nil {
			log.Fatal(UserCountDataErr.Error())
			return model.GetInvoiceDataReponse{}
		}

	default:
		log.Printf("Invalid Type: %d", reqVal.Type)
		return model.GetInvoiceDataReponse{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return model.GetInvoiceDataReponse{}
	}

	return response
}

func GenerateInvoiceDataService(db *gorm.DB, reqVal model.GenerateInvoiceReq, idValue int) (bool, string) {
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

	InsertInvoiceErr := tx.Exec(
		query.InsertInvoiceSQL,
		reqVal.SCId,
		reqVal.UserId,
		reqVal.FromId,
		reqVal.FromName,
		reqVal.FromPhone,
		reqVal.FromEmail,
		reqVal.FromPan,
		reqVal.FromGst,
		reqVal.FromAddress,
		reqVal.ToId,
		reqVal.ToName,
		reqVal.BillingFrom,
		reqVal.BillingTo,
		reqVal.ModeOfPayment,
		reqVal.UpiId,
		reqVal.AccountHolderName,
		reqVal.AccountNumber,
		reqVal.Bank,
		reqVal.Branch,
		reqVal.IFSC,
		reqVal.Quantity,
		reqVal.Amount,
		reqVal.Total,
		timeZone.GetPacificTime(),
		idValue,
		reqVal.ToAddress,
		reqVal.Signature,
	).Error
	if InsertInvoiceErr != nil {
		log.Fatal(InsertInvoiceErr)
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Invoice Generated"
}

func GetInvoiceHistoryService(db *gorm.DB, reqVal model.GetInvoiceHistoryReq, idValue int) ([]model.InvoiceHistory, []model.TakenDate) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	useQuery := query.GetInvoiceHistoryScancenterSQL

	if reqVal.Type == 2 {
		useQuery = query.GetInvoiceHistoryUserSQL
	} else if reqVal.Type != 1 && reqVal.Type != 2 {
		log.Printf("Invalid Type: %d", reqVal.Type)
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	//Get Invoice Data
	var invoiceHistory []model.InvoiceHistory
	InvoiceHistoryErr := tx.Raw(useQuery, reqVal.UserId).Scan(&invoiceHistory).Error
	if InvoiceHistoryErr != nil {
		log.Fatal(InvoiceHistoryErr)
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	for i, data := range invoiceHistory {
		if len(data.RefIHSignature) > 0 {
			DriversLicenseNoImgHelperData, viewErr := helper.ViewFile("./Assets/Files/" + data.RefIHSignature)
			if viewErr != nil {
				// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
				log.Fatalf("Failed to read DrivingLicense file: %v", viewErr)
			}
			if DriversLicenseNoImgHelperData != nil {
				invoiceHistory[i].RefIHSignatureFile = &model.FileData{
					Base64Data:  DriversLicenseNoImgHelperData.Base64Data,
					ContentType: DriversLicenseNoImgHelperData.ContentType,
				}
			}
		} else {
			invoiceHistory[i].RefIHSignatureFile = nil
		}
	}

	//Already Taken Date
	var invoiceHistoryTakenDate []model.TakenDate
	InvoiceHistoryTakenErr := tx.Raw(useQuery, reqVal.UserId).Scan(&invoiceHistoryTakenDate).Error
	if InvoiceHistoryTakenErr != nil {
		log.Fatal(InvoiceHistoryTakenErr)
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	return invoiceHistory, invoiceHistoryTakenDate
}
