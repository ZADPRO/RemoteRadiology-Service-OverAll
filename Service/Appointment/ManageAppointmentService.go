package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

func AddAppointmentService(db *gorm.DB, reqVal model.AddAppointmentReq, idValue int, roleIdValue int) (bool, string, int, int, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again", 0, 0, ""
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
		return false, "Invalid Scan Center Id", 0, 0, ""
	}

	var TotalCount []model.TotalCountModel

	// err := db.Raw(query.VerifyAppointment, reqVal.SCId, reqVal.AppointmentDate, reqVal.AppointmentStartTime, reqVal.AppointmentEndTime).Scan(&TotalCount).Error
	err := db.Raw(query.VerifyAppointment, FindScancenter[0].SCId, reqVal.AppointmentDate, idValue).Scan(&TotalCount).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch User Total Count: %v", err)
		return false, "Something went wrong, Try Again", 0, 0, ""
	}

	if TotalCount[0].TotalCount > 0 {
		return false, "Appointment Already Exists", 0, 0, ""
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
		return false, "Something went wrong, Try Again", 0, 0, ""
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
		return false, "Something went wrong, Try Again", 0, 0, ""
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again", 0, 0, ""
	}

	return true, "Successfully Appointment Created", Appointment.AppointmentId, Appointment.SCId, reqVal.SCId
}

func ViewPatientHistoryService(db *gorm.DB, idValue int) ([]model.ViewPatientHistoryModel, bool, string) {
	log := logger.InitLogger()

	var patientHistory []model.ViewPatientHistoryModel

	err := db.Raw(query.ViewPatientHistorySQL, idValue).Scan(&patientHistory).Error
	if err != nil {
		log.Printf("ERROR: Failed to View Patient History: %v", err)
		return []model.ViewPatientHistoryModel{}, false, ""
	}

	var ScanCenterMap []model.MapScanCenterPatientModel
	ScanCenterMapErr := db.Raw(query.ScanCenterMap, idValue).Scan(&ScanCenterMap).Error
	if ScanCenterMapErr != nil {
		log.Printf("ERROR: Failed to View Scan Center: %v", ScanCenterMapErr)
		return []model.ViewPatientHistoryModel{}, false, ""
	}

	if len(ScanCenterMap) > 0 {
		var ScanCenterConsultant []model.ScanCenterConsultantModel
		ScanCenterConsultantErr := db.Raw(query.ScanCenterConsultantSQL, ScanCenterMap[0].SCId).Scan(&ScanCenterConsultant).Error
		if ScanCenterConsultantErr != nil {
			log.Printf("ERROR: Failed to View Scan Center: %v", ScanCenterConsultantErr)
			return []model.ViewPatientHistoryModel{}, false, ""
		}

		var ConsultantStatus = false
		var ConsultantLink = ""
		if len(ScanCenterConsultant) > 0 {
			ConsultantStatus = ScanCenterConsultant[0].SCConsultantStatus
			ConsultantLink = ScanCenterConsultant[0].SCConsultantLink
		}

		return patientHistory, ConsultantStatus, ConsultantLink
	} else {
		return []model.ViewPatientHistoryModel{}, false, ""
	}

}

