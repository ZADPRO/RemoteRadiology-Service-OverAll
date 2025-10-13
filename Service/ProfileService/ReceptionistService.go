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

	var RadiologistData []model.GetReceptionistOne

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

	// fmt.Println(UserId)
	// fmt.Println(ScanCenterId)

	err := db.Raw(query.GetOneReceptionistSQL, UserId, ScanCenterId).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetReceptionistOne{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].SocialSecurityNo = hashdb.Decrypt(tech.SocialSecurityNo)
		RadiologistData[i].DrivingLicense = hashdb.Decrypt(tech.DrivingLicense)

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
		}

		// Driving License
		if isS3URL(RadiologistData[i].DrivingLicense) {
			log.Print("\n\n\nDriving License path => %v", RadiologistData[i].DrivingLicense)
			key := extractS3Key(RadiologistData[i].DrivingLicense) // "documents/ee84d792-71c8-4a39-9d6d-747316c9084b_20251009234001.pdf"
			presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
			if err == nil {
				RadiologistData[i].DrivingLicense = presignedURL
			}
		} else if len(RadiologistData[i].DrivingLicense) > 0 {
			data, err := helper.ViewFile("./Assets/Files/" + RadiologistData[i].DrivingLicense)
			if err == nil {
				RadiologistData[i].DrivingLicenseFile = &model.FileData{
					Base64Data:  data.Base64Data,
					ContentType: data.ContentType,
				}
			}
		}
	}

	return RadiologistData
}
