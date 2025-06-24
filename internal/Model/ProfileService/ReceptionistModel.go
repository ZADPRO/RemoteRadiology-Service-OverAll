package model

type GetReceptionistReq struct {
	ScanID int `json:"refSCId" binding:"required" mapstructure:"refSCId"`
}

type GetOneReceptionistReq struct {
	UserId int `json:"refUserId" mapstructure:"refUserId"`
	ScanID int `json:"refSCId" mapstructure:"refSCId"`
}

type GetReceptionist struct {
	RefUserId   uint   `json:"refUserId" gorm:"column:refUserId"`
	CustId      string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId      int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName   string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName    string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserStatus  string `json:"refUserStatus" gorm:"column:refUserStatus"`
	PhoneNumber string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email       string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
}

type GetReceptionistOne struct {
	RefUserId              uint      `json:"refUserId" gorm:"column:refUserId"`
	CustId                 string    `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId                 int       `json:"refRTId" gorm:"column:refRTId"`
	FirstName              string    `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName               string    `json:"refUserLastName" gorm:"column:refUserLastName"`
	ProfileImg             string    `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB                    string    `json:"refUserDOB" gorm:"column:refUserDOB"`
	UserStatus             string    `json:"refUserStatus" gorm:"column:refUserStatus"`
	Agreement              bool      `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
	PhoneNumberCountryCode string    `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	PhoneNumber            string    `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email                  string    `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	ProfileImgFile         *FileData `json:"profileImgFile" gorm:"-"`
	DrivingLicenseFile     *FileData `json:"drivingLicenseFile" gorm:"-"`
	SocialSecurityNo       string    `json:"refRDSSId" binding:"required" gorm:"column:refRDSSId"`
	DrivingLicense         string    `json:"refRDDrivingLicense" binding:"required" gorm:"column:refRDDrivingLicense"`
}

type GetReceptionistMapList struct {
	SCMID          int    `json:"refSCMId" gorm:"column:refSCMId"`
	SCID           int    `json:"refSCId" gorm:"column:refSCId"`
	RTId           int    `json:"refRTId" gorm:"column:refRTId"`
	Status         bool   `json:"refSCMStatus" gorm:"column:refSCMStatus"`
	ScanCenterName string `json:"refSCName" gorm:"column:refSCName"`
	RoleTpyeName   string `json:"refRTName" gorm:"column:refRTName"`
}

type GetReceptionistAudit struct {
	THID      int    `json:"refTHId" gorm:"column:refTHId"`
	TransId   int    `json:"transTypeId" gorm:"column:transTypeId"`
	THData    string `json:"refTHData" gorm:"column:refTHData"`
	THTime    string `json:"refTHTime" gorm:"column:refTHTime"`
	TransName string `json:"transTypeName" gorm:"column:transTypeName"`
	Username  string `json:"Username" gorm:"column:Username"`
}