func ViewTechnicianPatientQueueService(db *gorm.DB, idValue int, roleIdValue int) ([]model.ViewTechnicianPatientQueueModel, []model.StaffAvailableModel) {
	log := logger.InitLogger()

	var StaffAvailable []model.StaffAvailableModel
	var FinalStaffAvailable []model.StaffAvailableModel

	var patientQueue []model.ViewTechnicianPatientQueueModel

	var Scancenter []model.ScanCenterModel

	IdentifyScancentererr := db.Raw(query.IdentifyScanCenterWithUser, idValue).Scan(&Scancenter).Error
	if IdentifyScancentererr != nil {
		log.Printf("ERROR: Failed to Identify Scan Center: %v", IdentifyScancentererr)
		return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
	}

	if roleIdValue == 2 || roleIdValue == 3 || roleIdValue == 5 || roleIdValue == 8 {

		err := db.Raw(query.ViewTechnicianPatientQueueSQL, idValue, Scancenter[0].SCId).Scan(&patientQueue).Error
		if err != nil {
			log.Printf("ERROR: Failed to View Patient Queue: %v", err)
			return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
		}

		for i, data := range patientQueue {
			patientQueue[i].Username = hashdb.Decrypt(data.Username)

			var ReportUrgentStatus []model.ReportUrgentStatusModel

			GetReportUregentStatuserr := db.Raw(query.GetReportStatusSQL, data.AppointmentId).Scan(&ReportUrgentStatus).Error
			if GetReportUregentStatuserr != nil {
				log.Printf("ERROR: Failed to Get Report Status: %v", GetReportUregentStatuserr)
				return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			}

			if len(ReportUrgentStatus) > 0 {
				patientQueue[i].ReportStatus = hashdb.Decrypt(ReportUrgentStatus[0].ReportStatus)
			}

			// var dicom []model.GetDicomFile
			// log.Printf("Fetching Dicom for AppointmentId=%d, UserId=%d", data.AppointmentId, data.UserId)

			// dicomErr := db.Raw(query.ViewGetDicomFile, data.AppointmentId, data.UserId).Scan(&dicom).Error
			// if dicomErr != nil {
			// 	log.Printf("ERROR: Failed to fetch Dicom Files: %v", dicomErr)
			// 	return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			// }

			// patientQueue[i].DicomFiles = dicom
		}

		SuggestUserErr := db.Raw(query.GetUserWithScanDetails).Scan(&StaffAvailable).Error
		if SuggestUserErr != nil {
			log.Printf("ERROR: Failed to fetch Staff Available: %v", SuggestUserErr)
			return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
		}
		for i, data := range StaffAvailable {
			StaffAvailable[i].Username = hashdb.Decrypt(data.Username)
		}

		if roleIdValue == 1 || roleIdValue == 10 {
			FinalStaffAvailable = StaffAvailable
		} else if roleIdValue == 2 || roleIdValue == 3 || roleIdValue == 5 || roleIdValue == 8 {
			for _, data := range StaffAvailable {
				if data.RoleId == 10 || data.RoleId == 5 || data.RoleId == 8 || data.RoleId == 3 || data.RoleId == 2 {
					FinalStaffAvailable = append(FinalStaffAvailable, data)
				}
			}
		} else if roleIdValue == 6 || roleIdValue == 7 {
			for _, data := range StaffAvailable {
				if data.RoleId == 10 || data.RoleId == 6 || data.RoleId == 7 {
					FinalStaffAvailable = append(FinalStaffAvailable, data)
				}
			}
		}

		return patientQueue, FinalStaffAvailable

	} else {

		err := db.Raw(query.ViewAllPatientQueueSQL).Scan(&patientQueue).Error
		if err != nil {
			log.Printf("ERROR: Failed to View Patient Queue: %v", err)
			return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
		}

		for i, data := range patientQueue {
			patientQueue[i].Username = hashdb.Decrypt(data.Username)

			var ReportUrgentStatus []model.ReportUrgentStatusModel

			GetReportUregentStatuserr := db.Raw(query.GetReportStatusSQL, data.AppointmentId).Scan(&ReportUrgentStatus).Error
			if GetReportUregentStatuserr != nil {
				log.Printf("ERROR: Failed to Get Report Status: %v", GetReportUregentStatuserr)
				return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			}

			if len(ReportUrgentStatus) > 0 {
				patientQueue[i].ReportStatus = hashdb.Decrypt(ReportUrgentStatus[0].ReportStatus)
			}

			//Get the Private and Public in the Technician Intake
			var patientPrivatePublic []model.GetTechnicianIntakeData
			patientPrivatePublicErr := db.Raw(query.GetPatientPrivatePublicSQL, patientQueue[i].AppointmentId, 57).Scan(&patientPrivatePublic).Error
			if patientPrivatePublicErr != nil {
				log.Error(patientPrivatePublicErr)
			}

			for i, data := range patientPrivatePublic {
				patientPrivatePublic[i].Answer = hashdb.Decrypt(data.Answer)
			}

			if len(patientPrivatePublic) > 0 {
				patientQueue[i].PatientPrivatePublicStatus = patientPrivatePublic[0].Answer
			} else {
				patientQueue[i].PatientPrivatePublicStatus = ""
			}

			// var dicom []model.GetDicomFile
			// log.Printf("Fetching Dicom for AppointmentId=%d, UserId=%d", data.AppointmentId, data.UserId)

			// dicomErr := db.Raw(query.ViewGetDicomFile, data.AppointmentId, data.UserId).Scan(&dicom).Error
			// if dicomErr != nil {
			// 	log.Printf("ERROR: Failed to fetch Dicom Files: %v", dicomErr)
			// 	return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			// }

			// patientQueue[i].DicomFiles = dicom

			// //Get the Correct and Edit
			// if roleIdValue == 1 || roleIdValue == 6 {

			// 	var CorrectEditModel []model.GetCorrectEditModel

			// 	CorrectEditModelErr := db.Raw(query.CorrectEditStatusSQL, data.UserId, data.AppointmentId, idValue).Scan(&CorrectEditModel).Error
			// 	if CorrectEditModelErr != nil {
			// 		log.Printf("ERROR: Failed to Identify Scan Center: %v", CorrectEditModelErr)
			// 		return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			// 	}

			// 	if len(CorrectEditModel) > 0 {
			// 		patientQueue[i].GetCorrectEditModel.RHHandleCorrect = CorrectEditModel[0].RHHandleCorrect
			// 		patientQueue[i].GetCorrectEditModel.RHHandleEdit = CorrectEditModel[0].RHHandleEdit
			// 	} else {
			// 		patientQueue[i].GetCorrectEditModel.RHHandleCorrect = false
			// 		patientQueue[i].GetCorrectEditModel.RHHandleEdit = false
			// 	}

			// } else {
			// 	patientQueue[i].GetCorrectEditModel.RHHandleCorrect = false
			// 	patientQueue[i].GetCorrectEditModel.RHHandleEdit = false
			// }

		}

		SuggestUserErr := db.Raw(query.GetUserDetails).Scan(&StaffAvailable).Error
		if SuggestUserErr != nil {
			log.Printf("ERROR: Failed to fetch Staff Available: %v", SuggestUserErr)
			return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
		}
		for i, data := range StaffAvailable {
			StaffAvailable[i].Username = hashdb.Decrypt(data.Username)
		}

		if roleIdValue == 1 || roleIdValue == 10 || roleIdValue == 9 {
			FinalStaffAvailable = StaffAvailable
		} else if roleIdValue == 2 || roleIdValue == 3 || roleIdValue == 5 || roleIdValue == 8 {
			for _, data := range StaffAvailable {
				if data.RoleId == 10 || data.RoleId == 5 || data.RoleId == 8 {
					FinalStaffAvailable = append(FinalStaffAvailable, data)
				}
			}
		} else if roleIdValue == 6 || roleIdValue == 7 {
			for _, data := range StaffAvailable {
				if data.RoleId == 10 || data.RoleId == 6 || data.RoleId == 7 || data.RoleId == 9 {
					FinalStaffAvailable = append(FinalStaffAvailable, data)
				}
			}
		}

		return patientQueue, FinalStaffAvailable
	}

}

