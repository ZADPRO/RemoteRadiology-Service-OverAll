package service

import (
	s3Service "AuthenticationService/Service/S3"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

)

func GetAllPerformingProviderDataService(db *gorm.DB) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofWGPPSQL, 10).Scan(&RadiologistData).Error
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

func GetPerformingProviderDataService(db *gorm.DB, reqVal model.GetRadiologistreq, idValue int) []model.GetWGPPOne {
	log := logger.InitLogger()

	var RadiologistData []model.GetWGPPOne

	UserId := reqVal.Id

	if UserId == 0 {
		UserId = idValue
	}

	err := db.Raw(query.GetListofWGPPOneSQL, 10, UserId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetWGPPOne{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].MBBSRegNo = hashdb.Decrypt(tech.MBBSRegNo)
		RadiologistData[i].MDRegNo = hashdb.Decrypt(tech.MDRegNo)
		RadiologistData[i].Specialization = hashdb.Decrypt(tech.Specialization)
		RadiologistData[i].Pan = hashdb.Decrypt(tech.Pan)
		RadiologistData[i].Aadhar = hashdb.Decrypt(tech.Aadhar)
		RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)
		RadiologistData[i].DigitalSignature = hashdb.Decrypt(tech.DigitalSignature)

		fmt.Println("###################", RadiologistData[i])

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

		if len(hashdb.Decrypt(tech.Pan)) > 0 {
			decryptedPan := hashdb.Decrypt(tech.Pan)
			if isS3URL(decryptedPan) {
				key := extractS3Key(decryptedPan)
				presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
				if err != nil {
					log.Errorf("Failed to generate presigned URL for PAN: %v", err)
				} else {
					RadiologistData[i].Pan = presignedURL
				}
				RadiologistData[i].PanFile = nil
			} else {
				panFileHelperData, panFileErr := helper.ViewFile("./Assets/Files/" + decryptedPan)
				if panFileErr != nil {
					log.Errorf("Failed to read PAN file: %v", panFileErr)
				}
				if panFileHelperData != nil {
					RadiologistData[i].PanFile = &model.FileData{
						Base64Data:  panFileHelperData.Base64Data,
						ContentType: panFileHelperData.ContentType,
					}
				}
			}
		} else {
			RadiologistData[i].PanFile = nil
		}

		if len(hashdb.Decrypt(tech.Aadhar)) > 0 {
			decryptedAadhar := hashdb.Decrypt(tech.Aadhar)
			if isS3URL(decryptedAadhar) {
				key := extractS3Key(decryptedAadhar)
				presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
				if err != nil {
					log.Errorf("Failed to generate presigned URL for Aadhar: %v", err)
				} else {
					RadiologistData[i].Aadhar = presignedURL
				}
				RadiologistData[i].AadharFile = nil
			} else {
				aadharFileHelperData, aadharFileErr := helper.ViewFile("./Assets/Files/" + decryptedAadhar)
				if aadharFileErr != nil {
					log.Errorf("Failed to read Aadhar file: %v", aadharFileErr)
				}
				if aadharFileHelperData != nil {
					RadiologistData[i].AadharFile = &model.FileData{
						Base64Data:  aadharFileHelperData.Base64Data,
						ContentType: aadharFileHelperData.ContentType,
					}
				}
			}
		} else {
			RadiologistData[i].AadharFile = nil
		}

		if len(hashdb.Decrypt(tech.DrivingLicense)) > 0 {
			decryptedDL := hashdb.Decrypt(tech.DrivingLicense)
			if isS3URL(decryptedDL) {
				key := extractS3Key(decryptedDL)
				presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
				if err != nil {
					log.Errorf("Failed to generate presigned URL for Driving License: %v", err)
				} else {
					RadiologistData[i].DrivingLicense = presignedURL
				}
				RadiologistData[i].DrivingLicenseFile = nil
			} else {
				dlHelperData, viewErr := helper.ViewFile("./Assets/Files/" + decryptedDL)
				if viewErr != nil {
					log.Errorf("Failed to read DrivingLicense file: %v", viewErr)
				}
				if dlHelperData != nil {
					RadiologistData[i].DrivingLicenseFile = &model.FileData{
						Base64Data:  dlHelperData.Base64Data,
						ContentType: dlHelperData.ContentType,
					}
				}
			}
		} else {
			RadiologistData[i].DrivingLicenseFile = nil
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

		var MedicalLicenseSecurity []model.GetMedicalLicenseSecurityModel

		MedicalLicenseSecurityerr := db.Raw(query.GetMedicalLicenseSecuritySQL, reqVal.Id).Scan(&MedicalLicenseSecurity).Error
		if MedicalLicenseSecurityerr != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", MedicalLicenseSecurityerr)
			return []model.GetWGPPOne{}
		}

		RadiologistData[i].MedicalLicenseSecurity = make([]model.GetMedicalLicenseSecurityModel, 0, len(MedicalLicenseSecurity))

		for _, file := range MedicalLicenseSecurity {
			RadiologistData[i].MedicalLicenseSecurity = append(RadiologistData[i].MedicalLicenseSecurity, model.GetMedicalLicenseSecurityModel{
				MLSId:    file.MLSId,
				MLSState: hashdb.Decrypt(file.MLSState),
				MLSNo:    hashdb.Decrypt(file.MLSNo),
			})
		}
		var databaseCVFiles []model.GetCVFilesModel
		CVFileserr := db.Raw(query.GetCVFilesSQL, reqVal.Id).Scan(&databaseCVFiles).Error
		if CVFileserr != nil {
			log.Printf("ERROR: Failed to fetch CV files for user ID %d: %v", reqVal.Id, CVFileserr)
			return []model.GetWGPPOne{}
		}

		if len(databaseCVFiles) > 0 {
			RadiologistData[i].CVFiles = make([]model.GetCVFilesModel, 0, len(databaseCVFiles))
			for _, dbCvItem := range databaseCVFiles {
				processedCvItem := model.GetCVFilesModel{
					CVId:          dbCvItem.CVId,
					CVFileName:    hashdb.Decrypt(dbCvItem.CVFileName),
					CVOldFileName: hashdb.Decrypt(dbCvItem.CVOldFileName),
				}

				fileName := processedCvItem.CVFileName

				if strings.HasPrefix(fileName, "https://") {
					// File stored in S3 → Generate presigned URL
					key := extractS3Key(fileName)
					if presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 15*time.Minute); err == nil {
						processedCvItem.CVFileData = &model.FileData{
							Base64Data:  presignedURL,
							ContentType: "application/pdf",
						}
					} else {
						log.Errorf("Failed to generate presigned URL for CV: %v", err)
					}
				} else {
					// Local file → Use helper.ViewFile
					cvHelperFileData, cvFileReadErr := helper.ViewFile("./Assets/Files/" + fileName)
					if cvFileReadErr != nil {
						log.Printf("WARNING: Failed to read CV file %s: %v. Skipping file data.", fileName, cvFileReadErr)
					} else if cvHelperFileData != nil {
						processedCvItem.CVFileData = &model.FileData{
							Base64Data:  cvHelperFileData.Base64Data,
							ContentType: cvHelperFileData.ContentType,
						}
					}
				}
				RadiologistData[i].CVFiles = append(RadiologistData[i].CVFiles, processedCvItem)
			}
		}

		var databaseMalpracticeFiles []model.MalpracticeModel
		MalpracticeErr := db.Raw(query.GetMalpracticeFilesSQL, reqVal.Id).Scan(&databaseMalpracticeFiles).Error
		if MalpracticeErr != nil {
			log.Printf("ERROR: Failed to fetch License files for user ID %d: %v", reqVal.Id, MalpracticeErr)
			return []model.GetWGPPOne{}
		}

		if len(databaseMalpracticeFiles) > 0 {
			RadiologistData[i].MalpracticeInsuranceDetails = make([]model.MalpracticeModel, 0, len(databaseMalpracticeFiles))
			for _, dbLicenseItem := range databaseMalpracticeFiles {
				processedLicenseItem := model.MalpracticeModel{
					MPId:          dbLicenseItem.MPId,
					MPFileName:    hashdb.Decrypt(dbLicenseItem.MPFileName),
					MPOldFileName: hashdb.Decrypt(dbLicenseItem.MPOldFileName),
				}

				fileName := processedLicenseItem.MPFileName

				if strings.HasPrefix(fileName, "https://") {
					key := extractS3Key(fileName)
					if presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 15*time.Minute); err == nil {
						processedLicenseItem.MPFileData = &model.FileData{
							Base64Data:  presignedURL,
							ContentType: "application/pdf",
						}
					} else {
						log.Errorf("Failed to generate presigned URL for Malpractice file: %v", err)
					}
				} else {
					licenseHelperFileData, licenseFileReadErr := helper.ViewFile("./Assets/Files/" + fileName)
					if licenseFileReadErr != nil {
						log.Printf("WARNING: Failed to read License file %s: %v. Skipping file data.", fileName, licenseFileReadErr)
					} else if licenseHelperFileData != nil {
						processedLicenseItem.MPFileData = &model.FileData{
							Base64Data:  licenseHelperFileData.Base64Data,
							ContentType: licenseHelperFileData.ContentType,
						}
					}
				}

				RadiologistData[i].MalpracticeInsuranceDetails = append(RadiologistData[i].MalpracticeInsuranceDetails, processedLicenseItem)
			}
		}

		var databaseLicenseFiles []model.LicenseFilesModel
		LicenseErr := db.Raw(query.GetLicenseFilesSQL, reqVal.Id).Scan(&databaseLicenseFiles).Error
		if LicenseErr != nil {
			log.Printf("ERROR: Failed to fetch License files for user ID %d: %v", reqVal.Id, LicenseErr)
			return []model.GetWGPPOne{}
		}

		if len(databaseLicenseFiles) > 0 {
			RadiologistData[i].LicenseFiles = make([]model.LicenseFilesModel, 0, len(databaseLicenseFiles))
			for _, dbLicenseItem := range databaseLicenseFiles {
				processedLicenseItem := model.LicenseFilesModel{
					LId:          dbLicenseItem.LId,
					LFileName:    hashdb.Decrypt(dbLicenseItem.LFileName),
					LOldFileName: hashdb.Decrypt(dbLicenseItem.LOldFileName),
				}

				fileName := processedLicenseItem.LFileName

				if strings.HasPrefix(fileName, "https://") {
					key := extractS3Key(fileName)
					if presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 15*time.Minute); err == nil {
						processedLicenseItem.LFileData = &model.FileData{
							Base64Data:  presignedURL,
							ContentType: "application/pdf",
						}
					} else {
						log.Errorf("Failed to generate presigned URL for License file: %v", err)
					}
				} else {
					licenseHelperFileData, licenseFileReadErr := helper.ViewFile("./Assets/Files/" + fileName)
					if licenseFileReadErr != nil {
						log.Printf("WARNING: Failed to read License file %s: %v. Skipping file data.", fileName, licenseFileReadErr)
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

	}

	fmt.Println("%%%%%%%%%%%%%%%%", RadiologistData)
	return RadiologistData
}
