package model

type AddTrainingMaterialReq struct {
	FileName string `json:"fileName" binding:"required" mapstructure:"fileName"`
	Path     string `json:"filepath" binding:"required" mapstructure:"filepath"`
}

type CreateMaterialModel struct {
	TMId       int    `json:"refTMId" gorm:"primaryKey;autoIncrement;column:refTMId"`
	TMFileName string `json:"refTMFileName" gorm:"column:refTMFileName"`
	TMFilePath string `json:"refTMFilePath" gorm:"column:refTMFilePath"`
	TMStatus   bool   `json:"refTMStatus" gorm:"column:refTMStatus"`
}

func (CreateMaterialModel) TableName() string {
	return "public.TrainingMaterial"
}

type DeleteTrainingMaterialReq struct {
	Id int `json:"id" binding:"required"`
}
