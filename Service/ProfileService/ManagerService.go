package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	s3path "AuthenticationService/internal/Helper/S3"
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

		if len(RadiologistData[i].ProfileImg) > 0 {
			s3Key := "images/" + RadiologistData[i].ProfileImg
			log.Printf("\n\n\nS3KEY -> %v", s3Key)
			url, err := s3path.GetS3FileURL(s3Key)
			if err != nil {
				log.Errorf("Failed to generate S3 URL for profile image: %v", err)
			} else {
				RadiologistData[i].ProfileImgFile = &model.FileData{
					Base64Data:  "",
					ContentType: "image/jpeg",
				}
				RadiologistData[i].ProfileImgFile.Base64Data = url
			}
		} else {
			RadiologistData[i].ProfileImgFile = nil
		}

		if len(RadiologistData[i].Pan) > 0 {
			s3Key := "documents/" + RadiologistData[i].Pan
			url, err := s3path.GetS3FileURL(s3Key)
			if err != nil {
				log.Errorf("Failed to generate S3 URL for PAN file: %v", err)
			} else {
				RadiologistData[i].PanFile = &model.FileData{
					Base64Data:  url,
					ContentType: "application/pdf",
				}
			}
		} else {
			RadiologistData[i].PanFile = nil
		}

		if len(RadiologistData[i].Aadhar) > 0 {
			s3Key := "documents/" + RadiologistData[i].Aadhar
			url, err := s3path.GetS3FileURL(s3Key)
			if err != nil {
				log.Errorf("Failed to generate S3 URL for Aadhar file: %v", err)
			} else {
				RadiologistData[i].AadharFile = &model.FileData{
					Base64Data:  url,
					ContentType: "application/pdf",
				}
			}
		} else {
			RadiologistData[i].AadharFile = nil
		}

		if len(RadiologistData[i].DrivingLicense) > 0 {
			s3Key := "documents/" + RadiologistData[i].DrivingLicense
			url, err := s3path.GetS3FileURL(s3Key)
			if err != nil {
				log.Errorf("Failed to generate S3 URL for Driving License: %v", err)
			} else {
				RadiologistData[i].DrivingLicenseFile = &model.FileData{
					Base64Data:  url,
					ContentType: "application/pdf",
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
