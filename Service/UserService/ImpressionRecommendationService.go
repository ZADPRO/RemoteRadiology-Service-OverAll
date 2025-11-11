package service

import (
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/UserService"
	query "AuthenticationService/query/UserService"

	"gorm.io/gorm"
)

func GetImpressionRecommendationService(db *gorm.DB) ([]model.GetAllImpressionRecommendationCategory, []model.ImpressionRecommendationValModel) {
	log := logger.InitLogger()

	var CategoryData []model.GetAllImpressionRecommendationCategory

	err := db.Raw(query.GetCategoryDataSQL).Scan(&CategoryData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetAllImpressionRecommendationCategory{}, []model.ImpressionRecommendationValModel{}
	}

	GetAllImpressionRecommendation := []model.ImpressionRecommendationValModel{}

	errCheckNewImpressionRecommendation := db.Raw(query.GetAllImpressionRecommendationSQL).Scan(&GetAllImpressionRecommendation).Error
	if errCheckNewImpressionRecommendation != nil {
		log.Printf("ERROR: Failed to fetch Check New Impression Recommendation: %v", errCheckNewImpressionRecommendation)
		return []model.GetAllImpressionRecommendationCategory{}, []model.ImpressionRecommendationValModel{}
	}

	return CategoryData, GetAllImpressionRecommendation
}

func AddImpressionRecommendationService(db *gorm.DB, reqVal model.AddImpressionRecommendationReq) (bool, string) {
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

	//Check the Order Id and Impression Recommendation ID
	CheckNewImpressionRecommendation := []model.CheckImpressionRecommendationCategory{}

	errCheckNewImpressionRecommendation := tx.Raw(query.CheckNewImpressionRecommendationSQL, reqVal.ImpressionRecommendationId).Scan(&CheckNewImpressionRecommendation).Error
	if errCheckNewImpressionRecommendation != nil {
		log.Printf("ERROR: Failed to fetch Check New Impression Recommendation: %v", errCheckNewImpressionRecommendation)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if len(CheckNewImpressionRecommendation) > 0 {
		if CheckNewImpressionRecommendation[0].CheckStatus {
			return false, "Impression and Recommendation ID Already Exists"
		}
	} else {
		return false, "Something went wrong, Try Again"
	}

	//Insert New Impression and Recommendation
	errNewImpressionRecommendation := tx.Exec(
		query.InsertNewImpressionRecommendationSQL,
		reqVal.CategoryId,
		CheckNewImpressionRecommendation[0].Next_order_id,
		reqVal.SystemType,
		reqVal.ImpressionRecommendationId,
		reqVal.ImpressionShortDescription,
		reqVal.ImpressionLongDescription,
		reqVal.ImpressionTextColor,
		reqVal.RecommendationShortDescription,
		reqVal.RecommendationLongDescription,
		reqVal.RecommendationTextColor,
	).Error
	if errNewImpressionRecommendation != nil {
		log.Printf("ERROR: Failed to Add New Impression Recommendation: %v", errNewImpressionRecommendation)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Added!"

}

func UpdateImpressionRecommendationService(db *gorm.DB, reqVal model.UpdateImpressionRecommendationReq) (bool, string) {
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

	//Check the Order Id and Impression Recommendation ID
	CheckUpdateImpressionRecommendation := []model.CheckImpressionRecommendationCategory{}

	errCheckUpdateImpressionRecommendation := tx.Raw(query.CheckUpdateImpressionRecommendationSQL, reqVal.ImpressionRecommendationId, reqVal.Id).Scan(&CheckUpdateImpressionRecommendation).Error
	if errCheckUpdateImpressionRecommendation != nil {
		log.Printf("ERROR: Failed to fetch Check Update Impression Recommendation: %v", errCheckUpdateImpressionRecommendation)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if len(CheckUpdateImpressionRecommendation) > 0 {
		if CheckUpdateImpressionRecommendation[0].CheckStatus {
			return false, "Impression and Recommendation ID Already Exists"
		}
	} else {
		return false, "Something went wrong, Try Again"
	}

	//Update Impression and Recommendation
	errUpdateImpressionRecommendation := tx.Exec(
		query.UpdateImpressionRecommendationSQL,
		reqVal.CategoryId,
		reqVal.ImpressionRecommendationId,
		reqVal.ImpressionShortDescription,
		reqVal.ImpressionLongDescription,
		reqVal.ImpressionTextColor,
		reqVal.RecommendationShortDescription,
		reqVal.RecommendationLongDescription,
		reqVal.RecommendationTextColor,
		reqVal.Id,
	).Error
	if errUpdateImpressionRecommendation != nil {
		log.Printf("ERROR: Failed to Add New Impression Recommendation: %v", errUpdateImpressionRecommendation)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Updated!"

}

func DeleteImpressionRecommendationService(db *gorm.DB, reqVal model.UpdateImpressionRecommendationReq) (bool, string) {
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

	//Delete Impression and Recommendation
	errDeleteImpressionRecommendation := tx.Exec(
		query.DeleteImpressionRecommendationSQL,
		reqVal.Id,
	).Error
	if errDeleteImpressionRecommendation != nil {
		log.Printf("ERROR: Failed to Delete Impression Recommendation: %v", errDeleteImpressionRecommendation)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Deleted!"

}

func UpdateOrderImpressionRecommendationService(db *gorm.DB, reqVal model.UpdateOrderImpressionRecommendationReq) (bool, string) {
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

	//Update Order ID
	for _, item := range reqVal.OrderData {
		errUpdateOrder := tx.Exec(
			query.UpdateOrderImpressionRecommendationSQL,
			item.OrderId,
			item.Id,
		).Error
		if errUpdateOrder != nil {
			log.Printf("ERROR: Failed to Update Order Impression Recommendation: %v", errUpdateOrder)
			tx.Rollback()
			return false, "Something went wrong, Try Again"
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Ordered!"

}

func GetFooterReportService(db *gorm.DB) (bool, string, string) {
	log := logger.InitLogger()

	var ReportFooter []model.GetReportFooterModel

	err := db.Raw(query.GetReportFooterSQL).Scan(&ReportFooter).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return false, "Something went wrong, Try Again", ""
	}

	var FooterData = ""
	if len(ReportFooter) > 0 {
		FooterData = ReportFooter[0].RefFRContent
	}

	return true, "Succcessfully Fetched!", FooterData
}

func SaveFooterReportService(db *gorm.DB, reqVal model.SaveReportFooterReq) (bool, string) {
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

	//Update the Report Footer
	errUpdateReportFooter := tx.Exec(
		query.UpdateReportFooterSQL,
		reqVal.ReportText,
	).Error
	if errUpdateReportFooter != nil {
		log.Printf("ERROR: Failed to Update Report Footer: %v", errUpdateReportFooter)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Fetched!"
}
