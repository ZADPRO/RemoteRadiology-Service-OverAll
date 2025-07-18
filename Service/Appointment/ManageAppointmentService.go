package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"
	"encoding/json"
	"fmt"

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
	err := db.Raw(query.VerifyAppointment, FindScancenter[0].SCId, reqVal.AppointmentDate, idValue).Scan(&TotalCount).Error
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

func ViewTechnicianPatientQueueService(db *gorm.DB, idValue int, roleIdValue int) ([]model.ViewTechnicianPatientQueueModel, []model.StaffAvailableModel) {
	log := logger.InitLogger()

	var StaffAvailable []model.StaffAvailableModel

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

			var dicom []model.GetDicomFile
			log.Printf("Fetching Dicom for AppointmentId=%d, UserId=%d", data.AppointmentId, data.UserId)

			dicomErr := db.Raw(query.ViewGetDicomFile, data.AppointmentId, data.UserId).Scan(&dicom).Error
			if dicomErr != nil {
				log.Printf("ERROR: Failed to fetch Dicom Files: %v", dicomErr)
				return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			}

			patientQueue[i].DicomFiles = dicom
		}

		var suggestId []int
		switch roleIdValue {
		case 5:
			suggestId = []int{8}
		case 8:
			suggestId = []int{1, 5}
		}

		if len(suggestId) != 0 {

			if roleIdValue == 5 {
				SuggestUserErr := db.Raw(query.GetUserWithScanDetails, suggestId[0], Scancenter[0].SCId).Scan(&StaffAvailable).Error
				if SuggestUserErr != nil {
					log.Printf("ERROR: Failed to fetch Staff Available: %v", SuggestUserErr)
					return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
				}
				for i, data := range StaffAvailable {
					StaffAvailable[i].Username = hashdb.Decrypt(data.Username)
				}
				return patientQueue, StaffAvailable
			} else {
				SuggestUserErr := db.Raw(query.GetUserDetails, suggestId).Scan(&StaffAvailable).Error
				if SuggestUserErr != nil {
					log.Printf("ERROR: Failed to fetch Staff Available: %v", SuggestUserErr)
					return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
				}
				for i, data := range StaffAvailable {
					StaffAvailable[i].Username = hashdb.Decrypt(data.Username)
				}
				return patientQueue, StaffAvailable

			}

		} else {

			return patientQueue, []model.StaffAvailableModel{}
		}

	} else {

		err := db.Raw(query.ViewAllPatientQueueSQL).Scan(&patientQueue).Error
		if err != nil {
			log.Printf("ERROR: Failed to View Patient Queue: %v", err)
			return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
		}

		for i, data := range patientQueue {
			patientQueue[i].Username = hashdb.Decrypt(data.Username)

			var dicom []model.GetDicomFile
			log.Printf("Fetching Dicom for AppointmentId=%d, UserId=%d", data.AppointmentId, data.UserId)

			dicomErr := db.Raw(query.ViewGetDicomFile, data.AppointmentId, data.UserId).Scan(&dicom).Error
			if dicomErr != nil {
				log.Printf("ERROR: Failed to fetch Dicom Files: %v", dicomErr)
				return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			}

			patientQueue[i].DicomFiles = dicom
		}

		var suggestId []int
		switch roleIdValue {
		case 7:
			suggestId = []int{6}
		case 6:
			suggestId = []int{7, 6}
		}

		if len(suggestId) != 0 {

			// if roleIdValue == 7 {
			// 	SuggestUserErr := db.Raw(query.GetUserWithScanDetails, suggestId[0], Scancenter[0].SCId).Scan(&StaffAvailable).Error
			// 	if SuggestUserErr != nil {
			// 		log.Printf("ERROR: Failed to fetch Staff Available: %v", SuggestUserErr)
			// 		return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			// 	}
			// 	for i, data := range StaffAvailable {
			// 		StaffAvailable[i].Username = hashdb.Decrypt(data.Username)
			// 	}
			// 	return patientQueue, StaffAvailable
			// } else {

			fmt.Println(suggestId)

			SuggestUserErr := db.Raw(query.GetUserDetails, suggestId).Scan(&StaffAvailable).Error
			if SuggestUserErr != nil {
				log.Printf("ERROR: Failed to fetch Staff Available: %v", SuggestUserErr)
				return []model.ViewTechnicianPatientQueueModel{}, []model.StaffAvailableModel{}
			}
			for i, data := range StaffAvailable {
				StaffAvailable[i].Username = hashdb.Decrypt(data.Username)
			}
			return patientQueue, StaffAvailable

			// }

		} else {

			return patientQueue, []model.StaffAvailableModel{}
		}
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

	errExec := tx.Exec(query.InsertAdditionalFiles, idValue, reqVal.AppointmentId, true, string(filesJSON)).Error
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
			log.Fatalf("Failed to read profile image file: %v", viewErr)
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

	fmt.Println(reqVal.AssingUserId, reqVal.AppointmentId)

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
