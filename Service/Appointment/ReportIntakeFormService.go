package service

import (
	helper "AuthenticationService/internal/Helper/GetChanges"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	mailservice "AuthenticationService/internal/Helper/MailService"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	helperView "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

func CheckAccessService(db *gorm.DB, reqVal model.CheckAccessReq, idValue int, roleIdValue int) (bool, string, int, string) {
	log := logger.InitLogger()

	var result []model.AccessStatus

	fmt.Println("222222222222222222", roleIdValue, idValue, reqVal)

	if roleIdValue == 7 {
		err := db.Raw(query.ScribeCheckAccessSQL, idValue, reqVal.AppointmentId).Scan(&result).Error

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("&&&&&&&&&&&&&&&", result)

	} else {
		err := db.Raw(query.CheckAccessSQL, idValue, reqVal.AppointmentId).Scan(&result).Error

		if err != nil {
			log.Fatal(err)
		}
	}

	var message = "Another User Already Accessing it"

	if result[0].Status {
		message = "Report are Available for Edit"
	}

	return result[0].Status, message, result[0].RefAppointmentAccessId, result[0].CustID
}

func AssignGetReportService(db *gorm.DB, reqVal model.AssignGetReportReq, idValue int, roleIdValue int) (bool, string, []model.GetViewIntakeData, []model.GetTechnicianIntakeData, []model.GetReportIntakeData, []model.GetReportTextContent, []model.GetReportHistory, []model.GetReportComments, []model.GetOneUserAppointmentModel, []model.ReportFormateModel, []model.GetUserDetails, []model.PatientCustId, bool, *model.FileData, string, []model.AddAddendumModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again",
			[]model.GetViewIntakeData{},
			[]model.GetTechnicianIntakeData{},
			[]model.GetReportIntakeData{},
			[]model.GetReportTextContent{},
			[]model.GetReportHistory{},
			[]model.GetReportComments{},
			[]model.GetOneUserAppointmentModel{},
			[]model.ReportFormateModel{},
			[]model.GetUserDetails{},
			[]model.PatientCustId{},
			false,
			&model.FileData{},
			"",
			[]model.AddAddendumModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	//patientCustId
	var PatientUserDetails []model.PatientCustId
	PatientUserDetailsErr := db.Raw(query.PatientUserDetailsSQL, reqVal.PatientId).Scan(&PatientUserDetails).Error
	if PatientUserDetailsErr != nil {
		log.Fatal(PatientUserDetailsErr)
	}

	//GetUserDetails
	var UserDetails []model.GetUserDetails
	UserDetailsErr := db.Raw(query.GetUserDetailsSQL, idValue).Scan(&UserDetails).Error
	if UserDetailsErr != nil {
		log.Fatal(UserDetailsErr)
	}

	//Decrypt UserDetails
	for i, data := range UserDetails {
		UserDetails[i].FirstName = hashdb.Decrypt(data.FirstName)
		if len(UserDetails[i].Specialization) > 0 {
			UserDetails[i].Specialization = hashdb.Decrypt(data.Specialization)
		}
		if len(UserDetails[i].Department) > 0 {
			UserDetails[i].Department = hashdb.Decrypt(data.Department)
		}
	}

	checkAccessReq := model.CheckAccessReq{
		AppointmentId: reqVal.AppointmentId,
	}

	status, message, _, _ := CheckAccessService(db, checkAccessReq, idValue, roleIdValue)

	fmt.Println(status, message)

	if (status && !reqVal.ReadOnlyStatus) || (!status && reqVal.ReadOnlyStatus) || (status && reqVal.ReadOnlyStatus) {

		//Appointment Table
		var Appointment []model.AppointmentModel
		Appointmenterr := db.Raw(query.GetAppointmentSQL, reqVal.AppointmentId).Scan(&Appointment).Error
		if Appointmenterr != nil {
			log.Fatal(Appointmenterr)
		}

		if !reqVal.ReadOnlyStatus {

			var AppointementAccessIdVal = Appointment[0].AppointmentAccessId

			if roleIdValue == 7 {
				AppointementAccessIdVal = Appointment[0].AppointmentScribeAccessId
			}

			oldDataCat := map[string]interface{}{
				"Report Access ID": AppointementAccessIdVal,
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
					return false, "Something went wrong, Try Again",
						[]model.GetViewIntakeData{},
						[]model.GetTechnicianIntakeData{},
						[]model.GetReportIntakeData{},
						[]model.GetReportTextContent{},
						[]model.GetReportHistory{},
						[]model.GetReportComments{},
						[]model.GetOneUserAppointmentModel{},
						[]model.ReportFormateModel{},
						[]model.GetUserDetails{},
						[]model.PatientCustId{},
						false,
						&model.FileData{},
						"",
						[]model.AddAddendumModel{}
				}

				transData := 28
				errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(Appointment[0].UserId), int(idValue), string(ChangesDataJSON)).Error
				if errTrans != nil {
					log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
					tx.Rollback()
					return false, "Something went wrong, Try Again",
						[]model.GetViewIntakeData{},
						[]model.GetTechnicianIntakeData{},
						[]model.GetReportIntakeData{},
						[]model.GetReportTextContent{},
						[]model.GetReportHistory{},
						[]model.GetReportComments{},
						[]model.GetOneUserAppointmentModel{},
						[]model.ReportFormateModel{},
						[]model.GetUserDetails{},
						[]model.PatientCustId{},
						false,
						&model.FileData{},
						"",
						[]model.AddAddendumModel{}
				}

				var UpdateAccessSQL = query.UpdateAccessAppointment

				if roleIdValue == 7 {
					UpdateAccessSQL = query.ScribeUpdateAccessAppointment
				}

				categoryUpdate := tx.Exec(
					UpdateAccessSQL,
					true,
					idValue,
					reqVal.AppointmentId,
				).Error

				if categoryUpdate != nil {
					log.Printf("ERROR: Failed toCategory Id: %v\n", categoryUpdate)
					tx.Rollback()
					return false, "Something went wrong, Try Again",
						[]model.GetViewIntakeData{},
						[]model.GetTechnicianIntakeData{},
						[]model.GetReportIntakeData{},
						[]model.GetReportTextContent{},
						[]model.GetReportHistory{},
						[]model.GetReportComments{},
						[]model.GetOneUserAppointmentModel{},
						[]model.ReportFormateModel{},
						[]model.GetUserDetails{},
						[]model.PatientCustId{},
						false,
						&model.FileData{},
						"",
						[]model.AddAddendumModel{}
				}

				//List the Latest Report History
				var ReportHistory []model.GetReportHistory
				ListReportHistoryErr := db.Raw(query.CheckLatestReportSQL, reqVal.AppointmentId, reqVal.PatientId).Scan(&ReportHistory).Error
				if ListReportHistoryErr != nil {
					log.Fatal(ListReportHistoryErr)
				}

				if ListReportHistoryErr != nil {
					log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", errChange)
					tx.Rollback()
					return false, "Invalid User Accessing",
						[]model.GetViewIntakeData{},
						[]model.GetTechnicianIntakeData{},
						[]model.GetReportIntakeData{},
						[]model.GetReportTextContent{},
						[]model.GetReportHistory{},
						[]model.GetReportComments{},
						[]model.GetOneUserAppointmentModel{},
						[]model.ReportFormateModel{},
						[]model.GetUserDetails{},
						[]model.PatientCustId{},
						false,
						&model.FileData{},
						"",
						[]model.AddAddendumModel{}
				}

				if len(ReportHistory) > 0 {

					var starttime = ReportHistory[0].HandleEndTime

					if len(ReportHistory[0].HandleEndTime) == 0 {
						starttime = ReportHistory[1].HandleEndTime
					}

					fmt.Println("___________________>", starttime)

					//Insert the History
					ReportHistoryErr := tx.Exec(
						query.InsertReportHistorySQL,
						reqVal.PatientId,
						reqVal.AppointmentId,
						idValue,
						starttime,
					).Error
					if ReportHistoryErr != nil {
						log.Printf("ERROR: Failed to Insert Report History: %v\n", ReportHistoryErr)
						tx.Rollback()
						return false, "Something went wrong, Try Again",
							[]model.GetViewIntakeData{},
							[]model.GetTechnicianIntakeData{},
							[]model.GetReportIntakeData{},
							[]model.GetReportTextContent{},
							[]model.GetReportHistory{},
							[]model.GetReportComments{},
							[]model.GetOneUserAppointmentModel{},
							[]model.ReportFormateModel{},
							[]model.GetUserDetails{},
							[]model.PatientCustId{},
							false,
							&model.FileData{},
							"",
							[]model.AddAddendumModel{}
					}
				}

			}

		}

		if err := tx.Commit().Error; err != nil {
			log.Printf("ERROR: Failed to commit transaction: %v\n", err)
			tx.Rollback()
			return false, "Something went wrong, Try Again",
				[]model.GetViewIntakeData{},
				[]model.GetTechnicianIntakeData{},
				[]model.GetReportIntakeData{},
				[]model.GetReportTextContent{},
				[]model.GetReportHistory{},
				[]model.GetReportComments{},
				[]model.GetOneUserAppointmentModel{},
				[]model.ReportFormateModel{},
				[]model.GetUserDetails{},
				[]model.PatientCustId{},
				false,
				&model.FileData{},
				"",
				[]model.AddAddendumModel{}
		}

		var IntakeFormData []model.GetViewIntakeData
		var OneUserAppointment []model.GetOneUserAppointmentModel
		//Appointment Table
		ViewAppointmentErr := db.Raw(query.GetOneUserAppointment, reqVal.PatientId, reqVal.AppointmentId).Scan(&OneUserAppointment).Error
		if ViewAppointmentErr != nil {
			log.Fatal(ViewAppointmentErr)
		}

		//Decrypt Appointment Table
		for i, data := range OneUserAppointment {
			OneUserAppointment[i].SCName = hashdb.Decrypt(data.SCName)
		}

		//Intake Form Table
		IntakeFormDataerr := db.Raw(query.GetIntakeFormSQL, reqVal.AppointmentId).Scan(&IntakeFormData).Error
		if IntakeFormDataerr != nil {
			log.Fatal(IntakeFormDataerr)
		}

		//Decrypt Intake Form Table
		for i, data := range IntakeFormData {
			IntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}

		//Technician Intake Form Table
		var TechnicianIntakeFormData []model.GetTechnicianIntakeData
		TechnicianIntakeFormDataerr := db.Raw(query.GetTechnicianIntakeFormSQL, reqVal.AppointmentId).Scan(&TechnicianIntakeFormData).Error
		if TechnicianIntakeFormDataerr != nil {
			log.Fatal(TechnicianIntakeFormDataerr)
		}

		//Decrypt the Techncian Form Table
		for i, data := range TechnicianIntakeFormData {
			TechnicianIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}

		//Report Intake Form Table
		var ReportIntakeFormData []model.GetReportIntakeData
		ReportIntakeFormDataerr := db.Raw(query.GetReportIntakeFormSQL, reqVal.AppointmentId).Scan(&ReportIntakeFormData).Error
		if ReportIntakeFormDataerr != nil {
			log.Fatal(ReportIntakeFormDataerr)
		}

		//Decrypt Report Intake Form Table
		for i, data := range ReportIntakeFormData {
			ReportIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}

		//Report Text Content Table
		var ReportTextContentData []model.GetReportTextContent
		ReportTextContentDataerr := db.Raw(query.GetReporttextContent, reqVal.AppointmentId).Scan(&ReportTextContentData).Error
		if ReportTextContentDataerr != nil {
			log.Fatal(ReportTextContentDataerr)
		}

		//Decrypt Report Text Content Table
		for i, data := range ReportTextContentData {
			ReportTextContentData[i].TextContent = hashdb.Decrypt(data.TextContent)
		}

		//Report History Table
		var ReportHistoryData []model.GetReportHistory
		ReportHistoryDataerr := db.Raw(query.GetReportHistorySQL, reqVal.AppointmentId).Scan(&ReportHistoryData).Error
		if ReportHistoryDataerr != nil {
			log.Fatal(ReportHistoryDataerr)
		}

		//Decrypt Report History Table
		for i, data := range ReportHistoryData {
			ReportHistoryData[i].HandleUserName = hashdb.Decrypt(data.HandleUserName)
			ReportHistoryData[i].HandleContentText = hashdb.Decrypt(data.HandleContentText)
		}

		// Report Comment Table
		var ReportCommentsData []model.GetReportComments
		ReportCommentsDataerr := db.Raw(query.GetReportCommentsSQL, reqVal.AppointmentId).Scan(&ReportCommentsData).Error
		if ReportCommentsDataerr != nil {
			log.Fatal(ReportCommentsDataerr)
		}

		// Decrypt Report Comment Table
		for i, data := range ReportCommentsData {
			ReportCommentsData[i].Status = hashdb.Decrypt(data.Status)
			ReportCommentsData[i].Comments = hashdb.Decrypt(data.Comments)
		}

		//Get the Template all listed
		var ReportFormateList []model.ReportFormateModel
		ReportFormateListErr := db.Raw(query.GetReportFormateListSQL).Scan(&ReportFormateList).Error
		if ReportFormateListErr != nil {
			log.Fatal(ReportFormateListErr)
		}

		// Decrypt Report Formate List
		for i, data := range ReportFormateList {
			ReportFormateList[i].RFName = hashdb.Decrypt(data.RFName)
		}

		//Scan Center Profile Img
		var GetScanCenterImg []model.ScanCenterModel
		GetScanCenterImgErr := db.Raw(query.ScanCenterSQL, Appointment[0].SCId).Scan(&GetScanCenterImg).Error
		if GetScanCenterImgErr != nil {
			log.Fatal(GetScanCenterImgErr)
		}

		var ScanCenterProfileImg *model.FileData

		if len(GetScanCenterImg) > 0 {
			viewedFile, viewErr := helperView.ViewFile("./Assets/Profile/" + hashdb.Decrypt(GetScanCenterImg[0].ProfileImg))
			if viewErr != nil {
				log.Fatalf("Failed to read ScanCenter profile image: %v", viewErr)
			}

			ScanCenterProfileImg = &model.FileData{
				Base64Data:  viewedFile.Base64Data,
				ContentType: viewedFile.ContentType,
			}
		} else {
			ScanCenterProfileImg = &model.FileData{}
		}

		var EaseQTReportAccess = false

		//Get the Ease QT Report Access Status
		switch roleIdValue {
		case 1: //Master Admin
			EaseQTReportAccess = true
		case 2: //Scan Center Technician
			EaseQTReportAccess = false
		case 3: //Scan Center Manager
			EaseQTReportAccess = false
		case 4: //Patient
			EaseQTReportAccess = false
		case 5: //Scan Center Doctor

			var ReportStatus []model.DoctorReportAccessStatus
			err := db.Raw(query.DoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
			if err != nil {
				log.Fatal(err)
			}

			if len(ReportStatus) == 0 || ReportStatus[0].DDEaseQTReportAccess == nil {
				EaseQTReportAccess = false
				break
			}

			EaseQTReportAccess = *ReportStatus[0].DDEaseQTReportAccess

		case 6: //Junior Doctor
			EaseQTReportAccess = true
		case 7: //Scribe
			EaseQTReportAccess = true
		case 8: //Scan Center Reviewer
			var ReportStatus []model.CoDoctorReportAccessStatus
			err := db.Raw(query.CoDoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
			if err != nil {
				log.Fatal(err)
			}

			if len(ReportStatus) == 0 || ReportStatus[0].CDEaseQTReportAccess == nil {
				EaseQTReportAccess = false
				break
			}

			EaseQTReportAccess = *ReportStatus[0].CDEaseQTReportAccess
		case 9: //Manager
			EaseQTReportAccess = true
		case 10: //Performing Provider
			EaseQTReportAccess = true
		default:
			EaseQTReportAccess = false
		}
		// if roleIdValue == 1 || roleIdValue == 9 || roleIdValue == 6 || roleIdValue == 7 || roleIdValue == 10 {
		// 	EaseQTReportAccess = true
		// } else if roleIdValue == 5 { //Check the Scan Center Doctor

		// }

		return true, "Successfully Fetched", IntakeFormData, TechnicianIntakeFormData, ReportIntakeFormData, ReportTextContentData, ReportHistoryData, ReportCommentsData, OneUserAppointment, ReportFormateList, UserDetails, PatientUserDetails, EaseQTReportAccess, ScanCenterProfileImg, hashdb.Decrypt(GetScanCenterImg[0].SCAddress), ListAddendumService(db, reqVal.AppointmentId)

	} else {

		if err := tx.Commit().Error; err != nil {
			log.Printf("ERROR: Failed to commit transaction: %v\n", err)
			tx.Rollback()
			return false, "Something went wrong, Try Again",
				[]model.GetViewIntakeData{},
				[]model.GetTechnicianIntakeData{},
				[]model.GetReportIntakeData{},
				[]model.GetReportTextContent{},
				[]model.GetReportHistory{},
				[]model.GetReportComments{},
				[]model.GetOneUserAppointmentModel{},
				[]model.ReportFormateModel{},
				[]model.GetUserDetails{},
				[]model.PatientCustId{},
				false,
				&model.FileData{},
				"",
				[]model.AddAddendumModel{}
		}

		return status, message,
			[]model.GetViewIntakeData{},
			[]model.GetTechnicianIntakeData{},
			[]model.GetReportIntakeData{},
			[]model.GetReportTextContent{},
			[]model.GetReportHistory{},
			[]model.GetReportComments{},
			[]model.GetOneUserAppointmentModel{},
			[]model.ReportFormateModel{},
			[]model.GetUserDetails{},
			[]model.PatientCustId{},
			false,
			&model.FileData{},
			"",
			[]model.AddAddendumModel{}
	}

}

func AnswerReportIntakeService(db *gorm.DB, reqVal model.AnswerReportIntakeReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db

	// tx := db.Begin()
	// if tx.Error != nil {
	// 	log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
	// 	return false, "Something went wrong, Try Again"
	// }

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
	// 		tx.Rollback()
	// 	}
	// }()

	//Checking the Question ID is Available
	var ReportIntakeFormData []model.GetReportIntakeData

	ReportIntakeFormDataerr := db.Raw(query.GetReportIntakeFormQuestionSQL, reqVal.AppointmentId, reqVal.QuestionId).Scan(&ReportIntakeFormData).Error
	if ReportIntakeFormDataerr != nil {
		log.Fatal(ReportIntakeFormDataerr)
	}

	//If AVaiable Need to Update, else Create a QuestionID and Answer
	if len(ReportIntakeFormData) > 0 {

		//Find the If any Changes is Avaiable
		oldData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): hashdb.Decrypt(ReportIntakeFormData[0].Answer),
		}

		updatedData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): reqVal.Answer,
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

			//Insert Aduit Row for Answers Update
			transData := 30
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

			//Update the Answer with QuestionID
			UpdateTechnicianInputErr := tx.Exec(
				query.UpdateReportIntakeSQL,
				hashdb.Encrypt(reqVal.Answer),
				timeZone.GetPacificTime(),
				idValue,
				reqVal.QuestionId,
				reqVal.AppointmentId,
			).Error
			if UpdateTechnicianInputErr != nil {
				log.Printf("ERROR: Failed to Update Technician Input: %v\n", UpdateTechnicianInputErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

		}

	} else {

		//Inserting a new QuestionID and Answer
		InsertTechnicianInputErr := tx.Exec(
			query.InsertTechnicianIntakeSQL,
			reqVal.PatientId,
			reqVal.AppointmentId,
			reqVal.QuestionId,
			hashdb.Encrypt(reqVal.Answer),
			timeZone.GetPacificTime(),
			idValue,
		).Error
		if InsertTechnicianInputErr != nil {
			log.Printf("ERROR: Failed to Insert Technician Input: %v\n", InsertTechnicianInputErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		//Adding the Aduit Row Data
		oldData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): "",
		}

		updatedData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): reqVal.Answer,
		}

		ChangesData := helper.GetChanges(updatedData, oldData)

		var ChangesDataJSON []byte
		var errChange error
		ChangesDataJSON, errChange = json.Marshal(ChangesData)
		if errChange != nil {
			// Corrected log message
			log.Printf("ERROR: Failed to marshal ChangesData to JSON: %v\n", errChange)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		transData := 29

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

	}

	// if err := tx.Commit().Error; err != nil {
	// 	log.Printf("ERROR: Failed to commit transaction: %v\n", err)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	return true, "Successfully Changes Saved"
}

