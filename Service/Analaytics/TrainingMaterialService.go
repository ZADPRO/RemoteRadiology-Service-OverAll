package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Analaytics"
	query "AuthenticationService/query/Analaytics"

	"gorm.io/gorm"
)

func AddTrainingMaterialService(db *gorm.DB, reqVal model.AddTrainingMaterialReq, idValue int) (bool, string) {
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

	Material := model.CreateMaterialModel{
		TMFileName: hashdb.Encrypt(reqVal.FileName),
		TMFilePath: hashdb.Encrypt(reqVal.Path),
		TMStatus:   true,
	}

	MaterialErr := tx.Create(&Material).Error
	if MaterialErr != nil {
		log.Error(MaterialErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Added Successfully"
}

func ListTrainingMaterialService(db *gorm.DB, idValue int) []model.CreateMaterialModel {
	log := logger.InitLogger()

	var ListFiles []model.CreateMaterialModel

	ListFilesErr := db.Raw(query.ListTrainingFilesSQL).Scan(&ListFiles).Error
	if ListFilesErr != nil {
		log.Error(ListFilesErr)
		return []model.CreateMaterialModel{}
	}

	for i, data := range ListFiles {
		ListFiles[i].TMFileName = hashdb.Decrypt(data.TMFileName)
		ListFiles[i].TMFilePath = hashdb.Decrypt(data.TMFilePath)
	}

	return ListFiles
}

func DeleteTrainingMaterialService(db *gorm.DB, reqVal model.DeleteTrainingMaterialReq, idValue int) (bool, string) {
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

	DeleeErr := tx.Exec(query.DeleteTrainingSQL, reqVal.Id).Error
	if DeleeErr != nil {
		log.Error(DeleeErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Deleted Successfully"
}
