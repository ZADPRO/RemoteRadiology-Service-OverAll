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

func CheckAccessService(db *gorm.DB, reqVal model.CheckAccessReq, idValue int) (bool, string, int) {
	log := logger.InitLogger()

	var result []model.AccessStatus

	err := db.Raw(query.CheckAccessSQL, idValue, reqVal.AppointmentId).Scan(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	var message = "Another User Already Accessing it"

	if result[0].Status {
		message = "Report are Available for Edit"
	}

	return result[0].Status, message, result[0].RefAppointmentAccessId
}

func AssignGetReportService(db *gorm.DB, reqVal model.AssignGetReportReq, idValue int) (bool, string, []model.GetViewIntakeData, []model.GetTechnicianIntakeData, []model.GetReportIntakeData, []model.GetReportTextContent, []model.GetReportHistory, []model.GetReportComments, []model.GetOneUserAppointmentModel, []model.ReportFormateModel, []model.GetUserDetails, []model.PatientCustId) {
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
			[]model.PatientCustId{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	fmt.Println("Check1---------------------------")
	//patientCustId
	var PatientUserDetails []model.PatientCustId
	PatientUserDetailsErr := db.Raw(query.PatientUserDetailsSQL, reqVal.PatientId).Scan(&PatientUserDetails).Error
	if PatientUserDetailsErr != nil {
		log.Fatal(PatientUserDetailsErr)
	}

	fmt.Println("Check2---------------------------")
	//GetUserDetails
	var UserDetails []model.GetUserDetails
	UserDetailsErr := db.Raw(query.GetUserDetailsSQL, idValue).Scan(&UserDetails).Error
	if UserDetailsErr != nil {
		log.Fatal(UserDetailsErr)
	}

	fmt.Println("Check3---------------------------")
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

	fmt.Println("Check4---------------------------")
	checkAccessReq := model.CheckAccessReq{
		AppointmentId: reqVal.AppointmentId,
	}

	fmt.Println("Check5---------------------------")
	status, message, _ := CheckAccessService(db, checkAccessReq, idValue)

	fmt.Println(status, message)

	fmt.Println("Check6---------------------------")
	if (status && !reqVal.ReadOnlyStatus) || (!status && reqVal.ReadOnlyStatus) || (status && reqVal.ReadOnlyStatus) {

		fmt.Println("Check7---------------------------")
		//Appointment Table
		var Appointment []model.AppointmentModel
		Appointmenterr := db.Raw(query.GetAppointmentSQL, reqVal.AppointmentId).Scan(&Appointment).Error
		if Appointmenterr != nil {
			log.Fatal(Appointmenterr)
		}

		fmt.Println("Check8---------------------------")
		if !reqVal.ReadOnlyStatus {
			oldDataCat := map[string]interface{}{
				"Report Access ID": Appointment[0].AppointmentAccessId,
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
						[]model.PatientCustId{}
				}

				transData := 28
				fmt.Println("Check9---------------------------")
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
						[]model.PatientCustId{}
				}

				fmt.Println("Check10---------------------------")
				categoryUpdate := tx.Exec(
					query.UpdateAccessAppointment,
					true,
					idValue,
					reqVal.AppointmentId,
				).Error
				fmt.Println("Check11---------------------------")
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
						[]model.PatientCustId{}
				}
				fmt.Println("Check12--------------------------")
				//List the Latest Report History
				var ReportHistory []model.GetReportHistory
				ListReportHistoryErr := db.Raw(query.CheckLatestReportSQL, reqVal.AppointmentId, reqVal.PatientId).Scan(&ReportHistory).Error
				if ListReportHistoryErr != nil {
					log.Fatal(ListReportHistoryErr)
				}
				fmt.Println("Check13---------------------------")
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
						[]model.PatientCustId{}
				}

				fmt.Println("Check14---------------------------")
				if len(ReportHistory) > 0 {
					//Insert the History
					ReportHistoryErr := tx.Exec(
						query.InsertReportHistorySQL,
						reqVal.PatientId,
						reqVal.AppointmentId,
						idValue,
						ReportHistory[0].HandleEndTime,
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
							[]model.PatientCustId{}
					}
				}

			}

		}

		var IntakeFormData []model.GetViewIntakeData
		var OneUserAppointment []model.GetOneUserAppointmentModel
		fmt.Println("Check15---------------------------")
		//Appointment Table
		ViewAppointmentErr := db.Raw(query.GetOneUserAppointment, reqVal.PatientId, reqVal.AppointmentId).Scan(&OneUserAppointment).Error
		if ViewAppointmentErr != nil {
			log.Fatal(ViewAppointmentErr)
		}
		fmt.Println("Check16---------------------------")

		//Decrypt Appointment Table
		for i, data := range OneUserAppointment {
			OneUserAppointment[i].SCName = hashdb.Decrypt(data.SCName)
		}
		fmt.Println("Check17---------------------------")
		//Intake Form Table
		IntakeFormDataerr := db.Raw(query.GetIntakeFormSQL, reqVal.AppointmentId).Scan(&IntakeFormData).Error
		if IntakeFormDataerr != nil {
			log.Fatal(IntakeFormDataerr)
		}
		fmt.Println("Check18---------------------------")
		//Decrypt Intake Form Table
		for i, data := range IntakeFormData {
			IntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}
		fmt.Println("Check19---------------------------")
		//Technician Intake Form Table
		var TechnicianIntakeFormData []model.GetTechnicianIntakeData
		TechnicianIntakeFormDataerr := db.Raw(query.GetTechnicianIntakeFormSQL, reqVal.AppointmentId).Scan(&TechnicianIntakeFormData).Error
		if TechnicianIntakeFormDataerr != nil {
			log.Fatal(TechnicianIntakeFormDataerr)
		}
		fmt.Println("Check20---------------------------")
		//Decrypt the Techncian Form Table
		for i, data := range TechnicianIntakeFormData {
			TechnicianIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}
		fmt.Println("Check21---------------------------")
		//Report Intake Form Table
		var ReportIntakeFormData []model.GetReportIntakeData
		ReportIntakeFormDataerr := db.Raw(query.GetReportIntakeFormSQL, reqVal.AppointmentId).Scan(&ReportIntakeFormData).Error
		if ReportIntakeFormDataerr != nil {
			log.Fatal(ReportIntakeFormDataerr)
		}
		fmt.Println("Check22---------------------------")
		//Decrypt Report Intake Form Table
		for i, data := range ReportIntakeFormData {
			ReportIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}
		fmt.Println("Check23---------------------------")
		//Report Text Content Table
		var ReportTextContentData []model.GetReportTextContent
		ReportTextContentDataerr := db.Raw(query.GetReporttextContent, reqVal.AppointmentId).Scan(&ReportTextContentData).Error
		if ReportTextContentDataerr != nil {
			log.Fatal(ReportTextContentDataerr)
		}
		fmt.Println("Check24---------------------------")
		//Decrypt Report Text Content Table
		for i, data := range ReportTextContentData {
			ReportTextContentData[i].TextContent = hashdb.Decrypt(data.TextContent)
		}
		fmt.Println("Check25---------------------------")
		//Report History Table
		var ReportHistoryData []model.GetReportHistory
		ReportHistoryDataerr := db.Raw(query.GetReportHistorySQL, reqVal.AppointmentId).Scan(&ReportHistoryData).Error
		if ReportHistoryDataerr != nil {
			log.Fatal(ReportHistoryDataerr)
		}
		fmt.Println("Check26---------------------------")
		//Decrypt Report History Table
		for i, data := range ReportHistoryData {
			ReportHistoryData[i].HandleUserName = hashdb.Decrypt(data.HandleUserName)
			ReportHistoryData[i].HandleContentText = hashdb.Decrypt(data.HandleContentText)
		}
		fmt.Println("Check27---------------------------")
		// Report Comment Table
		var ReportCommentsData []model.GetReportComments
		ReportCommentsDataerr := db.Raw(query.GetReportCommentsSQL, reqVal.AppointmentId).Scan(&ReportCommentsData).Error
		if ReportCommentsDataerr != nil {
			log.Fatal(ReportCommentsDataerr)
		}
		fmt.Println("Check28---------------------------")
		// Decrypt Report Comment Table
		for i, data := range ReportCommentsData {
			ReportCommentsData[i].Status = hashdb.Decrypt(data.Status)
			ReportCommentsData[i].Comments = hashdb.Decrypt(data.Comments)
		}
		fmt.Println("Check29---------------------------")
		//Get the Template all listed
		var ReportFormateList []model.ReportFormateModel
		ReportFormateListErr := db.Raw(query.GetReportFormateListSQL).Scan(&ReportFormateList).Error
		if ReportFormateListErr != nil {
			log.Fatal(ReportFormateListErr)
		}
		fmt.Println("Check30--------------------------")
		// Decrypt Report Formate List
		for i, data := range ReportFormateList {
			ReportFormateList[i].RFName = hashdb.Decrypt(data.RFName)
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
				[]model.PatientCustId{}
		}

		return true, "Successfully Fetched", IntakeFormData, TechnicianIntakeFormData, ReportIntakeFormData, ReportTextContentData, ReportHistoryData, ReportCommentsData, OneUserAppointment, ReportFormateList, UserDetails, PatientUserDetails

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
				[]model.PatientCustId{}
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
			[]model.PatientCustId{}
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

	//Updating the End Time For the Report History
	ReportHistoryErr := tx.Exec(
		query.CompleteReportHistorySQL,
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

	return true, "Successfully Changes Saved"
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

	//Updating the TechnicianIntake Form
	for _, data := range reqVal.TechnicianIntakeForm {
		status, message := AnswerTechnicianIntakeService(tx, model.AnswerReportIntakeReq{
			PatientId:     reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			QuestionId:    data.QuestionId,
			Answer:        data.Answer,
		}, idValue)

		if !status {
			return status, message
		}
	}

	//Updating the PatientIntake Form
	for _, data := range reqVal.PatientIntakeForm {
		status, message := AnswerPatientIntakeService(tx, model.AnswerReportIntakeReq{
			PatientId:     reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			QuestionId:    data.QuestionId,
			Answer:        data.Answer,
		}, idValue)

		if !status {
			return status, message
		}
	}

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

	//Updating the Appointment Status
	UpdateAppointementErr := tx.Exec(
		query.CompleteReportAppointmentSQL,
		reqVal.MovedStatus,
		reqVal.Impression,
		reqVal.Recommendation,
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

	//Updating the End Time For the Report History
	ReportHistoryErr := tx.Exec(
		query.CompleteReportHistorySQL,
		reqVal.MovedStatus,
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
