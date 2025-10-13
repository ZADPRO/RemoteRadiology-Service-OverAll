package service

import (
	s3Service "AuthenticationService/Service/S3"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"context"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

)

func isS3URL(path string) bool {
	log.Printf("IS S3 URL checkin : %v", path)
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

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

	// Determine UserId and ScanCenterId if UserId not provided
	if reqVal.UserId == 0 {
		UserId = idValue
		var MappingData []model.Mapping
		if err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error; err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return []model.GetCoDoctorOne{}
		}
		if len(MappingData) > 0 {
			ScanCenterId = MappingData[0].SCId
		} else {
			ScanCenterId = 0
		}
	}

	// Fetch co-doctor data
	if err := db.Raw(query.GetListofCoDoctorOneSQL, UserId, ScanCenterId).Scan(&RadiologistData).Error; err != nil {
		log.Printf("ERROR: Failed to fetch co-doctors: %v", err)
		return []model.GetCoDoctorOne{}
	}

	// Helper function: return presigned URL if S3, else local file
	getFile := func(filePath, localDir string) (string, *model.FileData) {
		if filePath == "" {
			return "", nil
		}
		if isS3URL(filePath) {
			key := extractS3Key(filePath)
			presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
			if err != nil {
				log.Errorf("Failed to generate presigned URL for %s: %v", filePath, err)
				return "", nil
			}
			return presignedURL, &model.FileData{Base64Data: presignedURL, ContentType: "url"}
		}
		fileData, err := helper.ViewFile(localDir + "/" + filePath)
		if err != nil {
			log.Errorf("Failed to read local file %s: %v", filePath, err)
			return "", nil
		}
		return "", (*model.FileData)(fileData)
	}

	for i, tech := range RadiologistData {
		// Decrypt basic fields
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].SocialSecurityNo = hashdb.Decrypt(tech.SocialSecurityNo)
		RadiologistData[i].DriversLicenseNo = hashdb.Decrypt(tech.DriversLicenseNo)
		RadiologistData[i].Specialization = hashdb.Decrypt(tech.Specialization)
		RadiologistData[i].DigitalSignature = hashdb.Decrypt(tech.DigitalSignature)
		RadiologistData[i].NPI = hashdb.Decrypt(tech.NPI)

		// Profile Image
		RadiologistData[i].ProfileImg, RadiologistData[i].ProfileImgFile =
			getFile(RadiologistData[i].ProfileImg, "./Assets/Profile")

		// Driving License
		RadiologistData[i].DriversLicenseNo, RadiologistData[i].DriversLicenseFile =
			getFile(RadiologistData[i].DriversLicenseNo, "./Assets/Files")

		// Digital Signature
		RadiologistData[i].DigitalSignature, RadiologistData[i].DigitalSignatureFile =
			getFile(RadiologistData[i].DigitalSignature, "./Assets/Profile")

		// License Files
		var databaseLicenseFiles []model.LicenseFilesModel
		if err := db.Raw(query.GetLicenseFilesSQL, UserId).Scan(&databaseLicenseFiles).Error; err == nil && len(databaseLicenseFiles) > 0 {
			RadiologistData[i].LicenseFiles = make([]model.LicenseFilesModel, 0, len(databaseLicenseFiles))
			for _, lf := range databaseLicenseFiles {
				lf.LFileName = hashdb.Decrypt(lf.LFileName)
				lf.LOldFileName = hashdb.Decrypt(lf.LOldFileName)
				lf.LFileName, lf.LFileData = getFile(lf.LFileName, "./Assets/Files")
				RadiologistData[i].LicenseFiles = append(RadiologistData[i].LicenseFiles, lf)
			}
		}

		// Malpractice Files
		var databaseMalpracticeFiles []model.MalpracticeModel
		if err := db.Raw(query.GetMalpracticeFilesSQL, UserId).Scan(&databaseMalpracticeFiles).Error; err == nil && len(databaseMalpracticeFiles) > 0 {
			RadiologistData[i].MalpracticeInsuranceDetails = make([]model.MalpracticeModel, 0, len(databaseMalpracticeFiles))
			for _, mp := range databaseMalpracticeFiles {
				mp.MPFileName = hashdb.Decrypt(mp.MPFileName)
				mp.MPOldFileName = hashdb.Decrypt(mp.MPOldFileName)
				mp.MPFileName, mp.MPFileData = getFile(mp.MPFileName, "./Assets/Files")
				RadiologistData[i].MalpracticeInsuranceDetails = append(RadiologistData[i].MalpracticeInsuranceDetails, mp)
			}
		}
	}

	return RadiologistData
}
