package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/Analaytics"
	query "AuthenticationService/query/Analaytics"
	"fmt"

	"gorm.io/gorm"
)

func GetAmountService(db *gorm.DB) (bool, []model.AmountModel, []model.ScanCenterModel, []model.UserModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, []model.AmountModel{}, []model.ScanCenterModel{}, []model.UserModel{}
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
		log.Error(AmountValueErr.Error())
		return false, []model.AmountModel{}, []model.ScanCenterModel{}, []model.UserModel{}
	}

	//List of Scan Center
	var ScancenterData []model.ScanCenterModel
	ScancenterDataErr := tx.Raw(query.ListAllScanCenter).Scan(&ScancenterData).Error
	if ScancenterDataErr != nil {
		log.Error(ScancenterDataErr.Error())
		return false, []model.AmountModel{}, []model.ScanCenterModel{}, []model.UserModel{}
	}

	//List of User
	var UserData []model.UserModel
	UserDataErr := tx.Raw(query.ListUserSQL).Scan(&UserData).Error
	if UserDataErr != nil {
		log.Error(UserDataErr.Error())
		return false, []model.AmountModel{}, []model.ScanCenterModel{}, []model.UserModel{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, []model.AmountModel{}, []model.ScanCenterModel{}, []model.UserModel{}
	}

	return true, AmountModel, ScancenterData, UserData
}

