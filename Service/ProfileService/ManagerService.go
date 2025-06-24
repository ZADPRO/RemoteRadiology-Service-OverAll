package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"

	"gorm.io/gorm"
)

func GetAllManagerDataService(db *gorm.DB) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofManagerSQL, 9).Scan(&RadiologistData).Error
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

func GetManagerDataService(db *gorm.DB, reqVal model.GetRadiologistreq, idValue int) []model.GetManagerOne {
	log := logger.InitLogger()

	var RadiologistData []model.GetManagerOne

	UserId := reqVal.Id

	if UserId == 0 {
		UserId = idValue
	}

	err := db.Raw(query.GetOneListofManagerSQL, 9, UserId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetManagerOne{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].Pan = hashdb.Decrypt(tech.Pan)
		RadiologistData[i].Aadhar = hashdb.Decrypt(tech.Aadhar)
		RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)

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

		if len(hashdb.Decrypt(tech.Pan)) > 0 {
			panFileHelperData, panFileErr := helper.ViewFile("./Assets/Files/" + hashdb.Decrypt(tech.Pan))
			if panFileErr != nil {
				log.Fatalf("Failed to read PAN file: %v", panFileErr)
			}
			if panFileHelperData != nil {
				RadiologistData[i].PanFile = &model.FileData{
					Base64Data:  panFileHelperData.Base64Data,
					ContentType: panFileHelperData.ContentType,
				}
			}
		} else {
			RadiologistData[i].PanFile = nil
		}

		if len(hashdb.Decrypt(tech.Aadhar)) > 0 {
			aadharFileHelperData, aadharFileErr := helper.ViewFile("./Assets/Files/" + hashdb.Decrypt(tech.Aadhar))
			if aadharFileErr != nil {
				log.Fatalf("Failed to read Aadhar file: %v", aadharFileErr)
			}
			if aadharFileHelperData != nil {
				RadiologistData[i].AadharFile = &model.FileData{
					Base64Data:  aadharFileHelperData.Base64Data,
					ContentType: aadharFileHelperData.ContentType,
				}
			}
		} else {
			RadiologistData[i].AadharFile = nil
		}

		if len(hashdb.Decrypt(tech.DrivingLicense)) > 0 {
			drivingLicenseFileHelperData, drivingLicenseErr := helper.ViewFile("./Assets/Files/" + hashdb.Decrypt(tech.DrivingLicense))
			if drivingLicenseErr != nil {
				log.Fatalf("Failed to read Driving License file: %v", drivingLicenseErr)
			}
			if drivingLicenseFileHelperData != nil {
				RadiologistData[i].DrivingLicenseFile = &model.FileData{
					Base64Data:  drivingLicenseFileHelperData.Base64Data,
					ContentType: drivingLicenseFileHelperData.ContentType,
				}
			}
		} else {
			RadiologistData[i].DrivingLicenseFile = nil
		}

		var databaseECFiles []model.GetEducationCertificateFilesModel
		CVFileserr := db.Raw(query.GetECFilesSQL, reqVal.Id).Scan(&databaseECFiles).Error
		if CVFileserr != nil {
			log.Printf("ERROR: Failed to fetch CV files for user ID %d: %v", reqVal.Id, CVFileserr)
			return []model.GetManagerOne{}
		}

		if len(databaseECFiles) > 0 {
			RadiologistData[i].EducationCertificateFile = make([]model.GetEducationCertificateFilesModel, 0, len(databaseECFiles))
			for _, dbCvItem := range databaseECFiles {
				processedCvItem := model.GetEducationCertificateFilesModel{
					ECId:          dbCvItem.ECId,
					ECFileName:    hashdb.Decrypt(dbCvItem.ECFileName),
					ECOldFileName: hashdb.Decrypt(dbCvItem.ECOldFileName),
				}
				ECHelperFileData, ECFileReadErr := helper.ViewFile("./Assets/Files/" + processedCvItem.ECFileName)
				if ECFileReadErr != nil {
					log.Printf("WARNING: Failed to read CV file %s: %v. Skipping file data.", processedCvItem.ECFileName, ECFileReadErr)
					processedCvItem.ECFileData = nil
				} else if ECHelperFileData != nil {
					processedCvItem.ECFileData = &model.FileData{
						Base64Data:  ECHelperFileData.Base64Data,
						ContentType: ECHelperFileData.ContentType,
					}
				} else {
					processedCvItem.ECFileData = nil // Should ideally not happen if error is nil
				}
				RadiologistData[i].EducationCertificateFile = append(RadiologistData[i].EducationCertificateFile, processedCvItem)
			}
		}

	}

	return RadiologistData
}