func AnswerTechnicianIntakeService(db *gorm.DB, reqVal model.AnswerReportIntakeReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db
	// tx := db.Begin()
	// if tx.Error != nil {
	// 	log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
	// 	return false, "Something went wrong, Try Again"
	// }

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
	// 		tx.Rollback()
	// 	}
	// }()

	//Checking the Question ID is Available
	var TechnicianIntakeFormData []model.GetTechnicianIntakeData

	TechnicianIntakeFormDataErr := db.Raw(query.GetTechnicianIntakeFormQuestionSQL, reqVal.AppointmentId, reqVal.QuestionId).Scan(&TechnicianIntakeFormData).Error
	if TechnicianIntakeFormDataErr != nil {
		log.Fatal(TechnicianIntakeFormDataErr)
	}

	//If AVaiable Need to Update, else Create a QuestionID and Answer
	if len(TechnicianIntakeFormData) > 0 {

		//Find the If any Changes is Avaiable
		oldData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): hashdb.Decrypt(TechnicianIntakeFormData[0].Answer),
		}

		updatedData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): reqVal.Answer,
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

			//Insert Aduit Row for Answers Update
			transData := 27
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

			//Update the Answer with QuestionID
			UpdateTechnicianInputErr := tx.Exec(
				query.UpdateTechnicianIntakeSQL,
				hashdb.Encrypt(reqVal.Answer),
				timeZone.GetPacificTime(),
				idValue,
				reqVal.QuestionId,
				reqVal.AppointmentId,
			).Error
			if UpdateTechnicianInputErr != nil {
				log.Printf("ERROR: Failed to Update Technician Input: %v\n", UpdateTechnicianInputErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

		}

	} else {
		return false, "Invalid Question ID"
	}

	// if err := tx.Commit().Error; err != nil {
	// 	log.Printf("ERROR: Failed to commit transaction: %v\n", err)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	return true, "Successfully Changes Saved"
}

