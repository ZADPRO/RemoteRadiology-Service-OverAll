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

	var TotalCount []model.TotalCountModel

	err := db.Raw(query.VerifyAppointment, reqVal.SCId, reqVal.AppointmentDate, reqVal.AppointmentStartTime, reqVal.AppointmentEndTime).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return false, "Something went wrong, Try Again"
	}

	if TotalCount[0].TotalCount > 0 {
		return false, "Appointment Already Exists"
	}

	Appointment := model.CreateAppointmentModel{
		UserId:               idValue,
		SCId:                 reqVal.SCId,
		AppointmentDate:      reqVal.AppointmentDate,
		AppointmentStartTime: reqVal.AppointmentStartTime,
		AppointmentEndTime:   reqVal.AppointmentEndTime,
		AppointmentUrgency:   reqVal.AppointmentUrgency,
		AppointmentStatus:    true,
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
