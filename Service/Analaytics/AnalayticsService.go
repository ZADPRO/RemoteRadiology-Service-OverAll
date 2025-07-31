package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Analaytics"
	query "AuthenticationService/query/Analaytics"

	"gorm.io/gorm"
)

func AdminOverallOneAnalayticsService(db *gorm.DB, reqVal model.AdminOverallOneAnalyticsReq, roleIdValue int) model.AdminOverallAnalyticsResponse {
	log := logger.InitLogger()

	//Model Created
	var response model.AdminOverallAnalyticsResponse

	//6 Months Analaytics
	AdminOverallAnalayticsErr := db.Raw(query.AdminOverallAnalayticsSQL, reqVal.SCId, reqVal.SCId).Scan(&response.AdminScanCenterModel).Error
	if AdminOverallAnalayticsErr != nil {
		log.Fatal(AdminOverallAnalayticsErr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	//Particualr Month Scan Indications
	AdminOverallScanIndicatesAnalayticsErr := db.Raw(query.AdminOverallScanIndicatesAnalayticsSQL, reqVal.StartDate, reqVal.EndDate, reqVal.SCId, reqVal.SCId).Scan(&response.AdminOverallScanIndicatesAnalayticsModel).Error
	if AdminOverallScanIndicatesAnalayticsErr != nil {
		log.Fatal(AdminOverallScanIndicatesAnalayticsErr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	//List All the Scan Center
	err := db.Raw(query.GetAllScanCenter).Scan(&response.GetAllScaCenter).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return model.AdminOverallAnalyticsResponse{}
	}

	//Decrypt the Scan Center List Data
	for i, tech := range response.GetAllScaCenter {
		response.GetAllScaCenter[i].SCProfile = hashdb.Decrypt(tech.SCProfile)
		response.GetAllScaCenter[i].SCName = hashdb.Decrypt(tech.SCName)
		response.GetAllScaCenter[i].SCAddress = hashdb.Decrypt(tech.SCAddress)
		response.GetAllScaCenter[i].SCWebsite = hashdb.Decrypt(tech.SCWebsite)
	}

	//Choosing the UserListIds Based for Master Admin and Scan Center Admin
	var userListIds []int

	switch roleIdValue {
	case 3:
		userListIds = []int{2, 5, 8}
	case 1, 9:
		userListIds = []int{1, 2, 5, 6, 7, 8, 10}
	default:
		userListIds = []int{} // or keep as is since it's initialized empty
	}

	//List the All the User List with the Above Choosen User RoleIds
	var UserListIdsDataerr = db.Raw(query.UserListIdsSQL, userListIds).Scan(&response.UserListIdsData).Error
	if UserListIdsDataerr != nil {
		log.Fatal(UserListIdsDataerr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	//Impression and Recommentation
	ImpressionNRecommentationErr := db.Raw(query.ImpressionNRecommentationScanCenterSQL, reqVal.StartDate, reqVal.EndDate, reqVal.SCId, reqVal.SCId).Scan(&response.ImpressionModel).Error
	if ImpressionNRecommentationErr != nil {
		log.Fatal(ImpressionNRecommentationErr.Error())
		return model.AdminOverallAnalyticsResponse{}
	}

	return response
}

func UserAnalaytics(db *gorm.DB, reqVal model.OneUserReq, UserId int, roleIdValue int) model.OneUserReponse {
	log := logger.InitLogger()

	var response model.OneUserReponse

	//6 Months Analaytics
	AdminOverallAnalayticsErr := db.Raw(query.WellGreenUserAnalayticsSQL, UserId).Scan(&response.AdminScanCenterModel).Error
	if AdminOverallAnalayticsErr != nil {
		log.Fatal(AdminOverallAnalayticsErr.Error())
		return model.OneUserReponse{}
	}

	//Particualr Month Scan Indications
	AdminOverallScanIndicatesAnalayticsErr := db.Raw(query.WellGreenUserIndicatesAnalayticsSQL, UserId, reqVal.StartDate, reqVal.EndDate).Scan(&response.AdminOverallScanIndicatesAnalayticsModel).Error
	if AdminOverallScanIndicatesAnalayticsErr != nil {
		log.Fatal(AdminOverallScanIndicatesAnalayticsErr.Error())
		return model.OneUserReponse{}
	}

	//User Worked Timing
	UserWorkedTimingErr := db.Raw(query.UserWorkedTimingSQL, UserId).Scan(&response.UserAccessTimingModel).Error
	if UserWorkedTimingErr != nil {
		log.Fatal(UserWorkedTimingErr.Error())
		return model.OneUserReponse{}
	}

	//For Each Scan Center How many Count
	if reqVal.RoleId == 6 || reqVal.RoleId == 7 {
		ListScanAppointmentCountErr := db.Raw(query.ListScanAppointmentCountSQL, UserId, reqVal.StartDate, reqVal.EndDate).Scan(&response.ListScanAppointmentCountModel).Error
		if ListScanAppointmentCountErr != nil {
			log.Fatal(ListScanAppointmentCountErr.Error())
			return model.OneUserReponse{}
		}

		//Decrypt the Scan Center Name
		for i, data := range response.ListScanAppointmentCountModel {
			response.ListScanAppointmentCountModel[i].SCName = hashdb.Decrypt(data.SCName)
		}
	}

	//Total Correct and Edit
	TotalCorrectEditErr := db.Raw(query.TotalCorrectEditSQL, UserId, reqVal.StartDate, reqVal.EndDate).Scan(&response.TotalCorrectEdit).Error
	if TotalCorrectEditErr != nil {
		log.Fatal(TotalCorrectEditErr.Error())
		return model.OneUserReponse{}
	}

	//Impression and Recommentation
	ImpressionNRecommentationErr := db.Raw(query.ImpressionNRecommentationSQL, UserId, reqVal.StartDate, reqVal.EndDate).Scan(&response.ImpressionModel).Error
	if ImpressionNRecommentationErr != nil {
		log.Fatal(ImpressionNRecommentationErr.Error())
		return model.OneUserReponse{}
	}

	//Total TAT Timing
	TotalTATErr := db.Raw(query.TotalTATSQL, reqVal.StartDate, reqVal.EndDate, UserId).Scan(&response.DurationBucketModel).Error
	if TotalTATErr != nil {
		log.Fatal(TotalTATErr.Error())
		return model.OneUserReponse{}
	}

	return response

}

func OneUserService(db *gorm.DB, reqVal model.OneUserReq, idValue int, roleIdValue int) model.OneUserReponse {
	if roleIdValue == 1 || roleIdValue == 9 || roleIdValue == 3 {
		response := UserAnalaytics(db, reqVal, reqVal.UserId, roleIdValue)
		return response
	} else {
		response := UserAnalaytics(db, reqVal, idValue, roleIdValue)
		return response
	}
}
