package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Analaytics"
	query "AuthenticationService/query/Analaytics"
	"time"

	"gorm.io/gorm"
)

func AdminOverallAnalayticsService(db *gorm.DB) model.AdminOverallAnalyticsResponse {
	log := logger.InitLogger()

	var response model.AdminOverallAnalyticsResponse

	AdminOverallAnalayticsErr := db.Raw(query.AdminOverallAnalayticsSQL, 0, 0).Scan(&response.AdminScanCenterModel).Error
	if AdminOverallAnalayticsErr != nil {
		log.Fatal(AdminOverallAnalayticsErr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	current := time.Now()
	formatted := current.Format("2006-01")

	AdminOverallScanIndicatesAnalayticsErr := db.Raw(query.AdminOverallScanIndicatesAnalayticsSQL, formatted, 0, 0).Scan(&response.AdminOverallScanIndicatesAnalayticsModel).Error
	if AdminOverallScanIndicatesAnalayticsErr != nil {
		log.Fatal(AdminOverallScanIndicatesAnalayticsErr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	err := db.Raw(query.GetAllScanCenter).Scan(&response.GetAllScaCenter).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return model.AdminOverallAnalyticsResponse{}
	}

	for i, tech := range response.GetAllScaCenter {
		response.GetAllScaCenter[i].SCProfile = hashdb.Decrypt(tech.SCProfile)
		response.GetAllScaCenter[i].SCName = hashdb.Decrypt(tech.SCName)
		response.GetAllScaCenter[i].SCAddress = hashdb.Decrypt(tech.SCAddress)
		response.GetAllScaCenter[i].SCWebsite = hashdb.Decrypt(tech.SCWebsite)
	}

	return response
}

func AdminOverallOneAnalayticsService(db *gorm.DB, reqVal model.AdminOverallOneAnalyticsReq) model.AdminOverallAnalyticsResponse {
	log := logger.InitLogger()

	var response model.AdminOverallAnalyticsResponse

	AdminOverallAnalayticsErr := db.Raw(query.AdminOverallAnalayticsSQL, reqVal.SCId, reqVal.SCId).Scan(&response.AdminScanCenterModel).Error
	if AdminOverallAnalayticsErr != nil {
		log.Fatal(AdminOverallAnalayticsErr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	AdminOverallScanIndicatesAnalayticsErr := db.Raw(query.AdminOverallScanIndicatesAnalayticsSQL, reqVal.Monthyear, reqVal.SCId, reqVal.SCId).Scan(&response.AdminScanCenterModel).Error
	if AdminOverallScanIndicatesAnalayticsErr != nil {
		log.Fatal(AdminOverallScanIndicatesAnalayticsErr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	return response
}
