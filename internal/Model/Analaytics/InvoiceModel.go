package model

type AmountModel struct {
	TAId               int    `json:"refTAId" gorm:"column:refTAId"  mapstructure:"refTAId"`
	TASformEdit        string `json:"refTASformEdit" gorm:"column:refTASformEdit"  mapstructure:"refTASformEdit"`
	TASformCorrect     string `json:"refTASformCorrect" gorm:"column:refTASformCorrect"  mapstructure:"refTASformCorrect"`
	TADaformEdit       string `json:"refTADaformEdit" gorm:"column:refTADaformEdit"  mapstructure:"refTADaformEdit"`
	TADaformCorrect    string `json:"refTADaformCorrect" gorm:"column:refTADaformCorrect"  mapstructure:"refTADaformCorrect"`
	TADbformEdit       string `json:"refTADbformEdit" gorm:"column:refTADbformEdit"  mapstructure:"refTADbformEdit"`
	TADbformCorrect    string `json:"refTADbformCorrect" gorm:"column:refTADbformCorrect"  mapstructure:"refTADbformCorrect"`
	TADcformEdit       string `json:"refTADcformEdit" gorm:"column:refTADcformEdit"  mapstructure:"refTADcformEdit"`
	TADcformCorrect    string `json:"refTADcformCorrect" gorm:"column:refTADcformCorrect"  mapstructure:"refTADcformCorrect"`
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
	SFormEdit         int `json:"SFormEdit" gorm:"column:SFormEdit"`
	SFormCorrect      int `json:"SFormCorrect" gorm:"column:SFormCorrect"`
	DaFormEdit        int `json:"DaFormEdit" gorm:"column:DaFormEdit"`
	DaFormCorrect     int `json:"DaFormCorrect" gorm:"column:DaFormCorrect"`
	DbFormEdit        int `json:"DbFormEdit" gorm:"column:DbFormEdit"`
	DbFormCorrect     int `json:"DbFormCorrect" gorm:"column:DbFormCorrect"`
	DcFormEdit        int `json:"DcFormEdit" gorm:"column:DcFormEdit"`
	DcFormCorrect     int `json:"DcFormCorrect" gorm:"column:DcFormCorrect"`
}

type GenerateInvoiceReq struct {
	SCId                    int                  `json:"refSCId" mapstructure:"refSCId"`
	UserId                  int                  `json:"refUserId" mapstructure:"refUserId"`
	FromId                  int                  `json:"fromId" mapstructure:"fromId"`
	FromName                string               `json:"fromName" mapstructure:"fromName"`
	FromPhone               string               `json:"fromPhone" mapstructure:"fromPhone"`
	FromEmail               string               `json:"fromEmail" mapstructure:"fromEmail"`
	FromPan                 string               `json:"fromPan" mapstructure:"fromPan"`
	FromGst                 string               `json:"fromGst" mapstructure:"fromGst"`
	FromAddress             string               `json:"fromAddress" mapstructure:"fromAddress"`
	ToId                    int                  `json:"toId" mapstructure:"toId"`
	ToName                  string               `json:"toName" mapstructure:"toName"`
	BillingFrom             string               `json:"billingfrom" mapstructure:"billingfrom"`
	BillingTo               string               `json:"billingto" mapstructure:"billingto"`
	ModeOfPayment           string               `json:"modeofpayment" mapstructure:"modeofpayment"`
	UpiId                   string               `json:"upiId" mapstructure:"upiId"`
	AccountHolderName       string               `json:"accountHolderName" mapstructure:"accountHolderName"`
	AccountNumber           string               `json:"accountNumber" mapstructure:"accountNumber"`
	Bank                    string               `json:"bank" mapstructure:"bank"`
	Branch                  string               `json:"branch" mapstructure:"branch"`
	IFSC                    string               `json:"ifsc" mapstructure:"ifsc"`
	ToAddress               string               `json:"toAddress" mapstructure:"toAddress"`
	Signature               string               `json:"signature" mapstructure:"signature"`
	SformEditquantity       int                  `json:"refTASformEditquantity" mapstructure:"refTASformEditquantity"`
	SformEditamount         int                  `json:"refTASformEditamount" mapstructure:"refTASformEditamount"`
	SformCorrectquantity    int                  `json:"refTASformCorrectquantity" mapstructure:"refTASformCorrectquantity"`
	SformCorrectamount      int                  `json:"refTASformCorrectamount" mapstructure:"refTASformCorrectamount"`
	DaformEditquantity      int                  `json:"refTADaformEditquantity" mapstructure:"refTADaformEditquantity"`
	DaformEditamount        int                  `json:"refTADaformEditamount" mapstructure:"refTADaformEditamount"`
	DaformCorrectquantity   int                  `json:"refTADaformCorrectquantity" mapstructure:"refTADaformCorrectquantity"`
	DaformCorrectamount     int                  `json:"refTADaformCorrectamount" mapstructure:"refTADaformCorrectamount"`
	DbformEditquantity      int                  `json:"refTADbformEditquantity" mapstructure:"refTADbformEditquantity"`
	DbformEditamount        int                  `json:"refTADbformEditamount" mapstructure:"refTADbformEditamount"`
	DbformCorrectquantity   int                  `json:"refTADbformCorrectquantity" mapstructure:"refTADbformCorrectquantity"`
	DbformCorrectamount     int                  `json:"refTADbformCorrectamount" mapstructure:"refTADbformCorrectamount"`
	DcformEditquantity      int                  `json:"refTADcformEditquantity" mapstructure:"refTADcformEditquantity"`
	DcformEditamount        int                  `json:"refTADcformEditamount" mapstructure:"refTADcformEditamount"`
	DcformCorrectquantity   int                  `json:"refTADcformCorrectquantity" mapstructure:"refTADcformCorrectquantity"`
	DcformCorrectamount     int                  `json:"refTADcformCorrectamount" mapstructure:"refTADcformCorrectamount"`
	ScribeTotalcasequantity int                  `json:"refTADScribeTotalcasequantity" mapstructure:"refTADScribeTotalcasequantity"`
	ScribeTotalcaseamount   int                  `json:"refTADScribeTotalcaseamount" mapstructure:"refTADScribeTotalcaseamount"`
	OtherExpensiveName      string               `json:"otherExpensiveName" mapstructure:"otherExpensiveName"`
	OtherAmount             int                  `json:"otherAmount" mapstructure:"otherAmount"`
	ScanCenterTotalCase     int                  `json:"refScanCenterTotalCase" mapstructure:"refScanCenterTotalCase"`
	ScancentercaseAmount    int                  `json:"refScancentercaseAmount" mapstructure:"refScancentercaseAmount"`
	Total                   int                  `json:"total" mapstructure:"total"`
	OtherExpenses           []OtherExpensesModel `json:"otherExpenses" mapstructure:"otherExpenses"`
}

