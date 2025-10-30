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
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

func AddTechnicianIntakeFormService(db *gorm.DB, reqVal model.AddTechnicianIntakeFormReq, idValue int) (bool, string) {
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

	// Update Category Id
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
		// Insert Aduit Row for Category ID
		transData := 22

		errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.PatientId), int(idValue), string(ChangesDataJSON)).Error
		if errTrans != nil {
			log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		//Updating the Category ID
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

	//Store the Aduit Row for Answers
	for _, answer := range reqVal.UpdatedAnswers {

		PrevData := model.GetViewIntakeData{}
		errPrev := tx.Raw(query.GetIntakeDataSQL, answer.ITFId).Scan(&PrevData).Error
		if errPrev != nil {
			log.Printf("ERROR: Failed to Get Intake: %v\n", PrevData)
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

		if len(ChangesData) > 0 || answer.TechinicianStatus {
			var ChangesDataJSON []byte
			var errChange error
			ChangesDataJSON, errChange = json.Marshal(ChangesData)
			if errChange != nil {
				// Corrected log message
				log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", errChange)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			// combined := string(ChangesDataJSON) + "," + fmt.Sprintf(`"questionId": %d}`, answer.QuestionId)
			//Insert Aduit Row for Answers Update
			transData := 24

			// errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.PatientId), int(idValue), hashdb.Encrypt(string(ChangesDataJSON))).Error
			// if errTrans != nil {
			// 	log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
			// 	tx.Rollback()
			// 	return false, "Something went wrong, Try Again"
			// }

			errTrans := model.RefTransHistory{
				TransTypeId: transData,
				THData:      hashdb.Encrypt(string(ChangesDataJSON)),
				UserId:      reqVal.PatientId,
				THActionBy:  idValue,
			}

			errTransStatus := db.Create(&errTrans).Error
			if errTransStatus != nil {
				log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
				return false, "Something went wrong, Try Again"
			}

			fmt.Println(answer.ITFId, answer.TechinicianStatus)

		}

		//Updating Answers
		updatedIntakeErr := tx.Exec(
			query.UpdateIntakeDataSQL,
			answer.Answer,
			idValue,
			timeZone.GetPacificTime(),
			answer.TechinicianStatus,
			answer.ITFId,
		).Error

		if updatedIntakeErr != nil {
			log.Printf("ERROR: Failed to UpdatedIntake: %v\n", updatedIntakeErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

	}

	//Preparing the Encryption for the Technician Intake Form
	// for i, answer := range reqVal.TechnicianAnswers {
	// 	reqVal.TechnicianAnswers[i].Answer = hashdb.Encrypt(answer.Answer)
	// }

	// jsonAnswers, err := json.Marshal(reqVal.TechnicianAnswers)
	// if err != nil {
	// 	log.Printf("ERROR: Failed to marshal answers: %v\n", err)
	// 	tx.Rollback()
	// 	return false, "Invalid input format"
	// }

	// //Inserting the Technician Intake Form
	// InsertAnswer := tx.Exec(
	// 	query.TechnicianInsertAnswerSQL,
	// 	reqVal.PatientId,
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

	//Update and Insert the technician intake
	for _, answer := range reqVal.TechnicianAnswers {

		PrevData := model.GetViewTechnicianIntakeData{}
		errPrev := tx.Raw(query.GetTechnicianIntakeDataSQL, answer.QuestionId, reqVal.AppointmentId).Scan(&PrevData).Error
		if errPrev != nil {
			log.Printf("ERROR: Failed to Get Intake: %v\n", PrevData)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		if PrevData.QuestionId != 0 {

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

				// combined := string(ChangesDataJSON) + "," + fmt.Sprintf(`"questionId": %d}`, answer.QuestionId)
				//Insert Aduit Row for Answers Update
				transData := 27

				// errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.PatientId), int(idValue), hashdb.Encrypt(string(ChangesDataJSON))).Error
				// if errTrans != nil {
				// 	log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
				// 	tx.Rollback()
				// 	return false, "Something went wrong, Try Again"
				// }

				errTrans := model.RefTransHistory{
					TransTypeId: transData,
					THData:      hashdb.Encrypt(string(ChangesDataJSON)),
					UserId:      reqVal.PatientId,
					THActionBy:  idValue,
				}

				errTransStatus := db.Create(&errTrans).Error
				if errTransStatus != nil {
					log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
					return false, "Something went wrong, Try Again"
				}

				//Updating Answers
				updatedIntakeErr := tx.Exec(
					query.UpdateTechnicianIntakeDataSQL,
					hashdb.Encrypt(answer.Answer),
					idValue,
					timeZone.GetPacificTime(),
					answer.QuestionId,
					reqVal.AppointmentId,
				).Error

				if updatedIntakeErr != nil {
					log.Printf("ERROR: Failed to UpdatedIntake: %v\n", updatedIntakeErr)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}
			}
		} else {
			fmt.Println("++++++++++++++++++++")
			InsertTechErr := tx.Exec(
				query.InsertTechnicianIntakeDataSQL,
				reqVal.PatientId,
				reqVal.AppointmentId,
				answer.QuestionId,
				answer.Answer,
				timeZone.GetPacificTime(),
				idValue,
			).Error
			if InsertTechErr != nil {
				log.Printf("ERROR: Failed to Insert Answers: %v\n", InsertTechErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}
		}
	}

	// Storing the Technician Aduit Row for Answers

	// var techallChangeLogs []any

	// for _, answer := range reqVal.TechnicianAnswers {

	// 	oldData := map[string]interface{}{
	// 		fmt.Sprintf("%d", answer.QuestionId): "",
	// 	}

	// 	updatedData := map[string]interface{}{
	// 		fmt.Sprintf("%d", answer.QuestionId): answer.Answer,
	// 	}

	// 	ChangesData := helper.GetChanges(updatedData, oldData)

	// 	if len(ChangesData) > 0 {
	// 		var ChangesDataJSON []byte
	// 		var errChange error
	// 		ChangesDataJSON, errChange = json.Marshal(ChangesData)
	// 		if errChange != nil {
	// 			// Corrected log message
	// 			log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", errChange)
	// 			tx.Rollback()
	// 			return false, "Something went wrong, Try Again"
	// 		}

	// 		techallChangeLogs = append(techallChangeLogs, hashdb.Encrypt(string(ChangesDataJSON)))

	// 	}
	// }

	// technfinalJSON, techerr := json.Marshal(techallChangeLogs)
	// if techerr != nil {
	// 	log.Printf("ERROR: Failed to marshal allChangeLogs: %v\n", techerr)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	// transData := 26

	// errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(idValue), int(idValue), string(technfinalJSON)).Error
	// if errTrans != nil {
	// 	log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	if !reqVal.SaveStatus {
		reportStatus := model.RefTransHistory{
			TransTypeId: 26,
			THData:      "Technician Form Filled Successfully",
			UserId:      idValue,
			THActionBy:  idValue,
		}

		errreportStatus := db.Create(&reportStatus).Error
		if errreportStatus != nil {
			log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errreportStatus.Error())
			return false, "Something went wrong, Try Again"
		}

		//Updating a Appointment Status
		UpdateAppointmentStatuserr := tx.Exec(
			query.UpdateTechnicianAppointmentStatus,
			timeZone.GetPacificTimeDateOnly(),
			"reportformfill",
			reqVal.Priority,
			reqVal.ArtificatsLeft,
			reqVal.ArtificatsRight,
			reqVal.AppointmentId,
		).Error
		if UpdateAppointmentStatuserr != nil {
			log.Printf("ERROR: Failed to Update Appointment Status: %v\n", UpdateAppointmentStatuserr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		//Updating Report Hstory
		//Inserting the Audit for the Report Accessing
		HistoryoldDataCat := map[string]interface{}{
			"Report Access ID": idValue,
		}

		HistoryupdatedDataCat := map[string]interface{}{
			"Report Access ID": "",
		}

		HistoryChangesDataCat := helper.GetChanges(HistoryupdatedDataCat, HistoryoldDataCat)

		var HistoryChangesDataJSON []byte
		var HistoryerrChange error
		HistoryChangesDataJSON, HistoryerrChange = json.Marshal(HistoryChangesDataCat)
		if HistoryerrChange != nil {
			// Corrected log message
			log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", HistoryerrChange)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		HistorytransData := 28

		HistoryerrTrans := tx.Exec(query.InsertTransactionDataSQL, int(HistorytransData), int(reqVal.PatientId), int(idValue), string(HistoryChangesDataJSON)).Error
		if HistoryerrTrans != nil {
			log.Printf("ERROR: Failed to Transaction History: %v\n", HistoryerrTrans)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Technician Intake Form Created"
}

func ViewTechnicianIntakeFormService(db *gorm.DB, reqVal model.ViewTechnicianIntakeFormReq, idValue int) ([]model.GetViewIntakeData, []model.AduitModel, []model.TechIntakeModel, string, string) {
	log := logger.InitLogger()

	var ViewIntakeData []model.GetViewIntakeData

	err := db.Raw(query.ViewIntakeFormQuery, reqVal.PatientId, reqVal.AppointmentId).Scan(&ViewIntakeData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetViewIntakeData{}, []model.AduitModel{}, []model.TechIntakeModel{}, "", ""
	}

	// for i, data := range Aduit {
	// 	ViewIntakeData[i] = hashdb.Decrypt(strings.Trim(data.THData, `"`))
	// }

	var Aduit []model.AduitModel

	Aduiterr := db.Raw(query.GetAuditforIntakeForm).Scan(&Aduit).Error
	if Aduiterr != nil {
		log.Printf("ERROR: Failed to fetch Aduit: %v", Aduiterr)
		return []model.GetViewIntakeData{}, []model.AduitModel{}, []model.TechIntakeModel{}, "", ""
	}

	for i, data := range Aduit {
		Aduit[i].THData = hashdb.Decrypt(strings.Trim(data.THData, `"`))
	}

	var TechIntakeData []model.TechIntakeModel

	TechDataerr := db.Raw(query.GetTechIntakeForm, reqVal.AppointmentId, reqVal.PatientId).Scan(&TechIntakeData).Error
	if TechDataerr != nil {
		log.Printf("ERROR: Failed to fetch Tech Aduit: %v", TechDataerr)
		return []model.GetViewIntakeData{}, []model.AduitModel{}, []model.TechIntakeModel{}, "", ""
	}

	for i, data := range TechIntakeData {
		TechIntakeData[i].TITFAnswer = hashdb.Decrypt(data.TITFAnswer)
	}

	name := ""
	custId := ""

	var UserData []model.TechnicianModel

	UserDataerr := db.Raw(query.TechnicianUserSQL, reqVal.AppointmentId).Scan(&UserData).Error
	if UserDataerr != nil {
		log.Printf("ERROR: Failed to fetch Aduit: %v", UserDataerr)
		return []model.GetViewIntakeData{}, []model.AduitModel{}, []model.TechIntakeModel{}, "", ""
	}

	if len(UserData) > 0 {
		name = hashdb.Decrypt(UserData[0].FirstName)
		custId = UserData[0].CustId
	}

	// if len(Aduit) > 0 {

	// 	var UserData []model.TechnicianModel

	// 	UserDataerr := db.Raw(query.TechnicianUserSQL, Aduit[0].THActionBy).Scan(&UserData).Error
	// 	if UserDataerr != nil {
	// 		log.Printf("ERROR: Failed to fetch Aduit: %v", UserDataerr)
	// 		return []model.GetViewIntakeData{}, []model.AduitModel{}, []model.TechIntakeModel{}, "", ""
	// 	}

	// 	name = hashdb.Decrypt(UserData[0].FirstName)
	// 	custId = UserData[0].CustId

	// }

	// fmt.Println("$$$$$$$$$$$$$$$$$$$$", name, custId)

	return ViewIntakeData, Aduit, TechIntakeData, name, custId
}

func AssignTechnicianService(db *gorm.DB, reqVal model.ViewTechnicianIntakeFormReq, idValue int) (bool, string) {
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

	var ListUserData []model.ListUserModel

	ViewTechnicianErr := db.Raw(query.ListTechnicianSQL, reqVal.PatientId, reqVal.AppointmentId, idValue).Scan(&ListUserData).Error
	if ViewTechnicianErr != nil {
		log.Printf("ERROR: Failed to fetch Tech Aduit: %v", ViewTechnicianErr)
		return false, "Something went wrong, Try Again"
	}

	if len(ListUserData) > 0 {
		return true, "UserId Already Assigned"
	}

	oldDataCat := map[string]interface{}{
		"Report Access ID": "",
	}

	updatedDataCat := map[string]interface{}{
		"Report Access ID": idValue,
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

		transData := 28

		errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), reqVal.PatientId, int(idValue), string(ChangesDataJSON)).Error
		if errTrans != nil {
			log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

	}

	//Insert the History
	ReportHistoryErr := tx.Exec(
		query.TechInsertReportHistorySQL,
		reqVal.PatientId,
		reqVal.AppointmentId,
		idValue,
		timeZone.GetPacificTime(),
		"Technologist Form Fill",
	).Error
	if ReportHistoryErr != nil {
		log.Printf("ERROR: Failed to Insert Report History: %v\n", ReportHistoryErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Assigned"
}

func SaveDicomService(db *gorm.DB, reqVal model.SaveDicomReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, try again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	// --- 1️⃣ Check if DICOM already exists ---
	var existingDicoms []model.DicomFileModel
	if err := tx.Raw(query.ViewGetDicomFile, reqVal.AppointmentId, reqVal.PatientId).Scan(&existingDicoms).Error; err != nil {
		tx.Rollback()
		return false, "Something went wrong while checking DICOM files"
	}

	// --- 2️⃣ Get Patient Customer ID ---
	var patientCustId string
	if err := tx.Table(`"Users"`).
		Select(`"refUserCustId"`).
		Where(`"refUserId" = ?`, reqVal.PatientId).
		Scan(&patientCustId).Error; err != nil {
		log.Printf("ERROR: Failed to get patient custom ID: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, try again"
	}
	if patientCustId == "" {
		errName := tx.Table(`"Users"`).
			Select(`"refUserFirstName"`).
			Where(`"refUserId" = ?`, reqVal.PatientId).
			Scan(&patientCustId).Error
		if errName != nil || patientCustId == "" {
			log.Printf("ERROR: Patient custom ID not found for user ID: %d\n", reqVal.PatientId)
			tx.Rollback()
			return false, "Patient not found"
		}
	}

	// --- 3️⃣ Get Scan Center Customer ID ---
	type ScanCenterResult struct {
		RefSCCustId string `gorm:"column:refSCCustId"`
	}
	var scanCenter ScanCenterResult

	if err := tx.Table(`appointment."refAppointments" AS ra`).
		Joins(`JOIN public."ScanCenter" AS sc ON sc."refSCId" = ra."refSCId"`).
		Where(`ra."refAppointmentId" = ?`, reqVal.AppointmentId).
		Scan(&scanCenter).Error; err != nil || scanCenter.RefSCCustId == "" {
		log.Printf("ERROR: Failed to get scan center customer ID: %v\n", err)
		tx.Rollback()
		return false, "Scan center configuration error"
	}

	// --- 4️⃣ Save all DICOM files ---
	for _, file := range reqVal.DicomFiles {
		dicomFile := model.DicomFileModel{
			UserId:        reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			FileName:      file.FilesName, // already renamed and uploaded to S3
			Side:          file.Side,
			CreatedBy:     idValue,
			CreatedAt:     time.Now().In(timeZone.MustGetPacificLocation()),
		}

		if err := tx.Create(&dicomFile).Error; err != nil {
			log.Printf("ERROR: Failed to insert DICOM file: %v\n", err)
			tx.Rollback()
			return false, "Something went wrong while saving DICOM file"
		}
	}

	// --- 5️⃣ Update Report History if no previous DICOMs existed ---
	if len(existingDicoms) == 0 {
		if err := tx.Exec(
			query.CompleteReportHistorySQL,
			timeZone.GetPacificTime(),
			"Technologist Form Fill",
			"",
			reqVal.AppointmentId,
			idValue,
			reqVal.PatientId,
		).Error; err != nil {
			log.Printf("ERROR: Failed to update Report History: %v\n", err)
			tx.Rollback()
			return false, "Something went wrong while updating report history"
		}
	}

	// --- 6️⃣ Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, try again"
	}

	return true, "DICOM files saved successfully"
}

func ViewDicomService(db *gorm.DB, reqVal model.ViewTechnicianIntakeFormReq, idValue int) []model.DicomFileModel {
	log := logger.InitLogger()

	var DicomModel []model.DicomFileModel

	DicomErr := db.Raw(query.ViewGetDicomFile, reqVal.AppointmentId, reqVal.PatientId).Scan(&DicomModel).Error
	if DicomErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", DicomErr)
		return []model.DicomFileModel{}
	}

	return DicomModel
}

func DeleteDicomService(db *gorm.DB, reqVal model.DeleteDicomReq) (bool, string) {

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

	var DicomFiles []model.DicomFileModel
	DicomErr := db.Raw(query.GetListDicomSQL, reqVal.DFId).Scan(&DicomFiles).Error
	if DicomErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", DicomErr)
		return false, "Something went wrong, Try Again"
	}

	DeleteDicomErr := tx.Exec(
		query.DeleteDicomFileSQL,
		reqVal.DFId,
	).Error
	if DeleteDicomErr != nil {
		log.Error(DeleteDicomErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	uploadPath := "./Assets/Dicom/"

	for _, data := range DicomFiles {
		// Skip deletion if the file is an S3 URL
		if strings.HasPrefix(data.FileName, "http") {
			log.Printf("Skipping deletion for S3 file: %s", data.FileName)
			continue
		}

		filePath := filepath.Join(uploadPath, data.FileName)
		if err := os.Remove(filePath); err != nil {
			log.Error("File deletion failed:", err)
			return false, "Something went wrong, Try Again"
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Deleted"
}
