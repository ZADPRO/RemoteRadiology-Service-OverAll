package service

import (
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"
	"encoding/json"
	"errors"
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
	for i, answer := range reqVal.TechnicianAnswers {
		reqVal.TechnicianAnswers[i].Answer = hashdb.Encrypt(answer.Answer)
	}

	jsonAnswers, err := json.Marshal(reqVal.TechnicianAnswers)
	if err != nil {
		log.Printf("ERROR: Failed to marshal answers: %v\n", err)
		tx.Rollback()
		return false, "Invalid input format"
	}

	//Inserting the Technician Intake Form
	InsertAnswer := tx.Exec(
		query.TechnicianInsertAnswerSQL,
		reqVal.PatientId,
		reqVal.AppointmentId,
		idValue,
		timeZone.GetPacificTime(),
		string(jsonAnswers),
	).Error
	if InsertAnswer != nil {
		log.Printf("ERROR: Failed to Insert Answers: %v\n", InsertAnswer)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	// Storing the Technician Aduit Row for Answers

	var techallChangeLogs []any

	for _, answer := range reqVal.TechnicianAnswers {

		oldData := map[string]interface{}{
			fmt.Sprintf("%d", answer.QuestionId): "",
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

			techallChangeLogs = append(techallChangeLogs, hashdb.Encrypt(string(ChangesDataJSON)))

		}
	}

	technfinalJSON, techerr := json.Marshal(techallChangeLogs)
	if techerr != nil {
		log.Printf("ERROR: Failed to marshal allChangeLogs: %v\n", techerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	transData := 26

	errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(idValue), int(idValue), string(technfinalJSON)).Error
	if errTrans != nil {
		log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	reportStatus := model.RefTransHistory{
		TransTypeId: 25,
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
		"reportformfill",
		reqVal.Priority,
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

	if len(Aduit) > 0 {

		var UserData []model.TechnicianModel

		UserDataerr := db.Raw(query.TechnicianUserSQL, Aduit[0].THActionBy).Scan(&UserData).Error
		if UserDataerr != nil {
			log.Printf("ERROR: Failed to fetch Aduit: %v", UserDataerr)
			return []model.GetViewIntakeData{}, []model.AduitModel{}, []model.TechIntakeModel{}, "", ""
		}

		name = hashdb.Decrypt(UserData[0].FirstName)
		custId = UserData[0].CustId

	}

	fmt.Println("$$$$$$$$$$$$$$$$$$$$", name, custId)

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
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	// Get Patient Customer ID
	var patientCustId string
	errPatient := tx.Table("\"Users\"").Select("\"refUserCustId\"").Where("\"refUserId\" = ?", reqVal.PatientId).Scan(&patientCustId).Error
	if errPatient != nil {
		log.Printf("ERROR: Failed to get patient custom ID: %v\n", errPatient)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}
	if patientCustId == "" {
		log.Printf("ERROR: Patient custom ID not found for user ID: %d\n", reqVal.PatientId)
		tx.Rollback()
		return false, "Patient not found"
	}

	// Get Scan Center Customer ID from Technician ID
	// Define a struct to hold the result
	type ScanCenterResult struct {
		RefSCCustId string `gorm:"column:refSCCustId"`
	}

	var scanCenterCustId ScanCenterResult

	// Perform the join query
	errSC := tx.Table("appointment.\"refAppointments\" AS ra").
		Joins("JOIN public.\"ScanCenter\" AS sc ON sc.\"refSCId\" = ra.\"refSCId\"").
		Where("ra.\"refAppointmentId\" = ?", reqVal.AppointmentId).
		Scan(&scanCenterCustId).Error

	// Error handling
	if errSC != nil {
		if errors.Is(errSC, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: No active scan center mapping found for technician ID: %d\n", idValue)
			tx.Rollback()
			return false, "Technician not mapped to an active scan center"
		}
		log.Printf("ERROR: Failed to get scan center customer ID: %v\n", errSC)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if scanCenterCustId.RefSCCustId == "" {
		log.Printf("ERROR: Scan center customer ID is empty for technician ID: %d\n", idValue)
		tx.Rollback()
		return false, "Scan center configuration error"
	}

	//Handle Dicom File Store Process
	currentDate := timeZone.GetTimeWithFormate("02-01-2006")
	for i, file := range reqVal.DicomFiles {
		// Get the file extension
		ext := filepath.Ext(file.FilesName)
		if ext == "" {
			ext = ".zip" // Default to .zip if no extension
		}

		side := "R"

		if file.Side == "Left" {
			side = "L"
		}

		// Construct the new filename
		newFilename := fmt.Sprintf("%s_%s_%s_%s_%d%s",
			scanCenterCustId.RefSCCustId,
			strings.ToUpper(patientCustId),
			currentDate,
			side,
			i+1,
			ext,
		)

		oldPath := filepath.Join("./Assets/Dicom/", file.FilesName)
		newPath := filepath.Join("./Assets/Dicom/", newFilename)

		// Ensure old file exists
		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			log.Printf("ERROR: Source file does not exist: %s\n", oldPath)
			tx.Rollback()
			return false, "DICOM source file not found"
		}

		// Try renaming
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("ERROR: Failed to rename DICOM file from %s to %s: %v\n", oldPath, newPath, err)
			tx.Rollback()
			return false, "Failed to process DICOM file"
		}

		DicomFile := model.DicomFileModel{
			UserId:        reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			FileName:      newFilename,
			CreatedAt:     time.Now().In(timeZone.MustGetPacificLocation()),
			CreatedBy:     idValue,
			Side:          file.Side,
		}

		DicomFileerr := db.Create(&DicomFile).Error
		if DicomFileerr != nil {
			log.Error("DicomFile INSERT ERROR at Technician Intake: " + DicomFileerr.Error())
			return false, "Something went wrong, Try Again"
		}

	}

	//Updating the End Time For the Report History
	ReportHistoryErr := tx.Exec(
		query.CompleteReportHistorySQL,
		timeZone.GetPacificTime(),
		"technologistformfill",
		"",
		reqVal.AppointmentId,
		idValue,
		reqVal.PatientId,
	).Error
	if ReportHistoryErr != nil {
		log.Printf("ERROR: Failed to Update Report History: %v\n", ReportHistoryErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Dicom File Saved"
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
