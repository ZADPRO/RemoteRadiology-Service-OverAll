package service

import (
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"

	"gorm.io/gorm"
)

func AddAppointmentService(db *gorm.DB, reqVal model.AddAppointmentReq, idValue int) (bool, string) {
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

	var FindScancenter []model.ScanCenterModel

	findScancenterErr := db.Raw(query.FindScanCenterSQL, reqVal.SCId).Scan(&FindScancenter).Error
	if findScancenterErr != nil {
		log.Printf("ERROR: Failed to fetch Scan Centers: %v", findScancenterErr)
	}

	if len(FindScancenter) == 0 {
		return false, "Invalid Scan Center Id"
	}

	var TotalCount []model.TotalCountModel

	// err := db.Raw(query.VerifyAppointment, reqVal.SCId, reqVal.AppointmentDate, reqVal.AppointmentStartTime, reqVal.AppointmentEndTime).Scan(&TotalCount).Error
	err := db.Raw(query.VerifyAppointment, reqVal.SCId, reqVal.AppointmentDate).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return false, "Something went wrong, Try Again"
	}

	if TotalCount[0].TotalCount > 0 {
		return false, "Appointment Already Exists"
	}

	Appointment := model.CreateAppointmentModel{
		UserId:          idValue,
		SCId:            FindScancenter[0].SCId,
		AppointmentDate: reqVal.AppointmentDate,
		// AppointmentStartTime: reqVal.AppointmentStartTime,
		// AppointmentEndTime:   reqVal.AppointmentEndTime,
		// AppointmentUrgency: reqVal.AppointmentUrgency,
		AppointmentStatus:   true,
		AppointmentComplete: "fillform",
	}

	Appointmenterr := db.Create(&Appointment).Error
	if Appointmenterr != nil {
		log.Printf("ERROR: Failed to create Appointment: %v\n", Appointmenterr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	history := model.RefTransHistory{
		TransTypeId: 21,
		THData:      "Appointment Created Successfully",
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

	return true, "Successfully Appointment Created"
}

func ViewPatientHistoryService(db *gorm.DB, idValue int) []model.ViewPatientHistoryModel {
	log := logger.InitLogger()

	var patientHistory []model.ViewPatientHistoryModel

	err := db.Raw(query.ViewPatientHistorySQL, idValue).Scan(&patientHistory).Error
	if err != nil {
		log.Printf("ERROR: Failed to View Patient History: %v", err)
		return []model.ViewPatientHistoryModel{}
	}

	return patientHistory
}

func ViewTechnicianPatientQueueService(db *gorm.DB, idValue int) []model.ViewTechnicianPatientQueueModel {
	log := logger.InitLogger()

	var patientQueue []model.ViewTechnicianPatientQueueModel

	err := db.Raw(query.ViewTechnicianPatientQueueSQL, idValue).Scan(&patientQueue).Error
	if err != nil {
		log.Printf("ERROR: Failed to View Patient Queue: %v", err)
		return []model.ViewTechnicianPatientQueueModel{}
	}

	return patientQueue
}