type OtherExpensesModel struct {
	Name   string `json:"name" gorm:"column:refOIAName" mapstructure:"name"`
	Amount int    `json:"amount" gorm:"column:refOIAAmount" mapstructure:"amount"`
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
	RefIHSformEditquantity       int                  `json:"refIHSformEditquantity" gorm:"column:refIHSformEditquantity"`
	RefIHSformEditamount         int                  `json:"refIHSformEditamount" gorm:"column:refIHSformEditamount"`
	RefIHSformCorrectquantity    int                  `json:"refIHSformCorrectquantity" gorm:"column:refIHSformCorrectquantity"`
	RefIHSformCorrectamount      int                  `json:"refIHSformCorrectamount" gorm:"column:refIHSformCorrectamount"`
	RefIHDaformEditquantity      int                  `json:"refIHDaformEditquantity" gorm:"column:refIHDaformEditquantity"`
	RefIHDaformEditamount        int                  `json:"refIHDaformEditamount" gorm:"column:refIHDaformEditamount"`
	RefIHDaformCorrectquantity   int                  `json:"refIHDaformCorrectquantity" gorm:"column:refIHDaformCorrectquantity"`
	RefIHDaformCorrectamount     int                  `json:"refIHDaformCorrectamount" gorm:"column:refIHDaformCorrectamount"`
	RefIHDbformEditquantity      int                  `json:"refIHDbformEditquantity" gorm:"column:refIHDbformEditquantity"`
	RefIHDbformEditamount        int                  `json:"refIHDbformEditamount" gorm:"column:refIHDbformEditamount"`
	RefIHDbformCorrectquantity   int                  `json:"refIHDbformCorrectquantity" gorm:"column:refIHDbformCorrectquantity"`
	RefIHDbformCorrectamount     int                  `json:"refIHDbformCorrectamount" gorm:"column:refIHDbformCorrectamount"`
	RefIHDcformEditquantity      int                  `json:"refIHDcformEditquantity" gorm:"column:refIHDcformEditquantity"`
	RefIHDcformEditamount        int                  `json:"refIHDcformEditamount" gorm:"column:refIHDcformEditamount"`
	RefIHDcformCorrectquantity   int                  `json:"refIHDcformCorrectquantity" gorm:"column:refIHDcformCorrectquantity"`
	RefIHDcformCorrectamount     int                  `json:"refIHDcformCorrectamount" gorm:"column:refIHDcformCorrectamount"`
	RefIHScribeTotalcasequantity int                  `json:"refIHScribeTotalcasequantity" gorm:"column:refIHScribeTotalcasequantity"`
	RefIHScribeTotalcaseamount   int                  `json:"refIHScribeTotalcaseamount" gorm:"column:refIHScribeTotalcaseamount"`
	RefIHOtherExpensiveName      string               `json:"refIHOtherExpensiveName" gorm:"column:refIHOtherExpensiveName"`
	RefIHOtherAmount             int                  `json:"refIHOtherAmount" gorm:"column:refIHOtherAmount"`
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
