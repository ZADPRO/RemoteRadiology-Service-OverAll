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

func GetAllDoctorDataService(db *gorm.DB, reqVal model.GetReceptionistReq) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofDoctorSQL, 5, reqVal.ScanID).Scan(&RadiologistData).Error
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

func GetDoctorDataService(db *gorm.DB, reqVal model.GetOneReceptionistReq, idValue int) []model.GetDoctorOne {
	log := logger.InitLogger()

	var RadiologistData []model.GetDoctorOne

	UserId := reqVal.UserId
	ScanCenterId := reqVal.ScanID

	if reqVal.UserId == 0 {
		UserId = idValue
		var MappingData []model.Mapping
		err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return []model.GetDoctorOne{}
		}
		if len(MappingData) > 0 {
			ScanCenterId = MappingData[0].SCId
		} else {
			ScanCenterId = 0
		}
	}

	err := db.Raw(query.GetListofDoctorOneSQL, UserId, ScanCenterId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetDoctorOne{}
	}

	// Helper function for file access (S3 vs Local)
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
			return presignedURL, nil
		}
		fileData, viewErr := helper.ViewFile(localDir + "/" + filePath)
		if viewErr != nil {
			log.Errorf("Failed to read local file %s: %v", filePath, viewErr)
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

		// Medical License Security
		var MedicalLicenseSecurity []model.GetMedicalLicenseSecurityModel
		err := db.Raw(query.GetMedicalLicenseSecuritySQL, UserId).Scan(&MedicalLicenseSecurity).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch medical license security: %v", err)
			return []model.GetDoctorOne{}
		}
		RadiologistData[i].MedicalLicenseSecurity = make([]model.GetMedicalLicenseSecurityModel, 0, len(MedicalLicenseSecurity))
		for _, file := range MedicalLicenseSecurity {
			RadiologistData[i].MedicalLicenseSecurity = append(RadiologistData[i].MedicalLicenseSecurity, model.GetMedicalLicenseSecurityModel{
				MLSId:    file.MLSId,
				MLSState: hashdb.Decrypt(file.MLSState),
				MLSNo:    hashdb.Decrypt(file.MLSNo),
			})
		}

		// License Files
		var databaseLicenseFiles []model.LicenseFilesModel
		LicenseErr := db.Raw(query.GetLicenseFilesSQL, UserId).Scan(&databaseLicenseFiles).Error
		if LicenseErr != nil {
			log.Printf("ERROR: Failed to fetch License files for user ID %d: %v", UserId, LicenseErr)
			return []model.GetDoctorOne{}
		}
		if len(databaseLicenseFiles) > 0 {
			RadiologistData[i].LicenseFiles = make([]model.LicenseFilesModel, 0, len(databaseLicenseFiles))
			for _, dbLicenseItem := range databaseLicenseFiles {
				processedLicenseItem := model.LicenseFilesModel{
					LId:          dbLicenseItem.LId,
					LFileName:    hashdb.Decrypt(dbLicenseItem.LFileName),
					LOldFileName: hashdb.Decrypt(dbLicenseItem.LOldFileName),
				}
				processedLicenseItem.LFileName, processedLicenseItem.LFileData =
					getFile(processedLicenseItem.LFileName, "./Assets/Files")
				RadiologistData[i].LicenseFiles = append(RadiologistData[i].LicenseFiles, processedLicenseItem)
			}
		}

		// Malpractice Files
		var databaseMalpracticeFiles []model.MalpracticeModel
		MalpracticeErr := db.Raw(query.GetMalpracticeFilesSQL, UserId).Scan(&databaseMalpracticeFiles).Error
		if MalpracticeErr != nil {
			log.Printf("ERROR: Failed to fetch Malpractice files for user ID %d: %v", UserId, MalpracticeErr)
			return []model.GetDoctorOne{}
		}
		if len(databaseMalpracticeFiles) > 0 {
			RadiologistData[i].MalpracticeInsuranceDetails = make([]model.MalpracticeModel, 0, len(databaseMalpracticeFiles))
			for _, dbMalItem := range databaseMalpracticeFiles {
				processedMalItem := model.MalpracticeModel{
					MPId:          dbMalItem.MPId,
					MPFileName:    hashdb.Decrypt(dbMalItem.MPFileName),
					MPOldFileName: hashdb.Decrypt(dbMalItem.MPOldFileName),
				}
				processedMalItem.MPFileName, processedMalItem.MPFileData =
					getFile(processedMalItem.MPFileName, "./Assets/Files")
				RadiologistData[i].MalpracticeInsuranceDetails = append(RadiologistData[i].MalpracticeInsuranceDetails, processedMalItem)
			}
		}
	}

	return RadiologistData
}
