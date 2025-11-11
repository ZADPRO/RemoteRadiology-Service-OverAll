package model

type GetAllImpressionRecommendationCategory struct {
	RefIRCId    uint   `json:"refIRCId" gorm:"column:refIRCId"`
	RefIRCName  string `json:"refIRCName" gorm:"column:refIRCName"`
	RefIRCColor string `json:"refIRCColor" gorm:"column:refIRCColor"`
}

type CheckImpressionRecommendationCategory struct {
	Next_order_id int  `json:"next_order_id" gorm:"column:next_order_id"`
	CheckStatus   bool `json:"CheckStatus" gorm:"column:CheckStatus"`
}

type AddImpressionRecommendationReq struct {
	CategoryId                     int    `json:"categoryId" binding:"required" mapstructure:"categoryId"`
	SystemType                     string `json:"systemType" binding:"required" mapstructure:"systemType"`
	ImpressionRecommendationId     string `json:"impressionRecommendationId" binding:"required" mapstructure:"impressionRecommendationId"`
	ImpressionTextColor            string `json:"impressionTextColor" binding:"required" mapstructure:"impressionTextColor"`
	RecommendationTextColor        string `json:"recommedationTextColor" binding:"required" mapstructure:"recommedationTextColor"`
	ImpressionShortDescription     string `json:"impressionShortDescription" binding:"required" mapstructure:"impressionShortDescription"`
	RecommendationShortDescription string `json:"recommendationShortDescription" binding:"required" mapstructure:"recommendationShortDescription"`
	ImpressionLongDescription      string `json:"impressionLongDescription" binding:"required" mapstructure:"impressionLongDescription"`
	RecommendationLongDescription  string `json:"recommendationLongDescription" binding:"required" mapstructure:"recommendationLongDescription"`
}

type ImpressionRecommendationValModel struct {
	RefIRVId                      int    `json:"refIRVId" gorm:"column:refIRVId"`
	RefIRCId                      int    `json:"refIRCId" gorm:"column:refIRCId"`
	RefIRVOrderId                 int    `json:"refIRVOrderId" gorm:"column:refIRVOrderId"`
	RefIRVSystemType              string `json:"refIRVSystemType" gorm:"column:refIRVSystemType"`
	RefIRVCustId                  string `json:"refIRVCustId" gorm:"column:refIRVCustId"`
	RefIRVImpressionShortDesc     string `json:"refIRVImpressionShortDesc" gorm:"column:refIRVImpressionShortDesc"`
	RefIRVImpressionLongDesc      string `json:"refIRVImpressionLongDesc" gorm:"column:refIRVImpressionLongDesc"`
	RefIRVImpressionTextColor     string `json:"refIRVImpressionTextColor" gorm:"column:refIRVImpressionTextColor"`
	RefIRVRecommendationShortDesc string `json:"refIRVRecommendationShortDesc" gorm:"column:refIRVRecommendationShortDesc"`
	RefIRVRecommendationLongDesc  string `json:"refIRVRecommendationLongDesc" gorm:"column:refIRVRecommendationLongDesc"`
	RefIRVRecommendationTextColor string `json:"refIRVRecommendationTextColor" gorm:"column:refIRVRecommendationTextColor"`
	RefIRCName                    string `json:"refIRCName" gorm:"column:refIRCName"`
	RefIRCColor                   string `json:"refIRCColor" gorm:"column:refIRCColor"`
}

type UpdateOrderImpressionRecommendationReq struct {
	OrderData []OrderReq `json:"orderData" gorm:"column:orderData"`
}

type OrderReq struct {
	Id      int `json:"id" gorm:"column:id"`
	OrderId int `json:"orderId" gorm:"column:orderId"`
}

type UpdateImpressionRecommendationReq struct {
	Id                             int    `json:"id" binding:"required" mapstructure:"id"`
	CategoryId                     int    `json:"categoryId" binding:"required" mapstructure:"categoryId"`
	SystemType                     string `json:"systemType" binding:"required" mapstructure:"systemType"`
	ImpressionRecommendationId     string `json:"impressionRecommendationId" binding:"required" mapstructure:"impressionRecommendationId"`
	ImpressionTextColor            string `json:"impressionTextColor" binding:"required" mapstructure:"impressionTextColor"`
	RecommendationTextColor        string `json:"recommedationTextColor" binding:"required" mapstructure:"recommedationTextColor"`
	ImpressionShortDescription     string `json:"impressionShortDescription" binding:"required" mapstructure:"impressionShortDescription"`
	RecommendationShortDescription string `json:"recommendationShortDescription" binding:"required" mapstructure:"recommendationShortDescription"`
	ImpressionLongDescription      string `json:"impressionLongDescription" binding:"required" mapstructure:"impressionLongDescription"`
	RecommendationLongDescription  string `json:"recommendationLongDescription" binding:"required" mapstructure:"recommendationLongDescription"`
}

type GetReportFooterModel struct {
	RefFRId      int    `json:"refFRId" gorm:"column:refFRId"`
	RefFRContent string `json:"refFRContent" gorm:"column:refFRContent"`
}

type SaveReportFooterReq struct {
	ReportText string `json:"reportText" binding:"required" mapstructure:"reportText"`
}
