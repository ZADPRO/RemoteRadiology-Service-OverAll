package model

type GetUserModel struct {
	Id                      uint      `json:"refUserId" gorm:"column:refUserId"`
	CustId                  string    `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId                  int       `json:"refRTId" gorm:"column:refRTId"`
	Email                   string    `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	FirstName               string    `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName                string    `json:"refUserLastName" gorm:"column:refUserLastName"`
	DOB                     string    `json:"refUserDOB" gorm:"column:refUserDOB"`
	Gender                  string    `json:"refUserGender" gorm:"column:refUserGender"`
	CODOPhoneNo1CountryCode string    `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	CODOPhoneNo1            string    `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	UserProfileImg          string    `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	ProfileImgFile          *FileData `json:"profileImgFile" gorm:"-"`
}

type GetUserResModel struct {
	Id                      uint      `json:"refUserId"`
	CustId                  string    `json:"refUserCustId"`
	RoleId                  int       `json:"refRTId"`
	Email                   string    `json:"refCODOEmail"`
	FirstName               string    `json:"refUserFirstName"`
	LastName                string    `json:"refUserLastName"`
	DOB                     string    `json:"refUserDOB"`
	Gender                  string    `json:"refUserGender"`
	ScanCenterId            int       `json:"refSCId"`
	ScanCenterCustId        string    `json:"refSCCustId"`
	CODOPhoneNo1CountryCode string    `json:"refCODOPhoneNo1CountryCode"`
	CODOPhoneNo1            string    `json:"refCODOPhoneNo1"`
	ProfileImgFile          *FileData `json:"profileImgFile" gorm:"-"`
}