func AnswerPatientIntakeService(db *gorm.DB, reqVal model.AnswerReportIntakeReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db
	// tx := db.Begin()
	// if tx.Error != nil {
	// 	log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
	// 	return false, "Something went wrong, Try Again"
	// }

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
	// 		tx.Rollback()
	// 	}
	// }()

	//Checking the Question ID is Available
	var PatientIntakeFormData []model.GetViewIntakeData

	PatientIntakeFormDataErr := db.Raw(query.GetPatientIntakeFormQuestionSQL, reqVal.AppointmentId, reqVal.QuestionId).Scan(&PatientIntakeFormData).Error
	if PatientIntakeFormDataErr != nil {
		log.Fatal(PatientIntakeFormDataErr)
	}

	//If AVaiable Need to Update, else Create a QuestionID and Answer
	if len(PatientIntakeFormData) > 0 {

		//Find the If any Changes is Avaiable
		oldData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): hashdb.Decrypt(PatientIntakeFormData[0].Answer),
		}

		updatedData := map[string]interface{}{
			fmt.Sprintf("%d", reqVal.QuestionId): reqVal.Answer,
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

			//Insert Aduit Row for Answers Update
			transData := 24
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

			//Update the Answer with QuestionID
			UpdatePatientInputErr := tx.Exec(
				query.UpdatePatientIntakeSQL,
				hashdb.Encrypt(reqVal.Answer),
				timeZone.GetPacificTime(),
				idValue,
				reqVal.QuestionId,
				reqVal.AppointmentId,
			).Error
			if UpdatePatientInputErr != nil {
				log.Printf("ERROR: Failed to Update Technician Input: %v\n", UpdatePatientInputErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}

		}

	} else {
		return false, "Invalid Question ID"
	}

	// if err := tx.Commit().Error; err != nil {
	// 	log.Printf("ERROR: Failed to commit transaction: %v\n", err)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	return true, "Successfully Changes Saved"
}

