package model

type GetAllRadiologist struct {
	RefUserId   uint   `json:"refUserId" gorm:"column:refUserId"`
	CustId      string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId      int    `json:"refRTId" gorm:"column:refRTId"`
	FirstName   string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName    string `json:"refUserLastName" gorm:"column:refUserLastName"`
	UserStatus  string `json:"refUserStatus" gorm:"column:refUserStatus"`
	PhoneNumber string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email       string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
}

type GetRadiologistreq struct {
	Id int `json:"id" mapstructure:"id"`
}

type GetCVFilesModel struct {
	CVId          int       `json:"refCVID" gorm:"column:refCVID"`
	CVFileName    string    `json:"refCVFileName" gorm:"column:refCVFileName"`
	CVOldFileName string    `json:"refCVOldFileName" gorm:"column:refCVOldFileName"`
	CVFileData    *FileData `json:"cvFileData" gorm:"-"`
}

type GetEducationCertificateFilesModel struct {
	ECId          int       `json:"refECId" gorm:"column:refECId"`
	ECFileName    string    `json:"refECFileName" gorm:"column:refECFileName"`
	ECOldFileName string    `json:"refECOldFileName" gorm:"column:refECOldFileName"`
	ECFileData    *FileData `json:"educationCertificateFile" gorm:"-"`
}

type FileData struct {
	Base64Data  string `json:"base64Data"`  // base64-encoded file content
	ContentType string `json:"contentType"` // e.g., "image/jpeg"
}

type LicenseFilesModel struct {
	LId          int       `json:"refLId" gorm:"column:refLId"`
	LFileName    string    `json:"refLFileName" gorm:"column:refLFileName"`
	LOldFileName string    `json:"refLOldFileName" gorm:"column:refLOldFileName"`
	LFileData    *FileData `json:"lFileData" gorm:"-"`
}

type MalpracticeModel struct {
	MPId          int       `json:"refMPId" gorm:"column:refMPId"`
	MPFileName    string    `json:"refMPFileName" gorm:"column:refMPFileName"`
	MPOldFileName string    `json:"refMPOldFileName" gorm:"column:refMPOldFileName"`
	MPFileData    *FileData `json:"MPFileData" gorm:"-"`
}

type GetRadiologistOne struct {
	RefUserId                   uint                             `json:"refUserId" gorm:"column:refUserId"`
	CustId                      string                           `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleId                      int                              `json:"refRTId" gorm:"column:refRTId"`
	FirstName                   string                           `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	LastName                    string                           `json:"refUserLastName" gorm:"column:refUserLastName"`
	ProfileImg                  string                           `json:"refUserProfileImg" gorm:"column:refUserProfileImg"`
	DOB                         string                           `json:"refUserDOB" gorm:"column:refUserDOB"`
	UserStatus                  bool                             `json:"refUserStatus" gorm:"column:refUserStatus"`
	UserAgreement               bool                             `json:"refUserAgreementStatus" gorm:"column:refUserAgreementStatus"`
	PhoneNumberCountryCode      string                           `json:"refCODOPhoneNo1CountryCode" gorm:"column:refCODOPhoneNo1CountryCode"`
	PhoneNumber                 string                           `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email                       string                           `json:"refCODOEmail" gorm:"column:refCODOEmail"`
	MBBSRegNo                   string                           `json:"refRAMBBSRegNo" gorm:"column:refRAMBBSRegNo"`
	MDRegNo                     string                           `json:"refRAMDRegNo" gorm:"column:refRAMDRegNo"`
	Specialization              string                           `json:"refRASpecialization" gorm:"column:refRASpecialization"`
	Pan                         string                           `json:"refRAPan" gorm:"column:refRAPan"`
	Aadhar                      string                           `json:"refRAAadhar" gorm:"column:refRAAadhar"`
	DrivingLicense              string                           `json:"refRADrivingLicense" gorm:"column:refRADrivingLicense"`
	DigitalSignature            string                           `json:"refRADigitalSignature" gorm:"column:refRADigitalSignature"`
	ProfileImgFile              *FileData                        `json:"profileImgFile" gorm:"-"`
	PanFile                     *FileData                        `json:"panFile" gorm:"-"`
	AadharFile                  *FileData                        `json:"aadharFile" gorm:"-"`
	DrivingLicenseFile          *FileData                        `json:"drivingLicenseFile" gorm:"-"`
	DigitalSignatureFile        *FileData                        `json:"digitalSignatureFile" gorm:"-"`
	MedicalLicenseSecurity      []GetMedicalLicenseSecurityModel `json:"medicalLicenseSecurity" gorm:"-"`
	MalpracticeInsuranceDetails []MalpracticeModel               `json:"malpracticeinsureance_files" gorm:"-"`
	CVFiles                     []GetCVFilesModel                `json:"cvFiles" gorm:"-"`
	LicenseFiles                []LicenseFilesModel              `json:"licenseFiles" gorm:"-"`
}
