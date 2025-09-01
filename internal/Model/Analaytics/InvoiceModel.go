package model

type AmountModel struct {
	ScancenterAmount int `json:"refTAAmountScanCenter" gorm:"column:refTAAmountScanCenter"  mapstructure:"refTAAmountScanCenter"`
	UserAmount       int `json:"refTAAmountUser" gorm:"column:refTAAmountUser"  mapstructure:"refTAAmountUser"`
}

type ScanCenterModel struct {
	SCId     int    `json:"refSCId" gorm:"column:refSCId"`
	SCCustId string `json:"refSCCustId" gorm:"column:refSCCustId"`
	SCName   string `json:"refSCName" gorm:"column:refSCName"`
	Address  string `json:"refSCAddress" gorm:"column:refSCAddress"`
}

type UserModel struct {
	UserId     int    `json:"refUserId" gorm:"column:refUserId"`
	UserCustId string `json:"refUserCustId" gorm:"column:refUserCustId"`
}

type GetInvoiceDataReq struct {
	Type      int    `json:"type" mapstructure:"type"`
	UserId    int    `json:"userId" mapstructure:"userId"`
	Monthyear string `json:"monthnyear" mapstructure:"monthnyear"`
}

type GetCountScanCenterMonthModel struct {
	SCId              int `json:"refSCId" gorm:"column:refSCId"`
	TotalAppointments int `json:"total_appointments" gorm:"column:total_appointments"`
}

type GetUserModel struct {
	Name        string `json:"refUserFirstName" gorm:"column:refUserFirstName"`
	PhoneNumber string `json:"refCODOPhoneNo1" gorm:"column:refCODOPhoneNo1"`
	Email       string `json:"refCODOEmail" gorm:"column:refCODOEmail"`
}

type GetInvoiceDataReponse struct {
	AmountModel                              []AmountModel
	ScanCenterModel                          []ScanCenterModel
	GetCountScanCenterMonthModel             []GetCountScanCenterMonthModel
	GetUserModel                             []GetUserModel
	AdminOverallScanIndicatesAnalayticsModel []AdminOverallScanIndicatesAnalayticsModel
}

type GenerateInvoiceReq struct {
	SCId              int    `json:"refSCId" mapstructure:"refSCId"`
	UserId            int    `json:"refUserId" mapstructure:"refUserId"`
	FromId            int    `json:"fromId" mapstructure:"fromId"`
	FromName          string `json:"fromName" mapstructure:"fromName"`
	FromPhone         string `json:"fromPhone" mapstructure:"fromPhone"`
	FromEmail         string `json:"fromEmail" mapstructure:"fromEmail"`
	FromPan           string `json:"fromPan" mapstructure:"fromPan"`
	FromGst           string `json:"fromGst" mapstructure:"fromGst"`
	FromAddress       string `json:"fromAddress" mapstructure:"fromAddress"`
	ToId              int    `json:"toId" mapstructure:"toId"`
	ToName            string `json:"toName" mapstructure:"toName"`
	BillingFrom       string `json:"billingfrom" mapstructure:"billingfrom"`
	BillingTo         string `json:"billingto" mapstructure:"billingto"`
	ModeOfPayment     string `json:"modeofpayment" mapstructure:"modeofpayment"`
	UpiId             string `json:"upiId" mapstructure:"upiId"`
	AccountHolderName string `json:"accountHolderName" mapstructure:"accountHolderName"`
	AccountNumber     string `json:"accountNumber" mapstructure:"accountNumber"`
	Bank              string `json:"bank" mapstructure:"bank"`
	Branch            string `json:"branch" mapstructure:"branch"`
	IFSC              string `json:"ifsc" mapstructure:"ifsc"`
	Quantity          int    `json:"quantity" mapstructure:"quantity"`
	Amount            int    `json:"amount" mapstructure:"amount"`
	Total             int    `json:"total" mapstructure:"total"`
	ToAddress         string `json:"toAddress" mapstructure:"toAddress"`
	Signature         string `json:"signature" mapstructure:"signature"`
}

type GetInvoiceHistoryReq struct {
	Type   int `json:"type" mapstructure:"type"`
	UserId int `json:"id" mapstructure:"id"`
}

type FileData struct {
	Base64Data  string `json:"base64Data"`  // base64-encoded file content
	ContentType string `json:"contentType"` // e.g., "image/jpeg"
}

