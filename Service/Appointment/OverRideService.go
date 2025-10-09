package service

import (
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"

	"gorm.io/gorm"
)

func ListAllOverRideService(db *gorm.DB, reqVal model.OverrideListReq) []model.ListAllDataModel {
	log := logger.InitLogger()

	var ListAllData []model.ListAllDataModel

	err := db.Raw(query.ListAllOverRideSQL, reqVal.ScanCenterId).Scan(&ListAllData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.ListAllDataModel{}
	}

	return ListAllData
}

func WriteOverRideService(db *gorm.DB, reqVal model.WriteOverrideListReq, idValue int) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	OverRideerr := tx.Exec(
		query.UpdateOverRideSQL,
		reqVal.Status,
		reqVal.OverRideId,
	).Error

	if OverRideerr != nil {
		log.Printf("ERROR: Failed to update Override %d: %v\n", reqVal.OverRideId, OverRideerr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	history := model.RefTransHistory{
		TransTypeId: 23,
		THData:      "Updated Override as " + reqVal.Status,
		UserId:      reqVal.UserId,
		THActionBy:  idValue,
	}

	errhistory := db.Create(&history).Error
	if errhistory != nil {
		log.Error("LoginService INSERT ERROR at Trnasaction: " + errhistory.Error())
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Successfully Updated"
}
