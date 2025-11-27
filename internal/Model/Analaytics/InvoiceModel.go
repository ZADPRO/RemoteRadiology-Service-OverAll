package model

type AmountModel struct {
	TAId               int    `json:"refTAId" gorm:"column:refTAId"  mapstructure:"refTAId"`
	TASform            string `json:"refTASform" gorm:"column:refTASform"  mapstructure:"refTASform"`
	TADaform           string `json:"refTADaform" gorm:"column:refTADaform"  mapstructure:"refTADaform"`
	TADbform           string `json:"refTADbform" gorm:"column:refTADbform"  mapstructure:"refTADbform"`
	TADcform           string `json:"refTADcform" gorm:"column:refTADcform"  mapstructure:"refTADcform"`
	TAXform            string `json:"refTAXform" gorm:"column:refTAXform"  mapstructure:"refTAXform"`
	TAEditform         string `json:"refTAEditform" gorm:"column:refTAEditform"  mapstructure:"refTAEditform"`
	TADScribeTotalcase string `json:"refTADScribeTotalcase" gorm:"column:refTADScribeTotalcase"  mapstructure:"refTADScribeTotalcase"`
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
	RoleId      int    `json:"refRTId" gorm:"column:refRTId"`
	Pan         string `json:"refPan" gorm:"column:refPan"`
}

type GetInvoiceDataReponse struct {
	AmountModel                              []AmountModel
	ScanCenterModel                          []ScanCenterModel
	GetCountScanCenterMonthModel             []GetCountScanCenterMonthModel
	GetUserModel                             []GetUserModel
	AdminOverallScanIndicatesAnalayticsModel []UserTotalCase
}

type UserTotalCase struct {
	TotalAppointments int `json:"total_appointments" gorm:"column:total_appointments"`
	SForm             int `json:"SForm" gorm:"column:SForm"`
	DaForm            int `json:"DaForm" gorm:"column:DaForm"`
	DbForm            int `json:"DbForm" gorm:"column:DbForm"`
	DcForm            int `json:"DcForm" gorm:"column:DcForm"`
	XForm             int `json:"xForm" gorm:"column:xForm"`
	EditForm          int `json:"editForm" gorm:"column:editForm"`
}

