package model

type OverrideListReq struct {
	ScanCenterId int `json:"refSCId" binding:"required" mapstructure:"refSCId"`
}

type ListAllDataModel struct {
	AppointmentId        int    `json:"refAppointmentId" binding:"required"  gorm:"column:refAppointmentId"`
	UserId               int    `json:"refUserId" binding:"required"  gorm:"column:refUserId"`
	SCId                 int    `json:"refSCId" binding:"required"  gorm:"column:refSCId"`
	AppointmentDate      string `json:"refAppointmentDate" binding:"required"  gorm:"column:refAppointmentDate"`
	AppointmentStartTime string `json:"refAppointmentStartTime" binding:"required"  gorm:"column:refAppointmentStartTime"`
	AppointmentEndTime   string `json:"refAppointmentEndTime" binding:"required"  gorm:"column:refAppointmentEndTime"`
	AppointmentUrgency   bool   `json:"refAppointmentUrgency" binding:"required"  gorm:"column:refAppointmentUrgency"`
	AppointmentStatus    bool   `json:"refAppointmentStatus" binding:"required"  gorm:"column:refAppointmentStatus"`
	OVId                 int    `json:"refOVId" binding:"required"  gorm:"column:refOVId"`
	ApprovedStatus       string `json:"refApprovedStatus" binding:"required"  gorm:"column:refApprovedStatus"`
}

type WriteOverrideListReq struct {
	UserId     int    `json:"refUserId" binding:"required"  gorm:"column:refUserId"`
	OverRideId int    `json:"refOVId" binding:"required" mapstructure:"refOVId"`
	Status     string `json:"status" binding:"required" mapstructure:"status"`
}
