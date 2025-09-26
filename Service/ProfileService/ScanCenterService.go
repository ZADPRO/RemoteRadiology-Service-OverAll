package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"fmt"

	"gorm.io/gorm"
)

func GetAllScanCenterService(db *gorm.DB) []model.GetAllScaCenter {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllScaCenter

	err := db.Raw(query.GetAllScanCenter).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetAllScaCenter{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].SCProfile = hashdb.Decrypt(tech.SCProfile)
		RadiologistData[i].SCName = hashdb.Decrypt(tech.SCName)
		RadiologistData[i].SCAddress = hashdb.Decrypt(tech.SCAddress)
		RadiologistData[i].SCWebsite = hashdb.Decrypt(tech.SCWebsite)

		fmt.Println(hashdb.Decrypt(tech.SCProfile))
		fmt.Println(tech.SCId)

		if len(hashdb.Decrypt(tech.SCProfile)) > 0 {
			profileImgHelperData, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(tech.SCProfile))
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
	}

	return RadiologistData
}

func GetScanCenterService(db *gorm.DB, reqVal model.GetRadiologistreq, idValue int) []model.GetAllScaCenter {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllScaCenter

	ScancenterID := reqVal.Id

	if reqVal.Id == 0 {
		var MappingData []model.Mapping
		err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return []model.GetAllScaCenter{}
		}
		if len(MappingData) > 0 {
			ScancenterID = MappingData[0].SCId
		} else {
			ScancenterID = 0
		}
	}

	err := db.Raw(query.GetScanCenter, ScancenterID).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetAllScaCenter{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].SCProfile = hashdb.Decrypt(tech.SCProfile)
		RadiologistData[i].SCName = hashdb.Decrypt(tech.SCName)
		RadiologistData[i].SCAddress = hashdb.Decrypt(tech.SCAddress)
		RadiologistData[i].SCWebsite = hashdb.Decrypt(tech.SCWebsite)

		if len(hashdb.Decrypt(tech.SCProfile)) > 0 {
			profileImgHelperData, viewErr := helper.ViewFile("./Assets/Profile/" + hashdb.Decrypt(tech.SCProfile))
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

	}

	return RadiologistData
}