type GenerateInvoiceReq struct {
	SCId                     int                  `json:"refSCId" mapstructure:"refSCId"`
	UserId                   int                  `json:"refUserId" mapstructure:"refUserId"`
	FromId                   int                  `json:"fromId" mapstructure:"fromId"`
	FromName                 string               `json:"fromName" mapstructure:"fromName"`
	FromPhone                string               `json:"fromPhone" mapstructure:"fromPhone"`
	FromEmail                string               `json:"fromEmail" mapstructure:"fromEmail"`
	FromPan                  string               `json:"fromPan" mapstructure:"fromPan"`
	FromGst                  string               `json:"fromGst" mapstructure:"fromGst"`
	FromAddress              string               `json:"fromAddress" mapstructure:"fromAddress"`
	ToId                     int                  `json:"toId" mapstructure:"toId"`
	ToName                   string               `json:"toName" mapstructure:"toName"`
	BillingFrom              string               `json:"billingfrom" mapstructure:"billingfrom"`
	BillingTo                string               `json:"billingto" mapstructure:"billingto"`
	ModeOfPayment            string               `json:"modeofpayment" mapstructure:"modeofpayment"`
	UpiId                    string               `json:"upiId" mapstructure:"upiId"`
	AccountHolderName        string               `json:"accountHolderName" mapstructure:"accountHolderName"`
	AccountNumber            string               `json:"accountNumber" mapstructure:"accountNumber"`
	Bank                     string               `json:"bank" mapstructure:"bank"`
	Branch                   string               `json:"branch" mapstructure:"branch"`
	IFSC                     string               `json:"ifsc" mapstructure:"ifsc"`
	ToAddress                string               `json:"toAddress" mapstructure:"toAddress"`
	Signature                string               `json:"signature" mapstructure:"signature"`
	SFormquantity            int                  `json:"refIHSFormquantity" mapstructure:"refIHSFormquantity"`
	SFormamount              int                  `json:"refIHSFormamount" mapstructure:"refIHSFormamount"`
	DaFormquantity           int                  `json:"refIHDaFormquantity" mapstructure:"refIHDaFormquantity"`
	DaFormamount             int                  `json:"refIHDaFormamount" mapstructure:"refIHDaFormamount"`
	DbFormquantity           int                  `json:"refIHDbFormquantity" mapstructure:"refIHDbFormquantity"`
	DbFormamount             int                  `json:"refIHDbFormamount" mapstructure:"refIHDbFormamount"`
	DcFormquantity           int                  `json:"refIHDcFormquantity" mapstructure:"refIHDcFormquantity"`
	DcFormamount             int                  `json:"refIHDcFormamount" mapstructure:"refIHDcFormamount"`
	XFormquantity            int                  `json:"refIHxFormquantity" mapstructure:"refIHxFormquantity"`
	XFormamount              int                  `json:"refIHxFormamount" mapstructure:"refIHxFormamount"`
	Editquantity             int                  `json:"refIHEditquantity" mapstructure:"refIHEditquantity"`
	EditFormamount           int                  `json:"refIHEditFormamount" mapstructure:"refIHEditFormamount"`
	ScribeTotalcasequantity  int                  `json:"refTADScribeTotalcasequantity" mapstructure:"refTADScribeTotalcasequantity"`
	ScribeTotalcaseamount    int                  `json:"refTADScribeTotalcaseamount" mapstructure:"refTADScribeTotalcaseamount"`
	ScanCenterTotalCase      int                  `json:"refScanCenterTotalCase" mapstructure:"refScanCenterTotalCase"`
	ScancentercaseAmount     int                  `json:"refScancentercaseAmount" mapstructure:"refScancentercaseAmount"`
	Total                    int                  `json:"total" mapstructure:"total"`
	OtherExpenses            []OtherExpensesModel `json:"otherExpenses" mapstructure:"otherExpenses"`
	DeductibleExpenses       []OtherExpensesModel `json:"deductibleExpenses" mapstructure:"deductibleExpenses"`
	OtherExpensesAmount      int                  `json:"otherExpensesAmount" mapstructure:"otherExpensesAmount"`
	DeductibleExpensesAmount int                  `json:"deductibleExpensesAmount" mapstructure:"deductibleExpensesAmount"`
}

