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
	var ManagerData []model.GetManagerOne

	UserId := reqVal.Id
	if UserId == 0 {
		UserId = idValue
	}

	// Fetch manager details
	err := db.Raw(query.GetOneListofManagerSQL, 9, UserId).Scan(&ManagerData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch manager data: %v", err)
		return []model.GetManagerOne{}
	}

	// =============== Helper: Check if URL is from S3 ===============
	isS3URL := func(url string) bool {
		return strings.HasPrefix(url, "https://") && strings.Contains(url, "amazonaws.com")
	}

	// =============== Helper: Generate presigned or local file data ===============
	generateFileData := func(fileName, fileType, folder string) *model.FileData {
		if len(fileName) == 0 {
			return nil
		}

		// ✅ If already full S3 URL → return directly
		if isS3URL(fileName) {
			key := extractS3Key(fileName)
			presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
			if err != nil {
				log.Errorf("Failed to presign S3 URL for %s: %v", fileName, err)
				return &model.FileData{Base64Data: fileName, ContentType: fileType} // fallback
			}
			return &model.FileData{Base64Data: presignedURL, ContentType: fileType}
		}

		// ✅ Local file fallback
		localPath := fmt.Sprintf("./Assets/%s/%s", folder, fileName)
		fileData, err := helper.ViewFile(localPath)
		if err != nil {
			log.Warnf("Local file not found: %s (%v)", localPath, err)
			return nil
		}
		return &model.FileData{
			Base64Data:  fileData.Base64Data,
			ContentType: fileData.ContentType,
		}
	}

	// =============== Populate Data for Each Manager ===============
	for i, mgr := range ManagerData {
		ManagerData[i].FirstName = hashdb.Decrypt(mgr.FirstName)
		ManagerData[i].LastName = hashdb.Decrypt(mgr.LastName)
		ManagerData[i].ProfileImg = hashdb.Decrypt(mgr.ProfileImg)
		ManagerData[i].DOB = hashdb.Decrypt(mgr.DOB)
		ManagerData[i].Pan = hashdb.Decrypt(mgr.Pan)
		ManagerData[i].Aadhar = hashdb.Decrypt(mgr.Aadhar)
		ManagerData[i].DrivingLicense = hashdb.Decrypt(mgr.DrivingLicense)

		// ===== Profile Image (JPEG / PNG) =====
		ManagerData[i].ProfileImgFile = generateFileData(ManagerData[i].ProfileImg, "image/jpeg", "Profile")

		// ===== Documents (PDF) =====
		ManagerData[i].PanFile = generateFileData(ManagerData[i].Pan, "application/pdf", "Files")
		ManagerData[i].AadharFile = generateFileData(ManagerData[i].Aadhar, "application/pdf", "Files")
		ManagerData[i].DrivingLicenseFile = generateFileData(ManagerData[i].DrivingLicense, "application/pdf", "Files")

		// ===== Education Certificates =====
		var educationFiles []model.GetEducationCertificateFilesModel
		ecErr := db.Raw(query.GetECFilesSQL, UserId).Scan(&educationFiles).Error
		if ecErr != nil {
			log.Printf("ERROR: Failed to fetch Education Certificate files for user ID %d: %v", UserId, ecErr)
			continue
		}

		if len(educationFiles) > 0 {
			ManagerData[i].EducationCertificateFile = make([]model.GetEducationCertificateFilesModel, 0, len(educationFiles))
			for _, dbEC := range educationFiles {
				decryptedEC := model.GetEducationCertificateFilesModel{
					ECId:          dbEC.ECId,
					ECFileName:    hashdb.Decrypt(dbEC.ECFileName),
					ECOldFileName: hashdb.Decrypt(dbEC.ECOldFileName),
				}
				decryptedEC.ECFileData = generateFileData(decryptedEC.ECFileName, "application/pdf", "Files")
				ManagerData[i].EducationCertificateFile = append(ManagerData[i].EducationCertificateFile, decryptedEC)
			}
		}
	}

	return ManagerData
}