func UpdateAmountService(db *gorm.DB, reqVal model.AmountModel) (bool, string) {
	log := logger.InitLogger()

	fmt.Println(reqVal)

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

	UpdateErr := tx.Exec(
		query.UpdateAmountSQL,
		reqVal.TASform,
		reqVal.TADaform,
		reqVal.TADbform,
		reqVal.TADcform,
		reqVal.TAXform,
		reqVal.TAEditform,
		reqVal.TADScribeTotalcase,
		1,
	).Error
	if UpdateErr != nil {
		log.Error(UpdateErr)
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
		log.Error(AmountValueErr.Error())
		return model.GetInvoiceDataReponse{}
	}

	switch reqVal.Type {
	case 1:
		// Get Scan Center Data
		ScancenterDataErr := tx.Raw(query.ListOneScanCenter, reqVal.UserId).Scan(&response.ScanCenterModel).Error
		if ScancenterDataErr != nil {
			log.Error(ScancenterDataErr.Error())
			return model.GetInvoiceDataReponse{}
		}

		for i, data := range response.ScanCenterModel {
			response.ScanCenterModel[i].SCName = hashdb.Decrypt(data.SCName)
			response.ScanCenterModel[i].Address = hashdb.Decrypt(data.Address)
		}

		// Get Scan Center Total Count
		ScanCenterCountDataErr := tx.Raw(query.GetScanCenterCountPerMonthSQL, reqVal.UserId, reqVal.Monthyear).Scan(&response.GetCountScanCenterMonthModel).Error
		if ScanCenterCountDataErr != nil {
			log.Error(ScanCenterCountDataErr.Error())
			return model.GetInvoiceDataReponse{}
		}

	case 2:
		// Get User Details
		UserDataErr := tx.Raw(query.GetOneUserSQL, reqVal.UserId).Scan(&response.GetUserModel).Error
		if UserDataErr != nil {
			log.Error(UserDataErr.Error())
			return model.GetInvoiceDataReponse{}
		}

		// Decrypt User Details
		for i, data := range response.GetUserModel {
			response.GetUserModel[i].Name = hashdb.Decrypt(data.Name)
			response.GetUserModel[i].Pan = hashdb.Decrypt(data.Pan)
		}

		var Date = ""
		if len(reqVal.Monthyear) > 0 {
			Date = reqVal.Monthyear + "-01"
		}

		// Total Count for User
		UserCountDataErr := tx.Raw(query.WellGreenUserIndicatesAnalayticsInvoiceSQL, reqVal.UserId, Date).Scan(&response.AdminOverallScanIndicatesAnalayticsModel).Error
		if UserCountDataErr != nil {
			log.Error(UserCountDataErr.Error())
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

	fmt.Println("==================>", reqVal)

	var invoiceHistory []model.InvoiceHistory

	InsertInvoiceErr := tx.Raw(
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
		timeZone.GetPacificTime(),
		idValue,
		reqVal.ToAddress,
		reqVal.Signature,
		reqVal.SFormquantity,
		reqVal.SFormamount,
		reqVal.DaFormquantity,
		reqVal.DaFormamount,
		reqVal.DbFormquantity,
		reqVal.DbFormamount,
		reqVal.DcFormquantity,
		reqVal.DcFormamount,
		reqVal.XFormquantity,
		reqVal.XFormamount,
		reqVal.Editquantity,
		reqVal.EditFormamount,
		reqVal.ScribeTotalcasequantity,
		reqVal.ScribeTotalcaseamount,
		reqVal.OtherExpensesAmount,
		reqVal.DeductibleExpensesAmount,
		reqVal.ScanCenterTotalCase,
		reqVal.ScancentercaseAmount,
		reqVal.Total,
	).Scan(&invoiceHistory).Error
	if InsertInvoiceErr != nil {
		log.Error(InsertInvoiceErr)
		return false, "Something went wrong, Try Again"
	}

	if len(invoiceHistory) > 0 {
		for _, data := range reqVal.OtherExpenses {
			InsertOtherInvoiceErr := tx.Exec(query.InsertOtherInvoiceAmount, invoiceHistory[0].RefIHId, data.Name, data.Amount, "plus").Error
			if InsertOtherInvoiceErr != nil {
				log.Error(InsertOtherInvoiceErr)
				return false, "Something went wrong, Try Again"
			}
		}

		for _, data := range reqVal.DeductibleExpenses {
			InsertOtherInvoiceErr := tx.Exec(query.InsertOtherInvoiceAmount, invoiceHistory[0].RefIHId, data.Name, data.Amount, "minus").Error
			if InsertOtherInvoiceErr != nil {
				log.Error(InsertOtherInvoiceErr)
				return false, "Something went wrong, Try Again"
			}
		}
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
		log.Error(InvoiceHistoryErr)
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	for i, data := range invoiceHistory {
		invoiceOtherExpensesErr := tx.Raw(query.GetInvoiceOtherAmount, data.RefIHId).Scan(&invoiceHistory[i].OtherExpenses).Error
		if invoiceOtherExpensesErr != nil {
			log.Error(invoiceOtherExpensesErr)
			return []model.InvoiceHistory{}, []model.TakenDate{}
		}
		if len(data.RefIHSignature) > 0 {
			DriversLicenseNoImgHelperData, viewErr := helper.ViewFile("./Assets/Files/" + data.RefIHSignature)
			if viewErr != nil {
				// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
				log.Errorf("Failed to read DrivingLicense file: %v", viewErr)
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
		log.Error(InvoiceHistoryTakenErr)
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.InvoiceHistory{}, []model.TakenDate{}
	}

	return invoiceHistory, invoiceHistoryTakenDate
}

func GetInvoiceOverAllHistoryService(db *gorm.DB, reqVal model.GetInvoiceOverAllHistoryReq, roleIdValue int) []model.InvoiceHistory {
	log := logger.InitLogger()

	//Get OverAll Invoice
	var AllInvoiceModel []model.InvoiceHistory

	var FromDate = ""
	if len(reqVal.FromDate) > 0 {
		FromDate = reqVal.FromDate + " 00:00:00"
	}
	var ToDate = ""
	if len(reqVal.ToDate) > 0 {
		ToDate = reqVal.ToDate + " 23:59:59"
	}

	AllInvoiceErr := db.Raw(query.GetInvoiceOverAllHistorySQL, roleIdValue, FromDate, ToDate).Scan(&AllInvoiceModel).Error
	if AllInvoiceErr != nil {
		log.Error(AllInvoiceErr.Error())
		return []model.InvoiceHistory{}
	}

	for i, data := range AllInvoiceModel {
		if len(data.RefIHSignature) > 0 {
			DriversLicenseNoImgHelperData, viewErr := helper.ViewFile("./Assets/Files/" + data.RefIHSignature)
			if viewErr != nil {
				// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
				log.Errorf("Failed to read DrivingLicense file: %v", viewErr)
			}
			if DriversLicenseNoImgHelperData != nil {
				AllInvoiceModel[i].RefIHSignatureFile = &model.FileData{
					Base64Data:  DriversLicenseNoImgHelperData.Base64Data,
					ContentType: DriversLicenseNoImgHelperData.ContentType,
				}
			}
		} else {
			AllInvoiceModel[i].RefIHSignatureFile = nil
		}
	}

	return AllInvoiceModel

}
