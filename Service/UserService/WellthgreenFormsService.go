package service

import (
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/UserService"
	query "AuthenticationService/query/UserService"
	"fmt"

	"gorm.io/gorm"
)

func ListPatientConsentService(db *gorm.DB, reqVal model.ListPatientConsentReq) []model.ListPatientConsent {
	log := logger.InitLogger()

	var ListPatient []model.ListPatientConsent
	if err := db.Raw(query.ListPatientConsentSQL, reqVal.ScanCenterId).Scan(&ListPatient).Error; err != nil {
		log.Error("Error fetching WG brochure: " + err.Error())
		return []model.ListPatientConsent{}
	}

	return ListPatient
}

func ListPatientBrochureService(db *gorm.DB, reqVal model.ListPatientBrochureReq, wgDataForm int, scDataForm int) model.ListPatientBrochureRes {
	log := logger.InitLogger()

	var res model.ListPatientBrochureRes
	res.Status = true
	res.SCBrochureAccessStatus = true

	// WG Patient Brochure
	var wgBrochures []model.ListWGPatientBrochureModel
	if err := db.Raw(query.ListWGPatientBrochureSQL, wgDataForm).Scan(&wgBrochures).Error; err != nil {
		log.Error("Error fetching WG brochure: " + err.Error())
		res.Status = false
		return res
	}
	if len(wgBrochures) > 0 {
		res.WGPatientBrochure = wgBrochures[0].FDData
	}

	// SC Patient Brochure
	var scBrochures []model.ListWGPatientBrochureModel
	if err := db.Raw(query.ListSCPatientBrochureSQL, scDataForm, reqVal.ScancenterId).Scan(&scBrochures).Error; err != nil {
		log.Error("Error fetching SC brochure: " + err.Error())
		res.Status = false
		return res
	}
	fmt.Println(scBrochures, scDataForm, reqVal.ScancenterId)
	if len(scBrochures) > 0 {
		res.SCBrochureAccessStatus = scBrochures[0].FDAccessData
		res.SCPatientBrochure = scBrochures[0].FDData
	}

	return res
}

