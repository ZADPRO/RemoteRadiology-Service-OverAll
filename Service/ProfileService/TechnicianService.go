package service

import (
	s3Service "AuthenticationService/Service/S3"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"context"
	"time"

	"gorm.io/gorm"

)

func GetAllTechnicianDataService(db *gorm.DB, reqVal model.GetReceptionistReq) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofTechnicianSQL, 2, reqVal.ScanID).Scan(&RadiologistData).Error
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

func GetOneTechnicianDataService(db *gorm.DB, reqVal model.GetOneReceptionistReq, idValue int) []model.GetTechnicianOne {
	log := logger.InitLogger()

	var RadiologistData []model.GetTechnicianOne

	UserId := reqVal.UserId
	ScanCenterId := reqVal.ScanID

	if reqVal.UserId == 0 {
		UserId = idValue
		var MappingData []model.Mapping
		err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return []model.GetTechnicianOne{}
		}
		if len(MappingData) > 0 {
			ScanCenterId = MappingData[0].SCId
		} else {
			ScanCenterId = 0
		}
	}

	err := db.Raw(query.GetOneTechnicianSQL, UserId, ScanCenterId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetTechnicianOne{}
	}

	// for i, tech := range RadiologistData {
	// 	RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
	// 	RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
	// 	RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
	// 	RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
	// 	RadiologistData[i].SSNo = hashdb.Decrypt(tech.SSNo)
	// 	RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)
	// 	RadiologistData[i].DigitalSignature = hashdb.Decrypt(tech.DigitalSignature)

	// 	if len(hashdb.Decrypt(tech.ProfileImg)) > 0 {
	// 		profileImgHelperData, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(tech.ProfileImg))
	// 		if viewErr != nil {
	// 			// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
	// 			log.Errorf("Failed to read profile image file: %v", viewErr)
	// 		}
	// 		if profileImgHelperData != nil {
	// 			RadiologistData[i].ProfileImgFile = &model.FileData{
	// 				Base64Data:  profileImgHelperData.Base64Data,
	// 				ContentType: profileImgHelperData.ContentType,
	// 			}
	// 		}
	// 	} else {
	// 		RadiologistData[i].ProfileImgFile = nil
	// 	}

	// 	var databaseLicenseFiles []model.LicenseFilesModel
	// 	LicenseErr := db.Raw(query.GetLicenseFilesSQL, reqVal.UserId).Scan(&databaseLicenseFiles).Error
	// 	if LicenseErr != nil {
	// 		log.Printf("ERROR: Failed to fetch License files for user ID %d: %v", reqVal.UserId, LicenseErr)
	// 		return []model.GetTechnicianOne{}
	// 	}

	// 	if len(databaseLicenseFiles) > 0 {
	// 		RadiologistData[i].LicenseFiles = make([]model.LicenseFilesModel, 0, len(databaseLicenseFiles))
	// 		for _, dbLicenseItem := range databaseLicenseFiles {
	// 			processedLicenseItem := model.LicenseFilesModel{
	// 				LId:          dbLicenseItem.LId,
	// 				LFileName:    hashdb.Decrypt(dbLicenseItem.LFileName),
	// 				LOldFileName: hashdb.Decrypt(dbLicenseItem.LOldFileName),
	// 			}
	// 			licenseHelperFileData, licenseFileReadErr := helper.ViewFile("./Assets/Files/" + processedLicenseItem.LFileName)
	// 			if licenseFileReadErr != nil {
	// 				log.Printf("WARNING: Failed to read License file %s: %v. Skipping file data.", processedLicenseItem.LFileName, licenseFileReadErr)
	// 				processedLicenseItem.LFileData = nil
	// 			} else if licenseHelperFileData != nil {
	// 				processedLicenseItem.LFileData = &model.FileData{
	// 					Base64Data:  licenseHelperFileData.Base64Data,
	// 					ContentType: licenseHelperFileData.ContentType,
	// 				}
	// 			} else {
	// 				processedLicenseItem.LFileData = nil // Should ideally not happen if error is nil
	// 			}
	// 			RadiologistData[i].LicenseFiles = append(RadiologistData[i].LicenseFiles, processedLicenseItem)
	// 		}
	// 	}

	// 	if len(hashdb.Decrypt(tech.DrivingLicense)) > 0 {
	// 		DriversLicenseNoImgHelperData, viewErr := helper.ViewFile("./Assets/Files/" + hashdb.Decrypt(tech.DrivingLicense))
	// 		if viewErr != nil {
	// 			// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
	// 			log.Errorf("Failed to read DrivingLicense file: %v", viewErr)
	// 		}
	// 		if DriversLicenseNoImgHelperData != nil {
	// 			RadiologistData[i].DrivingLicenseFile = &model.FileData{
	// 				Base64Data:  DriversLicenseNoImgHelperData.Base64Data,
	// 				ContentType: DriversLicenseNoImgHelperData.ContentType,
	// 			}
	// 		}
	// 	} else {
	// 		RadiologistData[i].DrivingLicenseFile = nil
	// 	}

	// 	if len(hashdb.Decrypt(tech.DigitalSignature)) > 0 {
	// 		DigitalSignatureHelper, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(tech.DigitalSignature))
	// 		if viewErr != nil {
	// 			// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
	// 			log.Errorf("Failed to read DigitalSignature file: %v", viewErr)
	// 		}
	// 		if DigitalSignatureHelper != nil {
	// 			RadiologistData[i].DigitalSignatureFile = &model.FileData{
	// 				Base64Data:  DigitalSignatureHelper.Base64Data,
	// 				ContentType: DigitalSignatureHelper.ContentType,
	// 			}
	// 		}
	// 	} else {
	// 		RadiologistData[i].DigitalSignatureFile = nil
	// 	}

	// }

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].SSNo = hashdb.Decrypt(tech.SSNo)
		RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)
		RadiologistData[i].DigitalSignature = hashdb.Decrypt(tech.DigitalSignature)

		// ==================== PROFILE IMAGE (LOCAL) ====================
		if len(RadiologistData[i].ProfileImg) > 0 {
			profileImgHelperData, viewErr := helper.ViewFile("./Assets/Profile/" + RadiologistData[i].ProfileImg)
			if viewErr != nil {
				log.Errorf("Failed to read profile image file: %v", viewErr)
			} else if profileImgHelperData != nil {
				RadiologistData[i].ProfileImgFile = &model.FileData{
					Base64Data:  profileImgHelperData.Base64Data,
					ContentType: profileImgHelperData.ContentType,
				}
			}
		} else {
			RadiologistData[i].ProfileImgFile = nil
		}

		// ==================== LICENSE FILES ====================
		var databaseLicenseFiles []model.LicenseFilesModel
		LicenseErr := db.Raw(query.GetLicenseFilesSQL, reqVal.UserId).Scan(&databaseLicenseFiles).Error
		if LicenseErr != nil {
			log.Printf("ERROR: Failed to fetch License files for user ID %d: %v", reqVal.UserId, LicenseErr)
			return []model.GetTechnicianOne{}
		}

		if len(databaseLicenseFiles) > 0 {
			RadiologistData[i].LicenseFiles = make([]model.LicenseFilesModel, 0, len(databaseLicenseFiles))
			for _, dbLicenseItem := range databaseLicenseFiles {
				processedLicenseItem := model.LicenseFilesModel{
					LId:          dbLicenseItem.LId,
					LFileName:    hashdb.Decrypt(dbLicenseItem.LFileName),
					LOldFileName: hashdb.Decrypt(dbLicenseItem.LOldFileName),
				}

				// ✅ If file stored in S3 → generate presigned URL
				if isS3URL(processedLicenseItem.LFileName) {
					key := extractS3Key(processedLicenseItem.LFileName)
					presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
					if err != nil {
						log.Printf("WARNING: Failed to presign License file %s: %v", processedLicenseItem.LFileName, err)
					} else {
						processedLicenseItem.LFileName = presignedURL
					}
					processedLicenseItem.LFileData = nil // since file is remote
				} else {
					// ✅ Local file fallback (existing logic)
					licenseHelperFileData, licenseFileReadErr := helper.ViewFile("./Assets/Files/" + processedLicenseItem.LFileName)
					if licenseFileReadErr != nil {
						log.Printf("WARNING: Failed to read License file %s: %v. Skipping file data.", processedLicenseItem.LFileName, licenseFileReadErr)
						processedLicenseItem.LFileData = nil
					} else if licenseHelperFileData != nil {
						processedLicenseItem.LFileData = &model.FileData{
							Base64Data:  licenseHelperFileData.Base64Data,
							ContentType: licenseHelperFileData.ContentType,
						}
					}
				}

				RadiologistData[i].LicenseFiles = append(RadiologistData[i].LicenseFiles, processedLicenseItem)
			}
		}

		// ==================== DRIVING LICENSE ====================
		if isS3URL(RadiologistData[i].DrivingLicense) {
			key := extractS3Key(RadiologistData[i].DrivingLicense)
			presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
			if err == nil {
				RadiologistData[i].DrivingLicense = presignedURL
				RadiologistData[i].DrivingLicenseFile = nil
			} else {
				log.Errorf("Failed to generate presigned URL for Driving License: %v", err)
			}
		} else if len(RadiologistData[i].DrivingLicense) > 0 {
			DriversLicenseNoImgHelperData, viewErr := helper.ViewFile("./Assets/Files/" + RadiologistData[i].DrivingLicense)
			if viewErr != nil {
				log.Errorf("Failed to read DrivingLicense file: %v", viewErr)
			} else if DriversLicenseNoImgHelperData != nil {
				RadiologistData[i].DrivingLicenseFile = &model.FileData{
					Base64Data:  DriversLicenseNoImgHelperData.Base64Data,
					ContentType: DriversLicenseNoImgHelperData.ContentType,
				}
			}
		} else {
			RadiologistData[i].DrivingLicenseFile = nil
		}

		// ==================== DIGITAL SIGNATURE (LOCAL) ====================
		if len(RadiologistData[i].DigitalSignature) > 0 {
			DigitalSignatureHelper, viewErr := helper.ViewFile("./Assets/Profile/" + RadiologistData[i].DigitalSignature)
			if viewErr != nil {
				log.Errorf("Failed to read DigitalSignature file: %v", viewErr)
			} else if DigitalSignatureHelper != nil {
				RadiologistData[i].DigitalSignatureFile = &model.FileData{
					Base64Data:  DigitalSignatureHelper.Base64Data,
					ContentType: DigitalSignatureHelper.ContentType,
				}
			}
		} else {
			RadiologistData[i].DigitalSignatureFile = nil
		}
	}

	return RadiologistData
}
