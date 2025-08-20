package service

import (
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"

	"gorm.io/gorm"
)

func GetNotificationCountService(db *gorm.DB, idValue int) (bool, string, int) {
	log := logger.InitLogger()

	var Response model.ViewNotificationResponse

	var notificationErr = db.Raw(query.TotalReadCountNotificationQuery, idValue).Scan(&Response.TotalCount).Error
	if notificationErr != nil {
		log.Printf("ERROR: Failed to fetch Staff Available: %v", notificationErr)
		return false, "Something went wrong, Try Again", 0
	}

	return true, "Successfully Fetched", Response.TotalCount

}

func NotificationService(db *gorm.DB, reqVal model.ViewNotificationReq, idValue int) model.ViewNotificationResponse {
	log := logger.InitLogger()

	var Response model.ViewNotificationResponse

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.ViewNotificationResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var offset = reqVal.Offset - 10

	var notificationErr = db.Raw(query.ViewNotificationQuery, idValue, offset).Scan(&Response.Data).Error
	if notificationErr != nil {
		log.Printf("ERROR: Failed to fetch Staff Available: %v", notificationErr)
		return model.ViewNotificationResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	var TotalCountErr = db.Raw(query.TotalCountNotificationQuery, idValue).Scan(&Response.TotalCount).Error
	if TotalCountErr != nil {
		log.Printf("ERROR: Failed to fetch Staff Available: %v", TotalCountErr)
		return model.ViewNotificationResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return model.ViewNotificationResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	return model.ViewNotificationResponse{
		Status:     true,
		Message:    "Successfully Fetched",
		Data:       Response.Data,
		TotalCount: Response.TotalCount,
	}
}

func ReadStatusService(db *gorm.DB, reqVal model.ReadStatusReq, idValue int) model.ViewNotificationResponse {
	log := logger.InitLogger()

	var Response model.ViewNotificationResponse

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return model.ViewNotificationResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	var updateErr = db.Exec(query.UpdateReadMessageSQL, reqVal.Status, reqVal.Id).Error
	if updateErr != nil {
		log.Printf("ERROR: Failed to fetch Staff Available: %v", updateErr)
		return model.ViewNotificationResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return model.ViewNotificationResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	return model.ViewNotificationResponse{
		Status:     true,
		Message:    "Successfully Fetched",
		Data:       Response.Data,
		TotalCount: Response.TotalCount,
	}
}
