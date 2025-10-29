package model

type MigrationDicomFile struct {
	RefDFId          int    `json:"refDFId" binding:"required" gorm:"column:refDFId"`
	RefUserId        int    `json:"refUserId" binding:"required" gorm:"column:refUserId"`
	RefAppointmentId int    `json:"refAppointmentId" binding:"required" gorm:"column:refAppointmentId"`
	RefDFFilename    string `json:"refDFFilename" binding:"required" gorm:"column:refDFFilename"`
	RefDFSide        string `json:"refDFSide" binding:"required" gorm:"column:refDFSide"`
	IsMigrated       bool   `json:"isMigrated" binding:"required" gorm:"column:isMigrated"`
}