func AnswerTextContentService(db *gorm.DB, reqVal model.AnswerTextContentReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db
	// tx := db.Begin()
	// if tx.Error != nil {
	// 	log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
	// 	return false, "Something went wrong, Try Again"
	// }

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
	// 		tx.Rollback()
	// 	}
	// }()

	//Checking the Question ID is Available
	var ReportTextContent []model.GetReportTextContentModel

	ReportTextContentErr := db.Raw(query.GetTextContentSQL, reqVal.AppointmentId).Scan(&ReportTextContent).Error
	if ReportTextContentErr != nil {
		log.Fatal(ReportTextContentErr)
	}

	//If AVaiable Need to Update, else Create a QuestionID and Answer
	if len(ReportTextContent) > 0 {

		//Update the Text Content
		UpdateTextContentErr := tx.Exec(
			query.UpdateTextContentSQL,
			hashdb.Encrypt(reqVal.TextContent),
			timeZone.GetPacificTime(),
			idValue,
			reqVal.SyncStatus,
			reqVal.AppointmentId,
		).Error
		if UpdateTextContentErr != nil {
			log.Printf("ERROR: Failed to Update Text Content: %v\n", UpdateTextContentErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		transData := 32
		errTrans := model.RefTransHistory{
			TransTypeId: transData,
			THData:      "Text Content Updated",
			UserId:      reqVal.PatientId,
			THActionBy:  idValue,
		}

		errTransStatus := db.Create(&errTrans).Error
		if errTransStatus != nil {
			log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
			return false, "Something went wrong, Try Again"
		}

	} else {

		//Inserting a New Text Content
		InsertTextContentErr := tx.Exec(
			query.InsertTextContentSQL,
			reqVal.PatientId,
			reqVal.AppointmentId,
			hashdb.Encrypt(reqVal.TextContent),
			timeZone.GetPacificTime(),
			idValue,
			reqVal.SyncStatus,
		).Error
		if InsertTextContentErr != nil {
			log.Printf("ERROR: Failed to Insert Text Content: %v\n", InsertTextContentErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}

		transData := 31
		errTrans := model.RefTransHistory{
			TransTypeId: transData,
			THData:      "Text Content Created",
			UserId:      reqVal.PatientId,
			THActionBy:  idValue,
		}

		errTransStatus := db.Create(&errTrans).Error
		if errTransStatus != nil {
			log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
			return false, "Something went wrong, Try Again"
		}

	}

	// if err := tx.Commit().Error; err != nil {
	// 	log.Printf("ERROR: Failed to commit transaction: %v\n", err)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	return true, "Successfully Changes Saved"
}

func AddCommentsService(db *gorm.DB, reqVal model.AddCommentReq, idValue int) (bool, string) {
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

	//Adding Comments
	InsertCommentsErr := tx.Exec(
		query.InsertCommentsSQL,
		reqVal.PatientId,
		reqVal.AppointmentId,
		idValue,
		reqVal.AssignId,
		hashdb.Encrypt(reqVal.Status),
		hashdb.Encrypt(reqVal.Comments),
		timeZone.GetPacificTime(),
	).Error
	if InsertCommentsErr != nil {
		log.Printf("ERROR: Failed to Insert Comments: %v\n", InsertCommentsErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Changes Saved"
}

func CompleteReportService(db *gorm.DB, reqVal model.CompleteReportReq, idValue int) (bool, string) {
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

	//Updating the Appointment Status
	UpdateAppointementErr := tx.Exec(
		query.CompleteReportAppointmentSQL,
		reqVal.MovedStatus,
		false,
		nil,
		reqVal.AppointmentId,
		reqVal.PatientId,
	).Error
	if UpdateAppointementErr != nil {
		log.Printf("ERROR: Failed to Update Appointement: %v\n", UpdateAppointementErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Inserting the Audit for the Report Status
	ReportStatustransData := 25
	ReportStatusTransDataErr := model.RefTransHistory{
		TransTypeId: ReportStatustransData,
		THData:      "Report Finalized from " + reqVal.CurrentStatus,
		UserId:      reqVal.PatientId,
		THActionBy:  idValue,
	}

	errTransStatus := db.Create(&ReportStatusTransDataErr).Error
	if errTransStatus != nil {
		log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
		return false, "Something went wrong, Try Again"
	}

	//Inserting the Audit for the Report Accessing

	oldDataCat := map[string]interface{}{
		"Report Access ID": idValue,
	}

	updatedDataCat := map[string]interface{}{
		"Report Access ID": "",
	}

	ChangesDataCat := helper.GetChanges(updatedDataCat, oldDataCat)

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

	errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.PatientId), int(idValue), string(ChangesDataJSON)).Error
	if errTrans != nil {
		log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	// //Updating the End Time For the Report History
	// ReportHistoryErr := tx.Exec(
	// 	query.CompleteReportHistorySQL,
	// 	timeZone.GetPacificTime(),
	// 	reqVal.AppointmentId,
	// 	idValue,
	// 	reqVal.PatientId,
	// ).Error
	// if ReportHistoryErr != nil {
	// 	log.Printf("ERROR: Failed to Update Report History: %v\n", ReportHistoryErr)
	// 	tx.Rollback()
	// 	return false, "Something went wrong, Try Again"
	// }

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Changes Saved"
}

func AutosaveServicee(db *gorm.DB, reqVal model.AutoSubmitReportReq, idValue int, roleIdValue int) (bool, string, []model.GetReportIntakeData, []model.GetReportTextContent, []model.GetOneUserAppointmentModel, bool) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	fmt.Println("%%%%%%%%%%", reqVal.ChangedOneState)

	//Inserting and Upadating the Report Intake Form
	for _, data := range reqVal.ReportIntakeForm {
		status, message := AnswerReportIntakeService(tx, model.AnswerReportIntakeReq{
			PatientId:     reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			QuestionId:    data.QuestionId,
			Answer:        data.Answer,
		}, idValue)

		if !status {
			return status, message, []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Updating the Report Text Content
	if reqVal.ChangedOneState.ReportTextContent {

		var updateAutoerr = tx.Exec(query.UpdateAutosaveTextContentSQL,
			hashdb.Encrypt(reqVal.ReportTextContent),
			timeZone.GetPacificTime(),
			idValue,
			reqVal.AppointmentId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave Text Content: %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}

	}

	//Update Report Sync Status
	if reqVal.ChangedOneState.SyncStatus {

		var updateAutoerr = tx.Exec(query.UpdateAutosaveSyncStatusSQL,
			reqVal.SyncStatus,
			timeZone.GetPacificTime(),
			idValue,
			reqVal.AppointmentId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}

	}

	//Update the Impression
	if reqVal.ChangedOneState.Impression {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveImpressionSQL,
			reqVal.Impression,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the Recommendation
	if reqVal.ChangedOneState.Recommendation {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveRecommendationSQL,
			reqVal.Recommendation,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the ImpressionAddtional
	if reqVal.ChangedOneState.ImpressionAddtional {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveImpressionAddtionalSQL,
			reqVal.ImpressionAddtional,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the RecommendationAddtional
	if reqVal.ChangedOneState.RecommendationAddtional {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveRecommendationAddtionalSQL,
			reqVal.RecommendationAddtional,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the CommonImpressionRecommendation
	if reqVal.ChangedOneState.CommonImpressionRecommendation {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveCommonImpressionRecommendationSQL,
			reqVal.CommonImpressionRecommendation,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the ImpressionRight
	if reqVal.ChangedOneState.ImpressionRight {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveImpressionRightSQL,
			reqVal.ImpressionRight,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the RecommendationRight
	if reqVal.ChangedOneState.RecommendationRight {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveRecommendationRightSQL,
			reqVal.RecommendationRight,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the ImpressionAddtionalRight
	if reqVal.ChangedOneState.ImpressionAddtionalRight {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveImpressionAddtionalRightSQL,
			reqVal.ImpressionAddtionalRight,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the RecommendationAddtionalRight
	if reqVal.ChangedOneState.RecommendationAddtionalRight {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveRecommendationAddtionalRightSQL,
			reqVal.RecommendationAddtionalRight,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the CommonImpressionRecommendationRight
	if reqVal.ChangedOneState.CommonImpressionRecommendationRight {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveCommonImpressionRecommendationRightRightSQL,
			reqVal.CommonImpressionRecommendationRight,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the ArtificatsLeft
	if reqVal.ChangedOneState.ArtificatsLeft {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveArtificatsLeftSQL,
			reqVal.ArtificatsLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the ArtificatsLeft
	if reqVal.ChangedOneState.ArtificatsRight {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveArtificatsLeftSQL,
			reqVal.ArtificatsRight,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the PatientHistory
	if reqVal.ChangedOneState.PatientHistory {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosavePatientHistorySQL,
			reqVal.PatientHistory,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the BreastImplantsImagetext
	if reqVal.ChangedOneState.BreastImplantImageText {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveBreastImplantsImagetextSQL,
			reqVal.BreastImplantsImagetext,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the SymmetryImageText
	if reqVal.ChangedOneState.SymmetryImageText {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveSymmetryImageTextSQL,
			reqVal.SymmetryImageText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the BreastdensityImageText
	if reqVal.ChangedOneState.BreastDensityImageText {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveBreastDensityImageTextSQL,
			reqVal.BreastdensityImageText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the NippleAreolaImageText
	if reqVal.ChangedOneState.NippleAreolaImageText {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveNippleAreolaImageTextSQL,
			reqVal.NippleAreolaImageText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the GlandularImageText
	if reqVal.ChangedOneState.GlandularImageText {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveGlandularImageTextSQL,
			reqVal.GlandularImageText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the LymphnodesImageText
	if reqVal.ChangedOneState.LymphNodesImageText {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveLymphnodesImageTextSQL,
			reqVal.LymphnodesImageText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the BreastdensityImageTextLeft
	if reqVal.ChangedOneState.BreastDensityImageTextLeft {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveBreastDensityImageTextLeftSQL,
			reqVal.BreastdensityImageTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the NippleAreolaImageTextLeft
	if reqVal.ChangedOneState.NippleAreolaImageTextLeft {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveNippleAreolaImageTextLeftSQL,
			reqVal.NippleAreolaImageTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the GlandularImageTextLeft
	if reqVal.ChangedOneState.GlandularImageTextLeft {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveGlandularImageTextLeftSQL,
			reqVal.GlandularImageTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	//Update the LymphnodesImageTextLeft
	if reqVal.ChangedOneState.LymphNodesImageTextLeft {
		var updateAutoerr = tx.Exec(
			query.UpdateAutosaveLymphNodesImageTextLeftSQL,
			reqVal.LymphnodesImageTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateAutoerr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateAutoerr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false
	}

	//Report Intake Form Table
	var ReportIntakeFormData []model.GetReportIntakeData
	ReportIntakeFormDataerr := db.Raw(query.GetReportIntakeFormSQL, reqVal.AppointmentId).Scan(&ReportIntakeFormData).Error
	if ReportIntakeFormDataerr != nil {
		log.Fatal(ReportIntakeFormDataerr)
	}

	//Decrypt Report Intake Form Table
	for i, data := range ReportIntakeFormData {
		ReportIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
	}

	//Report Text Content Table
	var ReportTextContentData []model.GetReportTextContent
	ReportTextContentDataerr := db.Raw(query.GetReporttextContent, reqVal.AppointmentId).Scan(&ReportTextContentData).Error
	if ReportTextContentDataerr != nil {
		log.Fatal(ReportTextContentDataerr)
	}

	//Decrypt Report Text Content Table
	for i, data := range ReportTextContentData {
		ReportTextContentData[i].TextContent = hashdb.Decrypt(data.TextContent)
	}

	var OneUserAppointment []model.GetOneUserAppointmentModel
	//Appointment Table
	ViewAppointmentErr := db.Raw(query.GetOneUserAppointment, reqVal.PatientId, reqVal.AppointmentId).Scan(&OneUserAppointment).Error
	if ViewAppointmentErr != nil {
		log.Fatal(ViewAppointmentErr)
	}

	//Decrypt Appointment Table
	for i, data := range OneUserAppointment {
		OneUserAppointment[i].SCName = hashdb.Decrypt(data.SCName)
	}

	var EaseQTReportAccess = false

	//Get the Ease QT Report Access Status
	switch roleIdValue {
	case 1: //Master Admin
		EaseQTReportAccess = true
	case 2: //Scan Center Technician
		EaseQTReportAccess = false
	case 3: //Scan Center Manager
		EaseQTReportAccess = false
	case 4: //Patient
		EaseQTReportAccess = false
	case 5: //Scan Center Doctor

		var ReportStatus []model.DoctorReportAccessStatus
		err := db.Raw(query.DoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
		if err != nil {
			log.Fatal(err)
		}

		if len(ReportStatus) == 0 || ReportStatus[0].DDEaseQTReportAccess == nil {
			EaseQTReportAccess = false
			break
		}

		EaseQTReportAccess = *ReportStatus[0].DDEaseQTReportAccess

	case 6: //Junior Doctor
		EaseQTReportAccess = true
	case 7: //Scribe
		EaseQTReportAccess = true
	case 8: //Scan Center Reviewer
		var ReportStatus []model.CoDoctorReportAccessStatus
		err := db.Raw(query.CoDoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
		if err != nil {
			log.Fatal(err)
		}

		if len(ReportStatus) == 0 || ReportStatus[0].CDEaseQTReportAccess == nil {
			EaseQTReportAccess = false
			break
		}

		EaseQTReportAccess = *ReportStatus[0].CDEaseQTReportAccess
	case 9: //Manager
		EaseQTReportAccess = true
	case 10: //Performing Provider
		EaseQTReportAccess = true
	default:
		EaseQTReportAccess = false
	}

	return true, "Successfully Changes Saved", ReportIntakeFormData, ReportTextContentData, OneUserAppointment, EaseQTReportAccess
}

func SubmitReportService(db *gorm.DB, reqVal model.SubmitReportReq, idValue int, roleIdValue int) (bool, string) {
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

	//Inserting and Upadating the Report Intake Form
	for _, data := range reqVal.ReportIntakeForm {
		status, message := AnswerReportIntakeService(tx, model.AnswerReportIntakeReq{
			PatientId:     reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			QuestionId:    data.QuestionId,
			Answer:        data.Answer,
		}, idValue)

		if !status {
			return status, message
		}
	}

	// //Updating the TechnicianIntake Form
	// for _, data := range reqVal.TechnicianIntakeForm {
	// 	status, message := AnswerTechnicianIntakeService(tx, model.AnswerReportIntakeReq{
	// 		PatientId:     reqVal.PatientId,
	// 		AppointmentId: reqVal.AppointmentId,
	// 		QuestionId:    data.QuestionId,
	// 		Answer:        data.Answer,
	// 	}, idValue)

	// 	if !status {
	// 		return status, message
	// 	}
	// }

	// //Updating the PatientIntake Form
	// for _, data := range reqVal.PatientIntakeForm {
	// 	status, message := AnswerPatientIntakeService(tx, model.AnswerReportIntakeReq{
	// 		PatientId:     reqVal.PatientId,
	// 		AppointmentId: reqVal.AppointmentId,
	// 		QuestionId:    data.QuestionId,
	// 		Answer:        data.Answer,
	// 	}, idValue)

	// 	if !status {
	// 		return status, message
	// 	}
	// }

	//Updating the Report Text Content
	status, message := AnswerTextContentService(tx, model.AnswerTextContentReq{
		PatientId:     reqVal.PatientId,
		AppointmentId: reqVal.AppointmentId,
		TextContent:   reqVal.ReportTextContent,
		SyncStatus:    reqVal.SyncStatus,
	}, idValue)

	if !status {
		return status, message
	}

	var reportStatus = reqVal.CurrentStatus

	if reqVal.LeaveStatus {
		reportStatus = "Changes"
	}

	if roleIdValue == 7 {
		//Updating the Appointment Status
		UpdateAppointementErr := tx.Exec(
			query.ScribeCompleteReportAppointmentSQL,
			reqVal.MovedStatus,
			reqVal.Impression,
			reqVal.Recommendation,
			reqVal.ImpressionAddtional,
			reqVal.RecommendationAddtional,
			reqVal.CommonImpressionRecommendation,
			reqVal.ImpressionRight,
			reqVal.RecommendationRight,
			reqVal.ImpressionAddtionalRight,
			reqVal.RecommendationAddtionalRight,
			reqVal.CommonImpressionRecommendationRight,
			reqVal.ArtificatsLeft,
			reqVal.ArtificatsRight,
			reqVal.PatientHistory,
			reqVal.BreastImplantsImagetext,
			reqVal.SymmetryImageText,
			reqVal.BreastdensityImageText,
			reqVal.NippleAreolaImageText,
			reqVal.GlandularImageText,
			reqVal.LymphnodesImageText,
			reqVal.BreastdensityImageTextLeft,
			reqVal.NippleAreolaImageTextLeft,
			reqVal.GlandularImageTextLeft,
			reqVal.LymphnodesImageTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error
		if UpdateAppointementErr != nil {
			log.Printf("ERROR: Failed to Update Appointement: %v\n", UpdateAppointementErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	} else {
		//Updating the Appointment Status
		UpdateAppointementErr := tx.Exec(
			query.CompleteReportAppointmentSQL,
			reqVal.MovedStatus,
			reqVal.Impression,
			reqVal.Recommendation,
			reqVal.ImpressionAddtional,
			reqVal.RecommendationAddtional,
			reqVal.CommonImpressionRecommendation,
			reqVal.ImpressionRight,
			reqVal.RecommendationRight,
			reqVal.ImpressionAddtionalRight,
			reqVal.RecommendationAddtionalRight,
			reqVal.CommonImpressionRecommendationRight,
			reqVal.ArtificatsLeft,
			reqVal.ArtificatsRight,
			reqVal.PatientHistory,
			reqVal.BreastImplantsImagetext,
			reqVal.SymmetryImageText,
			reqVal.BreastdensityImageText,
			reqVal.NippleAreolaImageText,
			reqVal.GlandularImageText,
			reqVal.LymphnodesImageText,
			reqVal.BreastdensityImageTextLeft,
			reqVal.NippleAreolaImageTextLeft,
			reqVal.GlandularImageTextLeft,
			reqVal.LymphnodesImageTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error
		if UpdateAppointementErr != nil {
			log.Printf("ERROR: Failed to Update Appointement: %v\n", UpdateAppointementErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	//Inserting the Audit for the Report Status
	ReportStatustransData := 25
	ReportStatusTransDataErr := model.RefTransHistory{
		TransTypeId: ReportStatustransData,
		THData:      "Report Finalized from " + reportStatus,
		UserId:      reqVal.PatientId,
		THActionBy:  idValue,
	}

	errTransStatus := db.Create(&ReportStatusTransDataErr).Error
	if errTransStatus != nil {
		log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
		return false, "Something went wrong, Try Again"
	}

	//Inserting the Audit for the Report Accessing
	oldDataCat := map[string]interface{}{
		"Report Access ID": idValue,
	}

	updatedDataCat := map[string]interface{}{
		"Report Access ID": "",
	}

	ChangesDataCat := helper.GetChanges(updatedDataCat, oldDataCat)

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

	errTrans := tx.Exec(query.InsertTransactionDataSQL, int(transData), int(reqVal.PatientId), int(idValue), string(ChangesDataJSON)).Error
	if errTrans != nil {
		log.Printf("ERROR: Failed to Transaction History: %v\n", errTrans)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Updating the End Time For the Report History
	ReportHistoryErr := tx.Exec(
		query.CompleteReportHistorySQL,
		timeZone.GetPacificTime(),
		reportStatus,
		hashdb.Encrypt(reqVal.ReportTextContent),
		reqVal.AppointmentId,
		idValue,
		reqVal.PatientId,
	).Error
	if ReportHistoryErr != nil {
		log.Printf("ERROR: Failed to Update Report History: %v\n", ReportHistoryErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//totalCorrectEdit
	switch roleIdValue {
	case 1:
		var ListUserData []model.ListUserModel

		ListUserDataErr := db.Raw(query.ListUserDataSQL, reqVal.PatientId, reqVal.AppointmentId, 6).Scan(&ListUserData).Error
		if ListUserDataErr != nil {
			log.Fatal(ListUserDataErr.Error())
			return false, "Something went wrong, Try Again"
		}

		var correct = 0
		var edit = 0

		if reqVal.EditStatus {
			edit = 1
		} else {
			correct = 1
		}

		if len(ListUserData) > 0 {
			// for _, data := range ListUserData {
			UpdateChangesErr := tx.Exec(
				query.UpdateCorrectEditSQL,
				correct,
				edit,
				ListUserData[0].RHId,
			).Error
			if UpdateChangesErr != nil {
				log.Printf("ERROR: Failed to Update Report History: %v\n", UpdateChangesErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}
			// }
		}
	case 8:
		var ListUserData []model.ListUserModel

		ListUserDataErr := db.Raw(query.ListUserDataSQL, reqVal.PatientId, reqVal.AppointmentId, 1).Scan(&ListUserData).Error
		if ListUserDataErr != nil {
			log.Fatal(ListUserDataErr.Error())
			return false, "Something went wrong, Try Again"
		}

		var correct = 0
		var edit = 0

		if reqVal.EditStatus {
			edit = 1
		} else {
			correct = 1
		}

		if len(ListUserData) > 0 {
			// for _, data := range ListUserData {
			UpdateChangesErr := tx.Exec(
				query.UpdateCorrectEditSQL,
				correct,
				edit,
				ListUserData[0].RHId,
			).Error
			if UpdateChangesErr != nil {
				log.Printf("ERROR: Failed to Update Report History: %v\n", UpdateChangesErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}
			// }

		}

	}

	//Send Mail for the Patient
	if reqVal.PatientMailStatus {

		var PatientdataModel []model.Patientdata

		err := db.Raw(query.GetPatientData, reqVal.PatientId, reqVal.AppointmentId).Scan(&PatientdataModel).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return false, "Something went wrong, Try Again"
		}

		for i, data := range PatientdataModel {
			PatientdataModel[i].UserFirstName = hashdb.Decrypt(data.UserFirstName)
		}

		htmlContent := mailservice.PatientReportSignOff(PatientdataModel[0].UserFirstName, PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, PatientdataModel[0].SCCustId)

		subject := "Your Report Status"

		emailStatus := mailservice.MailService(PatientdataModel[0].Email, htmlContent, subject)

		if !emailStatus {
			log.Error("Sending Mail Meets Error")
			return false, "Something went wrong, Try Again"
		}
	}

	//Send Mail for the Scan Center Manager
	if reqVal.ManagerMailStatus {
		var PatientdataModel []model.Patientdata

		err := db.Raw(query.GetPatientData, reqVal.PatientId, reqVal.AppointmentId).Scan(&PatientdataModel).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return false, "Something went wrong, Try Again"
		}

		var ManagerModel []model.ManagerData

		Managererr := db.Raw(query.GetManagerData, 3, reqVal.AppointmentId).Scan(&ManagerModel).Error
		if Managererr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", Managererr)
			return false, "Something went wrong, Try Again"
		}

		for _, data := range ManagerModel {
			htmlContent := mailservice.PatientReportSignOff(PatientdataModel[0].UserFirstName, PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, data.SCCustId)

			subject := "Your Report Status"

			emailStatus := mailservice.MailService(data.Email, htmlContent, subject)

			if !emailStatus {
				log.Error("Sending Mail Meets Error")
				return false, "Something went wrong, Try Again"
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Report Submitted"
}

func UpdateRemarksService(db *gorm.DB, reqVal model.UpdateRemarkReq, idValue int) (bool, string) {
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

	//Updating Remarks
	UpdateRemarksErr := tx.Exec(
		query.UpdateReportRemarksSQL,
		reqVal.Remark,
		reqVal.AppointmentId,
		reqVal.PatientId,
	).Error
	if UpdateRemarksErr != nil {
		log.Printf("ERROR: Failed to Insert Comments: %v\n", UpdateRemarksErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//adding Remarks
	AddingRemarksErr := tx.Exec(
		query.InsertRemark,
		reqVal.AppointmentId,
		idValue,
		hashdb.Encrypt(reqVal.Remark),
		timeZone.GetPacificTime(),
	).Error
	if AddingRemarksErr != nil {
		log.Printf("ERROR: Failed to Insert Comments: %v\n", AddingRemarksErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Updating Audits
	transData := 34
	errTrans := model.RefTransHistory{
		TransTypeId: transData,
		THData:      hashdb.Encrypt(reqVal.Remark),
		UserId:      reqVal.PatientId,
		THActionBy:  idValue,
	}

	errTransStatus := db.Create(&errTrans).Error
	if errTransStatus != nil {
		log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Changes Saved"
}

func UploadReportFormateService(db *gorm.DB, reqVal model.UploadReportFormateReq, idValue int) (int, bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return 0, false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var insertedID int

	//Adding Template
	InsertReportTemplateErr := tx.Raw(
		query.InsertReportTemplate,
		hashdb.Encrypt(reqVal.Name),
		hashdb.Encrypt(reqVal.FormateTemplate),
		timeZone.GetPacificTime(),
		idValue,
	).Scan(&insertedID).Error
	if InsertReportTemplateErr != nil {
		log.Printf("ERROR: Failed to Insert Report Template: %v\n", InsertReportTemplateErr)
		tx.Rollback()
		return 0, false, "Something went wrong, Try Again"
	}

	//Addding audit For the Template
	transData := 35
	errTrans := model.RefTransHistory{
		TransTypeId: transData,
		THData:      "Report Template Added Successfully",
		UserId:      idValue,
		THActionBy:  idValue,
	}

	errTransStatus := db.Create(&errTrans).Error
	if errTransStatus != nil {
		log.Error("errreportStatus INSERT ERROR at Trnasaction: " + errTransStatus.Error())
		return 0, false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return 0, false, "Something went wrong, Try Again"
	}

	return insertedID, true, "Successfully Changes Saved"
}

func GetReportFormateService(db *gorm.DB, reqVal model.GetReportFormateReq, idValue int) []model.ReportTextFormateModel {
	log := logger.InitLogger()

	//Getting the Template
	var TemplateFormate []model.ReportTextFormateModel
	err := db.Raw(query.GetOneReportFormateListSQL, reqVal.Id).Scan(&TemplateFormate).Error
	if err != nil {
		log.Fatal(err)
	}

	for i, data := range TemplateFormate {
		TemplateFormate[i].RFText = hashdb.Decrypt(data.RFText)
		TemplateFormate[i].RFName = hashdb.Decrypt(data.RFName)
	}

	return TemplateFormate
}

func ListRemarkService(db *gorm.DB, reqVal model.ListRemarkReq) []model.ListRemarkModel {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return []model.ListRemarkModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var ListRewardModel []model.ListRemarkModel

	ListRewardErr := db.Raw(query.ListRemarkSQL, reqVal.AppointmentId).Scan(&ListRewardModel).Error
	if ListRewardErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", ListRewardErr)
		return []model.ListRemarkModel{}
	}

	for i, list := range ListRewardModel {
		ListRewardModel[i].RemarksMessage = hashdb.Decrypt(list.RemarksMessage)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return []model.ListRemarkModel{}
	}

	return ListRewardModel
}

func SendMailReportService(db *gorm.DB, reqVal model.SendMailReportReq) (bool, string) {
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

	//Update the Mail send Status in Appointment
	UpdateMailStatusErr := tx.Exec(
		query.UpdateMailStatusSQL,
		"sended",
		reqVal.AppointmentId,
		reqVal.PatientId,
	).Error
	if UpdateMailStatusErr != nil {
		log.Printf("ERROR: Failed to Update Mail Status: %v\n", UpdateMailStatusErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Send Mail for the Patient

	var PatientdataModel []model.Patientdata

	err := db.Raw(query.GetPatientData, reqVal.PatientId, reqVal.AppointmentId).Scan(&PatientdataModel).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return false, "Something went wrong, Try Again"
	}

	for i, data := range PatientdataModel {
		PatientdataModel[i].UserFirstName = hashdb.Decrypt(data.UserFirstName)
	}

	if reqVal.PatientMailStatus {

		htmlContent := mailservice.PatientReportSignOff(PatientdataModel[0].UserFirstName, PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, PatientdataModel[0].SCCustId)

		subject := "Your Report Status"

		emailStatus := mailservice.MailService(PatientdataModel[0].Email, htmlContent, subject)

		if !emailStatus {
			log.Error("Sending Mail Meets Error")
			return false, "Something went wrong, Try Again"
		}
	}

	//Send Mail for the Scan Center Manager
	if reqVal.ManagerMailStatus {
		// var PatientdataModel []model.Patientdata

		// err := db.Raw(query.GetPatientData, reqVal.PatientId, reqVal.AppointmentId).Scan(&PatientdataModel).Error
		// if err != nil {
		// 	log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		// 	return false, "Something went wrong, Try Again"
		// }

		var ManagerModel []model.ManagerData

		Managererr := db.Raw(query.GetManagerData, 3, reqVal.AppointmentId).Scan(&ManagerModel).Error
		if Managererr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", Managererr)
			return false, "Something went wrong, Try Again"
		}

		for _, data := range ManagerModel {
			htmlContent := mailservice.PatientReportSignOff(PatientdataModel[0].UserFirstName, PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, data.SCCustId)

			subject := "Your Report Status"

			emailStatus := mailservice.MailService(data.Email, htmlContent, subject)

			if !emailStatus {
				log.Error("Sending Mail Meets Error")
				return false, "Something went wrong, Try Again"
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Mail Sended !"

}

func DownloadReportService(db *gorm.DB, reqVal model.DownloadReportReq) model.GetViewIntakeData {
	log := logger.InitLogger()

	var PatientFile model.GetViewIntakeData

	err := db.Raw(query.DownloadReportSQL, reqVal.Id).Scan(&PatientFile).Error
	if err != nil {
		log.Fatal(err)
	}

	if len(PatientFile.Answer) > 0 {
		DriversLicenseNoImgHelperData, viewErr := helperView.ViewFile("./Assets/Files/" + hashdb.Decrypt(PatientFile.Answer))
		if viewErr != nil {
			// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
			log.Fatalf("Failed to read DrivingLicense file: %v", viewErr)
		}

		// for i, data := range PatientFile {
		PatientFile.File = &model.FileData{
			Base64Data:  DriversLicenseNoImgHelperData.Base64Data,
			ContentType: DriversLicenseNoImgHelperData.ContentType,
		}
	}

	// }

	return PatientFile
}

func ListAddendumService(db *gorm.DB, appointmentId int) []model.AddAddendumModel {
	log := logger.InitLogger()

	var ListAddendumModel []model.AddAddendumModel

	ListAddendumErr := db.Raw(query.ListAddendumSQL, appointmentId).Scan(&ListAddendumModel).Error
	if ListAddendumErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", ListAddendumErr)
		return []model.AddAddendumModel{}
	}

	return ListAddendumModel

}

func AddAddendumService(db *gorm.DB, reqVal model.AddAddendumReq, idValue int) (bool, string, []model.AddAddendumModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	InsertErr := tx.Exec(
		query.InsertAddedumSQL,
		reqVal.AppointmentId,
		idValue,
		reqVal.AddAddendumText,
		timeZone.GetPacificTime(),
	).Error
	if InsertErr != nil {
		log.Printf("ERROR: Failed to Insert Addendum: %v\n", InsertErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
	}

	return true, "Successfully Addedum Received", ListAddendumService(db, reqVal.AppointmentId)

}
