package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"fmt"

	"gorm.io/gorm"
)

func GetAllPatientService(db *gorm.DB, reqVal model.GetAllPatientReq) []model.GetAllRadiologist {
	log := logger.InitLogger()

	fmt.Println(reqVal.SCId)
	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetAllPatientList, reqVal.SCId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetAllRadiologist{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
	}

	return RadiologistData
}

func GetPatientDataService(db *gorm.DB, reqVal model.GetPatientReq, idValue int) []model.PatientOneModel {
	log := logger.InitLogger()

	var RadiologistData []model.PatientOneModel

	err := db.Raw(query.GetOnePatientList, reqVal.UserId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.PatientOneModel{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)

		if len(hashdb.Decrypt(tech.ProfileImg)) > 0 {
			profileImgHelperData, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(tech.ProfileImg))
			if viewErr != nil {
				// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
				log.Fatalf("Failed to read profile image file: %v", viewErr)
			}
			if profileImgHelperData != nil {
				RadiologistData[i].ProfileImgFile = &model.FileData{
					Base64Data:  profileImgHelperData.Base64Data,
					ContentType: profileImgHelperData.ContentType,
				}
			}
		} else {
			RadiologistData[i].ProfileImgFile = nil
		}

		var MappingData []model.PatientAppointmentModel

		Appointmenterr := db.Raw(query.GetAppointmentListSQL, reqVal.UserId).Scan(&MappingData).Error
		if Appointmenterr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", Appointmenterr)
			return []model.PatientOneModel{}
		}

		RadiologistData[i].Appointments = MappingData
	}

	return RadiologistData
}
