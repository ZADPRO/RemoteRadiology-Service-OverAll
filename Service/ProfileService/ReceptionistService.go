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

func extractS3Key(s3URL string) string {
	prefix := "https://easeqt-health-archive.s3.us-east-2.amazonaws.com/"
	return strings.TrimPrefix(s3URL, prefix)
}

func GetAllReceptionistDataService(db *gorm.DB, reqVal model.GetReceptionistReq) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofReceptionistSQL, 3, reqVal.ScanID).Scan(&RadiologistData).Error
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

func GetOneReceptionistDataService(db *gorm.DB, reqVal model.GetOneReceptionistReq, idValue int) []model.GetReceptionistOne {
	log := logger.InitLogger()

	var ReceptionistData []model.GetReceptionistOne

	UserId := reqVal.UserId
	ScanCenterId := reqVal.ScanID

	if reqVal.UserId == 0 {
		UserId = idValue
		var MappingData []model.Mapping
		err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return []model.GetReceptionistOne{}
		}
		if len(MappingData) > 0 {
			ScanCenterId = MappingData[0].SCId
		} else {
			ScanCenterId = 0
		}
	}

	// Fetch receptionist record
	err := db.Raw(query.GetOneReceptionistSQL, UserId, ScanCenterId).Scan(&ReceptionistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch receptionist: %v", err)
		return []model.GetReceptionistOne{}
	}

	for i, rec := range ReceptionistData {
		// ============ Decrypt all core fields ============
		ReceptionistData[i].FirstName = hashdb.Decrypt(rec.FirstName)
		ReceptionistData[i].LastName = hashdb.Decrypt(rec.LastName)
		ReceptionistData[i].ProfileImg = hashdb.Decrypt(rec.ProfileImg)
		ReceptionistData[i].DOB = hashdb.Decrypt(rec.DOB)
		ReceptionistData[i].SocialSecurityNo = hashdb.Decrypt(rec.SocialSecurityNo)
		ReceptionistData[i].DrivingLicense = hashdb.Decrypt(rec.DrivingLicense)

		if len(hashdb.Decrypt(rec.ProfileImg)) > 0 {
			profileImgHelperData, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(rec.ProfileImg))
			if viewErr != nil {
				log.Errorf("Failed to read profile image file: %v", viewErr)
			}
			if profileImgHelperData != nil {
				ReceptionistData[i].ProfileImgFile = &model.FileData{
					Base64Data: profileImgHelperData.Base64Data, ContentType: profileImgHelperData.ContentType,
				}
			}
		}

		// ============ DRIVING LICENSE HANDLING ============
		if isS3URL(ReceptionistData[i].DrivingLicense) {
			key := extractS3Key(ReceptionistData[i].DrivingLicense)
			presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
			if err != nil {
				log.Errorf("Failed to generate presigned URL for Driving License: %v", err)
			} else {
				ReceptionistData[i].DrivingLicense = presignedURL
				ReceptionistData[i].DrivingLicenseFile = nil
			}
		} else if len(ReceptionistData[i].DrivingLicense) > 0 {
			DriversLicenseFile, viewErr := helper.ViewFile("./Assets/Files/" + ReceptionistData[i].DrivingLicense)
			if viewErr != nil {
				log.Errorf("Failed to read Driving License file: %v", viewErr)
			} else if DriversLicenseFile != nil {
				ReceptionistData[i].DrivingLicenseFile = &model.FileData{
					Base64Data:  DriversLicenseFile.Base64Data,
					ContentType: DriversLicenseFile.ContentType,
				}
			}
		} else {
			ReceptionistData[i].DrivingLicenseFile = nil
		}
	}

	return ReceptionistData
}
