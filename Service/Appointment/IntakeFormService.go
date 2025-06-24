package service

import (
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"
	"encoding/json"
	"fmt"

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

	for i, answer := range reqVal.Answers {
		reqVal.Answers[i].Answer = hashdb.Encrypt(answer.Answer)
	}

	jsonAnswers, err := json.Marshal(reqVal.Answers)
	if err != nil {
		log.Printf("ERROR: Failed to marshal answers: %v\n", err)
		tx.Rollback()
		return false, "Invalid input format"
	}

	InsertAnswer := tx.Exec(
		query.InsertAnswerSQL,
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

	UpdateAppointmenterr := tx.Exec(
		query.UpdateAppointment,
		reqVal.CategoryId,
		reqVal.AppointmentId,
	).Error
	if UpdateAppointmenterr != nil {
		log.Printf("ERROR: Failed to Update Appointment: %v\n", UpdateAppointmenterr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if reqVal.OverrideRequest {
		override := model.OverrideRequestModel{
			UserId:         idValue,
			AppointmentId:  reqVal.AppointmentId,
			ApprovedStatus: "pending",
		}

		overrideerr := db.Create(&override).Error
		if overrideerr != nil {
			log.Error("LoginService INSERT ERROR at Trnasaction: " + overrideerr.Error())
			return false, "Something went wrong, Try Again"
		}

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

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Intake Form Created"
}

func ViewIntakeService(db *gorm.DB, reqVal model.ViewIntakeReq) []model.GetViewIntakeData {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return []model.GetViewIntakeData{}
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
		return []model.GetViewIntakeData{}
	}

	for i, data := range ViewIntakeData {
		ViewIntakeData[i].Answer = hashdb.Decrypt(data.Answer)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.GetViewIntakeData{}
	}

	return ViewIntakeData
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

	for _, answer := range reqVal.Answers {

		PrevData := model.GetViewIntakeData{}
		errPrev := tx.Raw(query.GetIntakeDataSQL, answer.ITFId).Scan(&PrevData).Error
		if errPrev != nil {
			log.Printf("ERROR: Failed to Get Intake: %v\n", PrevData)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		oldData := map[string]interface{}{
			fmt.Sprintf("Answers for %d", answer.QuestionId): hashdb.Decrypt(PrevData.Answer),
		}

		updatedData := map[string]interface{}{
			fmt.Sprintf("Answers for %d", answer.QuestionId): answer.Answer,
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

			updatedIntakeErr := tx.Exec(
				query.UpdateIntakeDataSQL,
				answer.Answer,
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

	return true, "Intakeform Updated from Technician"

}