type InvoiceHistory struct {
	RefIHId                int       `json:"refIHId" gorm:"column:refIHId;primaryKey;autoIncrement"`
	RefSCId                *int      `json:"refSCId" gorm:"column:refSCId"`     // nullable
	RefUserId              *int      `json:"refUserId" gorm:"column:refUserId"` // nullable
	RefIHFromId            int       `json:"refIHFromId" gorm:"column:refIHFromId"`
	RefIHFromName          string    `json:"refIHFromName" gorm:"column:refIHFromName"`
	RefIHFromPhoneNo       string    `json:"refIHFromPhoneNo" gorm:"column:refIHFromPhoneNo"`
	RefIHFromEmail         string    `json:"refIHFromEmail" gorm:"column:refIHFromEmail"`
	RefIHFromPan           string    `json:"refIHFromPan" gorm:"column:refIHFromPan"`
	RefIHFromGST           string    `json:"refIHFromGST" gorm:"column:refIHFromGST"`
	RefIHFromAddress       string    `json:"refIHFromAddress" gorm:"column:refIHFromAddress"`
	RefIHToId              int       `json:"refIHToId" gorm:"column:refIHToId"`
	RefIHToName            string    `json:"refIHToName" gorm:"column:refIHToName"`
	RefIHFromDate          string    `json:"refIHFromDate" gorm:"column:refIHFromDate"` // use time.Time if preferred
	RefIHToDate            string    `json:"refIHToDate" gorm:"column:refIHToDate"`     // same here
	RefIHModePayment       string    `json:"refIHModePayment" gorm:"column:refIHModePayment"`
	RefIHUPIId             string    `json:"refIHUPIId" gorm:"column:refIHUPIId"`
	RefIHAccountHolderName string    `json:"refIHAccountHolderName" gorm:"column:refIHAccountHolderName"`
	RefIHAccountNo         string    `json:"refIHAccountNo" gorm:"column:refIHAccountNo"`
	RefIHAccountBank       string    `json:"refIHAccountBank" gorm:"column:refIHAccountBank"`
	RefIHAccountBranch     string    `json:"refIHAccountBranch" gorm:"column:refIHAccountBranch"`
	RefIHAccountIFSC       string    `json:"refIHAccountIFSC" gorm:"column:refIHAccountIFSC"`
	RefIHQuantity          int       `json:"refIHQuantity" gorm:"column:refIHQuantity"`
	RefIHAmount            int       `json:"refIHAmount" gorm:"column:refIHAmount"`
	RefIHTotal             int       `json:"refIHTotal" gorm:"column:refIHTotal"`
	RefIHCreatedAt         string    `json:"refIHCreatedAt" gorm:"column:refIHCreatedAt"` // or time.Time
	RefIHCreatedBy         int       `json:"refIHCreatedBy" gorm:"column:refIHCreatedBy"`
	RefIHToAddress         string    `json:"refIHToAddress" gorm:"column:refIHToAddress"`
	RefIHSignature         string    `json:"refIHSignature" gorm:"column:refIHSignature"`
	RefIHSignatureFile     *FileData `json:"refIHSignatureFile" gorm:"-"`
}

type InvoiceHistoryOverAll struct {
	RefIHId                int       `json:"refIHId" gorm:"column:refIHId;primaryKey;autoIncrement"`
	RefSCId                *int      `json:"refSCId" gorm:"column:refSCId"`     // nullable
	RefUserId              *int      `json:"refUserId" gorm:"column:refUserId"` // nullable
	RefIHFromId            int       `json:"refIHFromId" gorm:"column:refIHFromId"`
	RefIHFromName          string    `json:"refIHFromName" gorm:"column:refIHFromName"`
	RefIHFromPhoneNo       string    `json:"refIHFromPhoneNo" gorm:"column:refIHFromPhoneNo"`
	RefIHFromEmail         string    `json:"refIHFromEmail" gorm:"column:refIHFromEmail"`
	RefIHFromPan           string    `json:"refIHFromPan" gorm:"column:refIHFromPan"`
	RefIHFromGST           string    `json:"refIHFromGST" gorm:"column:refIHFromGST"`
	RefIHFromAddress       string    `json:"refIHFromAddress" gorm:"column:refIHFromAddress"`
	RefIHToId              int       `json:"refIHToId" gorm:"column:refIHToId"`
	RefIHToName            string    `json:"refIHToName" gorm:"column:refIHToName"`
	RefIHFromDate          string    `json:"refIHFromDate" gorm:"column:refIHFromDate"` // use time.Time if preferred
	RefIHToDate            string    `json:"refIHToDate" gorm:"column:refIHToDate"`     // same here
	RefIHModePayment       string    `json:"refIHModePayment" gorm:"column:refIHModePayment"`
	RefIHUPIId             string    `json:"refIHUPIId" gorm:"column:refIHUPIId"`
	RefIHAccountHolderName string    `json:"refIHAccountHolderName" gorm:"column:refIHAccountHolderName"`
	RefIHAccountNo         string    `json:"refIHAccountNo" gorm:"column:refIHAccountNo"`
	RefIHAccountBank       string    `json:"refIHAccountBank" gorm:"column:refIHAccountBank"`
	RefIHAccountBranch     string    `json:"refIHAccountBranch" gorm:"column:refIHAccountBranch"`
	RefIHAccountIFSC       string    `json:"refIHAccountIFSC" gorm:"column:refIHAccountIFSC"`
	RefIHQuantity          int       `json:"refIHQuantity" gorm:"column:refIHQuantity"`
	RefIHAmount            int       `json:"refIHAmount" gorm:"column:refIHAmount"`
	RefIHTotal             int       `json:"refIHTotal" gorm:"column:refIHTotal"`
	RefIHCreatedAt         string    `json:"refIHCreatedAt" gorm:"column:refIHCreatedAt"` // or time.Time
	RefIHCreatedBy         int       `json:"refIHCreatedBy" gorm:"column:refIHCreatedBy"`
	RefIHToAddress         string    `json:"refIHToAddress" gorm:"column:refIHToAddress"`
	RefIHSignature         string    `json:"refIHSignature" gorm:"column:refIHSignature"`
	RefIHSignatureFile     *FileData `json:"refIHSignatureFile" gorm:"-"`
	RefUserCustId          string    `json:"refUserCustId" gorm:"column:refUserCustId"`
	RefSCCustId            string    `json:"refSCCustId" gorm:"column:refSCCustId"`
}

type TakenDate struct {
	RefIHFromDate string `json:"refIHFromDate" gorm:"column:refIHFromDate"`
}

type GetInvoiceOverAllHistoryReq struct {
	FromDate string `json:"fromDate" mapstructure:"fromDate"`
	ToDate   string `json:"toDate" mapstructure:"toDate"`
}