type OtherExpensesModel struct {
	Name   string `json:"name" gorm:"column:refOIAName" mapstructure:"name"`
	Amount int    `json:"amount" gorm:"column:refOIAAmount" mapstructure:"amount"`
	Type   string `json:"type" gorm:"column:refOIAAmountType" mapstructure:"type"`
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
	RefIHId                      int                  `json:"refIHId" gorm:"column:refIHId;primaryKey;autoIncrement"`
	RefSCId                      *int                 `json:"refSCId" gorm:"column:refSCId"`     // nullable
	RefUserId                    *int                 `json:"refUserId" gorm:"column:refUserId"` // nullable
	RefIHFromId                  int                  `json:"refIHFromId" gorm:"column:refIHFromId"`
	RefIHFromName                string               `json:"refIHFromName" gorm:"column:refIHFromName"`
	RefIHFromPhoneNo             string               `json:"refIHFromPhoneNo" gorm:"column:refIHFromPhoneNo"`
	RefIHFromEmail               string               `json:"refIHFromEmail" gorm:"column:refIHFromEmail"`
	RefIHFromPan                 string               `json:"refIHFromPan" gorm:"column:refIHFromPan"`
	RefIHFromGST                 string               `json:"refIHFromGST" gorm:"column:refIHFromGST"`
	RefIHFromAddress             string               `json:"refIHFromAddress" gorm:"column:refIHFromAddress"`
	RefIHToId                    int                  `json:"refIHToId" gorm:"column:refIHToId"`
	RefIHToName                  string               `json:"refIHToName" gorm:"column:refIHToName"`
	RefIHFromDate                string               `json:"refIHFromDate" gorm:"column:refIHFromDate"` // use time.Time if preferred
	RefIHToDate                  string               `json:"refIHToDate" gorm:"column:refIHToDate"`     // same here
	RefIHModePayment             string               `json:"refIHModePayment" gorm:"column:refIHModePayment"`
	RefIHUPIId                   string               `json:"refIHUPIId" gorm:"column:refIHUPIId"`
	RefIHAccountHolderName       string               `json:"refIHAccountHolderName" gorm:"column:refIHAccountHolderName"`
	RefIHAccountNo               string               `json:"refIHAccountNo" gorm:"column:refIHAccountNo"`
	RefIHAccountBank             string               `json:"refIHAccountBank" gorm:"column:refIHAccountBank"`
	RefIHAccountBranch           string               `json:"refIHAccountBranch" gorm:"column:refIHAccountBranch"`
	RefIHAccountIFSC             string               `json:"refIHAccountIFSC" gorm:"column:refIHAccountIFSC"`
	RefIHCreatedAt               string               `json:"refIHCreatedAt" gorm:"column:refIHCreatedAt"` // or time.Time
	RefIHCreatedBy               int                  `json:"refIHCreatedBy" gorm:"column:refIHCreatedBy"`
	RefIHToAddress               string               `json:"refIHToAddress" gorm:"column:refIHToAddress"`
	RefIHSignature               string               `json:"refIHSignature" gorm:"column:refIHSignature"`
	RefIHSignatureFile           *FileData            `json:"refIHSignatureFile" gorm:"-"`
	RefIHSFormquantity           int                  `json:"refIHSFormquantity" gorm:"column:refIHSFormquantity"`
	RefIHSFormamount             int                  `json:"refIHSFormamount" gorm:"column:refIHSFormamount"`
	RefIHDaFormquantity          int                  `json:"refIHDaFormquantity" gorm:"column:refIHDaFormquantity"`
	RefIHDaFormamount            int                  `json:"refIHDaFormamount" gorm:"column:refIHDaFormamount"`
	RefIHDbFormquantity          int                  `json:"refIHDbFormquantity" gorm:"column:refIHDbFormquantity"`
	RefIHDbFormamount            int                  `json:"refIHDbFormamount" gorm:"column:refIHDbFormamount"`
	RefIHDcFormquantity          int                  `json:"refIHDcFormquantity" gorm:"column:refIHDcFormquantity"`
	RefIHDcFormamount            int                  `json:"refIHDcFormamount" gorm:"column:refIHDcFormamount"`
	RefIHxFormquantity           int                  `json:"refIHxFormquantity" gorm:"column:refIHxFormquantity"`
	RefIHxFormamount             int                  `json:"refIHxFormamount" gorm:"column:refIHxFormamount"`
	RefIHEditquantity            int                  `json:"refIHEditquantity" gorm:"column:refIHEditquantity"`
	RefIHEditFormamount          int                  `json:"refIHEditFormamount" gorm:"column:refIHEditFormamount"`
	RefIHScribeTotalcasequantity int                  `json:"refIHScribeTotalcasequantity" gorm:"column:refIHScribeTotalcasequantity"`
	RefIHScribeTotalcaseamount   int                  `json:"refIHScribeTotalcaseamount" gorm:"column:refIHScribeTotalcaseamount"`
	RefIHScanCenterTotalCase     int                  `json:"refIHScanCenterTotalCase" gorm:"column:refIHScanCenterTotalCase"`
	RefIHScancentercaseAmount    int                  `json:"refIHScancentercaseAmount" gorm:"column:refIHScancentercaseAmount"`
	RefIHTotal                   int                  `json:"refIHTotal" gorm:"column:refIHTotal"`
	RefRTId                      int                  `json:"refRTId" gorm:"column:refRTId"`
	OtherExpenses                []OtherExpensesModel `json:"otherExpenses" gorm:"-" mapstructure:"otherExpenses"`
	RefUserCustId                string               `json:"refUserCustId" gorm:"column:refUserCustId"`
	RefSCCustId                  string               `json:"refSCCustId" gorm:"column:refSCCustId"`
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
