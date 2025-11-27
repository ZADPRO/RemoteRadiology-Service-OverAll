package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Analaytics"
	query "AuthenticationService/query/Analaytics"

	"gorm.io/gorm"
)

func GetDailyListService(db *gorm.DB, reqVal model.GetDailyListModel) (bool, []model.DailyListReponse) {
	log := logger.InitLogger()

	var DailyList []model.DailyListReponse

	DailyListErr := db.Raw(query.GetDailyListSQL, reqVal.FromDate, reqVal.ToDate).Scan(&DailyList).Error
	if DailyListErr != nil {
		log.Errorf("Error getting daily list: %v", DailyListErr)
		return false, []model.DailyListReponse{}
	}

	for i := range DailyList {
		DailyList[i].RefUserFirstName = hashdb.Decrypt(DailyList[i].RefUserFirstName)
	}

	return true, DailyList

}