func AddAddtionalFilesService(db *gorm.DB, reqVal model.AddAddtionalFilesReq, idValue int) (bool, string) {
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

	// Convert file slice to JSON string
	filesJSON, err := json.Marshal(reqVal.Files)
	if err != nil {
		log.Printf("ERROR: Failed to marshal files JSON: %v\n", err)
		tx.Rollback()
		return false, "Invalid file data"
	}

	errExec := tx.Exec(query.InsertAdditionalFiles, idValue, reqVal.AppointmentId, true, timeZone.GetPacificTime(), string(filesJSON)).Error
	if errExec != nil {
		log.Printf("ERROR: Failed to insert additional files: %v\n", errExec)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		return false, "Something went wrong, Try Again"
	}

	return true, "Added Successfully!"
}

func ViewAddtionalFilesService(db *gorm.DB, reqVal model.ViewAddtionalFileReq) []model.AdditionalFileUploadModel {
	log := logger.InitLogger()

	var ViewFile []model.AdditionalFileUploadModel

	err := db.Raw(query.ViewAddtionalFilesSQL, reqVal.UserId, reqVal.AppointmentId).Scan(&ViewFile).Error
	if err != nil {
		log.Printf("ERROR: Failed to View Patient Queue: %v", err)
		return []model.AdditionalFileUploadModel{}
	}

	for i, data := range ViewFile {
		ViewFileData, viewErr := helper.ViewFile("./Assets/Files/" + data.FileName)
		if viewErr != nil {
			// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
			log.Errorf("Failed to read profile image file: %v", viewErr)
		}
		if ViewFileData != nil {
			ViewFile[i].FileData = &model.FileData{
				Base64Data:  ViewFileData.Base64Data,
				ContentType: ViewFileData.ContentType,
			}
		}
	}

	return ViewFile
}

