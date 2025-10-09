package service

import (
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func AddIntakeFormService(db *gorm.DB, reqVal model.AddIntakeFormReq, idValue int) (bool, string) {
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

	// var allChangeLogs []any

	for _, answer := range reqVal.Answers {

		PrevData := []model.GetViewIntakeData{}
		errPrev := tx.Raw(query.GetIntakeAppointmentDataSQL, reqVal.AppointmentId, answer.QuestionId).Scan(&PrevData).Error
		if errPrev != nil {
			log.Printf("ERROR: Failed to Get Intake: %v\n", errPrev)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		value := ""
		if len(PrevData) != 0 {
			value = hashdb.Decrypt(PrevData[0].Answer)
		}

		oldData := map[string]interface{}{
			fmt.Sprintf("%d", answer.QuestionId): value,
		}

		updatedData := map[string]interface{}{
			fmt.Sprintf("%d", answer.QuestionId): answer.Answer,
		}

		ChangesData := helper.GetChanges(updatedData, oldData)

		if len(PrevData) == 0 {
			//Inserting New Question
			InsertIntakeErr := tx.Exec(
				query.InsertIntakeSQL,
				idValue,
				reqVal.AppointmentId,
				answer.QuestionId,
				hashdb.Encrypt(answer.Answer),
				timeZone.GetPacificTime(),
				idValue,
			).Error
			if InsertIntakeErr != nil {
				log.Printf("ERROR: Failed to InsertIntake: %v\n", InsertIntakeErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}
		} else {
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

				transData := 24

				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(idValue), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

				fmt.Println("Updateing the Intake Data", string(ChangesDataJSON))
				updatedIntakeErr := tx.Exec(
					query.UpdateCreateIntakeDataSQL,
					hashdb.Encrypt(answer.Answer),
					idValue,
					timeZone.GetPacificTime(),
					reqVal.AppointmentId,
					answer.QuestionId,
				).Error

				if updatedIntakeErr != nil {
					log.Printf("ERROR: Failed to UpdatedIntake: %v\n", updatedIntakeErr)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

			}
		}

	}

	// finalJSON, err := json.Marshal(allChangeLogs)
	// if err != nil {
	// 	log.Printf("ERROR: Failed to marshal allChangeLogs: %v\n", err)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	// transData := 23

	// errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(idValue), int(idValue), string(finalJSON)).Error
	// if errTrans != nil {
	// 	log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	// jsonAnswers, err := json.Marshal(reqVal.Answers)
	// if err != nil {
	// 	log.Printf("ERROR: Failed to marshal answers: %v\n", err)
	// 	tx.Rollback()
	// 	return false, "Invalid input format"
	// }

	// InsertAnswer := tx.Exec(
	// 	query.InsertAnswerSQL,
	// 	idValue,
	// 	reqVal.AppointmentId,
	// 	idValue,
	// 	timeZone.GetPacificTime(),
	// 	string(jsonAnswers),
	// ).Error
	// if InsertAnswer != nil {
	// 	log.Printf("ERROR: Failed to Insert Answers: %v\n", InsertAnswer)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	UpdateAppointmenterr := tx.Exec(
		query.UpdateAppointment,
		reqVal.CategoryId,
		reqVal.Consent,
		reqVal.AppointmentId,
	).Error
	if UpdateAppointmenterr != nil {
		log.Printf("ERROR: Failed to Update Appointment: %v\n", UpdateAppointmenterr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	var CheckOverride []model.OverrideRequestModel

	var CheckOverrideErr = tx.Raw(query.CheckOverride, reqVal.AppointmentId).Scan(&CheckOverride).Error
	if CheckOverrideErr != nil {
		log.Printf("ERROR: Failed to Check Override: %v\n", CheckOverrideErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	var UpdateAppointmentStatus = "technologistformfill"
	var oversideStatus = "success"

	if len(CheckOverride) == 1 && CheckOverride[0].ApprovedStatus == "approved" {

		UpdateAppointmentStatuserr := tx.Exec(
			query.UpdateAppointmentStatus,
			UpdateAppointmentStatus,
			reqVal.AppointmentId,
		).Error
		if UpdateAppointmentStatuserr != nil {
			log.Printf("ERROR: Failed to Update Appointment Status: %v\n", UpdateAppointmentStatuserr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		history := model.RefTransHistory{
			TransTypeId: 23,
			THData:      "Intake Form Created Successfully",
			UserId:      idValue,
			THActionBy:  idValue,
		}

		errhistory := db.Create(&history).Error
		if errhistory != nil {
			log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
			return false, "Something went wrong, Try Again"
		}

		ReportHistoryErr := tx.Exec(
			query.ReportHistorySQL,
			idValue,
			reqVal.AppointmentId,
			idValue,
			reqVal.PatientIntakeStartTime,
			timeZone.GetPacificTime(),
			"Patient Intake Override Form Fill",
		).Error
		if ReportHistoryErr != nil {
			log.Printf("ERROR: Failed to Insert Report History: %v\n", ReportHistoryErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

	} else {

		if reqVal.OverrideRequest {
			UpdateAppointmentStatus = "noteligible"
			oversideStatus = "pending"
		}
		override := model.OverrideRequestModel{
			UserId:         idValue,
			AppointmentId:  reqVal.AppointmentId,
			ApprovedStatus: oversideStatus,
		}

		overrideerr := db.Create(&override).Error
		if overrideerr != nil {
			log.Error("LoginService INSERT ERROR at Trnasaction: " + overrideerr.Error())
			return false, "Something went wrong, Try Again"
		}

		UpdateAppointmentStatuserr := tx.Exec(
			query.UpdateAppointmentStatus,
			UpdateAppointmentStatus,
			reqVal.AppointmentId,
		).Error
		if UpdateAppointmentStatuserr != nil {
			log.Printf("ERROR: Failed to Update Appointment Status: %v\n", UpdateAppointmentStatuserr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		if reqVal.OverrideRequest {
			history := model.RefTransHistory{
				TransTypeId: 23,
				THData:      "Applied Override Request",
				UserId:      idValue,
				THActionBy:  idValue,
			}

			errhistory := db.Create(&history).Error
			if errhistory != nil {
				log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
				return false, "Something went wrong, Try Again"
			}
		} else {
			history := model.RefTransHistory{
				TransTypeId: 23,
				THData:      "Intake Form Created Successfully",
				UserId:      idValue,
				THActionBy:  idValue,
			}

			errhistory := db.Create(&history).Error
			if errhistory != nil {
				log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
				return false, "Something went wrong, Try Again"
			}
		}

		reportStatus := model.RefTransHistory{
			TransTypeId: 25,
			THData:      "Patient Intake Filled Successfully",
			UserId:      idValue,
			THActionBy:  idValue,
		}

		errreportStatus := db.Create(&reportStatus).Error
		if errreportStatus != nil {
			log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errreportStatus.Error())
			return false, "Something went wrong, Try Again"
		}

		ReportHistoryErr := tx.Exec(
			query.ReportHistorySQL,
			idValue,
			reqVal.AppointmentId,
			idValue,
			reqVal.PatientIntakeStartTime,
			timeZone.GetPacificTime(),
			"Patient Intake Form Fill",
		).Error
		if ReportHistoryErr != nil {
			log.Printf("ERROR: Failed to Insert Report History: %v\n", ReportHistoryErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Intake Form Created"
}

func ViewIntakeService(db *gorm.DB, reqVal model.ViewIntakeReq) ([]model.GetViewIntakeData, []model.AduitModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return []model.GetViewIntakeData{}, []model.AduitModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var ViewIntakeData []model.GetViewIntakeData

	err := db.Raw(query.ViewIntakeFormQuery, reqVal.UserId, reqVal.AppointmentId).Scan(&ViewIntakeData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetViewIntakeData{}, []model.AduitModel{}
	}

	var Aduit []model.AduitModel

	Aduiterr := db.Raw(query.GetAuditforIntakeForm).Scan(&Aduit).Error
	if Aduiterr != nil {
		log.Printf("ERROR: Failed to fetch Aduit: %v", Aduiterr)
		return []model.GetViewIntakeData{}, []model.AduitModel{}
	}

	for i, data := range Aduit {
		Aduit[i].THData = hashdb.Decrypt(strings.Trim(data.THData, `"`))
	}

	for i, data := range ViewIntakeData {
		ViewIntakeData[i].Answer = hashdb.Decrypt(data.Answer)
		// if data.QuestionId == 128 || data.QuestionId == 137 || data.QuestionId == 138 || data.QuestionId == 143 || data.QuestionId == 148 || data.QuestionId == 153 || data.QuestionId == 158 || data.QuestionId == 165 {
		// 	if len(hashdb.Decrypt(data.Answer)) > 0 {
		// 		FilesData, viewErr := helperfile.ViewFile("./Assets/Files/" + hashdb.Decrypt(data.Answer))
		// 		if viewErr != nil {
		// 			// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
		// 			log.Errorf("Failed to read profile image file: %v", viewErr)
		// 		}
		// 		if FilesData != nil {
		// 			ViewIntakeData[i].File = &model.FileData{
		// 				Base64Data:  FilesData.Base64Data,
		// 				ContentType: FilesData.ContentType,
		// 			}
		// 		}
		// 	} else {
		// 		ViewIntakeData[i].File = nil
		// 	}
		// }
	}

	var PatientData []model.PatientCustId

	PatientDataErr := tx.Raw(query.GetPatientData, reqVal.UserId, reqVal.AppointmentId).Scan(&PatientData).Error
	if PatientDataErr != nil {
		log.Printf("ERROR: Failed to fetch Patient Data: %v", PatientDataErr)
		return []model.GetViewIntakeData{}, []model.AduitModel{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.GetViewIntakeData{}, []model.AduitModel{}
	}

	return ViewIntakeData, Aduit
}

func VerifyIntakeFormService(db *gorm.DB, reqVal model.VerifyIntakeFormReq) (bool, []model.OverrideRequestModel) {
	log := logger.InitLogger()

	var OverrideData []model.OverrideRequestModel

	OverrideDataerr := db.Raw(query.GetVerifyIntakeFormQuery, reqVal.AppointmentId).Scan(&OverrideData).Error
	if OverrideDataerr != nil {
		log.Printf("ERROR: Failed to fetch Get Verify Data: %v", OverrideDataerr)
		return false, []model.OverrideRequestModel{}
	}

	if len(OverrideData) == 0 {
		return false, []model.OverrideRequestModel{}
	} else {
		return true, OverrideData
	}
}

func UpdateIntakeFormService(db *gorm.DB, reqVal model.UpdateIntakeFormReq, idValue int) (bool, string) {
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

	PrevDataCat := model.GetCategoryIdModel{}
	errPrevCat := tx.Raw(query.GetCategoryId, reqVal.AppointmentId, reqVal.PatientId).Scan(&PrevDataCat).Error
	if errPrevCat != nil {
		log.Printf("ERROR: Failed to Get Category: %v\n", PrevDataCat)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	oldDataCat := map[string]interface{}{
		"Category ID": PrevDataCat.CategoryId,
	}

	updatedDataCat := map[string]interface{}{
		"Category ID": reqVal.CategoryId,
	}

	ChangesDataCat := helper.GetChanges(updatedDataCat, oldDataCat)

	if len(ChangesDataCat) > 0 {
		var ChangesDataJSON []byte
		var errChange error
		ChangesDataJSON, errChange = json.Marshal(ChangesDataCat)
		if errChange != nil {
			// Corrected log message
			log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", errChange)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		// combined := string(ChangesDataJSON) + "," + fmt.Sprintf(`"questionId": %d}`, answer.QuestionId)

		transData := 24

		errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.PatientId), int(idValue), string(ChangesDataJSON)).Error
		if errTrans != nil {
			log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		fmt.Println("Updateing the Category Data", string(ChangesDataJSON))

		categoryUpdate := tx.Exec(
			query.UpdateCategoryId,
			reqVal.CategoryId,
			reqVal.AppointmentId,
		).Error
		if categoryUpdate != nil {
			log.Printf("ERROR: Failed toCategory Id: %v\n", categoryUpdate)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	for _, answer := range reqVal.Answers {
		PrevData := model.GetViewIntakeData{}
		errPrev := tx.Raw(query.GetIntakeDataSQL, answer.ITFId).Scan(&PrevData).Error
		if errPrev != nil {
			log.Printf("ERROR: Failed to Get Intake: %v\n", errPrev)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		oldData := map[string]interface{}{
			fmt.Sprintf("%d", answer.QuestionId): hashdb.Decrypt(PrevData.Answer),
		}

		updatedData := map[string]interface{}{
			fmt.Sprintf("%d", answer.QuestionId): answer.Answer,
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

			transData := 24

			errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.UserId), int(idValue), string(ChangesDataJSON)).Error
			if errTrans != nil {
				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			fmt.Println("Updateing the Intake Data", string(ChangesDataJSON))
			updatedIntakeErr := tx.Exec(
				query.UpdateIntakeDataSQL,
				hashdb.Encrypt(answer.Answer),
				idValue,
				timeZone.GetPacificTime(),
				answer.VerifiedTechnician,
				answer.ITFId,
			).Error

			if updatedIntakeErr != nil {
				log.Printf("ERROR: Failed to UpdatedIntake: %v\n", updatedIntakeErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

		}

	}

	reportStatus := model.RefTransHistory{
		TransTypeId: 25,
		THData:      "Patient Intake Updated Successfully",
		UserId:      reqVal.UserId,
		THActionBy:  idValue,
	}

	errreportStatus := db.Create(&reportStatus).Error
	if errreportStatus != nil {
		log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errreportStatus.Error())
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Intakeform Updated from Technician"

}

func GetReportDataService(db *gorm.DB, reqVal model.GetViewReportReq) []model.PatientResponse {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return []model.PatientResponse{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var TextContentModel []model.PatientResponse

	err := db.Raw(query.GetTextContent, reqVal.AppointmentId).Scan(&TextContentModel).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.PatientResponse{}
	}

	for i, data := range TextContentModel {
		TextContentModel[i].RTCText = hashdb.Decrypt(data.RTCText)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.PatientResponse{}
	}

	return TextContentModel
}

func GetConsentDataService(db *gorm.DB, reqVal model.GetViewReportReq) []model.GetViewConsentResponse {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return []model.GetViewConsentResponse{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var TextContentModel []model.GetViewConsentResponse

	err := db.Raw(query.GetAppointmentConsent, reqVal.AppointmentId).Scan(&TextContentModel).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetViewConsentResponse{}
	}

	// for i, data := range TextContentModel {
	// 	TextContentModel[i].RTCText = hashdb.Decrypt(data.RTCText)
	// }

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.GetViewConsentResponse{}
	}

	return TextContentModel
}

func AllowOverrideService(db *gorm.DB, reqVal model.GetOverRideReportReq, idValue int) (bool, string) {
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

	UpdateAppointmentStatuserr := tx.Exec(
		query.UpdateAppointmentStatus,
		"fillform",
		reqVal.AppointmentId,
	).Error
	if UpdateAppointmentStatuserr != nil {
		log.Printf("ERROR: Failed to Update Appointment Status: %v\n", UpdateAppointmentStatuserr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	UpdateOverideERr := tx.Exec(
		query.UpdateOverrideSQL,
		"approved",
		idValue,
		timeZone.GetPacificTime(),
		reqVal.AppointmentId,
	)

	if UpdateOverideERr.Error != nil {
		log.Printf("ERROR: Failed to Update Override: %v\n", UpdateOverideERr.Error)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succesfully Override Request Approved"
}
