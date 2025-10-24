package service

import (
	s3Service "AuthenticationService/Service/S3"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GetAllScribeDataService(db *gorm.DB) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofScribeSQL, 7).Scan(&RadiologistData).Error
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

func GetScribeDataService(db *gorm.DB, reqVal model.GetRadiologistreq, idValue int) []model.GetScribeOne {
	log := logger.InitLogger()

	var RadiologistData []model.GetScribeOne
	UserId := reqVal.Id
	if UserId == 0 {
		UserId = idValue
	}

	// Fetch scribe data from DB
	err := db.Raw(query.GetListofScribeOneSQL, 7, UserId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scribe data: %v", err)
		return []model.GetScribeOne{}
	}

	// Helper: Check if URL is S3
	isS3URL := func(url string) bool {
		return strings.HasPrefix(url, "https://easeqt-health-archive.s3")
	}

	// Helper: Generate FileData (S3 presigned URL or local file)
	getFile := func(filePath, contentType, localDir string) *model.FileData {
		if filePath == "" {
			return nil
		}

		// S3 file
		if isS3URL(filePath) {
			key := extractS3Key(filePath)
			presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
			if err != nil {
				log.Errorf("Failed to generate presigned URL for %s: %v", filePath, err)
				return nil
			}
			return &model.FileData{
				Base64Data:  presignedURL,
				ContentType: contentType,
			}
		}

		// Local file
		fileData, err := helper.ViewFile(localDir + "/" + filePath)
		if err != nil {
			log.Errorf("Failed to read local file %s: %v", filePath, err)
			return nil
		}

		return (*model.FileData)(fileData)
	}

	for i, tech := range RadiologistData {
		// Decrypt all personal fields
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].Pan = hashdb.Decrypt(tech.Pan)
		RadiologistData[i].Aadhar = hashdb.Decrypt(tech.Aadhar)
		RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)

		// Profile image
		RadiologistData[i].ProfileImgFile = getFile(RadiologistData[i].ProfileImg, "image/jpeg", "Profile")

		// Documents
		RadiologistData[i].PanFile = getFile(RadiologistData[i].Pan, "application/pdf", "Files")
		RadiologistData[i].AadharFile = getFile(RadiologistData[i].Aadhar, "application/pdf", "Files")
		RadiologistData[i].DrivingLicenseFile = getFile(RadiologistData[i].DrivingLicense, "application/pdf", "Files")

		// Education Certificates
		var databaseECFiles []model.GetEducationCertificateFilesModel
		ECFilesErr := db.Raw(query.GetECFilesSQL, UserId).Scan(&databaseECFiles).Error
		if ECFilesErr != nil {
			log.Printf("ERROR: Failed to fetch EC files for user ID %d: %v", UserId, ECFilesErr)
			return []model.GetScribeOne{}
		}

		if len(databaseECFiles) > 0 {
			RadiologistData[i].EducationCertificateFile = make([]model.GetEducationCertificateFilesModel, 0, len(databaseECFiles))
			for _, dbCvItem := range databaseECFiles {
				processedCvItem := model.GetEducationCertificateFilesModel{
					ECId:          dbCvItem.ECId,
					ECFileName:    hashdb.Decrypt(dbCvItem.ECFileName),
					ECOldFileName: hashdb.Decrypt(dbCvItem.ECOldFileName),
				}

				// Generate FileData
				processedCvItem.ECFileData = getFile(processedCvItem.ECFileName, "application/pdf", "Files")
				RadiologistData[i].EducationCertificateFile = append(RadiologistData[i].EducationCertificateFile, processedCvItem)
			}
		}
	}

	return RadiologistData
}
