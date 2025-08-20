package model

type ViewNotificationReq struct {
	Offset int `json:"offset" binding:"required" mapstructure:"offset"`
}

type ViewNotificationResponse struct {
	Status     bool                `json:"status"`
	Message    string              `json:"message"`
	Data       []NotificationTable `json:"data"`
	TotalCount int                 `json:"totalCount" gorm:"column:total_count"`
}

type NotificationTable struct {
	NId           int    `json:"refNId" gorm:"column:refNId"`
	UserId        int    `json:"refUserId" gorm:"column:refUserId"`
	NMessage      string `json:"refNMessage" gorm:"column:refNMessage"`
	NReadStatus   bool   `json:"refNReadStatus" gorm:"column:refNReadStatus"`
	Nstatus       bool   `json:"refNstatus" gorm:"column:refNstatus"`
	AppointmentId int    `json:"refAppointmentId" gorm:"column:refAppointmentId"`
	NAssignedBy   int    `json:"refNAssignedBy" gorm:"column:refNAssignedBy"`
	NCreatedAt    string `json:"refNCreatedAt" gorm:"column:refNCreatedAt"`
}

type ReadStatusReq struct {
	Id     int  `json:"id" binding:"required" mapstructure:"id"`
	Status bool `json:"status" mapstructure:"status"`
}
