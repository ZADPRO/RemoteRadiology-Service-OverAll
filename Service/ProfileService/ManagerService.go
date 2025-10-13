package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	s3path "AuthenticationService/internal/Helper/S3"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"strings"

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

	// Helper function to determine if a string is an S3 URL
	isS3URL := func(url string) bool {
		return strings.HasPrefix(url, "https://easeqt-health-archive.s3")
	}

	// Helper to generate FileData
	generateFileData := func(fileName, fileType, folder string) *model.FileData {
		if len(fileName) == 0 {
			return nil
		}

		// If it's already a valid S3 URL
		if isS3URL(fileName) {
			return &model.FileData{
				Base64Data:  fileName,
				ContentType: fileType,
			}
		}

		s3Key := folder + "/" + fileName
		url, err := s3path.GetS3FileURL(s3Key)
		if err != nil {
			log.Errorf("Failed to generate S3 URL for %s: %v", fileName, err)
			return nil
		}
		return &model.FileData{
			Base64Data:  url,
			ContentType: fileType,
		}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].Pan = hashdb.Decrypt(tech.Pan)
		RadiologistData[i].Aadhar = hashdb.Decrypt(tech.Aadhar)
		RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)

		// Profile image (image type, no token needed if S3 URL)
		RadiologistData[i].ProfileImgFile = generateFileData(RadiologistData[i].ProfileImg, "image/jpeg", "images")

		// Documents (PDF type)
		RadiologistData[i].PanFile = generateFileData(RadiologistData[i].Pan, "application/pdf", "documents")
		RadiologistData[i].AadharFile = generateFileData(RadiologistData[i].Aadhar, "application/pdf", "documents")
		RadiologistData[i].DrivingLicenseFile = generateFileData(RadiologistData[i].DrivingLicense, "application/pdf", "documents")

		// Education Certificates
		var databaseECFiles []model.GetEducationCertificateFilesModel
		CVFileserr := db.Raw(query.GetECFilesSQL, UserId).Scan(&databaseECFiles).Error
		if CVFileserr != nil {
			log.Printf("ERROR: Failed to fetch CV files for user ID %d: %v", UserId, CVFileserr)
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

				// Generate presigned URL for document
				processedCvItem.ECFileData = generateFileData(processedCvItem.ECFileName, "application/pdf", "documents")
				RadiologistData[i].EducationCertificateFile = append(RadiologistData[i].EducationCertificateFile, processedCvItem)
			}
		}
	}

	return RadiologistData
}