func AssignUserService(db *gorm.DB, reqVal model.AssignUserReq, idValue int) (bool, string) {
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

	var GetUserDetails model.GetAllUserDetailsModel

	// fmt.Println(reqVal.AssingUserId, reqVal.AppointmentId)
	GetUserDeatilsErr := tx.Raw(
		query.GetAllUserDetailsSQL,
		reqVal.AssingUserId,
		reqVal.PatientId,
		reqVal.AppointmentId,
		idValue,
	).Scan(&GetUserDetails).Error
	if GetUserDeatilsErr != nil {
		log.Printf("ERROR: Failed to update assign user: %v\n", GetUserDeatilsErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Updating the Assigning User
	errExec := tx.Exec(
		query.UpdateAssignUser,
		reqVal.AssingUserId,
		reqVal.AppointmentId,
	).Error
	if errExec != nil {
		log.Printf("ERROR: Failed to update assign user: %v\n", errExec)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	message := fmt.Sprintf(
		"A report has been assigned to you by %s. The patient ID is %s, and the appointment is scheduled for %s.",
		GetUserDetails.User_Id,
		GetUserDetails.Patient_Id,
		GetUserDetails.Appointment_date,
	)

	notificationExec := tx.Exec(
		query.InsertNotificationSQL,
		reqVal.AssingUserId,
		message,
		reqVal.AppointmentId,
		idValue,
		timeZone.GetPacificTime(),
		false,
		true,
	).Error
	if notificationExec != nil {
		log.Printf("ERROR: Failed to update assign user: %v\n", notificationExec)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Adding Audit Page
	history := model.RefTransHistory{
		TransTypeId: 33,
		THData:      "Report was Assigned for " + reqVal.AssingUserCustId,
		UserId:      reqVal.PatientId,
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

	return true, "Added Successfully!"

}

func ListMessagesService(db *gorm.DB, idValue int) []model.Notification {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return []model.Notification{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var Notification []model.Notification

	NotificationErr := db.Raw(query.NotificationSQL, idValue).Scan(&Notification).Error
	if NotificationErr != nil {
		log.Error(NotificationErr)
		return []model.Notification{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.Notification{}
	}

	return Notification

}

func ListAuditLogService(db *gorm.DB) []model.RefAuditTransHistory {
	log := logger.InitLogger()

	var AuditTransHitory []model.RefAuditTransHistory

	AuditTransHitoryErr := db.Raw(query.GetAuditReportStatusSQL).Scan(&AuditTransHitory).Error
	if AuditTransHitoryErr != nil {
		log.Error(AuditTransHitoryErr)
		return []model.RefAuditTransHistory{}
	}

	return AuditTransHitory

}
