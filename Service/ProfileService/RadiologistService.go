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

func GetAllRadiologistDataService(db *gorm.DB) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofRadiologistSQL, 6).Scan(&RadiologistData).Error
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

func GetRadiologistDataService(db *gorm.DB, reqVal model.GetRadiologistreq, idValue int) []model.GetRadiologistOne {
	log := logger.InitLogger()

	var RadiologistData []model.GetRadiologistOne
	UserId := reqVal.Id
	if UserId == 0 {
		UserId = idValue
	}

	err := db.Raw(query.GetListofRadiologistOneSQL, 6, UserId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch radiologist data: %v", err)
		return []model.GetRadiologistOne{}
	}

	// Helper function: returns presigned URL if S3, else local file
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
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg) // leave as-is
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].MBBSRegNo = hashdb.Decrypt(tech.MBBSRegNo)
		RadiologistData[i].MDRegNo = hashdb.Decrypt(tech.MDRegNo)
		RadiologistData[i].Specialization = hashdb.Decrypt(tech.Specialization)
		RadiologistData[i].Pan = hashdb.Decrypt(tech.Pan)
		RadiologistData[i].Aadhar = hashdb.Decrypt(tech.Aadhar)
		RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)
		RadiologistData[i].DigitalSignature = hashdb.Decrypt(tech.DigitalSignature) // leave as-is

		// Files: PAN, Aadhar, Driving License (use getFile)

		if len(hashdb.Decrypt(tech.Pan)) > 0 {
			panFileHelperData, panFileErr := helper.ViewFile("./Assets/Files/" + hashdb.Decrypt(tech.Pan))
			if panFileErr != nil {
				log.Errorf("Failed to read PAN file: %v", panFileErr)
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
				log.Errorf("Failed to read Aadhar file: %v", aadharFileErr)
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

		RadiologistData[i].DrivingLicense, RadiologistData[i].DrivingLicenseFile = getFile(RadiologistData[i].DrivingLicense, "./Assets/Files")

		// Medical License Security
		var MedicalLicenseSecurity []model.GetMedicalLicenseSecurityModel
		if err := db.Raw(query.GetMedicalLicenseSecuritySQL, UserId).Scan(&MedicalLicenseSecurity).Error; err != nil {
			log.Printf("ERROR: Failed to fetch medical license security: %v", err)
			return []model.GetRadiologistOne{}
		}
		RadiologistData[i].MedicalLicenseSecurity = make([]model.GetMedicalLicenseSecurityModel, 0, len(MedicalLicenseSecurity))
		for _, file := range MedicalLicenseSecurity {
			RadiologistData[i].MedicalLicenseSecurity = append(RadiologistData[i].MedicalLicenseSecurity, model.GetMedicalLicenseSecurityModel{
				MLSId:    file.MLSId,
				MLSState: hashdb.Decrypt(file.MLSState),
				MLSNo:    hashdb.Decrypt(file.MLSNo),
			})
		}

		// CV Files
		var databaseCVFiles []model.GetCVFilesModel
		if err := db.Raw(query.GetCVFilesSQL, UserId).Scan(&databaseCVFiles).Error; err == nil && len(databaseCVFiles) > 0 {
			RadiologistData[i].CVFiles = make([]model.GetCVFilesModel, 0, len(databaseCVFiles))
			for _, dbCvItem := range databaseCVFiles {
				dbCvItem.CVFileName = hashdb.Decrypt(dbCvItem.CVFileName)
				dbCvItem.CVOldFileName = hashdb.Decrypt(dbCvItem.CVOldFileName)
				dbCvItem.CVFileName, dbCvItem.CVFileData = getFile(dbCvItem.CVFileName, "./Assets/Files")
				RadiologistData[i].CVFiles = append(RadiologistData[i].CVFiles, dbCvItem)
			}
		}

		// Malpractice Files
		var databaseMalpracticeFiles []model.MalpracticeModel
		if err := db.Raw(query.GetMalpracticeFilesSQL, UserId).Scan(&databaseMalpracticeFiles).Error; err == nil && len(databaseMalpracticeFiles) > 0 {
			RadiologistData[i].MalpracticeInsuranceDetails = make([]model.MalpracticeModel, 0, len(databaseMalpracticeFiles))
			for _, dbLicenseItem := range databaseMalpracticeFiles {
				dbLicenseItem.MPFileName = hashdb.Decrypt(dbLicenseItem.MPFileName)
				dbLicenseItem.MPOldFileName = hashdb.Decrypt(dbLicenseItem.MPOldFileName)
				dbLicenseItem.MPFileName, dbLicenseItem.MPFileData = getFile(dbLicenseItem.MPFileName, "./Assets/Files")
				RadiologistData[i].MalpracticeInsuranceDetails = append(RadiologistData[i].MalpracticeInsuranceDetails, dbLicenseItem)
			}
		}

		// License Files
		var databaseLicenseFiles []model.LicenseFilesModel
		if err := db.Raw(query.GetLicenseFilesSQL, UserId).Scan(&databaseLicenseFiles).Error; err == nil && len(databaseLicenseFiles) > 0 {
			RadiologistData[i].LicenseFiles = make([]model.LicenseFilesModel, 0, len(databaseLicenseFiles))
			for _, dbLicenseItem := range databaseLicenseFiles {
				dbLicenseItem.LFileName = hashdb.Decrypt(dbLicenseItem.LFileName)
				dbLicenseItem.LOldFileName = hashdb.Decrypt(dbLicenseItem.LOldFileName)
				dbLicenseItem.LFileName, dbLicenseItem.LFileData = getFile(dbLicenseItem.LFileName, "./Assets/Files")
				RadiologistData[i].LicenseFiles = append(RadiologistData[i].LicenseFiles, dbLicenseItem)
			}
		}
	}

	return RadiologistData
}
