package model

type GetMedicalLicenseSecurityModel struct {
	MLSId    int    `json:"refMLSId"  gorm:"column:refMLSId"`
	MLSState string `json:"refMLSState" gorm:"column:refMLSState"`
	MLSNo    string `json:"refMLSNo" gorm:"column:refMLSNo"`
	MLStatus string `json:"refMLStatus" gorm:"column:refMLStatus"`
}

type GetDoctorOne struct {
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
	SocialSecurityNo            string                           `json:"refDDSocialSecurityNo" gorm:"column:refDDSocialSecurityNo"`
	DriversLicenseNo            string                           `json:"drivers_license" gorm:"column:refDDDrivingLicense"`
	Specialization              string                           `json:"Specialization" gorm:"column:refDDSpecialization"`
	DigitalSignature            string                           `json:"digital_signature" gorm:"column:refDDDigitalSignature"`
	DDEaseQTReportAccess        bool                             `json:"refDDEaseQTReportAccess" gorm:"column:refDDEaseQTReportAccess"`
	NPI                         string                           `json:"refDDNPI" gorm:"column:refDDNPI"`
	ProfileImgFile              *FileData                        `json:"profileImgFile" gorm:"-"`
	DriversLicenseFile          *FileData                        `json:"driversLicenseFile"  gorm:"-"`
	DigitalSignatureFile        *FileData                        `json:"digitalSignatureFile"  gorm:"-"`
	MedicalLicenseSecurity      []GetMedicalLicenseSecurityModel `json:"medicalLicenseSecurity" gorm:"-"`
	MalpracticeInsuranceDetails []MalpracticeModel               `json:"malpracticeinsureance_files" gorm:"-"`
	LicenseFiles                []LicenseFilesModel              `json:"licenseFiles" gorm:"-"`
}
