package model

type GetCoDoctorOne struct {
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
	SocialSecurityNo            string                           `json:"refCDSocialSecurityNo" gorm:"column:refCDSocialSecurityNo"`
	DriversLicenseNo            string                           `json:"drivers_license" gorm:"column:refCDDrivingLicense"`
	Specialization              string                           `json:"Specialization" gorm:"column:refCDSpecialization"`
	DigitalSignature            string                           `json:"digital_signature" gorm:"column:refCDDigitalSignature"`
	CDEaseQTReportAccess        bool                             `json:"refCDEaseQTReportAccess" gorm:"column:refCDEaseQTReportAccess"`
	CDNAsystemReportAccess      bool                             `json:"refCDNAsystemReportAccess" gorm:"column:refCDNAsystemReportAccess"`
	NPI                         string                           `json:"refCDNPI" gorm:"column:refCDNPI"`
	ProfileImgFile              *FileData                        `json:"profileImgFile" gorm:"-"`
	DriversLicenseFile          *FileData                        `json:"driversLicenseFile"  gorm:"-"`
	DigitalSignatureFile        *FileData                        `json:"digitalSignatureFile"  gorm:"-"`
	MedicalLicenseSecurity      []GetMedicalLicenseSecurityModel `json:"medicalLicenseSecurity" gorm:"-"`
	MalpracticeInsuranceDetails []MalpracticeModel               `json:"malpracticeinsureance_files" gorm:"-"`
	LicenseFiles                []LicenseFilesModel              `json:"licenseFiles" gorm:"-"`
}
