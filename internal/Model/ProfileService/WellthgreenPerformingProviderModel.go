package model

type GetWGPPOne struct {
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
	MBBSRegNo                   string                           `json:"refWGPPMBBSRegNo" gorm:"column:refWGPPMBBSRegNo"`
	MDRegNo                     string                           `json:"refWGPPMDRegNo" gorm:"column:refWGPPMDRegNo"`
	Specialization              string                           `json:"refWGPPSpecialization" gorm:"column:refWGPPSpecialization"`
	Pan                         string                           `json:"refWGPPPan" gorm:"column:refWGPPPan"`
	Aadhar                      string                           `json:"refWGPPAadhar" gorm:"column:refWGPPAadhar"`
	DrivingLicense              string                           `json:"refWGPPDrivingLicense" gorm:"column:refWGPPDrivingLicense"`
	DigitalSignature            string                           `json:"refWGPPDigitalSignature" gorm:"column:refWGPPDigitalSignature"`
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
