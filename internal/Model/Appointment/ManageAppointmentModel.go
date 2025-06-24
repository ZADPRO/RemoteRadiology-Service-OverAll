package model

type AddAppointmentReq struct {
	SCId                 int    `json:"refSCId" binding:"required" mapstructure:"refSCId"`
	AppointmentDate      string `json:"refAppointmentDate" binding:"required" mapstructure:"refAppointmentDate"`
	AppointmentStartTime string `json:"refAppointmentStartTime" binding:"required" mapstructure:"refAppointmentStartTime"`
	AppointmentEndTime   string `json:"refAppointmentEndTime" binding:"required" mapstructure:"refAppointmentEndTime"`
	AppointmentUrgency   bool   `json:"refAppointmentUrgency" mapstructure:"refAppointmentUrgency"`
}

type CreateAppointmentModel struct {
	AppointmentId        int    `json:"refAppointmentId" gorm:"primaryKey;autoIncrement;column:refAppointmentId"`
	UserId               int    `json:"refUserId" gorm:"column:refUserId"`
	SCId                 int    `json:"refSCId" gorm:"column:refSCId"`
	AppointmentDate      string `json:"refAppointmentDate" gorm:"column:refAppointmentDate"`
	AppointmentStartTime string `json:"refAppointmentStartTime" gorm:"column:refAppointmentStartTime"`
	AppointmentEndTime   string `json:"refAppointmentEndTime" gorm:"column:refAppointmentEndTime"`
	AppointmentUrgency   bool   `json:"refAppointmentUrgency" gorm:"column:refAppointmentUrgency"`
	AppointmentStatus    bool   `json:"refAppointmentStatus" gorm:"column:refAppointmentStatus"`
}

func (CreateAppointmentModel) TableName() string {
	return "appointment.refAppointments"
}

type TotalCountModel struct {
	TotalCount int `json:"TotalCount" gorm:"column:TotalCount"`
}

type RefTransHistory struct {
	TransTypeId int    `json:"transTypeId" gorm:"column:transTypeId"`
	THData      string `json:"refTHData" gorm:"column:refTHData"`
	UserId      int    `json:"refUserId" gorm:"column:refUserId"`
	THActionBy  int    `json:"refTHActionBy" gorm:"column:refTHActionBy"`
}
