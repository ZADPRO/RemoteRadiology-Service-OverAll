package service

import (
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
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
		transData := 24

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
			transData := 24

			errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.PatientId), int(idValue), hashdb.Encrypt(string(ChangesDataJSON))).Error
			if errTrans != nil {
				log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

			//Updating Answers
			updatedIntakeErr := tx.Exec(
				query.UpdateIntakeDataSQL,
				answer.Answer,
				idValue,
				answer.TechinicianStatus,
				answer.ITFId,
			).Error

			if updatedIntakeErr != nil {
				log.Printf("ERROR: Failed to UpdatedIntake: %v\n", updatedIntakeErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

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
		idValue,
		reqVal.AppointmentId,
		idValue,
		string(jsonAnswers),
	).Error
	if InsertAnswer != nil {
		log.Printf("ERROR: Failed to Insert Answers: %v\n", InsertAnswer)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

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
	var scanCenterCustId string
	errSC := tx.Table("map.refScanCenterMap as scm").
		Select("sc.\"refSCCustId\"").
		Joins("join \"public\".\"ScanCenter\" as sc on sc.refSCId = scm.refSCId").
		Where("scm.\"refUserId\" = ? AND scm.\"refSCMStatus\" = ?", idValue, true).
		First(&scanCenterCustId).Error

	if errSC != nil {
		if errSC == gorm.ErrRecordNotFound {
			log.Printf("ERROR: No active scan center mapping found for technician ID: %d\n", idValue)
			tx.Rollback()
			return false, "Technician not mapped to an active scan center"
		}
		log.Printf("ERROR: Failed to get scan center custom ID: %v\n", errSC)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}
	if scanCenterCustId == "" {
		log.Printf("ERROR: Scan center custom ID is empty for technician ID: %d\n", idValue)
		tx.Rollback()
		return false, "Scan center configuration error"
	}

	//Handle Dicom File Store Process
	currentDate := time.Now().Format("02-01-2006")
	for i, file := range reqVal.DicomFiles {
		// Get the file extension
		ext := filepath.Ext(file.FilesName)
		if ext == "" {
			ext = ".zip" // Default to .zip if no extension
		}

		// Construct the new filename
		newFilename := fmt.Sprintf("%s_%s_%s_%d%s",
			scanCenterCustId,
			currentDate,
			strings.ToUpper(patientCustId),
			i+1,
			ext)

		oldPath := filepath.Join("./Assets/Dicom/", file.FilesName)
		newPath := filepath.Join("./Assets/Dicom/", newFilename)

		// Rename the file
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("ERROR: Failed to rename DICOM file from %s to %s: %v\n", oldPath, newPath, err)
			tx.Rollback()
			return false, "Failed to process DICOM file"
		}

		DicomFile := model.DicomFileModel{
			UserId:        reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			FileName:      newFilename,
			CreatedAt:     time.Now(),
		}

		DicomFileerr := db.Create(&DicomFile).Error
		if DicomFileerr != nil {
			log.Error("DicomFile INSERT ERROR at Technician Intake: " + DicomFileerr.Error())
			return false, "Something went wrong, Try Again"
		}

	}

	//Updating a Appointment Status
	UpdateAppointmentStatuserr := tx.Exec(
		query.UpdateTechnicianAppointmentStatus,
		"doctorreview",
		reqVal.Priority,
		reqVal.AppointmentId,
	).Error
	if UpdateAppointmentStatuserr != nil {
		log.Printf("ERROR: Failed to Update Appointment Status: %v\n", UpdateAppointmentStatuserr)
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

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Technician Intake Form Created"
}