func UpdatePatientBrochureService(db *gorm.DB, reqVal model.UpdatePatientBroucherReq, idValue int, wgDataForm int, scDataForm int, wgTransactionId int, scTransactionId int) model.UpdatePatientBroucherResponse {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.UpdatePatientBroucherResponse{
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

	if reqVal.ScancenterId == 0 { //Wellthgreen Patient Brochure need to Handle
		var wgBrochures []model.ListWGPatientBrochureModel
		if err := tx.Raw(query.ListWGPatientBrochureSQL, wgDataForm).Scan(&wgBrochures).Error; err != nil {
			log.Error("Error fetching WG brochure: " + err.Error())
			return model.UpdatePatientBroucherResponse{
				Status:  false,
				Message: "Something Went Wrong, Try Again",
			}
		}

		if len(wgBrochures) > 0 { //Need to Update the Brochure
			UpdateWGBrochureErr := tx.Exec(query.UpdateWGPatientBrochureSQL, reqVal.Brochure, wgDataForm).Error
			if UpdateWGBrochureErr != nil {
				log.Error("Error updating WG brochure: " + UpdateWGBrochureErr.Error())
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}

			errTrans := tx.Exec(query.InsertTransactionDataSQL, wgTransactionId, wgDataForm, int(idValue), `["Wellthgreen Patient Brochure Updated!"]`).Error
			if errTrans != nil {
				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
				tx.Rollback()
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}

		} else { //Need to Add the Brochure
			AddWGBrochureErr := tx.Exec(query.InsertWGPatientBrochureSQL, wgDataForm, reqVal.Brochure).Error
			if AddWGBrochureErr != nil {
				log.Error("Error adding WG brochure: " + AddWGBrochureErr.Error())
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}

			errTrans := tx.Exec(query.InsertTransactionDataSQL, wgTransactionId, wgDataForm, int(idValue), `["Wellthgreen Patient Brochure Added!"]`).Error
			if errTrans != nil {
				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
				tx.Rollback()
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}
		}

	} else { //Scan Center Patient Brochure need to Handle
		var scBrochures []model.ListWGPatientBrochureModel
		if err := tx.Raw(query.ListSCPatientBrochureSQL, scDataForm, reqVal.ScancenterId).Scan(&scBrochures).Error; err != nil {
			log.Error("Error fetching SC brochure: " + err.Error())
			return model.UpdatePatientBroucherResponse{
				Status:  false,
				Message: "Something Went Wrong, Try Again",
			}
		}

		if len(scBrochures) > 0 { //Need to Update the Brochure
			UpdatescBrochureErr := tx.Exec(
				query.UpdatescPatientBrochureSQL,
				reqVal.Brochure,
				reqVal.AccessStatus,
				scDataForm,
				reqVal.ScancenterId,
			).Error
			if UpdatescBrochureErr != nil {
				log.Error("Error updating SC brochure: " + UpdatescBrochureErr.Error())
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}

			errTrans := tx.Exec(query.InsertTransactionDataSQL, scTransactionId, reqVal.ScancenterId, int(idValue), `["Scan Center Patient Brochure Updated!"]`).Error
			if errTrans != nil {
				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
				tx.Rollback()
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}
		} else { // Need to Add the Brochure
			AddscBrochureErr := tx.Exec(
				query.InsertscPatientBrochureSQL,
				scDataForm,
				reqVal.ScancenterId,
				reqVal.Brochure,
				reqVal.AccessStatus,
			).Error
			if AddscBrochureErr != nil {
				log.Error("Error adding SC brochure: " + AddscBrochureErr.Error())
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}

			errTrans := tx.Exec(query.InsertTransactionDataSQL, scTransactionId, reqVal.ScancenterId, int(idValue), `["Scan Center Patient Brochure Updated!"]`).Error
			if errTrans != nil {
				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
				tx.Rollback()
				return model.UpdatePatientBroucherResponse{
					Status:  false,
					Message: "Something Went Wrong, Try Again",
				}
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return model.UpdatePatientBroucherResponse{
			Status:  false,
			Message: "Something Went Wrong, Try Again",
		}
	}

	return model.UpdatePatientBroucherResponse{
		Status:  true,
		Message: "Successfully Updated the Brochure",
	}
}

// func ListPatientConsentService(db *gorm.DB, reqVal model.ListPatientBrochureReq) model.ListPatientBrochureRes {
// 	log := logger.InitLogger()

// 	var res model.ListPatientBrochureRes
// 	res.Status = true

// 	// WG Patient Brochure
// 	var wgBrochures []model.ListWGPatientBrochureModel
// 	if err := db.Raw(query.ListWGPatientBrochureSQL, 3).Scan(&wgBrochures).Error; err != nil {
// 		log.Error("Error fetching WG brochure: " + err.Error())
// 		res.Status = false
// 		return res
// 	}
// 	if len(wgBrochures) > 0 {
// 		res.WGPatientBrochure = wgBrochures[0].FDData
// 	}

// 	// SC Patient Brochure
// 	var scBrochures []model.ListWGPatientBrochureModel
// 	if err := db.Raw(query.ListSCPatientBrochureSQL, 4, reqVal.ScancenterId).Scan(&scBrochures).Error; err != nil {
// 		log.Error("Error fetching SC brochure: " + err.Error())
// 		res.Status = false
// 		return res
// 	}
// 	if len(scBrochures) > 0 {
// 		res.SCBrochureAccessStatus = scBrochures[0].FDAccessData
// 		res.SCPatientBrochure = scBrochures[0].FDData
// 	}

// 	return res
// }

// func UpdatePatientConsentService(db *gorm.DB, reqVal model.UpdatePatientBroucherReq, idValue int) model.UpdatePatientBroucherResponse {
// 	log := logger.InitLogger()

// 	tx := db.Begin()
// 	if tx.Error != nil {
// 		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
// 		return model.UpdatePatientBroucherResponse{
// 			Status:  false,
// 			Message: "Something Went Wrong, Try Again",
// 		}
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
// 			tx.Rollback()
// 		}
// 	}()

// 	if reqVal.ScancenterId == 0 { //Wellthgreen Patient Brochure need to Handle
// 		var wgBrochures []model.ListWGPatientBrochureModel
// 		if err := tx.Raw(query.ListWGPatientBrochureSQL, 3).Scan(&wgBrochures).Error; err != nil {
// 			log.Error("Error fetching WG brochure: " + err.Error())
// 			return model.UpdatePatientBroucherResponse{
// 				Status:  false,
// 				Message: "Something Went Wrong, Try Again",
// 			}
// 		}

// 		if len(wgBrochures) > 0 { //Need to Update the Brochure
// 			UpdateWGBrochureErr := tx.Exec(query.UpdateWGPatientBrochureSQL, reqVal.Brochure, 3).Error
// 			if UpdateWGBrochureErr != nil {
// 				log.Error("Error updating WG brochure: " + UpdateWGBrochureErr.Error())
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}

// 			errTrans := tx.Exec(query.InsertTransactionDataSQL, 38, 3, int(idValue), "Wellthgreen Patient Consent Updated!").Error
// 			if errTrans != nil {
// 				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
// 				tx.Rollback()
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}

// 		} else { //Need to Add the Brochure
// 			AddWGBrochureErr := tx.Exec(query.InsertWGPatientBrochureSQL, 3, reqVal.Brochure).Error
// 			if AddWGBrochureErr != nil {
// 				log.Error("Error adding WG brochure: " + AddWGBrochureErr.Error())
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}

