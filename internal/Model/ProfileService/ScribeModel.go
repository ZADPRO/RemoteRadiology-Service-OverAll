package model

type GetScribeOne struct {
	RefUserId                uint                                `json:"refUserId" gorm:"column:refUserId"`
	CustId                   string                              `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId                   int                                 `json:"refRTId" gorm:"column:refRTId"`
	FirstName                string                              `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName                 string                              `json:"refUserLastName" gorm:"column:refUserLastName"`
	ProfileImg               string                              `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB                      string                              `json:"refUserDOB" gorm:"column:refUserDOB"`
	UserStatus               bool                                `json:"refUserStatus" gorm:"column:refUserStatus"`
	UserAgreement            bool                                `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
	PhoneNumberCountryCode   string                              `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	PhoneNumber              string                              `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email                    string                              `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	DrivingLicense           string                              `json:"refSDDrivingLicense" gorm:"column:refSDDrivingLicense"`
	Pan                      string                              `json:"refSDPan" gorm:"column:refSDPan"`
	Aadhar                   string                              `json:"refSDAadhar" gorm:"column:refSDAadhar"`
	ProfileImgFile           *FileData                           `json:"profileImgFile" gorm:"-"`
	PanFile                  *FileData                           `json:"panFile" gorm:"-"`
	DrivingLicenseFile       *FileData                           `json:"drivingLicenseFile" gorm:"-"`
	AadharFile               *FileData                           `json:"aadharFile" gorm:"-"`
	EducationCertificateFile []GetEducationCertificateFilesModel `json:"educationCertificateFiles" gorm:"-"`
}
