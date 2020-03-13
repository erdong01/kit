package model

// 广告模型
type Advertisement struct {
	AdvertisementId            int64                            `gorm:"primary_key" json:"advertisement_id,omitempty"`
	AdvertisementNo            int64                            `json:"advertisement_no,omitempty"`
	AdvertisementApplicationId int64                            `json:"advertisement_application_id,omitempty"`
	AdvertisementIndex         int64                            `json:"advertisement_index,omitempty"`
	LinkType                   int64                            `json:"link_type,omitempty"`
	AdvertisementImgUrl        string                           `json:"advertisement_img_url,omitempty"`
	AdvertisementLinkUrl       string                           `json:"advertisement_link_url,omitempty"`
	ApplicationUrl             string                           `json:"application_url,omitempty"`
	Status                     int64                            `json:"status,omitempty"`
	ApplicationPosition        AdvertisementApplicationPosition // 产品广告位置 一对一
	BaseModel
}

// 表名
func (Advertisement) TableName() string {
	return "rxt_advertisement"
}
