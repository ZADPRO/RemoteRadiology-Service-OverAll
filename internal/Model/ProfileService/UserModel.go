package model

type GetUserModel struct {
	Id        uint   `json:"refUserId" gorm:"column:refUserId"`
	CustId    string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId    int    `json:"refRTId" gorm:"column:refRTId"`
	Email     string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	FirstName string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName  string `json:"refUserLastName" gorm:"column:refUserLastName"`
}

type GetUserResModel struct {
	Id           uint   `json:"refUserId"`
	CustId       string `json:"refUserCustId"`
	RoleId       int    `json:"refRTId"`
	Email        string `json:"refCODOEmail"`
	FirstName    string `json:"refUserFirstName"`
	LastName     string `json:"refUserLastName"`
	ScanCenterId int    `json:"refSCId"`
}
