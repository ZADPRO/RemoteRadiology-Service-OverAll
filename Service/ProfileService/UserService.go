package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"

	"gorm.io/gorm"
)

func GetUserService(db *gorm.DB, idValue int) model.GetUserResModel {
	log := logger.InitLogger()

	var UserData []model.GetUserModel

	err := db.Raw(query.GetUserModel, idValue).Scan(&UserData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch user data for id %d: %v", idValue, err)
		return model.GetUserResModel{}
	}

	for i, tech := range UserData {
		UserData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		UserData[i].LastName = hashdb.Decrypt(tech.LastName)
	}

	if UserData[0].RoleId != 9 && UserData[0].RoleId != 7 && UserData[0].RoleId != 6 {

		ScanCenterId := 0

		var MappingData []model.Mapping
		err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return model.GetUserResModel{}
		}
		if len(MappingData) > 0 {
			ScanCenterId = MappingData[0].SCId
		} else {
			ScanCenterId = 0
		}

		return model.GetUserResModel{
			Id:           UserData[0].Id,
			CustId:       UserData[0].CustId,
			RoleId:       UserData[0].RoleId,
			Email:        UserData[0].Email,
			FirstName:    UserData[0].FirstName,
			LastName:     UserData[0].LastName,
			ScanCenterId: ScanCenterId,
		}
	} else {
		return model.GetUserResModel{
			Id:           UserData[0].Id,
			CustId:       UserData[0].CustId,
			RoleId:       UserData[0].RoleId,
			Email:        UserData[0].Email,
			FirstName:    UserData[0].FirstName,
			LastName:     UserData[0].LastName,
			ScanCenterId: 0,
		}
	}
}
