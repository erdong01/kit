package model

// 产品广告位置模型
type AdvertisementApplicationPosition struct {
	AdvertisementApplicationPositionId int64  `gorm:"primary_key" json:"advertisement_application_position_id,omitempty"`
	ApplicationId                      int8   `json:"application_id,omitempty"`
	PositionName                       string `json:"position_name,omitempty"`
	BaseModel
}

// 表名
func (AdvertisementApplicationPosition) TableName() string {
	return "rxt_advertisement_application_position"
}