// 			errTrans := tx.Exec(query.InsertTransactionDataSQL, 38, 3, int(idValue), "Wellthgreen Patient Consent Added!").Error
// 			if errTrans != nil {
// 				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
// 				tx.Rollback()
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}
// 		}

// 	} else { //Scan Center Patient Brochure need to Handle
// 		var scBrochures []model.ListWGPatientBrochureModel
// 		if err := tx.Raw(query.ListSCPatientBrochureSQL, 4, reqVal.ScancenterId).Scan(&scBrochures).Error; err != nil {
// 			log.Error("Error fetching SC brochure: " + err.Error())
// 			return model.UpdatePatientBroucherResponse{
// 				Status:  false,
// 				Message: "Something Went Wrong, Try Again",
// 			}
// 		}

// 		if len(scBrochures) > 0 { //Need to Update the Brochure
// 			UpdatescBrochureErr := tx.Exec(
// 				query.UpdatescPatientBrochureSQL,
// 				reqVal.Brochure,
// 				reqVal.AccessStatus,
// 				4,
// 				reqVal.ScancenterId,
// 			).Error
// 			if UpdatescBrochureErr != nil {
// 				log.Error("Error updating SC brochure: " + UpdatescBrochureErr.Error())
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}

// 			errTrans := tx.Exec(query.InsertTransactionDataSQL, 39, reqVal.ScancenterId, int(idValue), "Scan Center Patient Consent Updated!").Error
// 			if errTrans != nil {
// 				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
// 				tx.Rollback()
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}
// 		} else { // Need to Add the Brochure
// 			AddscBrochureErr := tx.Exec(
// 				query.InsertscPatientBrochureSQL,
// 				4,
// 				reqVal.ScancenterId,
// 				reqVal.Brochure,
// 				reqVal.AccessStatus,
// 			).Error
// 			if AddscBrochureErr != nil {
// 				log.Error("Error adding SC brochure: " + AddscBrochureErr.Error())
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}

// 			errTrans := tx.Exec(query.InsertTransactionDataSQL, 39, reqVal.ScancenterId, int(idValue), "Scan Center Patient Consent Added!").Error
// 			if errTrans != nil {
// 				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
// 				tx.Rollback()
// 				return model.UpdatePatientBroucherResponse{
// 					Status:  false,
// 					Message: "Something Went Wrong, Try Again",
// 				}
// 			}
// 		}

// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
// 		tx.Rollback()
// 		return model.UpdatePatientBroucherResponse{
// 			Status:  false,
// 			Message: "Something Went Wrong, Try Again",
// 		}
// 	}

// 	return model.UpdatePatientBroucherResponse{
// 		Status:  true,
// 		Message: "Successfully Updated the Brochure",
// 	}
// }
