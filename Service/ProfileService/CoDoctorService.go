package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"

	"gorm.io/gorm"
)

func GetAllCoDoctorDataService(db *gorm.DB, reqVal model.GetReceptionistReq) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofCoDoctorSQL, 8, reqVal.ScanID).Scan(&RadiologistData).Error
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

func GetDoctorCoDataService(db *gorm.DB, reqVal model.GetOneReceptionistReq, idValue int) []model.GetCoDoctorOne {
	log := logger.InitLogger()

	var RadiologistData []model.GetCoDoctorOne

	UserId := reqVal.UserId
	ScanCenterId := reqVal.ScanID

	if reqVal.UserId == 0 {
		UserId = idValue
		var MappingData []model.Mapping
		err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return []model.GetCoDoctorOne{}
		}
		if len(MappingData) > 0 {
			ScanCenterId = MappingData[0].SCId
		} else {
			ScanCenterId = 0
		}
	}

	err := db.Raw(query.GetListofCoDoctorOneSQL, UserId, ScanCenterId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetCoDoctorOne{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].SocialSecurityNo = hashdb.Decrypt(tech.SocialSecurityNo)
		RadiologistData[i].DriversLicenseNo = hashdb.Decrypt(tech.DriversLicenseNo)
		RadiologistData[i].Specialization = hashdb.Decrypt(tech.Specialization)
		RadiologistData[i].DigitalSignature = hashdb.Decrypt(tech.DigitalSignature)
		RadiologistData[i].NPI = hashdb.Decrypt(tech.NPI)

		if len(hashdb.Decrypt(tech.DriversLicenseNo)) > 0 {
			DriversLicenseNoImgHelperData, viewErr := helper.ViewFile("./Assets/Files/" + hashdb.Decrypt(tech.DriversLicenseNo))
			if viewErr != nil {
				// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
				log.Errorf("Failed to read DrivingLicense file: %v", viewErr)
			}
			if DriversLicenseNoImgHelperData != nil {
				RadiologistData[i].DriversLicenseFile = &model.FileData{
					Base64Data:  DriversLicenseNoImgHelperData.Base64Data,
					ContentType: DriversLicenseNoImgHelperData.ContentType,
				}
			}
		} else {
			RadiologistData[i].DriversLicenseFile = nil
		}

		if len(hashdb.Decrypt(tech.DigitalSignature)) > 0 {
			DigitalSignatureHelper, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(tech.DigitalSignature))
			if viewErr != nil {
				// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
				log.Errorf("Failed to read DigitalSignature file: %v", viewErr)
			}
			if DigitalSignatureHelper != nil {
				RadiologistData[i].DigitalSignatureFile = &model.FileData{
					Base64Data:  DigitalSignatureHelper.Base64Data,
					ContentType: DigitalSignatureHelper.ContentType,
				}
			}
		} else {
			RadiologistData[i].DigitalSignatureFile = nil
		}

		if len(hashdb.Decrypt(tech.ProfileImg)) > 0 {
			profileImgHelperData, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(tech.ProfileImg))
			if viewErr != nil {
				// Consider if Fatalf is appropriate or if logging a warning and setting to nil is better
				log.Errorf("Failed to read profile image file: %v", viewErr)
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

		var MedicalLicenseSecurity []model.GetMedicalLicenseSecurityModel

		MedicalLicenseSecurityerr := db.Raw(query.GetMedicalLicenseSecuritySQL, reqVal.UserId).Scan(&MedicalLicenseSecurity).Error
		if MedicalLicenseSecurityerr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", MedicalLicenseSecurityerr)
			return []model.GetCoDoctorOne{}
		}

		RadiologistData[i].MedicalLicenseSecurity = make([]model.GetMedicalLicenseSecurityModel, 0, len(MedicalLicenseSecurity))

		for _, file := range MedicalLicenseSecurity {
			RadiologistData[i].MedicalLicenseSecurity = append(RadiologistData[i].MedicalLicenseSecurity, model.GetMedicalLicenseSecurityModel{
				MLSId:    file.MLSId,
				MLSState: hashdb.Decrypt(file.MLSState),
				MLSNo:    hashdb.Decrypt(file.MLSNo),
			})
		}

		var databaseLicenseFiles []model.LicenseFilesModel
		LicenseErr := db.Raw(query.GetLicenseFilesSQL, reqVal.UserId).Scan(&databaseLicenseFiles).Error
		if LicenseErr != nil {
			log.Printf("ERROR: Failed to fetch License files for user ID %d: %v", reqVal.UserId, LicenseErr)
			return []model.GetCoDoctorOne{}
		}

		if len(databaseLicenseFiles) > 0 {
			RadiologistData[i].LicenseFiles = make([]model.LicenseFilesModel, 0, len(databaseLicenseFiles))
			for _, dbLicenseItem := range databaseLicenseFiles {
				processedLicenseItem := model.LicenseFilesModel{
					LId:          dbLicenseItem.LId,
					LFileName:    hashdb.Decrypt(dbLicenseItem.LFileName),
					LOldFileName: hashdb.Decrypt(dbLicenseItem.LOldFileName),
				}
				licenseHelperFileData, licenseFileReadErr := helper.ViewFile("./Assets/Files/" + processedLicenseItem.LFileName)
				if licenseFileReadErr != nil {
					log.Printf("WARNING: Failed to read License file %s: %v. Skipping file data.", processedLicenseItem.LFileName, licenseFileReadErr)
					processedLicenseItem.LFileData = nil
				} else if licenseHelperFileData != nil {
					processedLicenseItem.LFileData = &model.FileData{
						Base64Data:  licenseHelperFileData.Base64Data,
						ContentType: licenseHelperFileData.ContentType,
					}
				} else {
					processedLicenseItem.LFileData = nil // Should ideally not happen if error is nil
				}
				RadiologistData[i].LicenseFiles = append(RadiologistData[i].LicenseFiles, processedLicenseItem)
			}
		}

		var databaseMalpracticeFiles []model.MalpracticeModel
		MalpracticeErr := db.Raw(query.GetMalpracticeFilesSQL, reqVal.UserId).Scan(&databaseMalpracticeFiles).Error
		if MalpracticeErr != nil {
			log.Printf("ERROR: Failed to fetch License files for user ID %d: %v", reqVal.UserId, MalpracticeErr)
			return []model.GetCoDoctorOne{}
		}

		if len(databaseMalpracticeFiles) > 0 {
			RadiologistData[i].MalpracticeInsuranceDetails = make([]model.MalpracticeModel, 0, len(databaseMalpracticeFiles))
			for _, dbLicenseItem := range databaseMalpracticeFiles {
				processedLicenseItem := model.MalpracticeModel{
					MPId:          dbLicenseItem.MPId,
					MPFileName:    hashdb.Decrypt(dbLicenseItem.MPFileName),
					MPOldFileName: hashdb.Decrypt(dbLicenseItem.MPOldFileName),
				}
				licenseHelperFileData, licenseFileReadErr := helper.ViewFile("./Assets/Files/" + processedLicenseItem.MPFileName)
				if licenseFileReadErr != nil {
					log.Printf("WARNING: Failed to read License file %s: %v. Skipping file data.", processedLicenseItem.MPFileName, licenseFileReadErr)
					processedLicenseItem.MPFileData = nil
				} else if licenseHelperFileData != nil {
					processedLicenseItem.MPFileData = &model.FileData{
						Base64Data:  licenseHelperFileData.Base64Data,
						ContentType: licenseHelperFileData.ContentType,
					}
				} else {
					processedLicenseItem.MPFileData = nil // Should ideally not happen if error is nil
				}
				RadiologistData[i].MalpracticeInsuranceDetails = append(RadiologistData[i].MalpracticeInsuranceDetails, processedLicenseItem)
			}
		}

	}

	return RadiologistData
}
