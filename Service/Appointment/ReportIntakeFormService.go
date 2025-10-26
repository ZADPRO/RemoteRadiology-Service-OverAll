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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

func CheckAccessService(db *gorm.DB, reqVal model.CheckAccessReq, idValue int, roleIdValue int) (bool, string, int, string) {
	log := logger.InitLogger()

	var result []model.AccessStatus

	if roleIdValue == 7 {
		err := db.Raw(query.ScribeCheckAccessSQL, idValue, reqVal.AppointmentId).Scan(&result).Error

		if err != nil {
			log.Error(err)
		}

	} else {
		err := db.Raw(query.CheckAccessSQL, idValue, reqVal.AppointmentId).Scan(&result).Error

		if err != nil {
			log.Error(err)
		}
	}

	var message = "Another User Already Accessing it"

	if result[0].Status {
		message = "Report are Available for Edit"
	}

	return result[0].Status, message, result[0].RefAppointmentAccessId, result[0].CustID
}

func AssignGetReportService(db *gorm.DB, reqVal model.AssignGetReportReq, idValue int, roleIdValue int) (bool, string, []model.GetViewIntakeData, []model.GetTechnicianIntakeData, []model.GetReportIntakeData, []model.GetReportTextContent, []model.GetReportHistory, []model.GetReportComments, []model.GetOneUserAppointmentModel, []model.ReportFormateModel, []model.GetUserDetails, []model.PatientCustId, bool, *model.FileData, string, []model.AddAddendumModel, []model.GetOldReport, bool, string, string, string, []model.ListAllSignatureModel) {
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
			[]model.AddAddendumModel{},
			[]model.GetOldReport{},
			false,
			"", "", "",
			[]model.ListAllSignatureModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var PerformingProviderName = ""
	var VerifyingProviderName = ""

	//Performing Provider Name
	var PerformingProviderStatus = []string{"Reviewed 1", "Signed Off"}
	var PerformingProviderData []model.PatientCustId
	PerformingProviderErr := db.Raw(query.UserIdentifyRole, 10, reqVal.AppointmentId, PerformingProviderStatus).Scan(&PerformingProviderData).Error
	if PerformingProviderErr != nil {
		log.Error(PerformingProviderErr)
	}

	if len(PerformingProviderData) > 0 {
		PerformingProviderName = hashdb.Decrypt(PerformingProviderData[0].UserFirstName)
	}

	//Verifying Provider Name
	var VerifyingProviderStatus = []string{"Reviewed 2"}
	var VerifyingProviderData []model.PatientCustId
	VerifyingProviderDataErr := db.Raw(query.UserIdentifyRole, 8, reqVal.AppointmentId, VerifyingProviderStatus).Scan(&VerifyingProviderData).Error
	if VerifyingProviderDataErr != nil {
		log.Error(VerifyingProviderDataErr)
	}

	if len(VerifyingProviderData) > 0 {
		VerifyingProviderName = hashdb.Decrypt(VerifyingProviderData[0].UserFirstName)
	}

	//patientCustId
	var PatientUserDetails []model.PatientCustId
	PatientUserDetailsErr := db.Raw(query.PatientUserDetailsSQL, reqVal.PatientId).Scan(&PatientUserDetails).Error
	if PatientUserDetailsErr != nil {
		log.Error(PatientUserDetailsErr)
	}

	for i, data := range PatientUserDetails {
		PatientUserDetails[i].UserFirstName = hashdb.Decrypt(data.UserFirstName)
		PatientUserDetails[i].UserDOB = hashdb.Decrypt(data.UserDOB)
		PatientUserDetails[i].UserGender = hashdb.Decrypt(data.UserGender)
	}

	//GetUserDetails
	var UserDetails []model.GetUserDetails
	UserDetailsErr := db.Raw(query.GetUserDetailsSQL, idValue).Scan(&UserDetails).Error
	if UserDetailsErr != nil {
		log.Error(UserDetailsErr)
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

	if (status && !reqVal.ReadOnlyStatus) || (!status && reqVal.ReadOnlyStatus) || (status && reqVal.ReadOnlyStatus) {

		//Appointment Table
		var Appointment []model.AppointmentModel
		Appointmenterr := db.Raw(query.GetAppointmentSQL, reqVal.AppointmentId).Scan(&Appointment).Error
		if Appointmenterr != nil {
			log.Error(Appointmenterr)
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
						[]model.AddAddendumModel{},
						[]model.GetOldReport{},
						false,
						"", "", "",
						[]model.ListAllSignatureModel{}
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
						[]model.AddAddendumModel{},
						[]model.GetOldReport{},
						false,
						"", "", "",
						[]model.ListAllSignatureModel{}
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
						[]model.AddAddendumModel{},
						[]model.GetOldReport{},
						false,
						"", "", "",
						[]model.ListAllSignatureModel{}

				}

				//List the Latest Report History
				var ReportHistory []model.GetReportHistory
				ListReportHistoryErr := db.Raw(query.CheckLatestReportSQL, reqVal.AppointmentId, reqVal.PatientId).Scan(&ReportHistory).Error
				if ListReportHistoryErr != nil {
					log.Error(ListReportHistoryErr)
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
						[]model.AddAddendumModel{},
						[]model.GetOldReport{},
						false,
						"", "", "",
						[]model.ListAllSignatureModel{}
				}

				if len(ReportHistory) > 0 {

					var starttime = ReportHistory[0].HandleEndTime

					if len(ReportHistory[0].HandleEndTime) == 0 {
						starttime = ReportHistory[0].HandleEndTime
					}

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
							[]model.AddAddendumModel{},
							[]model.GetOldReport{},
							false,
							"", "", "",
							[]model.ListAllSignatureModel{}
					}
				} else {
					var starttime = timeZone.GetPacificTime()

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
							[]model.AddAddendumModel{},
							[]model.GetOldReport{},
							false,
							"", "", "",
							[]model.ListAllSignatureModel{}
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
				[]model.AddAddendumModel{},
				[]model.GetOldReport{},
				false,
				"", "", "",
				[]model.ListAllSignatureModel{}
		}

		var IntakeFormData []model.GetViewIntakeData
		var OneUserAppointment []model.GetOneUserAppointmentModel
		//Appointment Table
		ViewAppointmentErr := db.Raw(query.GetOneUserAppointment, reqVal.PatientId, reqVal.AppointmentId).Scan(&OneUserAppointment).Error
		if ViewAppointmentErr != nil {
			log.Error(ViewAppointmentErr)
		}

		//Decrypt Appointment Table
		for i, data := range OneUserAppointment {
			OneUserAppointment[i].SCName = hashdb.Decrypt(data.SCName)
		}

		//Intake Form Table
		IntakeFormDataerr := db.Raw(query.GetIntakeFormSQL, reqVal.AppointmentId).Scan(&IntakeFormData).Error
		if IntakeFormDataerr != nil {
			log.Error(IntakeFormDataerr)
		}

		//Decrypt Intake Form Table
		for i, data := range IntakeFormData {
			IntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}

		//Technician Intake Form Table
		var TechnicianIntakeFormData []model.GetTechnicianIntakeData
		TechnicianIntakeFormDataerr := db.Raw(query.GetTechnicianIntakeFormSQL, reqVal.AppointmentId).Scan(&TechnicianIntakeFormData).Error
		if TechnicianIntakeFormDataerr != nil {
			log.Error(TechnicianIntakeFormDataerr)
		}

		//Decrypt the Techncian Form Table
		for i, data := range TechnicianIntakeFormData {
			TechnicianIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}

		//Report Intake Form Table
		var ReportIntakeFormData []model.GetReportIntakeData
		ReportIntakeFormDataerr := db.Raw(query.GetReportIntakeFormSQL, reqVal.AppointmentId).Scan(&ReportIntakeFormData).Error
		if ReportIntakeFormDataerr != nil {
			log.Error(ReportIntakeFormDataerr)
		}

		//Decrypt Report Intake Form Table
		for i, data := range ReportIntakeFormData {
			ReportIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
		}

		//Report Text Content Table
		var ReportTextContentData []model.GetReportTextContent
		ReportTextContentDataerr := db.Raw(query.GetReporttextContent, reqVal.AppointmentId).Scan(&ReportTextContentData).Error
		if ReportTextContentDataerr != nil {
			log.Error(ReportTextContentDataerr)
		}

		//Decrypt Report Text Content Table
		for i, data := range ReportTextContentData {
			ReportTextContentData[i].TextContent = hashdb.Decrypt(data.TextContent)
		}

		//Report History Table
		var ReportHistoryData []model.GetReportHistory
		if roleIdValue == 1 || roleIdValue == 6 || roleIdValue == 7 || roleIdValue == 9 || roleIdValue == 10 {
			ReportHistoryDataerr := db.Raw(query.GetReportHistoryFullSQL, reqVal.AppointmentId).Scan(&ReportHistoryData).Error
			if ReportHistoryDataerr != nil {
				log.Error(ReportHistoryDataerr)
			}
		} else {
			ReportHistoryDataerr := db.Raw(query.GetReportHistorySQL, reqVal.AppointmentId).Scan(&ReportHistoryData).Error
			if ReportHistoryDataerr != nil {
				log.Error(ReportHistoryDataerr)
			}
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
			log.Error(ReportCommentsDataerr)
		}

		// Decrypt Report Comment Table
		for i, data := range ReportCommentsData {
			ReportCommentsData[i].Status = hashdb.Decrypt(data.Status)
			ReportCommentsData[i].Comments = hashdb.Decrypt(data.Comments)
		}

		var ReportFormateList []model.ReportFormateModel

		if roleIdValue == 1 {
			//Get the Template all listed
			ReportFormateListErr := db.Raw(query.GetReportFormateAllListSQL).Scan(&ReportFormateList).Error
			if ReportFormateListErr != nil {
				log.Error(ReportFormateListErr)
			}
		} else {
			//Get the Template all listed
			ReportFormateListErr := db.Raw(query.GetReportFormateListSQL, idValue).Scan(&ReportFormateList).Error
			if ReportFormateListErr != nil {
				log.Error(ReportFormateListErr)
			}
		}

		// Decrypt Report Formate List
		for i, data := range ReportFormateList {
			ReportFormateList[i].RFName = hashdb.Decrypt(data.RFName)
		}

		//Scan Center Profile Img
		// Scan Center Profile Img
		var GetScanCenterImg []model.ScanCenterModel
		GetScanCenterImgErr := db.Raw(query.ScanCenterSQL, Appointment[0].SCId).Scan(&GetScanCenterImg).Error
		if GetScanCenterImgErr != nil {
			log.Error(GetScanCenterImgErr)
		}

		var ScanCenterProfileImg *model.FileData

		if len(GetScanCenterImg) > 0 {
			profilePath := hashdb.Decrypt(GetScanCenterImg[0].ProfileImg)

			// ✅ If it's a remote (S3/public) URL
			if strings.HasPrefix(profilePath, "http://") || strings.HasPrefix(profilePath, "https://") {
				resp, err := http.Get(profilePath)
				if err != nil {
					log.Errorf("Failed to fetch S3 image: %v", err)
					ScanCenterProfileImg = &model.FileData{}
				} else {
					defer resp.Body.Close()
					imgBytes, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Errorf("Failed to read S3 image response: %v", err)
						ScanCenterProfileImg = &model.FileData{}
					} else {
						base64Str := base64.StdEncoding.EncodeToString(imgBytes)
						contentType := resp.Header.Get("Content-Type") // get MIME type from response
						ScanCenterProfileImg = &model.FileData{
							Base64Data:  base64Str,
							ContentType: contentType,
						}
					}
				}
			} else {
				log.Infof("🗂️ Detected Local file for ScanCenter profile: %s", profilePath)

				viewedFile, viewErr := helperView.ViewFile("./Assets/Profile/" + profilePath)
				if viewErr != nil {
					log.Errorf("❌ Failed to read ScanCenter profile image: %v", viewErr)
					ScanCenterProfileImg = &model.FileData{}
				} else {
					ScanCenterProfileImg = &model.FileData{
						Base64Data:  viewedFile.Base64Data,
						ContentType: viewedFile.ContentType,
					}
				}
			}

		} else {
			ScanCenterProfileImg = &model.FileData{}
		}

		var EaseQTReportAccess = false
		var NASystemReportAccess = false

		//Get the Ease QT Report Access Status
		switch roleIdValue {
		case 1: //Master Admin
			EaseQTReportAccess = true
			NASystemReportAccess = true
		case 2: //Scan Center Technician
			EaseQTReportAccess = false
			NASystemReportAccess = true
		case 3: //Scan Center Manager
			EaseQTReportAccess = false
			NASystemReportAccess = true
		case 4: //Patient
			EaseQTReportAccess = false
			NASystemReportAccess = true
		case 5: //Scan Center Doctor

			var ReportStatus []model.DoctorReportAccessStatus
			err := db.Raw(query.DoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
			if err != nil {
				log.Error(err)
			}

			if len(ReportStatus) == 0 || ReportStatus[0].DDEaseQTReportAccess == nil {
				EaseQTReportAccess = false
				break
			}

			EaseQTReportAccess = *ReportStatus[0].DDEaseQTReportAccess
			NASystemReportAccess = *ReportStatus[0].DDNAsystemReportAccess

		case 6: //Junior Doctor
			EaseQTReportAccess = true
			NASystemReportAccess = true
		case 7: //Scribe
			EaseQTReportAccess = true
			NASystemReportAccess = true
		case 8: //Scan Center Reviewer
			var ReportStatus []model.CoDoctorReportAccessStatus
			err := db.Raw(query.CoDoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
			if err != nil {
				log.Error(err)
			}

			if len(ReportStatus) == 0 || ReportStatus[0].CDEaseQTReportAccess == nil {
				EaseQTReportAccess = false
				break
			}

			EaseQTReportAccess = *ReportStatus[0].CDEaseQTReportAccess
			NASystemReportAccess = *ReportStatus[0].CDNAsystemReportAccess

		case 9: //Manager
			EaseQTReportAccess = true
			NASystemReportAccess = true
		case 10: //Performing Provider
			EaseQTReportAccess = true
			NASystemReportAccess = true
		default:
			EaseQTReportAccess = false
			NASystemReportAccess = false
		}
		// if roleIdValue == 1 || roleIdValue == 9 || roleIdValue == 6 || roleIdValue == 7 || roleIdValue == 10 {
		// 	EaseQTReportAccess = true
		// } else if roleIdValue == 5 { //Check the Scan Center Doctor

		// }

		//Get the Other Old Reports
		var oldReportData []model.GetOldReport
		oldReportErr := db.Raw(query.GetOldReportsSQL, reqVal.PatientId, reqVal.AppointmentId).Scan(&oldReportData).Error
		if oldReportErr != nil {
			log.Error(oldReportErr)
		}

		//Get the Private and Public in the Technician Intake
		var patientPrivatePublic []model.GetTechnicianIntakeData
		patientPrivatePublicErr := db.Raw(query.GetPatientPrivatePublicSQL, reqVal.AppointmentId, 57).Scan(&patientPrivatePublic).Error
		if patientPrivatePublicErr != nil {
			log.Error(patientPrivatePublicErr)
		}

		for i, data := range patientPrivatePublic {
			patientPrivatePublic[i].Answer = hashdb.Decrypt(data.Answer)
		}

		var patientPrivatePublicStatus = ""

		if len(patientPrivatePublic) > 0 {
			patientPrivatePublicStatus = patientPrivatePublic[0].Answer
		}

		return true, "Successfully Fetched", IntakeFormData, TechnicianIntakeFormData, ReportIntakeFormData, ReportTextContentData, ReportHistoryData, ReportCommentsData, OneUserAppointment, ReportFormateList, UserDetails, PatientUserDetails, EaseQTReportAccess, ScanCenterProfileImg, hashdb.Decrypt(GetScanCenterImg[0].SCAddress), ListAddendumService(db, reqVal.AppointmentId), oldReportData, NASystemReportAccess, patientPrivatePublicStatus, PerformingProviderName, VerifyingProviderName, ListAllSignatureService(db, reqVal.AppointmentId)

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
				[]model.AddAddendumModel{},
				[]model.GetOldReport{},
				false,
				"", "", "",
				[]model.ListAllSignatureModel{}
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
			[]model.AddAddendumModel{},
			[]model.GetOldReport{},
			false,
			"", "", "",
			[]model.ListAllSignatureModel{}
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
		log.Error(ReportIntakeFormDataerr)
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
		log.Error(TechnicianIntakeFormDataErr)
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
		log.Error(PatientIntakeFormDataErr)
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
		log.Error(ReportTextContentErr)
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

func AutosaveServicee(db *gorm.DB, reqVal model.AutoSubmitReportReq, idValue int, roleIdValue int) (bool, string, []model.GetReportIntakeData, []model.GetReportTextContent, []model.GetOneUserAppointmentModel, bool, bool, string, string, []model.ListAllSignatureModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var PerformingProviderName = ""
	var VerifyingProviderName = ""

	//Performing Provider Name
	var PerformingProviderStatus = []string{"Reviewed 1", "Signed Off"}
	var PerformingProviderData []model.PatientCustId
	PerformingProviderErr := db.Raw(query.UserIdentifyRole, 10, reqVal.AppointmentId, PerformingProviderStatus).Scan(&PerformingProviderData).Error
	if PerformingProviderErr != nil {
		log.Error(PerformingProviderErr)
	}

	if len(PerformingProviderData) > 0 {
		PerformingProviderName = hashdb.Decrypt(PerformingProviderData[0].UserFirstName)
	}

	//Verifying Provider Name
	var VerifyingProviderStatus = []string{"Reviewed 2"}
	var VerifyingProviderData []model.PatientCustId
	VerifyingProviderDataErr := db.Raw(query.UserIdentifyRole, 8, reqVal.AppointmentId, VerifyingProviderStatus).Scan(&VerifyingProviderData).Error
	if VerifyingProviderDataErr != nil {
		log.Error(VerifyingProviderDataErr)
	}

	if len(VerifyingProviderData) > 0 {
		VerifyingProviderName = hashdb.Decrypt(VerifyingProviderData[0].UserFirstName)
	}

	//Inserting and Upadating the Report Intake Form
	for _, data := range reqVal.ReportIntakeForm {
		status, message := AnswerReportIntakeService(tx, model.AnswerReportIntakeReq{
			PatientId:     reqVal.PatientId,
			AppointmentId: reqVal.AppointmentId,
			QuestionId:    data.QuestionId,
			Answer:        data.Answer,
		}, idValue)

		if !status {
			return status, message, []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}

		var updateSyncErr = tx.Exec(
			query.UpdateAutosavePatientHistorySyncSQL,
			false,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
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
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the BreastImplant SyncStatus
	if reqVal.ChangedOneState.BreastImplantSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveBreastImplantSyncSQL,
			reqVal.ReportSyncStatus.BreastImplantSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the Symmetry SyncStatus
	if reqVal.ChangedOneState.SymmetrySyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveSymmetrySyncSQL,
			reqVal.ReportSyncStatus.SymmetrySyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the BreastDensity SyncStatus
	if reqVal.ChangedOneState.BreastDensitySyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveBreastDensitySyncSQL,
			reqVal.ReportSyncStatus.BreastDensitySyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the NippleAreola SyncStatus
	if reqVal.ChangedOneState.NippleAreolaSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveNippleAreolaSyncSQL,
			reqVal.ReportSyncStatus.NippleAreolaSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the Glandular SyncStatus
	if reqVal.ChangedOneState.GlandularSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveGlandularSyncSQL,
			reqVal.ReportSyncStatus.GlandularSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the LymphNodes SyncStatus
	if reqVal.ChangedOneState.LymphNodesSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveLymphNodesSyncSQL,
			reqVal.ReportSyncStatus.LymphNodesSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the Lesions SyncStatus
	if reqVal.ChangedOneState.LesionsSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveLesionsSyncSQL,
			reqVal.ReportSyncStatus.LesionsSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the ComparisonPrior SyncStatus
	if reqVal.ChangedOneState.ComparisonPriorSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveComparisonPriorSyncSQL,
			reqVal.ReportSyncStatus.ComparisonPriorSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the BreastDensityLeft SyncStatus
	if reqVal.ChangedOneState.BreastDensityLeftSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveBreastDensityLeftSyncSQL,
			reqVal.ReportSyncStatus.BreastDensityLeftSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the NippleAreolaLeft SyncStatus
	if reqVal.ChangedOneState.NippleAreolaLeftSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveNippleAreolaLeftSyncSQL,
			reqVal.ReportSyncStatus.NippleAreolaLeftSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the GlandularLeft SyncStatus
	if reqVal.ChangedOneState.GlandularLeftSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveGlandularLeftSyncSQL,
			reqVal.ReportSyncStatus.GlandularLeftSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the LymphNodesLeft SyncStatus
	if reqVal.ChangedOneState.LymphNodesLeftSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveLymphNodesLeftSyncSQL,
			reqVal.ReportSyncStatus.LymphNodesLeftSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the LesionsLeft SyncStatus
	if reqVal.ChangedOneState.LesionsLeftSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveLesionsLeftSyncSQL,
			reqVal.ReportSyncStatus.LesionsLeftSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the ComparisonPriorLeft SyncStatus
	if reqVal.ChangedOneState.ComparisonPriorLeftSyncStatus {
		var updateSyncErr = tx.Exec(
			query.UpdateAutosaveComparisonPriorLeftSyncSQL,
			reqVal.ReportSyncStatus.ComparisonPriorLeftSyncStatus,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if updateSyncErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", updateSyncErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the BreastImplantReportText
	if reqVal.ChangedOneState.BreastImplantReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveBreastImplantReportTextSyncSQL,
			reqVal.AutoReportText.BreastImplantReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the SymmetryReportText
	if reqVal.ChangedOneState.SymmetryReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveSymmetryReportTextSyncSQL,
			reqVal.AutoReportText.SymmetryReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the BreastDensityReportText
	if reqVal.ChangedOneState.BreastDensityReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveBreastDensityReportTextSyncSQL,
			reqVal.AutoReportText.BreastDensityReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the NippleAreolaReportText
	if reqVal.ChangedOneState.NippleAreolaReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveNippleAreolaReportTextSyncSQL,
			reqVal.AutoReportText.NippleAreolaReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the LesionsReportText
	if reqVal.ChangedOneState.LesionsReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveLesionsReportTextTextSyncSQL,
			reqVal.AutoReportText.LesionsReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the ComparisonPriorReportText
	if reqVal.ChangedOneState.ComparisonPriorReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveComparisonPriorReportTextSyncSQL,
			reqVal.AutoReportText.ComparisonPriorReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the GrandularAndDuctalTissueReportText
	if reqVal.ChangedOneState.GrandularAndDuctalTissueReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveGrandularAndDuctalTissueReportTextSyncSQL,
			reqVal.AutoReportText.GrandularAndDuctalTissueReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the LymphNodesReportText
	if reqVal.ChangedOneState.LymphNodesReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveLymphNodesReportTextSyncSQL,
			reqVal.AutoReportText.LymphNodesReportText,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the BreastDensityReportTextLeft
	if reqVal.ChangedOneState.BreastDensityLeftReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveBreastDensityReportTextLeftSyncSQL,
			reqVal.AutoReportText.BreastDensityReportTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the NippleAreolaReportTextLeft
	if reqVal.ChangedOneState.NippleAreolaLeftReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveNippleAreolaReportTextLeftSyncSQL,
			reqVal.AutoReportText.NippleAreolaReportTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the LesionsReportTextLeft
	if reqVal.ChangedOneState.LesionsLeftReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveLesionsReportTextLeftSyncSQL,
			reqVal.AutoReportText.LesionsReportTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the ComparisonPriorReportTextLeft
	if reqVal.ChangedOneState.ComparisonPriorLeftReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveComparisonPriorReportTextLeftSyncSQL,
			reqVal.AutoReportText.ComparisonPriorReportTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the GrandularAndDuctalTissueReportTextLeft
	if reqVal.ChangedOneState.GrandularAndDuctalTissueLeftReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveGrandularAndDuctalTissueReportTextLeftSyncSQL,
			reqVal.AutoReportText.GrandularAndDuctalTissueReportTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	//Updte the LymphNodesReportTextLeft
	if reqVal.ChangedOneState.LymphNodesLeftReportText {
		var reportTextErr = tx.Exec(
			query.UpdateAutosaveLymphNodesReportTextLeftSyncSQL,
			reqVal.AutoReportText.LymphNodesReportTextLeft,
			reqVal.AppointmentId,
			reqVal.PatientId,
		).Error

		if reportTextErr != nil {
			log.Printf("ERROR: Failed to Update Autosave %v\n", reportTextErr)
			tx.Rollback()
			return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again", []model.GetReportIntakeData{}, []model.GetReportTextContent{}, []model.GetOneUserAppointmentModel{}, false, false, "", "", []model.ListAllSignatureModel{}
	}

	//Report Intake Form Table
	var ReportIntakeFormData []model.GetReportIntakeData
	ReportIntakeFormDataerr := db.Raw(query.GetReportIntakeFormSQL, reqVal.AppointmentId).Scan(&ReportIntakeFormData).Error
	if ReportIntakeFormDataerr != nil {
		log.Error(ReportIntakeFormDataerr)
	}

	//Decrypt Report Intake Form Table
	for i, data := range ReportIntakeFormData {
		ReportIntakeFormData[i].Answer = hashdb.Decrypt(data.Answer)
	}

	//Report Text Content Table
	var ReportTextContentData []model.GetReportTextContent
	ReportTextContentDataerr := db.Raw(query.GetReporttextContent, reqVal.AppointmentId).Scan(&ReportTextContentData).Error
	if ReportTextContentDataerr != nil {
		log.Error(ReportTextContentDataerr)
	}

	//Decrypt Report Text Content Table
	for i, data := range ReportTextContentData {
		ReportTextContentData[i].TextContent = hashdb.Decrypt(data.TextContent)
	}

	var OneUserAppointment []model.GetOneUserAppointmentModel
	//Appointment Table
	ViewAppointmentErr := db.Raw(query.GetOneUserAppointment, reqVal.PatientId, reqVal.AppointmentId).Scan(&OneUserAppointment).Error
	if ViewAppointmentErr != nil {
		log.Error(ViewAppointmentErr)
	}

	//Decrypt Appointment Table
	for i, data := range OneUserAppointment {
		OneUserAppointment[i].SCName = hashdb.Decrypt(data.SCName)
	}

	var EaseQTReportAccess = false
	var NASystemReportAccess = false

	//Get the Ease QT Report Access Status
	switch roleIdValue {
	case 1: //Master Admin
		EaseQTReportAccess = true
		NASystemReportAccess = true
	case 2: //Scan Center Technician
		EaseQTReportAccess = false
		NASystemReportAccess = true
	case 3: //Scan Center Manager
		EaseQTReportAccess = false
		NASystemReportAccess = true
	case 4: //Patient
		EaseQTReportAccess = false
		NASystemReportAccess = true
	case 5: //Scan Center Doctor

		var ReportStatus []model.DoctorReportAccessStatus
		err := db.Raw(query.DoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
		if err != nil {
			log.Error(err)
		}

		if len(ReportStatus) == 0 || ReportStatus[0].DDEaseQTReportAccess == nil {
			EaseQTReportAccess = false
			break
		}

		EaseQTReportAccess = *ReportStatus[0].DDEaseQTReportAccess
		NASystemReportAccess = *ReportStatus[0].DDNAsystemReportAccess

	case 6: //Junior Doctor
		EaseQTReportAccess = true
		NASystemReportAccess = true
	case 7: //Scribe
		EaseQTReportAccess = true
		NASystemReportAccess = true
	case 8: //Scan Center Reviewer
		var ReportStatus []model.CoDoctorReportAccessStatus
		err := db.Raw(query.CoDoctorReportAccessSQL, idValue).Scan(&ReportStatus).Error
		if err != nil {
			log.Error(err)
		}

		if len(ReportStatus) == 0 || ReportStatus[0].CDEaseQTReportAccess == nil {
			EaseQTReportAccess = false
			break
		}

		EaseQTReportAccess = *ReportStatus[0].CDEaseQTReportAccess
		NASystemReportAccess = *ReportStatus[0].CDNAsystemReportAccess

	case 9: //Manager
		EaseQTReportAccess = true
		NASystemReportAccess = true
	case 10: //Performing Provider
		EaseQTReportAccess = true
		NASystemReportAccess = true
	default:
		EaseQTReportAccess = false
		NASystemReportAccess = false
	}

	return true, "Successfully Changes Saved", ReportIntakeFormData, ReportTextContentData, OneUserAppointment, EaseQTReportAccess, NASystemReportAccess, PerformingProviderName, VerifyingProviderName, ListAllSignatureService(db, reqVal.AppointmentId)
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

	var reportStatus = reqVal.MovedStatus

	if reqVal.LeaveStatus {
		reportStatus = "Changes"
	}

	var MovedStatus = reqVal.MovedStatus

	var TextContent = reqVal.ReportTextContent

	if reqVal.MovedStatus == "Signed Off" {
		//Get Addedum
		var UserCount []model.AddedumCountModel

		UserCountErr := tx.Raw(query.AddedumCountSQL, reqVal.AppointmentId).Scan(&UserCount).Error
		if UserCountErr != nil {
			log.Error(UserCountErr)
		}

		if UserCount[0].Count > 0 {
			MovedStatus = "Signed Off (A)"
		}

	}

	if roleIdValue == 7 {
		//Updating the Appointment Status
		UpdateAppointementErr := tx.Exec(
			query.ScribeCompleteReportAppointmentSQL,
			MovedStatus,
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
			MovedStatus,
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
		hashdb.Encrypt(TextContent),
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

		ListUserDataErr := db.Raw(query.ListUserDataSQL, reqVal.PatientId, reqVal.AppointmentId, []int{6}).Scan(&ListUserData).Error
		if ListUserDataErr != nil {
			log.Error(ListUserDataErr.Error())
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

		ListUserDataErr := db.Raw(query.ListUserDataSQL, reqVal.PatientId, reqVal.AppointmentId, []int{1, 10}).Scan(&ListUserData).Error
		if ListUserDataErr != nil {
			log.Error(ListUserDataErr.Error())
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
			//Handler User
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

			//Update the Handler User
			UpdateChangesUserErr := tx.Exec(
				query.UpdateHandlerCorrectEditIdSQL,
				correct,
				edit,
				idValue,
				reqVal.AppointmentId,
			).Error
			if UpdateChangesUserErr != nil {
				log.Printf("ERROR: Failed to Update Report History: %v\n", UpdateChangesUserErr)
				tx.Rollback()
				return false, "Something went wrong, Try Again"
			}
			// }

		}

	}

	var AddedumContent []string

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

		AddedumContent = append(AddedumContent, "The written report was emailed to "+PatientdataModel[0].Email+" on "+timeZone.GetPacificTime())

		htmlContent := mailservice.PatientReportSignOff(hashdb.Decrypt(PatientdataModel[0].UserFirstName), PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, PatientdataModel[0].SCCustId)

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
			if data.UserStatus {
				AddedumContent = append(AddedumContent, "The written report was emailed to "+data.Email+" on "+timeZone.GetPacificTime())

				htmlContent := mailservice.ManagerReportSignOff(hashdb.Decrypt(PatientdataModel[0].UserFirstName), PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, data.SCCustId)

				subject := "Report Status"

				emailStatus := mailservice.MailService(data.Email, htmlContent, subject)

				if !emailStatus {
					log.Error("Sending Mail Meets Error")
					return false, "Something went wrong, Try Again"
				}
			}
		}

	}

	if reqVal.PatientMailStatus || reqVal.ManagerMailStatus {

		//Get User Cust Id
		var UserId []model.AddAddendumModel
		UserIdErr := db.Raw(query.GetUserId, idValue).Scan(&UserId).Error
		if UserIdErr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", UserIdErr)
		}

		var AddedumList = ListAddendumService(db, reqVal.AppointmentId)

		var userCustId = ""
		if len(UserId) > 0 {
			userCustId = UserId[0].CustId
		}

		if len(AddedumList) > 0 {
			TextContent += "<br/>" + timeZone.GetPacificTime() + " - " + userCustId +
				"<p>" + strings.Join(AddedumContent, ". ") + "</p>"
		} else {
			TextContent += "<br/><p><strong>ADDENDUM:</strong></p><br/><p>" +
				timeZone.GetPacificTime() + " - " + userCustId +
				"</p><p>" + strings.Join(AddedumContent, ". ") + "</p>"
		}

		InsertErr := tx.Exec(
			query.InsertAddedumSQL,
			reqVal.AppointmentId,
			idValue,
			strings.Join(AddedumContent, ". "),
			timeZone.GetPacificTime(),
		).Error
		if InsertErr != nil {
			log.Printf("ERROR: Failed to Insert Addendum: %v\n", InsertErr)
			tx.Rollback()
		}

		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$----------------------->")

	}

	//Updating the Report Text Content
	status, message := AnswerTextContentService(tx, model.AnswerTextContentReq{
		PatientId:     reqVal.PatientId,
		AppointmentId: reqVal.AppointmentId,
		TextContent:   TextContent,
		SyncStatus:    reqVal.SyncStatus,
	}, idValue)

	if !status {
		return status, message
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

func UploadReportFormateService(db *gorm.DB, reqVal model.UploadReportFormateReq, idValue int) (int, string, bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return 0, "", false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var insertedID []model.ReportFormateCreateModel

	//Adding Template
	InsertReportTemplateErr := tx.Raw(
		query.InsertReportTemplate,
		hashdb.Encrypt(reqVal.Name),
		hashdb.Encrypt(reqVal.FormateTemplate),
		timeZone.GetPacificTime(),
		idValue,
		reqVal.AccessStatus,
	).Scan(&insertedID).Error
	if InsertReportTemplateErr != nil {
		log.Printf("ERROR: Failed to Insert Report Template: %v\n", InsertReportTemplateErr)
		tx.Rollback()
		return 0, "", false, "Something went wrong, Try Again"
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
		return 0, "", false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return 0, "", false, "Something went wrong, Try Again"
	}

	return insertedID[0].RFId, insertedID[0].RefUserCustId, true, "Successfully Changes Saved"
}

func DeleteReportFormateService(db *gorm.DB, reqVal model.DeleteReportFormateReq, idValue int) (bool, string) {
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

	//Delete Template
	DeleteReportTemplateErr := tx.Exec(
		query.DeleteReportTemplateSQL,
		reqVal.Id,
	).Error
	if DeleteReportTemplateErr != nil {
		log.Printf("ERROR: Failed to Delete Report Template: %v\n", DeleteReportTemplateErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Addding audit For the Template
	transData := 36
	errTrans := model.RefTransHistory{
		TransTypeId: transData,
		THData:      "Report Template Deleted Successfully",
		UserId:      idValue,
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

func UpdateReportFormateService(db *gorm.DB, reqVal model.UpdateReportFormateReq, idValue int) (bool, string) {
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

	//Update Template
	UpdateReportTemplateErr := tx.Exec(
		query.UpdateReportTemplateSQL,
		reqVal.AccessStatus,
		reqVal.Id,
	).Error
	if UpdateReportTemplateErr != nil {
		log.Printf("ERROR: Failed to Update Report Template: %v\n", UpdateReportTemplateErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	//Addding audit For the Template
	transData := 36
	errTrans := model.RefTransHistory{
		TransTypeId: transData,
		THData:      "Report Template Update Successfully",
		UserId:      idValue,
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

func GetReportFormateService(db *gorm.DB, reqVal model.GetReportFormateReq, idValue int) []model.ReportTextFormateModel {
	log := logger.InitLogger()

	//Getting the Template
	var TemplateFormate []model.ReportTextFormateModel
	err := db.Raw(query.GetOneReportFormateListSQL, reqVal.Id).Scan(&TemplateFormate).Error
	if err != nil {
		log.Error(err)
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

func SendMailReportService(db *gorm.DB, reqVal model.SendMailReportReq, idValue int) (bool, string) {
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

	var AddedumContent []string

	if reqVal.PatientMailStatus {

		AddedumContent = append(AddedumContent, "The written report was emailed to "+PatientdataModel[0].Email+" on "+timeZone.GetPacificTime())

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
			if data.UserStatus {
				AddedumContent = append(AddedumContent, "The written report was emailed to "+data.Email+" on "+timeZone.GetPacificTime())
				htmlContent := mailservice.ManagerReportSignOff(hashdb.Decrypt(PatientdataModel[0].UserFirstName), PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, data.SCCustId)

				subject := "Report Status"

				emailStatus := mailservice.MailService(data.Email, htmlContent, subject)

				if !emailStatus {
					log.Error("Sending Mail Meets Error")
					return false, "Something went wrong, Try Again"
				}
			}
		}

	}

	if reqVal.PatientMailStatus || reqVal.ManagerMailStatus {

		//Get PatientContent
		var PatientContent []model.GetReportTextContent
		PatientContenterr := db.Raw(query.GetPatientContent, reqVal.AppointmentId).Scan(&PatientContent).Error
		if PatientContenterr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", PatientContenterr)
		}

		//Get User Cust Id
		var UserId []model.AddAddendumModel
		UserIdErr := db.Raw(query.GetUserId, idValue).Scan(&UserId).Error
		if UserIdErr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", UserIdErr)
		}

		var AddedumList = ListAddendumService(db, reqVal.AppointmentId)
		var TextContent = ""

		var userCustId = ""
		if len(UserId) > 0 {
			userCustId = UserId[0].CustId
		}

		if len(AddedumList) > 0 {
			TextContent += "<br/>" + timeZone.GetPacificTime() + " - " + userCustId +
				"<p>" + strings.Join(AddedumContent, ". ") + "</p>"
		} else {
			TextContent += "<br/><p><strong>ADDENDUM:</strong></p><br/><p>" +
				timeZone.GetPacificTime() + " - " + userCustId +
				"</p><p>" + strings.Join(AddedumContent, ". ") + "</p>"
		}

		InsertErr := tx.Exec(
			query.InsertAddedumSQL,
			reqVal.AppointmentId,
			idValue,
			strings.Join(AddedumContent, ". "),
			timeZone.GetPacificTime(),
		).Error
		if InsertErr != nil {
			log.Printf("ERROR: Failed to Insert Addendum: %v\n", InsertErr)
			tx.Rollback()
		}

		//Update textContent
		UpdateTextContentErr := tx.Exec(
			query.UpdateReportTextContentSQL,
			hashdb.Encrypt(hashdb.Decrypt(PatientContent[0].TextContent)+TextContent),
			reqVal.AppointmentId,
		).Error
		if UpdateTextContentErr != nil {
			log.Printf("ERROR: Failed to Update Report Text Content: %v\n", UpdateTextContentErr)
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
		log.Error(err)
	}

	if len(PatientFile.Answer) > 0 {
		DriversLicenseNoImgHelperData, viewErr := helperView.ViewFile("./Assets/Files/" + hashdb.Decrypt(PatientFile.Answer))
		if viewErr != nil {
			// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
			log.Errorf("Failed to read DrivingLicense file: %v", viewErr)
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

	//Get PatientContent
	var PatientContent []model.GetReportTextContent
	PatientContenterr := db.Raw(query.GetPatientContent, reqVal.AppointmentId).Scan(&PatientContent).Error
	if PatientContenterr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", PatientContenterr)
	}

	UpdateAppointmentErr := tx.Exec(
		query.UpdateReportAppointmentSQL,
		"Signed Off (A)",
		reqVal.AppointmentId,
	).Error

	if UpdateAppointmentErr != nil {
		log.Printf("ERROR: Failed to Update Appointment: %v\n", UpdateAppointmentErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
	}

	var AddedumContent []string

	//Send Mail for the Patient
	if reqVal.PatientMailStatus {

		var PatientdataModel []model.Patientdata

		err := db.Raw(query.GetPatientData, reqVal.PatientId, reqVal.AppointmentId).Scan(&PatientdataModel).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
		}

		for i, data := range PatientdataModel {
			PatientdataModel[i].UserFirstName = hashdb.Decrypt(data.UserFirstName)
		}

		AddedumContent = append(AddedumContent, "The written report was emailed to "+PatientdataModel[0].Email+" on "+timeZone.GetPacificTime())

		htmlContent := mailservice.PatientReportSignOff(hashdb.Decrypt(PatientdataModel[0].UserFirstName), PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, PatientdataModel[0].SCCustId)

		subject := "Your Report Status"

		emailStatus := mailservice.MailService(PatientdataModel[0].Email, htmlContent, subject)

		if !emailStatus {
			log.Error("Sending Mail Meets Error")
			return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
		}
	}

	//Send Mail for the Scan Center Manager
	if reqVal.ManagerMailStatus {
		var PatientdataModel []model.Patientdata

		err := db.Raw(query.GetPatientData, reqVal.PatientId, reqVal.AppointmentId).Scan(&PatientdataModel).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
		}

		var ManagerModel []model.ManagerData

		Managererr := db.Raw(query.GetManagerData, 3, reqVal.AppointmentId).Scan(&ManagerModel).Error
		if Managererr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", Managererr)
			return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
		}

		for _, data := range ManagerModel {
			if data.UserStatus {
				AddedumContent = append(AddedumContent, "The written report was emailed to "+data.Email+" on "+timeZone.GetPacificTime())

				htmlContent := mailservice.ManagerReportSignOff(hashdb.Decrypt(PatientdataModel[0].UserFirstName), PatientdataModel[0].CustId, PatientdataModel[0].AppointmentDate, data.SCCustId)

				subject := "Report Status"

				emailStatus := mailservice.MailService(data.Email, htmlContent, subject)

				if !emailStatus {
					log.Error("Sending Mail Meets Error")
					return false, "Something went wrong, Try Again", []model.AddAddendumModel{}
				}
			}
		}
	}

	//Get User Cust Id
	var UserId []model.AddAddendumModel
	UserIdErr := db.Raw(query.GetUserId, idValue).Scan(&UserId).Error
	if UserIdErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", UserIdErr)
	}

	var AddedumList = ListAddendumService(db, reqVal.AppointmentId)
	var TextContent = ""
	if len(PatientContent) > 0 {
		TextContent = hashdb.Decrypt(PatientContent[0].TextContent)
	}

	var userCustId = ""
	if len(UserId) > 0 {
		userCustId = UserId[0].CustId
	}

	if len(AddedumList) > 0 {
		TextContent += "<br/>" + timeZone.GetPacificTime() + " - " + userCustId + "" + reqVal.AddAddendumText + "<p>" + strings.Join(AddedumContent, ". ") + "</p>"
	} else {
		TextContent += "<br/><p><strong>ADDENDUM:</strong></p><br/><p>" + timeZone.GetPacificTime() + " - " + userCustId + "</p>" + reqVal.AddAddendumText + "<p>" + strings.Join(AddedumContent, ". ") + "</p>"
	}

	//Update textContent
	UpdateTextContentErr := tx.Exec(
		query.UpdateReportTextContentSQL,
		hashdb.Encrypt(TextContent),
		reqVal.AppointmentId,
	).Error
	if UpdateTextContentErr != nil {
		log.Printf("ERROR: Failed to Update Report Text Content: %v\n", UpdateTextContentErr)
	}

	InsertErr := tx.Exec(
		query.InsertAddedumSQL,
		reqVal.AppointmentId,
		idValue,
		reqVal.AddAddendumText+"<p>"+strings.Join(AddedumContent, ". ")+"</p>",
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

	return true, "Successfully Addedum Received", AddedumList

}

func ListOldReportService(db *gorm.DB, reqVal model.ListOldReportReq, idValue int) (bool, string, []model.ListOldReportModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again", []model.ListOldReportModel{}
	}

	var ListOldReportModel []model.ListOldReportModel

	ListAddendumErr := tx.Raw(query.ListOldReportSQL, reqVal.AppointmentId, reqVal.PatientId, reqVal.CategoryId).Scan(&ListOldReportModel).Error
	if ListAddendumErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", ListAddendumErr)
		return false, "Something went wrong, Try Again", []model.ListOldReportModel{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again", []model.ListOldReportModel{}
	}

	return true, "Successfully Changes Saved", ListOldReportModel

}

func DeleteOldReportService(db *gorm.DB, reqVal model.DeleteOldReportModel) (bool, string) {
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

	var DicomFiles []model.ListOldReportModel
	DicomErr := tx.Raw(query.GetParticularOldReport, reqVal.ORId).Scan(&DicomFiles).Error
	if DicomErr != nil {
		log.Printf("ERROR: Failed to fetch old report records: %v", DicomErr)
		return false, "Something went wrong, Try Again"
	}

	DeleteDicomErr := tx.Exec(query.DeleteOldReportSQL, false, reqVal.ORId).Error
	if DeleteDicomErr != nil {
		log.Error(DeleteDicomErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	uploadPath := "./Assets/Files/"

	for _, data := range DicomFiles {
		filePath := data.ORFilename

		// If it's an S3 or external URL, skip deletion
		if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
			log.Infof("Skipping deletion for S3/external file: %s", filePath)
			continue
		}

		// Otherwise, delete from local storage
		fullPath := filepath.Join(uploadPath, filePath)
		if err := os.Remove(fullPath); err != nil {
			log.Error("Local file deletion failed:", err)
			// Don’t fail entire operation for one file
			continue
		}

		log.Infof("Deleted local file: %s", fullPath)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Deleted"
}

func ListAllSignatureService(db *gorm.DB, appointmentId int) []model.ListAllSignatureModel {
	log := logger.InitLogger()

	var ListAllSignatureModel []model.ListAllSignatureModel
	ListAllSignatureErr := db.Raw(query.ListAllSignatureSQL, appointmentId).Scan(&ListAllSignatureModel).Error
	if ListAllSignatureErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", ListAllSignatureErr)
		return []model.ListAllSignatureModel{}
	}

	return ListAllSignatureModel

}

func InsertSignatureService(db *gorm.DB, reqVal model.AddSignatureReq) (bool, string, []model.ListAllSignatureModel) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again", []model.ListAllSignatureModel{}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var InsertSignatureServiceErr = tx.Exec(query.InsertSignatureSQL, reqVal.AppointmentId, reqVal.PatientId, reqVal.AddSignatureText, timeZone.GetPacificTime()).Error
	if InsertSignatureServiceErr != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", InsertSignatureServiceErr)
		return false, "Something went wrong, Try Again", []model.ListAllSignatureModel{}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again", []model.ListAllSignatureModel{}
	}

	return true, "Successfully Deleted", ListAllSignatureService(db, reqVal.AppointmentId)

}
